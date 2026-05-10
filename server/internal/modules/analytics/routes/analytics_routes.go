package routes

import (
	"github.com/aarcsx/krishaksh-backend/internal/core/middleware"
	"github.com/aarcsx/krishaksh-backend/internal/modules/analytics/handler"
	"github.com/gin-gonic/gin"
)

func RegisterAnalyticsRoutes(r *gin.RouterGroup, h *handler.AnalyticsHandler) {
	g := r.Group("/analytics")
	g.Use(middleware.RequireAuth(), middleware.RequireRole("ADMIN"))
	g.GET("/dashboard", h.Dashboard)
}
