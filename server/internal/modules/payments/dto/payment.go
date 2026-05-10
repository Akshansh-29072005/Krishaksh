package dto

import "github.com/google/uuid"

type CreatePaymentOrderRequest struct {
	OrderID uuid.UUID `json:"order_id" binding:"required"`
}
