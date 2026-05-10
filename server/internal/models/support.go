package models

import (
	"time"

	"github.com/google/uuid"
)

// SupportTicket represents a farmer's help request
type SupportTicket struct {
	ID          uuid.UUID  `db:"id"`
	UserID      uuid.UUID  `db:"user_id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Status      string     `db:"status"`      // open, under_review, awaiting_user, resolved, escalated, closed
	Priority    string     `db:"priority"`    // low, medium, high, critical
	AssignedTo  *uuid.UUID `db:"assigned_to"` // admin user_id
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	ResolvedAt  *time.Time `db:"resolved_at"`
}

// SupportMessage is a single chat-like entry in a ticket thread
type SupportMessage struct {
	ID         uuid.UUID `db:"id"`
	TicketID   uuid.UUID `db:"ticket_id"`
	SenderID   uuid.UUID `db:"sender_id"`
	SenderRole string    `db:"sender_role"` // farmer | admin
	Body       string    `db:"body"`
	CreatedAt  time.Time `db:"created_at"`
}

// SupportAttachment holds file metadata uploaded within a ticket (image or voice note)
type SupportAttachment struct {
	ID            uuid.UUID  `db:"id"`
	TicketID      uuid.UUID  `db:"ticket_id"`
	MessageID     *uuid.UUID `db:"message_id"` // optional – link to a specific message
	UploaderID    uuid.UUID  `db:"uploader_id"`
	FileURL       string     `db:"file_url"`
	FileType      string     `db:"file_type"` // image | voice
	FileSizeBytes int64      `db:"file_size_bytes"`
	CreatedAt     time.Time  `db:"created_at"`
}
