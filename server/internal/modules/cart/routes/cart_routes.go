package routes

import (
	"github.com/aarcsx/krisho-backend/internal/core/middleware"
	"github.com/aarcsx/krisho-backend/internal/modules/cart/handler"
	"github.com/gin-gonic/gin"
)

func RegisterCartRoutes(router *gin.RouterGroup, h *handler.CartHandler) {
	g := router.Group("/cart")
	g.Use(middleware.RequireAuth())
	{
		g.POST("/items", h.AddItem)
		g.PATCH("/items/:id", h.UpdateItem)
		g.GET("", h.GetCart)
		g.DELETE("/items/:id", h.RemoveItem)
	}
}
