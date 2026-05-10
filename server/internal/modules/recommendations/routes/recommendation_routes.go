package routes

import (
	"github.com/aarcsx/krishaksh-backend/internal/core/middleware"
	"github.com/aarcsx/krishaksh-backend/internal/modules/recommendations/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRecommendationRoutes(router *gin.RouterGroup, h *handler.RecommendationHandler) {
	recs := router.Group("/recommendations")
	recs.Use(middleware.RequireAuth())
	{
		recs.GET("/:scan_id", h.GetForScan)
	}
}
