package services

import (
	"fmt"
	"log"

	"github.com/fabianherreracruz/bre-b-pse-app/internal/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/twilio/twilio-go"
	tilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type NotificationService struct {
	cfg           *config.Config
	tilioClient  *twilio.RestClient
	sendgridClient *sendgrid.Client
}

func NewNotificationService(cfg *config.Config) *NotificationService {
	tilioClient := twilio.NewRestClient()
	sendgridClient := sendgrid.NewSendClient(cfg.SendGridAPIKey)

	return &NotificationService{
		cfg:            cfg,
		tilioClient:   tilioClient,
		sendgridClient: sendgridClient,
	}
}

// SendEmail envía un correo electrónico
func (n *NotificationService) SendEmail(to, subject, htmlContent string) error {
	from := mail.NewEmail("Recaudos BRE-B PSE", n.cfg.SendGridFromEmail)
	toMail := mail.NewEmail("", to)
	message := mail.NewSingleEmail(from, subject, toMail, subject, htmlContent)

	response, err := n.sendgridClient.Send(message)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	if response.StatusCode > 299 {
		return fmt.Errorf("SendGrid error: status %d", response.StatusCode)
	}

	log.Printf("✅ Email sent to %s", to)
	return nil
}

// SendSMS envía un SMS
func (n *NotificationService) SendSMS(to, message string) error {
	params := &tilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(n.cfg.TwilioPhoneNumber)
	params.SetBody(message)

	resp, err := n.tilioClient.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Error sending SMS: %v", err)
		return err
	}

	if resp.Sid == nil {
		return fmt.Errorf("failed to send SMS")
	}

	log.Printf("✅ SMS sent to %s (SID: %s)", to, *resp.Sid)
	return nil
}

// SendWhatsApp envía un mensaje por WhatsApp
func (n *NotificationService) SendWhatsApp(to, message string) error {
	params := &tilioApi.CreateMessageParams{}
	params.SetTo(fmt.Sprintf("whatsapp:%s", to))
	params.SetFrom(n.cfg.TwilioWhatsAppNumber)
	params.SetBody(message)

	resp, err := n.tilioClient.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Error sending WhatsApp: %v", err)
		return err
	}

	if resp.Sid == nil {
		return fmt.Errorf("failed to send WhatsApp message")
	}

	log.Printf("✅ WhatsApp sent to %s (SID: %s)", to, *resp.Sid)
	return nil
}

// Envía notificaciones por múltiples canales
func (n *NotificationService) SendMultiChannel(to, phone, subject, message string) error {
	// Email
	if to != "" {
		if err := n.SendEmail(to, subject, message); err != nil {
			log.Printf("Failed to send email: %v", err)
		}
	}

	// SMS
	if phone != "" {
		if err := n.SendSMS(phone, message); err != nil {
			log.Printf("Failed to send SMS: %v", err)
		}
	}

	// WhatsApp
	if phone != "" {
		if err := n.SendWhatsApp(phone, message); err != nil {
			log.Printf("Failed to send WhatsApp: %v", err)
		}
	}

	return nil
}

// GenerarMensajeRecaudoExitoso genera el mensaje para recaudo exitoso
func GenerarMensajeRecaudoExitoso(nombre string, monto float64, referencia string) string {
	return fmt.Sprintf(
		`¡Hola %s!

Tu recaudo ha sido procesado exitosamente.

Detalles:
- Monto: $%.2f
- Referencia: %s
- Estado: Completado

Gracias por usar nuestro servicio.`,
		nombre, monto, referencia,
	)
}

// GenerarMensajeRecaudoFallido genera el mensaje para recaudo fallido
func GenerarMensajeRecaudoFallido(nombre string, monto float64, referencia string, razon string) string {
	return fmt.Sprintf(
		`¡Hola %s!

Lamentablemente, tu recaudo no pudo ser procesado.

Detalles:
- Monto: $%.2f
- Referencia: %s
- Razón: %s
- Estado: Fallido

Por favor, intenta nuevamente o contacta al soporte.`,
		nombre, monto, referencia, razon,
	)
}

// GenerarMensajeSplitPago genera el mensaje del split de pago
func GenerarMensajeSplitPago(beneficiario string, monto float64, fecha string) string {
	return fmt.Sprintf(
		`¡Hola %s!

Recibiste un pago del recaudo BRE-B PSE.

Detalles:
- Monto: $%.2f
- Fecha: %s
- Estado: Procesado

El dinero será transferido a tu cuenta en las próximas 24-48 horas.

Gracias.`,
		beneficiario, monto, fecha,
	)
}
