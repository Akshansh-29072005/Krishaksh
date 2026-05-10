package repository

import (
	"context"

	"github.com/aarcsx/krisho-backend/internal/database"
	"github.com/aarcsx/krisho-backend/internal/models"
	"github.com/google/uuid"
)

type NotificationRepository interface {
	SaveNotification(ctx context.Context, n *models.Notification) error
	GetUserNotifications(ctx context.Context, userID uuid.UUID) ([]*models.Notification, error)
	MarkDelivered(ctx context.Context, id uuid.UUID) error
	UpsertDeviceToken(ctx context.Context, token *models.DeviceToken) error
	GetDeviceTokens(ctx context.Context, userID uuid.UUID) ([]string, error)
}

type notifRepoImpl struct {
	db *database.DB
}

func NewNotificationRepository(db *database.DB) NotificationRepository {
	return &notifRepoImpl{db: db}
}

func (r *notifRepoImpl) SaveNotification(ctx context.Context, n *models.Notification) error {
	q := `INSERT INTO notifications (id, user_id, type, title, body, is_read, payload, created_at)
	      VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Pool.Exec(ctx, q, n.ID, n.UserID, n.Type, n.Title, n.Body, n.IsRead, n.Payload, n.CreatedAt)
	return err
}

func (r *notifRepoImpl) GetUserNotifications(ctx context.Context, userID uuid.UUID) ([]*models.Notification, error) {
	q := `SELECT id, user_id, type, title, body, is_read, delivered_at, payload, created_at
	      FROM notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT 50`
	rows, err := r.db.Pool.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifs []*models.Notification
	for rows.Next() {
		n := &models.Notification{}
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body, &n.IsRead, &n.DeliveredAt, &n.Payload, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifs = append(notifs, n)
	}
	return notifs, nil
}

func (r *notifRepoImpl) MarkDelivered(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `UPDATE notifications SET delivered_at = NOW(), is_read = true WHERE id = $1`, id)
	return err
}

func (r *notifRepoImpl) UpsertDeviceToken(ctx context.Context, token *models.DeviceToken) error {
	// On conflict for same user+token just refresh updated_at
	q := `INSERT INTO device_tokens (id, user_id, token, platform, created_at, updated_at)
	      VALUES ($1, $2, $3, $4, $5, $6)
	      ON CONFLICT (user_id, token) DO UPDATE SET updated_at = EXCLUDED.updated_at`
	_, err := r.db.Pool.Exec(ctx, q, token.ID, token.UserID, token.Token, token.Platform, token.CreatedAt, token.UpdatedAt)
	return err
}

func (r *notifRepoImpl) GetDeviceTokens(ctx context.Context, userID uuid.UUID) ([]string, error) {
	q := `SELECT token FROM device_tokens WHERE user_id = $1`
	rows, err := r.db.Pool.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		tokens = append(tokens, t)
	}
	return tokens, nil
}
