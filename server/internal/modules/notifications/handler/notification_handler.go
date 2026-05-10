package handler

import (
	"net/http"
	"time"

	"github.com/aarcsx/krisho-backend/internal/core/response"
	"github.com/aarcsx/krisho-backend/internal/models"
	notifRepo "github.com/aarcsx/krisho-backend/internal/modules/notifications/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	repo notifRepo.NotificationRepository
}

func NewNotificationHandler(repo notifRepo.NotificationRepository) *NotificationHandler {
	return &NotificationHandler{repo: repo}
}

func (h *NotificationHandler) RegisterDevice(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uuid.UUID)

	var req struct {
		Token    string `json:"token" binding:"required"`
		Platform string `json:"platform" binding:"required"` // android | ios
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "token and platform are required")
		return
	}

	token := &models.DeviceToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     req.Token,
		Platform:  req.Platform,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.repo.UpsertDeviceToken(c.Request.Context(), token); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to register device")
		return
	}

	response.Success(c, http.StatusOK, "Device registered for push notifications", nil)
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uuid.UUID)

	notifs, err := h.repo.GetUserNotifications(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch notifications")
		return
	}

	response.Success(c, http.StatusOK, "Notifications fetched", notifs)
}
