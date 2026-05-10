package handler

import (
	"net/http"

	"github.com/aarcsx/krishaksh-backend/internal/core/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetMe responds with the details of the currently logged-in user through context.
func (h *UserHandler) GetMe(c *gin.Context) {
	// Extract details injected by the RequireAuth middleware
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	role, _ := c.Get("role")
	email, _ := c.Get("email")

	userID := userIDRaw.(uuid.UUID)

	// In a real application, you'd fetch deep user details from the UserRepository using userID.
	// We'll mock the response structure as requested.
	response.Success(c, http.StatusOK, "User details fetched successfully", gin.H{
		"id":    userID,
		"role":  role,
		"email": email,
	})
}
