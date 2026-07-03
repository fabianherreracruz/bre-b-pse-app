package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fabianherreracruz/bre-b-pse-app/internal/models"
	"github.com/fabianherreracruz/bre-b-pse-app/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentHandler struct {
	db                *gorm.DB
	epaycoService     *services.EPaycoService
	notifyService     *services.NotificationService
	splitService      *services.SplitService
}

type CreatePaymentRequest struct {
	Amount        float64                   `json:"amount" binding:"required"`
	Description   string                   `json:"description" binding:"required"`
	CustomerEmail string                   `json:"customer_email" binding:"required"`
	CustomerPhone string                   `json:"customer_phone" binding:"required"`
	Splits        []CreateSplitRequest      `json:"splits" binding:"required"`
}

type CreateSplitRequest struct {
	BeneficiarioID string  `json:"beneficiario_id" binding:"required"`
	Percentage     float64 `json:"percentage" binding:"required"`
}

type PaymentResponse struct {
	ID         uuid.UUID `json:"id"`
	PaymentURL string    `json:"payment_url"`
	Reference  string    `json:"reference"`
	Status     string    `json:"status"`
}

func NewPaymentHandler(
	db *gorm.DB,
	epaycoService *services.EPaycoService,
	notifyService *services.NotificationService,
	splitService *services.SplitService,
) *PaymentHandler {
	return &PaymentHandler{
		db:            db,
		epaycoService: epaycoService,
		notifyService: notifyService,
		splitService:  splitService,
	}
}

// CreatePayment crea un nuevo recaudo
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener usuario del contexto
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Crear recaudo
	recaudoID := uuid.New()
	reference := fmt.Sprintf("REC-%d-%s", recaudoID, uuid.New().String()[:8])

	recaudo := models.Recaudo{
		ID:            recaudoID,
		RecaudadorID:  userID.(uuid.UUID),
		ReferenceCode: reference,
		Amount:        req.Amount,
		Currency:      "COP",
		PaymentMethod: "PSE",
		Status:        "pendiente",
		Description:   req.Description,
		CustomerEmail: req.CustomerEmail,
		CustomerPhone: req.CustomerPhone,
	}

	if err := h.db.Create(&recaudo).Error; err != nil {
		log.Printf("Error creating recaudo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create payment"})
		return
	}

	// Crear splits
	var splitConfigs []services.SplitConfig
	for _, split := range req.Splits {
		beneficiarioID, err := uuid.Parse(split.BeneficiarioID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid beneficiary id"})
			return
		}

		// Obtener cuenta ePayco del beneficiario
		var beneficiario models.User
		if err := h.db.Where("id = ?", beneficiarioID).First(&beneficiario).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "beneficiary not found"})
			return
		}

		splitConfigs = append(splitConfigs, services.SplitConfig{
			BeneficiarioID: beneficiarioID,
			Percentage:     split.Percentage,
			EPaycoAccount:  beneficiario.EPaycoAccount,
		})
	}

	// Calcular y guardar splits
	splits, err := h.splitService.CalcularSplits(recaudoID, req.Amount, splitConfigs)
	if err != nil {
		log.Printf("Error calculating splits: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.splitService.GuardarSplits(splits); err != nil {
		log.Printf("Error saving splits: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save splits"})
		return
	}

	// Crear pago en ePayco
	paymentReq := services.CreatePaymentRequest{
		Reference:   reference,
		Amount:      req.Amount,
		Currency:    "COP",
		Email:       req.CustomerEmail,
		Phone:       req.CustomerPhone,
		Description: req.Description,
		Sandbox:     true,
	}

	paymentResp, err := h.epaycoService.CreatePayment(paymentReq)
	if err != nil {
		log.Printf("Error creating payment in ePayco: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process payment"})
		return
	}

	// Guardar referencia de ePayco
	h.db.Model(&recaudo).Update("epayco_reference", paymentResp.EPaycoReference)

	// Responder
	response := PaymentResponse{
		ID:         recaudoID,
		PaymentURL: paymentResp.PaymentURL,
		Reference:  reference,
		Status:     "pendiente",
	}

	c.JSON(http.StatusCreated, response)
}

// VerifyPayment verifica el estado de un pago
func (h *PaymentHandler) VerifyPayment(c *gin.Context) {
	reference := c.Param("reference")

	// Obtener recaudo
	var recaudo models.Recaudo
	if err := h.db.Where("reference_code = ?", reference).First(&recaudo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}

	// Verificar con ePayco
	paymentData, err := h.epaycoService.VerifyPayment(reference)
	if err != nil {
		log.Printf("Error verifying payment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify payment"})
		return
	}

	c.JSON(http.StatusOK, paymentData)
}

// WebhookPayment maneja webhooks de ePayco
func (h *PaymentHandler) WebhookPayment(c *gin.Context) {
	// Obtener referencia del webhook
	reference := c.PostForm("reference")
	status := c.PostForm("x_response")
	approvalCode := c.PostForm("x_approval_code")

	log.Printf("Webhook received - Reference: %s, Status: %s", reference, status)

	// Obtener recaudo
	var recaudo models.Recaudo
	if err := h.db.Where("reference_code = ?", reference).First(&recaudo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment not found"})
		return
	}

	// Actualizar estado
	if status == "1" || status == "aceptada" {
		recaudo.Status = "completado"

		// Procesar splits
		splits, err := h.splitService.ObtenerSplits(recaudo.ID)
		if err != nil {
			log.Printf("Error getting splits: %v", err)
		}

		for _, split := range splits {
			// Crear pago de split en ePayco
			var beneficiario models.User
			if err := h.db.Where("id = ?", split.BeneficiarioID).First(&beneficiario).Error; err != nil {
				log.Printf("Error fetching beneficiary: %v", err)
				continue
			}

			paymentReq := services.CreatePaymentRequest{
				Reference:   fmt.Sprintf("SPLIT-%s", split.ID.String()),
				Amount:      split.Amount,
				Currency:    "COP",
				Email:       beneficiario.Email,
				Phone:       beneficiario.Phone,
				Description: fmt.Sprintf("Split payment from %s", reference),
				Sandbox:     true,
			}

			if _, err := h.epaycoService.CreatePayment(paymentReq); err != nil {
				log.Printf("Error creating split payment: %v", err)
				continue
			}

			// Marcar split como pagado
			h.splitService.MarcarSplitPagado(split.ID)

			// Notificar beneficiario
			mensaje := services.GenerarMensajeSplitPago(
				beneficiario.FirstName,
				split.Amount,
				recaudo.CreatedAt.String(),
			)
			h.notifyService.SendMultiChannel(beneficiario.Email, beneficiario.Phone, "Split Payment", mensaje)
		}

		// Notificar cliente de recaudo exitoso
		mensaje := services.GenerarMensajeRecaudoExitoso("Cliente", recaudo.Amount, reference)
		h.notifyService.SendMultiChannel(recaudo.CustomerEmail, recaudo.CustomerPhone, "Recaudo Exitoso", mensaje)

	} else {
		recaudo.Status = "fallido"

		// Notificar cliente de recaudo fallido
		mensaje := services.GenerarMensajeRecaudoFallido("Cliente", recaudo.Amount, reference, "Pago rechazado")
		h.notifyService.SendMultiChannel(recaudo.CustomerEmail, recaudo.CustomerPhone, "Recaudo Fallido", mensaje)
	}

	// Guardar aprobación
	if approvalCode != "" {
		h.db.Model(&recaudo).Update("epayco_reference", approvalCode)
	}

	// Actualizar recaudo
	if err := h.db.Save(&recaudo).Error; err != nil {
		log.Printf("Error updating recaudo: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetPaymentStatus obtiene el estado de un pago
func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	recaudoID := c.Param("id")

	var recaudo models.Recaudo
	if err := h.db.Where("id = ?", recaudoID).First(&recaudo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}

	c.JSON(http.StatusOK, recaudo)
}
