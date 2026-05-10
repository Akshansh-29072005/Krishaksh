package repository

import (
	"context"

	"github.com/aarcsx/krisho-backend/internal/database"
	"github.com/aarcsx/krisho-backend/internal/models"
	"github.com/google/uuid"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, p *models.Payment) error
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Payment, error)
	GetByProviderOrderID(ctx context.Context, providerOrderID string) (*models.Payment, error)
	MarkPaid(ctx context.Context, paymentID uuid.UUID, providerPaymentID string) error
	SaveTransaction(ctx context.Context, tr *models.Transaction) error
	UpsertWebhookEvent(ctx context.Context, event *models.WebhookEvent) error
	MarkWebhookProcessed(ctx context.Context, provider, eventID string, errMsg *string) error
}

type paymentRepoImpl struct{ db *database.DB }

func NewPaymentRepository(db *database.DB) PaymentRepository { return &paymentRepoImpl{db: db} }

func (r *paymentRepoImpl) CreatePayment(ctx context.Context, p *models.Payment) error {
	_, err := r.db.Pool.Exec(ctx, `INSERT INTO payments (id, order_id, provider, provider_order_id, provider_payment_id, status, amount, currency, attempts, metadata, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,NOW(),NOW())`, p.ID, p.OrderID, p.Provider, p.ProviderOrderID, p.ProviderPaymentID, p.Status, p.Amount, p.Currency, p.Attempts, p.Metadata)
	return err
}
func (r *paymentRepoImpl) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Payment, error) {
	p := &models.Payment{}
	err := r.db.Pool.QueryRow(ctx, `SELECT id, order_id, provider, provider_order_id, provider_payment_id, status, amount, currency, attempts, metadata, created_at, updated_at FROM payments WHERE order_id=$1 ORDER BY created_at DESC LIMIT 1`, orderID).Scan(&p.ID, &p.OrderID, &p.Provider, &p.ProviderOrderID, &p.ProviderPaymentID, &p.Status, &p.Amount, &p.Currency, &p.Attempts, &p.Metadata, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}
func (r *paymentRepoImpl) GetByProviderOrderID(ctx context.Context, providerOrderID string) (*models.Payment, error) {
	p := &models.Payment{}
	err := r.db.Pool.QueryRow(ctx, `SELECT id, order_id, provider, provider_order_id, provider_payment_id, status, amount, currency, attempts, metadata, created_at, updated_at FROM payments WHERE provider_order_id=$1`, providerOrderID).Scan(&p.ID, &p.OrderID, &p.Provider, &p.ProviderOrderID, &p.ProviderPaymentID, &p.Status, &p.Amount, &p.Currency, &p.Attempts, &p.Metadata, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}
func (r *paymentRepoImpl) MarkPaid(ctx context.Context, paymentID uuid.UUID, providerPaymentID string) error {
	_, err := r.db.Pool.Exec(ctx, `UPDATE payments SET status='paid', provider_payment_id=$1, updated_at=NOW() WHERE id=$2`, providerPaymentID, paymentID)
	return err
}
func (r *paymentRepoImpl) SaveTransaction(ctx context.Context, tr *models.Transaction) error {
	_, err := r.db.Pool.Exec(ctx, `INSERT INTO transactions (id, payment_id, event_type, provider_event_id, amount, currency, status, raw_payload, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW())`, tr.ID, tr.PaymentID, tr.EventType, tr.ProviderEventID, tr.Amount, tr.Currency, tr.Status, tr.RawPayload)
	return err
}
func (r *paymentRepoImpl) UpsertWebhookEvent(ctx context.Context, event *models.WebhookEvent) error {
	_, err := r.db.Pool.Exec(ctx, `INSERT INTO webhook_events (id, provider, event_id, event_type, signature, payload, processed, retries, created_at) VALUES ($1,$2,$3,$4,$5,$6,false,0,NOW()) ON CONFLICT (provider,event_id) DO NOTHING`, event.ID, event.Provider, event.EventID, event.EventType, event.Signature, event.Payload)
	return err
}
func (r *paymentRepoImpl) MarkWebhookProcessed(ctx context.Context, provider, eventID string, errMsg *string) error {
	if errMsg == nil {
		_, err := r.db.Pool.Exec(ctx, `UPDATE webhook_events SET processed=true, processed_at=NOW(), last_error=NULL WHERE provider=$1 AND event_id=$2`, provider, eventID)
		return err
	}
	_, err := r.db.Pool.Exec(ctx, `UPDATE webhook_events SET retries=retries+1, last_error=$3 WHERE provider=$1 AND event_id=$2`, provider, eventID, *errMsg)
	return err
}
