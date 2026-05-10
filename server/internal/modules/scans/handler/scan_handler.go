package handler

import (
	"net/http"

	"github.com/aarcsx/krisho-backend/internal/core/response"
	"github.com/aarcsx/krisho-backend/internal/modules/scans/dto"
	"github.com/aarcsx/krisho-backend/internal/modules/scans/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ScanHandler struct {
	service service.ScanService
}

func NewScanHandler(s service.ScanService) *ScanHandler {
	return &ScanHandler{service: s}
}

func (h *ScanHandler) UploadImagePresigned(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uuid.UUID)

	contentType := c.Query("content_type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	url, key, err := h.service.GetPresignedUploadURL(c.Request.Context(), userID, contentType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to generate signed url")
		return
	}

	response.Success(c, http.StatusOK, "Upload URL generated", gin.H{
		"presigned_url": url,
		"image_key":     key,
	})
}

func (h *ScanHandler) CreateScanPostUpload(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uuid.UUID)

	var req dto.CreateScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format")
		return
	}

	// This instantly returns standard response 200/202, since the processing goes into the Asynq worker seamlessly.
	scan, err := h.service.ProcessUpload(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to initialize scan pipeline")
		return
	}

	response.Success(c, http.StatusAccepted, "Image scan queued successfully", scan)
}

func (h *ScanHandler) GetScanHistory(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uuid.UUID)

	scans, err := h.service.GetScanHistory(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch scan history")
		return
	}

	response.Success(c, http.StatusOK, "History fetched globally", scans)
}

func (h *ScanHandler) GetScanDetails(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uuid.UUID)

	scanIDParam := c.Param("id")
	scanID, err := uuid.Parse(scanIDParam)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid scan ID configuration")
		return
	}

	scan, err := h.service.GetScanDetails(c.Request.Context(), userID, scanID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Scan target absent or prohibited")
		return
	}

	response.Success(c, http.StatusOK, "Success", scan)
}
