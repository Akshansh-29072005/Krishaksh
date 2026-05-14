package models

import (
	"time"

	"github.com/google/uuid"
)

// SupportTicket represents a farmer's help request
type SupportTicket struct {
	ID                uuid.UUID  `db:"id" json:"id"`
	UserID            uuid.UUID  `db:"user_id" json:"user_id"`
	Title             string     `db:"title" json:"title"`
	Description       string     `db:"description" json:"description"`
	Status            string     `db:"status" json:"status"`                     // open, under_review, awaiting_user, resolved, escalated, closed
	Priority          string     `db:"priority" json:"priority"`                 // low, medium, high, critical
	AssignedTo        *uuid.UUID `db:"assigned_to" json:"assigned_to,omitempty"` // admin user_id
	CallbackRequested bool       `db:"callback_requested" json:"callback_requested"`
	CallbackStatus    string     `db:"callback_status" json:"callback_status"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at" json:"updated_at"`
	ResolvedAt        *time.Time `db:"resolved_at" json:"resolved_at,omitempty"`
}

// SupportMessage is a single chat-like entry in a ticket thread
type SupportMessage struct {
	ID         uuid.UUID `db:"id" json:"id"`
	TicketID   uuid.UUID `db:"ticket_id" json:"ticket_id"`
	SenderID   uuid.UUID `db:"sender_id" json:"sender_id"`
	SenderRole string    `db:"sender_role" json:"sender_role"` // farmer | admin
	Body       string    `db:"body" json:"body"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

// SupportAttachment holds file metadata uploaded within a ticket (image or voice note)
type SupportAttachment struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	TicketID      uuid.UUID  `db:"ticket_id" json:"ticket_id"`
	MessageID     *uuid.UUID `db:"message_id" json:"message_id,omitempty"` // optional – link to a specific message
	UploaderID    uuid.UUID  `db:"uploader_id" json:"uploader_id"`
	FileURL       string     `db:"file_url" json:"file_url"`
	FileType      string     `db:"file_type" json:"file_type"` // image | voice
	FileSizeBytes int64      `db:"file_size_bytes" json:"file_size_bytes"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
}
