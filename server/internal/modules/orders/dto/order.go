package dto

type CreateOrderRequest struct {
	ShippingMetadata map[string]interface{} `json:"shipping_metadata"`
	Notes            *string                `json:"notes"`
}
