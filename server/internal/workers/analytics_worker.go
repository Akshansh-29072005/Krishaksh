package workers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aarcsx/krisho-backend/internal/modules/analytics/repository"
	"github.com/aarcsx/krisho-backend/internal/observability"
	"github.com/aarcsx/krisho-backend/pkg/queue"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type AnalyticsWorker struct {
	repo repository.AnalyticsRepository
}

func NewAnalyticsWorker(r repository.AnalyticsRepository) *AnalyticsWorker {
	return &AnalyticsWorker{repo: r}
}

func (w *AnalyticsWorker) HandleAnalyticsEventTask(ctx context.Context, t *asynq.Task) error {
	var payload queue.AnalyticsEventPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		observability.M.Inc("worker_errors_total:analytics")
		return err
	}
	parse := func(s string) *uuid.UUID {
		if s == "" {
			return nil
		}
		id, err := uuid.Parse(s)
		if err != nil {
			return nil
		}
		return &id
	}
	e := repository.Event{
		EventType:  payload.EventType,
		ActorUser:  parse(payload.ActorUser),
		CompanyID:  parse(payload.CompanyID),
		EntityID:   parse(payload.EntityID),
		OrderID:    parse(payload.OrderID),
		PaymentID:  parse(payload.PaymentID),
		Metadata:   payload.Metadata,
		Value:      payload.Value,
		OccurredAt: time.Now(),
	}
	if err := w.repo.InsertEvent(ctx, e); err != nil {
		observability.M.Inc("worker_errors_total:analytics")
		return err
	}
	observability.M.Inc("analytics_events_processed_total")
	return nil
}
