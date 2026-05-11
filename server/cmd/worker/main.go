package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/hibiken/asynq"

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
	cfg := config.LoadConfig()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("config validation failed: %v", err)
	}
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	defer db.Close()

	if os.Getenv("RUN_DB_MIGRATIONS") == "true" {
		if err := database.RunMigrations(context.Background(), db, "db/migrations"); err != nil {
			log.Fatalf("migration failed: %v", err)
		}
	}

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

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.RedisAddr},
		asynq.Config{Concurrency: 20, Queues: map[string]int{"scans": 10, "payments": 4, "refunds": 2, "analytics": 4}},
	)
	mux := asynq.NewServeMux()
	mux.HandleFunc(queue.TypeScanAnalyze, workers.NewScanWorker(scanRepository, aiSvc).HandleScanAnalyzeTask)
	mux.HandleFunc(queue.TypePaymentEvent, workers.NewPaymentWorker(paymentRepository, orderRepository).HandlePaymentEventTask)
	mux.HandleFunc(queue.TypeRefundStart, workers.NewPaymentWorker(paymentRepository, orderRepository).HandleRefundPlaceholderTask)
	mux.HandleFunc(queue.TypeAnalyticsEvt, workers.NewAnalyticsWorker(analyticsRepository).HandleAnalyticsEventTask)

	go func() {
		observability.InitLogger().Info("worker_server_started", "redis", cfg.RedisAddr)
		if err := srv.Run(mux); err != nil {
			log.Fatalf("asynq run failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	srv.Shutdown()
	<-ctx.Done()
	observability.InitLogger().Info("worker_server_stopped")
}
