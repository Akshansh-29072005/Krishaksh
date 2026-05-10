package routes

import (
	"github.com/aarcsx/krisho-backend/internal/core/middleware"
	"github.com/aarcsx/krisho-backend/internal/modules/admin/handler"
	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(r *gin.RouterGroup, h *handler.AdminHandler) {
	g := r.Group("/admin")
	g.Use(middleware.RequireAuth(), middleware.RequireRole("ADMIN"))
	g.GET("/users", h.ListUsers)
	g.PATCH("/users/:id/role", h.ChangeRole)
	g.PATCH("/users/:id/suspend", h.SuspendUser)
	g.POST("/diseases", h.CreateDisease)
	g.PATCH("/diseases/:id", h.UpdateDisease)
	g.DELETE("/diseases/:id", h.DeleteDisease)
	g.POST("/products", h.CreateProduct)
	g.PATCH("/products/:id", h.UpdateProduct)
	g.DELETE("/products/:id", h.DeleteProduct)
	g.GET("/campaigns", h.ListCampaigns)
	g.PATCH("/campaigns/:id", h.UpdateCampaign)
	g.GET("/support/tickets", h.ListSupport)
	g.PATCH("/support/tickets/:id", h.UpdateSupport)
}
