package routes

import (
	"github.com/aarcsx/krisho-backend/internal/core/middleware"
	"github.com/aarcsx/krisho-backend/internal/modules/orders/handler"
	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(router *gin.RouterGroup, h *handler.OrderHandler) {
	g := router.Group("/orders")
	g.Use(middleware.RequireAuth())
	{
		g.POST("", h.CreateOrder)
		g.GET("", h.ListOrders)
		g.GET(":id", h.GetOrder)
	}
}
