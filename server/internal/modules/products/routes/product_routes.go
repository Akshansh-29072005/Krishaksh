package routes

import (
	"github.com/aarcsx/krishaksh-backend/internal/modules/products/handler"
	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(router *gin.RouterGroup, h *handler.ProductHandler) {
	products := router.Group("/products")
	{
		products.GET("", h.GetAll)
		products.GET("/:id", h.GetByID)
	}
}
