package dto

import "time"

type CreateTicketRequest struct {
	Title       string `json:"title" binding:"required,min=5"`
	Description string `json:"description" binding:"required,min=10"`
	Priority    string `json:"priority"` // optional: defaults to "medium"
}

type SendMessageRequest struct {
	Body string `json:"body" binding:"required"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"` // under_review, awaiting_user, resolved, escalated, closed
}

type TicketResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MessageResponse struct {
	ID         string    `json:"id"`
	SenderRole string    `json:"sender_role"`
	Body       string    `json:"body"`
	CreatedAt  time.Time `json:"created_at"`
}
