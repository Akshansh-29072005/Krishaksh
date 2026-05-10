package handler

import (
	"net/http"

	"github.com/aarcsx/krisho-backend/internal/core/response"
	adRepo "github.com/aarcsx/krisho-backend/internal/modules/advertisements/repository"
	"github.com/gin-gonic/gin"
)

type AdvertisementHandler struct {
	repo adRepo.AdvertisementRepository
}

func NewAdvertisementHandler(repo adRepo.AdvertisementRepository) *AdvertisementHandler {
	return &AdvertisementHandler{repo: repo}
}

func (h *AdvertisementHandler) GetFeatured(c *gin.Context) {
	region := c.Query("region")
	var regionPtr *string
	if region != "" {
		regionPtr = &region
	}

	ads, err := h.repo.GetFeaturedAds(c.Request.Context(), regionPtr)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch advertisements")
		return
	}

	response.Success(c, http.StatusOK, "Featured ads fetched", ads)
}
