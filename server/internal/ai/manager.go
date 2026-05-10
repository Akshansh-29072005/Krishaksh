package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aarcsx/krishaksh-backend/internal/ai/providers"
	"github.com/aarcsx/krishaksh-backend/internal/observability"
)

type Manager struct {
	providers          []providers.VisionProvider
	prompts            *PromptManager
	promptVersion      string
	confidenceMin      float64
	maxRetriesProvider int
}

func NewManager(p []providers.VisionProvider, prompts *PromptManager, promptVersion string, confidenceMin float64) *Manager {
	if confidenceMin <= 0 {
		confidenceMin = 0.65
	}
	return &Manager{providers: p, prompts: prompts, promptVersion: promptVersion, confidenceMin: confidenceMin, maxRetriesProvider: 2}
}

func (m *Manager) Infer(ctx context.Context, req InferenceRequest) (*InferenceResult, error) {
	prompt, err := m.prompts.Get(m.promptVersion)
	if err != nil {
		return nil, err
	}
	imgBytes, mime, err := FetchAndPreprocess(ctx, req.ImageURL)
	if err != nil {
		return nil, err
	}
	var lastErr error
	for _, p := range m.providers {
		for attempt := 1; attempt <= m.maxRetriesProvider; attempt++ {
			start := time.Now()
			pctx := ctx
			if req.Timeout > 0 {
				var cancel context.CancelFunc
				pctx, cancel = context.WithTimeout(ctx, req.Timeout)
				defer cancel()
			}
			raw, err := p.Infer(pctx, prompt, req.CropType, imgBytes, mime)
			if err != nil {
				lastErr = err
				observability.InitLogger().Warn("ai_provider_attempt_failed", "provider", p.Name(), "attempt", attempt, "error", err.Error(), "trace_id", req.TraceID, "request_id", req.RequestID)
				continue
			}
			parsed, err := ParseAndValidate(raw)
			if err != nil {
				lastErr = err
				observability.InitLogger().Warn("ai_response_validation_failed", "provider", p.Name(), "attempt", attempt, "error", err.Error(), "trace_id", req.TraceID, "request_id", req.RequestID)
				continue
			}
			best := parsed.Predictions[0]
			for _, pred := range parsed.Predictions[1:] {
				if pred.Confidence > best.Confidence {
					best = pred
				}
			}
			if strings.TrimSpace(parsed.CropType) == "" {
				parsed.CropType = req.CropType
			}
			result := &InferenceResult{
				Provider:          p.Name(),
				Response:          parsed,
				PrimaryPrediction: best,
				Uncertain:         best.Confidence < m.confidenceMin,
				Latency:           time.Since(start),
				PromptVersion:     m.promptVersion,
				RawResponse:       raw,
				InputBytes:        len(imgBytes),
				OutputBytes:       len(raw),
			}
			observability.M.Inc("ai_inference_success_total:" + p.Name())
			observability.InitLogger().Info("ai_inference_success", "provider", p.Name(), "confidence", best.Confidence, "uncertain", result.Uncertain, "latency_ms", result.Latency.Milliseconds(), "input_bytes", result.InputBytes, "output_bytes", result.OutputBytes, "trace_id", req.TraceID, "request_id", req.RequestID)
			return result, nil
		}
		observability.M.Inc("ai_provider_fallback_total")
	}
	observability.M.Inc("ai_inference_failure_total")
	if lastErr == nil {
		lastErr = fmt.Errorf("all providers failed")
	}
	return nil, lastErr
}
