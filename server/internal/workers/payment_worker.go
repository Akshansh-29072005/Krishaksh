package workers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aarcsx/krisho-backend/internal/models"
	orderRepo "github.com/aarcsx/krisho-backend/internal/modules/orders/repository"
	paymentRepo "github.com/aarcsx/krisho-backend/internal/modules/payments/repository"
	"github.com/aarcsx/krisho-backend/internal/observability"
	"github.com/aarcsx/krisho-backend/pkg/queue"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type PaymentWorker struct {
	payments paymentRepo.PaymentRepository
	orders   orderRepo.OrderRepository
}

func NewPaymentWorker(p paymentRepo.PaymentRepository, o orderRepo.OrderRepository) *PaymentWorker {
	return &PaymentWorker{payments: p, orders: o}
}

func (w *PaymentWorker) ProcessPaymentEventTask(ctx context.Context, payload queue.PaymentEventPayload) error {
	observability.InitLogger().Info("worker_payment_event_received", "event_id", payload.EventID, "event_type", payload.EventType)
	var env struct {
		Payload struct {
			Payment struct {
				Entity struct {
					ID      string `json:"id"`
					OrderID string `json:"order_id"`
					Amount  int64  `json:"amount"`
					Status  string `json:"status"`
				} `json:"entity"`
			} `json:"payment"`
		} `json:"payload"`
	}
	if err := json.Unmarshal(payload.RawPayload, &env); err != nil {
		observability.M.Inc("worker_errors_total:payment")
		return err
	}
	payment, err := w.payments.GetByProviderOrderID(ctx, env.Payload.Payment.Entity.OrderID)
	if err != nil {
		return err
	}
	amt := float64(env.Payload.Payment.Entity.Amount) / 100
	providerEventID := payload.EventID
	tr := &models.Transaction{ID: uuid.New(), PaymentID: payment.ID, EventType: payload.EventType, ProviderEventID: &providerEventID, Amount: &amt, Status: env.Payload.Payment.Entity.Status, RawPayload: payload.RawPayload}
	if err := w.payments.SaveTransaction(ctx, tr); err != nil {
		return err
	}
	if payload.EventType == "payment.captured" {
		if err := w.payments.MarkPaid(ctx, payment.ID, env.Payload.Payment.Entity.ID); err != nil {
			return err
		}
		_, err = w.orders.SetOrderStatus(ctx, payment.OrderID, []string{"payment_pending", "pending"}, "paid")
		if err != nil {
			return err
		}
	}
	if err := w.payments.MarkWebhookProcessed(ctx, "razorpay", payload.EventID, nil); err != nil {
		observability.M.Inc("worker_errors_total:payment")
		return err
	}
	observability.InitLogger().Info("worker_payment_event_processed", "event_id", payload.EventID, "payment_id", payment.ID.String())
	return nil
}

func (w *PaymentWorker) HandlePaymentEventTask(ctx context.Context, t *asynq.Task) error {
	var payload queue.PaymentEventPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		observability.M.Inc("worker_errors_total:payment")
		return err
	}
	return w.ProcessPaymentEventTask(ctx, payload)
}

func (w *PaymentWorker) ProcessRefundPlaceholderTask(ctx context.Context, payload map[string]string) error {
	return fmt.Errorf("refund processor not implemented yet")
}

func (w *PaymentWorker) HandleRefundPlaceholderTask(ctx context.Context, t *asynq.Task) error {
	var payload map[string]string
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	return w.ProcessRefundPlaceholderTask(ctx, payload)
}
