package response

import (
	"github.com/gin-gonic/gin"
)

// APIResponse represents the standard structure for all API JSON responses.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success responds with a 200 OK and standardized payload
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error responds with a standard error payload
func Error(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, APIResponse{
		Success: false,
		Message: message,
	})
}
