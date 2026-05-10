package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aarcsx/krisho-backend/internal/models"
	notifRepo "github.com/aarcsx/krisho-backend/internal/modules/notifications/repository"
	"github.com/aarcsx/krisho-backend/internal/observability"
	"github.com/aarcsx/krisho-backend/pkg/fcm"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const TypeNotificationSend = "notification:send"

type NotificationPayload struct {
	UserID string
	Type   string
	Title  string
	Body   string
	Data   map[string]string
}

type NotificationWorker struct {
	repo      notifRepo.NotificationRepository
	fcmClient fcm.FCMClient
}

func NewNotificationWorker(repo notifRepo.NotificationRepository, fcmClient fcm.FCMClient) *NotificationWorker {
	return &NotificationWorker{repo: repo, fcmClient: fcmClient}
}

func EnqueueNotification(client *asynq.Client, payload NotificationPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("notification marshal error: %w", err)
	}
	task := asynq.NewTask(TypeNotificationSend, data, asynq.MaxRetry(3), asynq.Queue("notifications"))
	_, err = client.Enqueue(task)
	return err
}

func (w *NotificationWorker) HandleNotificationTask(ctx context.Context, t *asynq.Task) error {
	var payload NotificationPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		observability.M.Inc("worker_errors_total:notification")
		return fmt.Errorf("invalid notification payload: %w", err)
	}
	observability.InitLogger().Info("worker_notification_started", "user_id", payload.UserID, "type", payload.Type)

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		observability.M.Inc("worker_errors_total:notification")
		return fmt.Errorf("invalid user ID in notification: %w", err)
	}

	// 1. Persist notification record in DB
	notif := &models.Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      payload.Type,
		Title:     payload.Title,
		Body:      payload.Body,
		IsRead:    false,
		CreatedAt: time.Now(),
	}
	if err := w.repo.SaveNotification(ctx, notif); err != nil {
		log.Printf("[NotifWorker] Failed to save notification: %v", err)
	}

	// 2. Fetch all device tokens for this user
	tokens, err := w.repo.GetDeviceTokens(ctx, userID)
	if err != nil || len(tokens) == 0 {
		log.Printf("[NotifWorker] No device tokens for user %s — skipping FCM", userID)
		return nil
	}

	// 3. Dispatch FCM push to all registered devices
	for _, token := range tokens {
		fcmPayload := fcm.FCMPayload{
			Title: payload.Title,
			Body:  payload.Body,
			Token: token,
			Data:  payload.Data,
		}
		if err := w.fcmClient.Send(ctx, fcmPayload); err != nil {
			log.Printf("[NotifWorker] FCM send failed for token: %v", err)
			// Don't return error — retry specific token next iteration, not entire job
		}
	}

	// 4. Mark delivered in DB
	_ = w.repo.MarkDelivered(ctx, notif.ID)
	observability.InitLogger().Info("worker_notification_delivered", "notification_id", notif.ID.String(), "user_id", userID.String())
	return nil
}
