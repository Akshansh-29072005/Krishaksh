package workers

import (
	"context"
	"encoding/json"
	"fmt"

	internalAI "github.com/aarcsx/krisho-backend/internal/ai"
	"github.com/aarcsx/krisho-backend/internal/modules/scans/repository"
	"github.com/aarcsx/krisho-backend/internal/observability"
	"github.com/aarcsx/krisho-backend/pkg/queue"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type ScanWorker struct {
	repo      repository.ScanRepository
	aiService *internalAI.Service
}

func NewScanWorker(repo repository.ScanRepository, aiService *internalAI.Service) *ScanWorker {
	return &ScanWorker{repo: repo, aiService: aiService}
}

func (w *ScanWorker) HandleScanAnalyzeTask(ctx context.Context, t *asynq.Task) error {
	var payload queue.ScanAnalyzePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json unmarshal fallback err: %w", err)
	}

	scanID, err := uuid.Parse(payload.ScanID)
	if err != nil {
		return fmt.Errorf("invalid scan id explicitly: %w", err)
	}

	traceID := ""
	requestID := ""
	urlPreview := payload.ImageURL
	if len(urlPreview) > 60 {
		urlPreview = urlPreview[:60] + "..."
	}

	observability.InitLogger().Info("worker_scan_started", 
		"scan_id", scanID.String(), 
		"crop_type", payload.CropType,
		"image_url_preview", urlPreview,
	)
	_ = w.repo.UpdateScanStatus(ctx, scanID, "PROCESSING")

	result, err := w.aiService.Analyze(ctx, payload.ImageURL, payload.CropType, traceID, requestID)
	if err != nil {
		errMsg := err.Error()
		status := "FAILED"
		provider := "none"
		meta := map[string]interface{}{"degraded_mode": true, "reason": "provider_failure"}
		_ = w.repo.UpdateAIInference(ctx, scanID, status, &provider, nil, &errMsg, meta)
		observability.M.Inc("worker_errors_total:scan")
		observability.InitLogger().Error("worker_scan_failed", "scan_id", scanID.String(), "error", errMsg)
		return err
	}

	status := "COMPLETED"
	if result.Uncertain {
		status = "UNCERTAIN"
	}
	provider := result.Provider
	confidence := result.PrimaryPrediction.Confidence
	meta := map[string]interface{}{
		"prompt_version":   result.PromptVersion,
		"notes":            result.Response.Notes,
		"uncertain":        result.Uncertain,
		"predictions":      result.Response.Predictions,
		"latency_ms":       result.Latency.Milliseconds(),
		"preprocess_bytes": result.InputBytes,
		"output_bytes":     result.OutputBytes,
	}
	if err := w.repo.UpdateAIInference(ctx, scanID, status, &provider, &confidence, nil, meta); err != nil {
		return err
	}

	observability.InitLogger().Info("worker_scan_completed", "scan_id", scanID.String(), "provider", provider, "confidence", confidence, "status", status)
	return nil
}
