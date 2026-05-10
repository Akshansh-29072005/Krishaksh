package models

import (
	"time"

	"github.com/google/uuid"
)

// Notification stores a push notification record per user
type Notification struct {
	ID          uuid.UUID  `db:"id"`
	UserID      uuid.UUID  `db:"user_id"`
	Type        string     `db:"type"` // scan_completed, ticket_updated, disease_alert, weather_alert, promotion, order_update
	Title       string     `db:"title"`
	Body        string     `db:"body"`
	IsRead      bool       `db:"is_read"`
	DeliveredAt *time.Time `db:"delivered_at"`
	Payload     *string    `db:"payload"` // JSON string for deep-link data
	CreatedAt   time.Time  `db:"created_at"`
}

// DeviceToken stores FCM push tokens per device
type DeviceToken struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Token     string    `db:"token"`
	Platform  string    `db:"platform"` // android | ios
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
