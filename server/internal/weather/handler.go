package weather

import (
	"net/http"

	"github.com/aarcsx/krisho-backend/internal/core/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(v1 *gin.RouterGroup) {
	v1.GET("/weather", h.Get)
}

func (h *Handler) Get(c *gin.Context) {
	lat := c.Query("lat")
	lon := c.Query("lon")
	if lat == "" || lon == "" {
		response.Error(c, http.StatusBadRequest, "lat and lon are required")
		return
	}
	data, source, err := h.service.Get(c.Request.Context(), lat, lon)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Keep cache/debug source available while preserving the standard API envelope.
	response.Success(c, http.StatusOK, source, data)
}
