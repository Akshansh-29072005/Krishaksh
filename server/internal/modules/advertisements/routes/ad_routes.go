package routes

import (
	"github.com/aarcsx/krishaksh-backend/internal/modules/advertisements/handler"
	"github.com/gin-gonic/gin"
)

func RegisterAdvertisementRoutes(router *gin.RouterGroup, h *handler.AdvertisementHandler) {
	ads := router.Group("/ads")
	{
		ads.GET("/featured", h.GetFeatured)
	}
}
