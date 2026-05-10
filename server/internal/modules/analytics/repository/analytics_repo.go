package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aarcsx/krishaksh-backend/internal/database"
	"github.com/google/uuid"
)

type Event struct {
	EventType  string
	ActorUser  *uuid.UUID
	CompanyID  *uuid.UUID
	EntityID   *uuid.UUID
	OrderID    *uuid.UUID
	PaymentID  *uuid.UUID
	Metadata   map[string]interface{}
	Value      *float64
	OccurredAt time.Time
}

type AnalyticsRepository interface {
	InsertEvent(ctx context.Context, e Event) error
	DashboardMetrics(ctx context.Context, since time.Time) (map[string]int64, error)
	CompanyMetrics(ctx context.Context, companyID uuid.UUID, since time.Time) (map[string]int64, error)
}

type analyticsRepoImpl struct{ db *database.DB }

func NewAnalyticsRepository(db *database.DB) AnalyticsRepository { return &analyticsRepoImpl{db: db} }

func (r *analyticsRepoImpl) InsertEvent(ctx context.Context, e Event) error {
	meta, _ := json.Marshal(e.Metadata)
	_, err := r.db.Pool.Exec(ctx, `INSERT INTO analytics_events (event_type, actor_user_id, company_id, entity_id, order_id, payment_id, metadata, value, occurred_at, created_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,NOW())`, e.EventType, e.ActorUser, e.CompanyID, e.EntityID, e.OrderID, e.PaymentID, meta, e.Value, e.OccurredAt)
	return err
}

func (r *analyticsRepoImpl) DashboardMetrics(ctx context.Context, since time.Time) (map[string]int64, error) {
	rows, err := r.db.Pool.Query(ctx, `SELECT event_type, COUNT(*) FROM analytics_events WHERE occurred_at >= $1 GROUP BY event_type`, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]int64{}
	for rows.Next() {
		var k string
		var v int64
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		out[k] = v
	}
	return out, nil
}

func (r *analyticsRepoImpl) CompanyMetrics(ctx context.Context, companyID uuid.UUID, since time.Time) (map[string]int64, error) {
	rows, err := r.db.Pool.Query(ctx, `SELECT event_type, COUNT(*) FROM analytics_events WHERE occurred_at >= $1 AND company_id = $2 GROUP BY event_type`, since, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]int64{}
	for rows.Next() {
		var k string
		var v int64
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		out[k] = v
	}
	return out, nil
}
