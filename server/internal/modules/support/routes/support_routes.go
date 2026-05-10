package routes

import (
	"github.com/aarcsx/krishaksh-backend/internal/core/middleware"
	"github.com/aarcsx/krishaksh-backend/internal/modules/support/handler"
	"github.com/gin-gonic/gin"
)

func RegisterSupportRoutes(router *gin.RouterGroup, h *handler.SupportHandler) {
	support := router.Group("/support")
	support.Use(middleware.RequireAuth())
	{
		support.POST("/tickets", h.CreateTicket)
		support.GET("/tickets", h.GetTickets)
		support.GET("/tickets/:id", h.GetTicket)
		support.POST("/tickets/:id/messages", h.SendMessage)
	}
}
