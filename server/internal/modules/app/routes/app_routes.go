package routes

import (
    "github.com/aarcsx/krisho-backend/internal/modules/app/handler"
    "github.com/gin-gonic/gin"
)

func RegisterAppRoutes(router *gin.RouterGroup, h *handler.AppHandler) {
    router.GET("/crops", h.GetCrops)
    router.GET("/app-config", h.GetAppConfig)
}
