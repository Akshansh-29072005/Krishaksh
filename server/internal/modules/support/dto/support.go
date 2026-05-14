package dto

import "time"

type CreateTicketRequest struct {
	Title           string `json:"title" binding:"required,min=5"`
	Description     string `json:"description" binding:"required,min=10"`
	Priority        string `json:"priority"` // optional: defaults to "medium"
	RequestCallback bool   `json:"request_callback,omitempty"`
}

type SendMessageRequest struct {
	Body string `json:"body" binding:"required"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"` // under_review, awaiting_user, resolved, escalated, closed
}

type TicketResponse struct {
	ID                string    `json:"id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	Status            string    `json:"status"`
	Priority          string    `json:"priority"`
	CallbackRequested bool      `json:"callback_requested"`
	CallbackStatus    string    `json:"callback_status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Body              string    `json:"body"`
}
