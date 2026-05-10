package routes

import (
	"github.com/aarcsx/krishaksh-backend/internal/core/middleware"
	"github.com/aarcsx/krishaksh-backend/internal/modules/users/handler"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.RouterGroup, userHandler *handler.UserHandler) {
	userGroup := router.Group("/users")

	// Protect these routes to ensure a valid JWT is fully present
	userGroup.Use(middleware.RequireAuth())
	{
		userGroup.GET("/me", userHandler.GetMe)

		// Example of RBAC logic
		userGroup.GET("/admin-only", middleware.RequireRole("ADMIN"), func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Welcome Admin!"})
		})
	}
}
