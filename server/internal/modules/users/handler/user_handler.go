package handler

import (
	"net/http"

	"github.com/aarcsx/krisho-backend/internal/core/response"
	"github.com/aarcsx/krisho-backend/internal/modules/auth/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	repo repository.AuthRepository
}

func NewUserHandler(repo repository.AuthRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// GetMe responds with the details of the currently logged-in user through context.
func (h *UserHandler) GetMe(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	userID := userIDRaw.(uuid.UUID)
	role, _ := c.Get("role")

	user, err := h.repo.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to load user profile")
		return
	}

	response.Success(c, http.StatusOK, "User details fetched successfully", gin.H{
		"id":    user.ID,
		"role":  role,
		"email": user.Email,
		"name":  user.FullName,
	})
}
