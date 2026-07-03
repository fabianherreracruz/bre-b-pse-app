package services

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/fabianherreracruz/bre-b-pse-app/internal/config"
)

type EPaycoService struct {
	cfg *config.Config
}

// Estructura para crear un pago
type CreatePaymentRequest struct {
	Reference   string  `json:"reference"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	Description string  `json:"description"`
	Sandbox     bool    `json:"sandbox"`
}

// Respuesta de ePayco
type EPaycoResponse struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// PaymentResponse estructura de respuesta de pago
type PaymentResponse struct {
	Reference       string `json:"reference"`
	SessionID       string `json:"session_id"`
	URL             string `json:"url"`
	PaymentURL      string `json:"payment_url"`
	ProcessedCode   string `json:"processed_code"`
	ApprovalCode    string `json:"approval_code"`
	TransactionID   string `json:"transaction_id"`
	EPaycoReference string `json:"epayco_reference"`
}

func NewEPaycoService(cfg *config.Config) *EPaycoService {
	return &EPaycoService{cfg: cfg}
}

// CreatePayment crea un pago en ePayco
func (e *EPaycoService) CreatePayment(req CreatePaymentRequest) (*PaymentResponse, error) {
	// Generar firma
	signature := e.generateSignature(req.Reference, fmt.Sprintf("%.2f", req.Amount))

	payload := map[string]interface{}{
		"clientId":   e.cfg.EPaycoClientID,
		"reference":  req.Reference,
		"amount":     req.Amount,
		"currency":   req.Currency,
		"email":      req.Email,
		"phone":      req.Phone,
		"description": req.Description,
		"signature":  signature,
		"sandbox":    req.Sandbox,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling payload: %v", err)
		return nil, err
	}

	// Hacer request a ePayco
	resp, err := http.Post(
		"https://api.epayco.co/payment/create/transaction",
		"application/json",
		bytes.NewBuffer(jsonPayload),
	)
	if err != nil {
		log.Printf("Error calling ePayco API: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return nil, err
	}

	var epaycoResp EPaycoResponse
	if err := json.Unmarshal(body, &epaycoResp); err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}

	if epaycoResp.Status != 1 {
		return nil, fmt.Errorf("ePayco error: %s", epaycoResp.Message)
	}

	// Parsear respuesta
	dataBytes, _ := json.Marshal(epaycoResp.Data)
	var paymentResp PaymentResponse
	if err := json.Unmarshal(dataBytes, &paymentResp); err != nil {
		log.Printf("Error parsing payment response: %v", err)
		return nil, err
	}

	return &paymentResp, nil
}

// VerifyPayment verifica un pago en ePayco
func (e *EPaycoService) VerifyPayment(reference string) (map[string]interface{}, error) {
	url := fmt.Sprintf(
		"https://api.epayco.co/payment/query/reference_%s/?public_key=%s",
		reference,
		e.cfg.EPaycoAPIKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error verifying payment: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}

	return result, nil
}

// generateSignature genera la firma para autenticación con ePayco
func (e *EPaycoService) generateSignature(reference, amount string) string {
	message := reference + amount + e.cfg.EPaycoPrivateKey
	hash := sha256.Sum256([]byte(message))
	return hex.EncodeToString(hash[:])
}

// ValidateWebhook valida que el webhook sea de ePayco
func (e *EPaycoService) ValidateWebhook(signature string, data string) bool {
	expectedSignature := e.generateSignature(data, "")
	return signature == expectedSignature
}
