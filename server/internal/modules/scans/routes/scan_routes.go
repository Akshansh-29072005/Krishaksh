package routes

import (
	"github.com/aarcsx/krishaksh-backend/internal/core/middleware"
	"github.com/aarcsx/krishaksh-backend/internal/modules/scans/handler"
	"github.com/gin-gonic/gin"
)

func RegisterScanRoutes(router *gin.RouterGroup, scanHandler *handler.ScanHandler) {
	scans := router.Group("/scans")
	scans.Use(middleware.RequireAuth())
	{
		// 1. Ask backend for AWS S3 upload url natively
		scans.GET("/upload-url", scanHandler.UploadImagePresigned)

		// 2. Client finishes uploading internally to AWS then pings this to start scanning
		scans.POST("", scanHandler.CreateScanPostUpload)

		scans.GET("/history", scanHandler.GetScanHistory)
		scans.GET("/:id", scanHandler.GetScanDetails)
	}
}
