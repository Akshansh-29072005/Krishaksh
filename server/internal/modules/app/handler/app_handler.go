package handler

import (
    "net/http"

    "github.com/aarcsx/krisho-backend/internal/core/response"
    "github.com/aarcsx/krisho-backend/internal/models"
    appService "github.com/aarcsx/krisho-backend/internal/modules/app/service"
    "github.com/gin-gonic/gin"
)

type AppHandler struct {
    service *appService.AppService
}

func NewAppHandler(service *appService.AppService) *AppHandler {
    return &AppHandler{service: service}
}

func (h *AppHandler) GetCrops(c *gin.Context) {
    crops, err := h.service.GetActiveCrops(c.Request.Context())
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to fetch crop metadata")
        return
    }
    if crops == nil {
        crops = []*models.Crop{}
    }
    response.Success(c, http.StatusOK, "Crop list fetched", crops)
}

func (h *AppHandler) GetAppConfig(c *gin.Context) {
    config, err := h.service.GetAppConfig(c.Request.Context())
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to fetch app config")
        return
    }
    response.Success(c, http.StatusOK, "App config fetched", config)
}
