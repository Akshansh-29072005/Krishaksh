package weather

import (
	"net/http"

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
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "lat and lon are required"})
		return
	}
	data, source, err := h.service.Get(c.Request.Context(), lat, lon)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "source": source, "data": data})
}
