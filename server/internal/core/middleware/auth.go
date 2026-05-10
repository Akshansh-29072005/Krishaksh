package middleware

import (
	"net/http"
	"strings"

	"github.com/aarcsx/krishaksh-backend/internal/core/response"
	"github.com/aarcsx/krishaksh-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

// RequireAuth blocks endpoints if no valid JWT is present
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Inject user context variables seamlessly into Gin
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// RequireRole uses RBAC to restrict authenticated routes strictly to allowed roles
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "User context not found")
			return
		}

		role, ok := roleVal.(string)
		if !ok {
			response.Error(c, http.StatusInternalServerError, "Invalid role format in context")
			return
		}

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		response.Error(c, http.StatusForbidden, "You do not have permission to access this resource")
	}
}
