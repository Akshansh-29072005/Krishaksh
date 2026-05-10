package dto

import "github.com/google/uuid"

type AddCartItemRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1,max=100"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1,max=100"`
}
