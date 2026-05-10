package routes

import (
	"github.com/aarcsx/krishaksh-backend/internal/core/middleware"
	"github.com/aarcsx/krishaksh-backend/internal/modules/payments/handler"
	"github.com/gin-gonic/gin"
)

func RegisterPaymentRoutes(router *gin.RouterGroup, h *handler.PaymentHandler) {
	auth := router.Group("/payments")
	auth.Use(middleware.RequireAuth())
	auth.POST("/create-order", h.CreateOrder)

	public := router.Group("/payments")
	public.POST("/webhook", h.Webhook)
}
