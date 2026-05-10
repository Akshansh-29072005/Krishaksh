package handler

import (
	"net/http"

	"github.com/aarcsx/krisho-backend/internal/core/response"
	"github.com/aarcsx/krisho-backend/internal/modules/orders/dto"
	"github.com/aarcsx/krisho-backend/internal/modules/orders/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct{ service service.OrderService }

func NewOrderHandler(s service.OrderService) *OrderHandler { return &OrderHandler{service: s} }

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	uid := c.MustGet("user_id").(uuid.UUID)
	var req dto.CreateOrderRequest
	_ = c.ShouldBindJSON(&req)
	ord, err := h.service.CreateFromCart(c.Request.Context(), uid, req.ShippingMetadata, req.Notes)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "order created", ord)
}
func (h *OrderHandler) ListOrders(c *gin.Context) {
	uid := c.MustGet("user_id").(uuid.UUID)
	orders, err := h.service.GetOrders(c.Request.Context(), uid)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "orders fetched", orders)
}
func (h *OrderHandler) GetOrder(c *gin.Context) {
	uid := c.MustGet("user_id").(uuid.UUID)
	oid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	ord, err := h.service.GetOrder(c.Request.Context(), uid, oid)
	if err != nil {
		response.Error(c, http.StatusNotFound, "order not found")
		return
	}
	response.Success(c, http.StatusOK, "order fetched", ord)
}
