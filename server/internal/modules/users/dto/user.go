package dto

type UpdatePhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
}
