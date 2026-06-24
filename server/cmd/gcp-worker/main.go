package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	internalAI "github.com/aarcsx/krisho-backend/internal/ai"
	aiProviders "github.com/aarcsx/krisho-backend/internal/ai/providers"
	"github.com/aarcsx/krisho-backend/internal/config"
	"github.com/aarcsx/krisho-backend/internal/database"
	analyticsRepo "github.com/aarcsx/krisho-backend/internal/modules/analytics/repository"
	orderRepo "github.com/aarcsx/krisho-backend/internal/modules/orders/repository"
	paymentRepo "github.com/aarcsx/krisho-backend/internal/modules/payments/repository"
	scanRepo "github.com/aarcsx/krisho-backend/internal/modules/scans/repository"
	"github.com/aarcsx/krisho-backend/internal/observability"
	"github.com/aarcsx/krisho-backend/internal/workers"
	"github.com/aarcsx/krisho-backend/pkg/queue"
)

func main() {
	observability.InitLogger()
	log.Printf("GCP Worker initializing...")

	cfg := config.LoadConfig()
	log.Printf("Config loaded: GCP Project=%s, Queue=%s, Worker URL=%s",
		os.Getenv("GCP_PROJECT_ID"), os.Getenv("CLOUD_TASKS_QUEUE"), os.Getenv("CLOUD_TASKS_WORKER_URL"))

	if err := cfg.Validate(); err != nil {
		log.Fatalf("config validation failed: %v", err)
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	defer db.Close()
	log.Printf("Connected to database successfully")

	scanRepository := scanRepo.NewScanRepository(db)
	orderRepository := orderRepo.NewOrderRepository(db)
	paymentRepository := paymentRepo.NewPaymentRepository(db)
	analyticsRepository := analyticsRepo.NewAnalyticsRepository(db)

	promptMgr := internalAI.NewPromptManager("internal/ai/prompts")
	aiMgr := internalAI.NewManager(
		[]aiProviders.VisionProvider{
			aiProviders.NewGeminiProvider(strings.TrimSpace(os.Getenv("GEMINI_API_KEY"))),
		},
		promptMgr,
		"vision_v1",
		0.65,
	)
	aiSvc := internalAI.NewService(aiMgr)

	scanWorker := workers.NewScanWorker(scanRepository, aiSvc)
	paymentWorker := workers.NewPaymentWorker(paymentRepository, orderRepository)
	analyticsWorker := workers.NewAnalyticsWorker(analyticsRepository)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		log.Printf("INCOMING REQUEST: Method=%s, Path=%s, ClientIP=%s",
			c.Request.Method, c.Request.URL.Path, c.ClientIP())
		c.Next()
	})

	r.POST("/tasks/scan", func(c *gin.Context) {
		log.Printf("Received /tasks/scan request!")

		var payload queue.ScanAnalyzePayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			log.Printf("ERROR: Invalid payload for /tasks/scan: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Processing scan task: ScanID=%s, CropType=%s, ImageURL=%s",
			payload.ScanID, payload.CropType, payload.ImageURL)

		ctx := context.Background()
		if err := scanWorker.ProcessScanTask(ctx, payload); err != nil {
			log.Printf("ERROR: Failed to process scan task: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Printf("SUCCESS: Scan task completed for ScanID=%s", payload.ScanID)

		c.Status(http.StatusOK)
	})

	r.POST("/tasks/payment", func(c *gin.Context) {
		var payload queue.PaymentEventPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := context.Background()
		if err := paymentWorker.ProcessPaymentEventTask(ctx, payload); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	})

	r.POST("/tasks/refund", func(c *gin.Context) {
		var payload map[string]string
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := context.Background()
		if err := paymentWorker.ProcessRefundPlaceholderTask(ctx, payload); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	})

	r.POST("/tasks/analytics", func(c *gin.Context) {
		var payload queue.AnalyticsEventPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := context.Background()
		if err := analyticsWorker.ProcessAnalyticsEventTask(ctx, payload); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	})

	r.GET("/health", func(c *gin.Context) {
		log.Printf("Health check received!")
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "gcp-worker"})
	})

	r.GET("/", func(c *gin.Context) {
		log.Printf("Root handler hit!")
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "gcp-worker is running!"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("GCP Worker server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
