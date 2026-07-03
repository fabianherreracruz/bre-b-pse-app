package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// User modelo para usuarios/recaudadores
type User struct {
	ID            uuid.UUID      `db:"id" json:"id"`
	Email         string         `db:"email" json:"email"`
	Phone         string         `db:"phone" json:"phone"`
	Password      string         `db:"password" json:"-"`
	FirstName     string         `db:"first_name" json:"first_name"`
	LastName      string         `db:"last_name" json:"last_name"`
	Role          string         `db:"role" json:"role"` // admin, recaudador, beneficiario
	Status        string         `db:"status" json:"status"`
	EPaycoAccount string         `db:"epayco_account" json:"epayco_account"`
	CreatedAt     time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at" json:"updated_at"`
}

// Recaudo modelo para recaudos
type Recaudo struct {
	ID              uuid.UUID      `db:"id" json:"id"`
	RecaudadorID    uuid.UUID      `db:"recaudador_id" json:"recaudador_id"`
	ReferenceCode   string         `db:"reference_code" json:"reference_code"`
	Amount          float64        `db:"amount" json:"amount"`
	Currency        string         `db:"currency" json:"currency"`
	PaymentMethod   string         `db:"payment_method" json:"payment_method"` // PSE, tarjeta, etc
	Status          string         `db:"status" json:"status"`                 // pendiente, completado, fallido
	Description     string         `db:"description" json:"description"`
	EPaycoReference string         `db:"epayco_reference" json:"epayco_reference"`
	CustomerEmail   string         `db:"customer_email" json:"customer_email"`
	CustomerPhone   string         `db:"customer_phone" json:"customer_phone"`
	CreatedAt       time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at" json:"updated_at"`
	CompletedAt     sql.NullTime   `db:"completed_at" json:"completed_at"`
}

// Split modelo para distribución de pagos
type Split struct {
	ID            uuid.UUID      `db:"id" json:"id"`
	RecaudoID     uuid.UUID      `db:"recaudo_id" json:"recaudo_id"`
	BeneficiarioID uuid.UUID     `db:"beneficiario_id" json:"beneficiario_id"`
	Percentage    float64        `db:"percentage" json:"percentage"`
	Amount        float64        `db:"amount" json:"amount"`
	Status        string         `db:"status" json:"status"` // pendiente, pagado
	PaidAt        sql.NullTime   `db:"paid_at" json:"paid_at"`
	CreatedAt     time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at" json:"updated_at"`
}

// Notificacion modelo para rastrear notificaciones
type Notificacion struct {
	ID        uuid.UUID      `db:"id" json:"id"`
	RecaudoID uuid.UUID      `db:"recaudo_id" json:"recaudo_id"`
	UserID    uuid.UUID      `db:"user_id" json:"user_id"`
	Type      string         `db:"type" json:"type"` // email, sms, whatsapp
	Recipient string         `db:"recipient" json:"recipient"`
	Status    string         `db:"status" json:"status"` // pendiente, enviado, fallido
	Message   string         `db:"message" json:"message"`
	Error     sql.NullString `db:"error" json:"error"`
	SentAt    sql.NullTime   `db:"sent_at" json:"sent_at"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
}

// AuditLog para auditoría
type AuditLog struct {
	ID        uuid.UUID      `db:"id" json:"id"`
	UserID    uuid.UUID      `db:"user_id" json:"user_id"`
	Action    string         `db:"action" json:"action"`
	Entity    string         `db:"entity" json:"entity"`
	EntityID  uuid.UUID      `db:"entity_id" json:"entity_id"`
	OldValues string         `db:"old_values" json:"old_values"`
	NewValues string         `db:"new_values" json:"new_values"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
}
