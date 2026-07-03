package services

import (
	"fmt"
	"log"

	"github.com/fabianherreracruz/bre-b-pse-app/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SplitService struct {
	db *gorm.DB
}

// SplitConfig configuración para un split
type SplitConfig struct {
	BeneficiarioID uuid.UUID
	Percentage     float64
	EPaycoAccount  string
}

func NewSplitService(db *gorm.DB) *SplitService {
	return &SplitService{db: db}
}

// CalcularSplits calcula la distribución de un pago
func (s *SplitService) CalcularSplits(recaudoID uuid.UUID, totalAmount float64, splitConfigs []SplitConfig) ([]models.Split, error) {
	var splits []models.Split

	for _, config := range splitConfigs {
		amount := (totalAmount * config.Percentage) / 100
		
		split := models.Split{
			ID:            uuid.New(),
			RecaudoID:     recaudoID,
			BeneficiarioID: config.BeneficiarioID,
			Percentage:    config.Percentage,
			Amount:        amount,
			Status:        "pendiente",
		}

		splits = append(splits, split)
	}

	// Validar que suma sea 100%
	totalPercentage := 0.0
	for _, split := range splits {
		totalPercentage += split.Percentage
	}

	if totalPercentage != 100.0 {
		return nil, fmt.Errorf("total percentage must be 100, got %.2f", totalPercentage)
	}

	return splits, nil
}

// GuardarSplits guarda los splits en la base de datos
func (s *SplitService) GuardarSplits(splits []models.Split) error {
	for _, split := range splits {
		if err := s.db.Create(&split).Error; err != nil {
			log.Printf("Error saving split: %v", err)
			return err
		}
	}
	return nil
}

// ObtenerSplits obtiene los splits de un recaudo
func (s *SplitService) ObtenerSplits(recaudoID uuid.UUID) ([]models.Split, error) {
	var splits []models.Split
	if err := s.db.Where("recaudo_id = ?", recaudoID).Find(&splits).Error; err != nil {
		log.Printf("Error fetching splits: %v", err)
		return nil, err
	}
	return splits, nil
}

// MarcarSplitPagado marca un split como pagado
func (s *SplitService) MarcarSplitPagado(splitID uuid.UUID) error {
	if err := s.db.Model(&models.Split{}).Where("id = ?", splitID).Updates(map[string]interface{}{
		"status": "pagado",
		"paid_at": gorm.Expr("NOW()"),
	}).Error; err != nil {
		log.Printf("Error updating split: %v", err)
		return err
	}
	return nil
}

// ObtenerSplitsPendientes obtiene los splits pendientes de pago
func (s *SplitService) ObtenerSplitsPendientes() ([]models.Split, error) {
	var splits []models.Split
	if err := s.db.Where("status = ?", "pendiente").Find(&splits).Error; err != nil {
		log.Printf("Error fetching pending splits: %v", err)
		return nil, err
	}
	return splits, nil
}

// ProcessarSplitsPago procesa automáticamente los pagos de splits
func (s *SplitService) ProcessarSplitsPago(epaycoService *EPaycoService) error {
	splits, err := s.ObtenerSplitsPendientes()
	if err != nil {
		return err
	}

	for _, split := range splits {
		// Obtener usuario beneficiario
		var beneficiario models.User
		if err := s.db.Where("id = ?", split.BeneficiarioID).First(&beneficiario).Error; err != nil {
			log.Printf("Error fetching beneficiary: %v", err)
			continue
		}

		// Crear pago en ePayco
		reference := fmt.Sprintf("SPLIT-%s", split.ID.String())
		paymentReq := CreatePaymentRequest{
			Reference:   reference,
			Amount:      split.Amount,
			Currency:    "COP",
			Email:       beneficiario.Email,
			Phone:       beneficiario.Phone,
			Description: fmt.Sprintf("Split payment - %.2f", split.Amount),
			Sandbox:     true,
		}

		_, err := epaycoService.CreatePayment(paymentReq)
		if err != nil {
			log.Printf("Error creating payment for split %s: %v", split.ID, err)
			continue
		}

		// Marcar como pagado
		if err := s.MarcarSplitPagado(split.ID); err != nil {
			log.Printf("Error marking split as paid: %v", err)
			continue
		}

		log.Printf("✅ Split payment processed: %s (%.2f)", split.ID, split.Amount)
	}

	return nil
}
