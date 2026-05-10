package handler

import (
	"net/http"
	"strconv"

	"github.com/aarcsx/krisho-backend/internal/core/response"
	"github.com/aarcsx/krisho-backend/internal/modules/analytics/service"
	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct{ service service.AnalyticsService }

func NewAnalyticsHandler(s service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: s}
}

func (h *AnalyticsHandler) Dashboard(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	m, err := h.service.Dashboard(c.Request.Context(), days)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to load dashboard")
		return
	}
	response.Success(c, http.StatusOK, "dashboard metrics", m)
}
