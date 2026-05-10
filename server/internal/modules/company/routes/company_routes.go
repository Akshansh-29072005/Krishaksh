package routes

import (
	"github.com/aarcsx/krisho-backend/internal/core/middleware"
	"github.com/aarcsx/krisho-backend/internal/modules/company/handler"
	"github.com/gin-gonic/gin"
)

func RegisterCompanyRoutes(r *gin.RouterGroup, h *handler.CompanyHandler) {
	g := r.Group("/company")
	g.Use(middleware.RequireAuth(), middleware.RequireRole("VENDOR", "ADMIN"))
	g.GET("/dashboard", h.Dashboard)
}
