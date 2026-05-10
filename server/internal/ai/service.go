package ai

import (
	"context"
	"time"
)

type Service struct{ mgr *Manager }

func NewService(mgr *Manager) *Service { return &Service{mgr: mgr} }

func (s *Service) Analyze(ctx context.Context, imageURL, cropType, traceID, requestID string) (*InferenceResult, error) {
	return s.mgr.Infer(ctx, InferenceRequest{
		ImageURL:  imageURL,
		CropType:  cropType,
		Timeout:   22 * time.Second,
		TraceID:   traceID,
		RequestID: requestID,
	})
}
