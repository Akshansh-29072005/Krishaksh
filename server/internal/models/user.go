package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a farmer, vendor, or admin
type User struct {
	ID          uuid.UUID `db:"id"`
	GoogleID    *string   `db:"google_id"`
	FullName    string    `db:"full_name"`
	Email       *string   `db:"email"`
	PhoneNumber *string   `db:"phone_number"`
	DeviceToken *string   `db:"device_token"`
	RoleID      int       `db:"role_id"`
	Village     *string   `db:"village"`
	Language    string    `db:"language"` // e.g. "en", "hi"
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Role manages RBAC system
type Role struct {
	ID          int      `db:"id"`
	Name        string   `db:"name"` // FARMER, VENDOR, ADMIN
	Permissions []string `db:"permissions"`
}

// RefreshToken stores stateful device session tokens
type RefreshToken struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
	Revoked   bool      `db:"revoked"`
}
