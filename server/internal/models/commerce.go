package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CartItem struct {
	ID        uuid.UUID `db:"id"`
	CartID    uuid.UUID `db:"cart_id"`
	ProductID uuid.UUID `db:"product_id"`
	Quantity  int       `db:"quantity"`
	UnitPrice float64   `db:"unit_price"`
	LineTotal float64   `db:"line_total"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Product   *Product  `db:"-"`
}

type Order struct {
	ID               uuid.UUID       `db:"id"`
	UserID           uuid.UUID       `db:"user_id"`
	Status           string          `db:"status"`
	Currency         string          `db:"currency"`
	Subtotal         float64         `db:"subtotal"`
	TaxAmount        float64         `db:"tax_amount"`
	ShippingAmount   float64         `db:"shipping_amount"`
	DiscountAmount   float64         `db:"discount_amount"`
	GrandTotal       float64         `db:"grand_total"`
	ShippingMetadata json.RawMessage `db:"shipping_metadata"`
	PaymentID        *uuid.UUID      `db:"payment_id"`
	Notes            *string         `db:"notes"`
	CreatedAt        time.Time       `db:"created_at"`
	UpdatedAt        time.Time       `db:"updated_at"`
	Items            []*OrderItem    `db:"-"`
}

type OrderItem struct {
	ID        uuid.UUID `db:"id"`
	OrderID   uuid.UUID `db:"order_id"`
	ProductID uuid.UUID `db:"product_id"`
	Quantity  int       `db:"quantity"`
	UnitPrice float64   `db:"unit_price"`
	LineTotal float64   `db:"line_total"`
	CreatedAt time.Time `db:"created_at"`
	Product   *Product  `db:"-"`
}

type Payment struct {
	ID                uuid.UUID       `db:"id"`
	OrderID           uuid.UUID       `db:"order_id"`
	Provider          string          `db:"provider"`
	ProviderOrderID   string          `db:"provider_order_id"`
	ProviderPaymentID *string         `db:"provider_payment_id"`
	Status            string          `db:"status"`
	Amount            float64         `db:"amount"`
	Currency          string          `db:"currency"`
	Attempts          int             `db:"attempts"`
	Metadata          json.RawMessage `db:"metadata"`
	CreatedAt         time.Time       `db:"created_at"`
	UpdatedAt         time.Time       `db:"updated_at"`
}

type Transaction struct {
	ID              uuid.UUID       `db:"id"`
	PaymentID       uuid.UUID       `db:"payment_id"`
	EventType       string          `db:"event_type"`
	ProviderEventID *string         `db:"provider_event_id"`
	Amount          *float64        `db:"amount"`
	Currency        *string         `db:"currency"`
	Status          string          `db:"status"`
	RawPayload      json.RawMessage `db:"raw_payload"`
	CreatedAt       time.Time       `db:"created_at"`
}

type WebhookEvent struct {
	ID          uuid.UUID       `db:"id"`
	Provider    string          `db:"provider"`
	EventID     string          `db:"event_id"`
	EventType   string          `db:"event_type"`
	Signature   *string         `db:"signature"`
	Payload     json.RawMessage `db:"payload"`
	Processed   bool            `db:"processed"`
	ProcessedAt *time.Time      `db:"processed_at"`
	Retries     int             `db:"retries"`
	LastError   *string         `db:"last_error"`
	CreatedAt   time.Time       `db:"created_at"`
}
