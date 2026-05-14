package handler

import (
	"net/http"
	"strings"

	"github.com/aarcsx/krisho-backend/internal/core/response"
	"github.com/aarcsx/krisho-backend/internal/modules/auth/repository"
	"github.com/aarcsx/krisho-backend/internal/modules/users/dto"
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
		"phone": user.PhoneNumber,
		"crops": []string{},
	})
}

func (h *UserHandler) UpdatePhone(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}

	userID := userIDRaw.(uuid.UUID)
	var req dto.UpdatePhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid phone request")
		return
	}

	phone := strings.TrimSpace(req.Phone)
	if phone == "" {
		response.Error(c, http.StatusBadRequest, "Phone number is required")
		return
	}

	if err := h.repo.UpdatePhoneNumber(c.Request.Context(), userID, phone); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update phone number")
		return
	}

	user, err := h.repo.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to load updated user profile")
		return
	}

	response.Success(c, http.StatusOK, "Phone number updated successfully", gin.H{
		"id":    user.ID,
		"role":  c.GetString("role"),
		"email": user.Email,
		"name":  user.FullName,
		"phone": user.PhoneNumber,
		"crops": []string{},
	})
}
