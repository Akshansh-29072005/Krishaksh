package handler

import (
	"io"
	"net/http"

	"github.com/aarcsx/krishaksh-backend/internal/core/response"
	"github.com/aarcsx/krishaksh-backend/internal/modules/payments/dto"
	"github.com/aarcsx/krishaksh-backend/internal/modules/payments/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentHandler struct{ service service.PaymentService }

func NewPaymentHandler(s service.PaymentService) *PaymentHandler { return &PaymentHandler{service: s} }

func (h *PaymentHandler) CreateOrder(c *gin.Context) {
	var req dto.CreatePaymentOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}
	uid := c.MustGet("user_id").(uuid.UUID)
	out, err := h.service.CreateRazorpayOrder(c.Request.Context(), uid, req.OrderID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "payment order created", out)
}
func (h *PaymentHandler) Webhook(c *gin.Context) {
	raw, _ := io.ReadAll(c.Request.Body)
	sig := c.GetHeader("X-Razorpay-Signature")
	if err := h.service.HandleWebhook(c.Request.Context(), sig, raw); err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "webhook accepted", nil)
}
