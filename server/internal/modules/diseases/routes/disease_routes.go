package routes

import (
	"github.com/aarcsx/krishaksh-backend/internal/modules/diseases/handler"
	"github.com/gin-gonic/gin"
)

func RegisterDiseaseRoutes(router *gin.RouterGroup, h *handler.DiseaseHandler) {
	diseases := router.Group("/diseases")
	{
		diseases.GET("", h.GetAll)
		diseases.GET("/:id", h.GetByID)
	}
}
