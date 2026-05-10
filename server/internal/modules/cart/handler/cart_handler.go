package handler

import (
	"net/http"

	"github.com/aarcsx/krisho-backend/internal/core/response"
	"github.com/aarcsx/krisho-backend/internal/modules/cart/dto"
	"github.com/aarcsx/krisho-backend/internal/modules/cart/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CartHandler struct{ service service.CartService }

func NewCartHandler(s service.CartService) *CartHandler { return &CartHandler{service: s} }

func (h *CartHandler) AddItem(c *gin.Context) {
	var req dto.AddCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}
	uid := c.MustGet("user_id").(uuid.UUID)
	if err := h.service.AddItem(c.Request.Context(), uid, req.ProductID, req.Quantity); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "cart updated", nil)
}
func (h *CartHandler) GetCart(c *gin.Context) {
	uid := c.MustGet("user_id").(uuid.UUID)
	out, err := h.service.GetCart(c.Request.Context(), uid)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "cart fetched", out)
}
func (h *CartHandler) RemoveItem(c *gin.Context) {
	uid := c.MustGet("user_id").(uuid.UUID)
	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.service.RemoveItem(c.Request.Context(), uid, itemID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "item removed", nil)
}

func (h *CartHandler) UpdateItem(c *gin.Context) {
	uid := c.MustGet("user_id").(uuid.UUID)
	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req dto.AddCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}
	if err := h.service.UpdateItemQuantity(c.Request.Context(), uid, itemID, req.ProductID, req.Quantity); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "item quantity updated", nil)
}
