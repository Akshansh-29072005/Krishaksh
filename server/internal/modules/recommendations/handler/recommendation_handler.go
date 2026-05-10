package handler

import (
	"net/http"

	"github.com/aarcsx/krishaksh-backend/internal/core/response"
	"github.com/aarcsx/krishaksh-backend/internal/modules/recommendations/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RecommendationHandler struct {
	service service.RecommendationService
}

func NewRecommendationHandler(s service.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{service: s}
}

func (h *RecommendationHandler) GetForScan(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uuid.UUID)

	scanID, err := uuid.Parse(c.Param("scan_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid scan ID")
		return
	}

	result, err := h.service.GetForScan(c.Request.Context(), userID, scanID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Recommendation data not found")
		return
	}

	response.Success(c, http.StatusOK, "Recommendations fetched", result)
}
