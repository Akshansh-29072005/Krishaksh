package queue

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aarcsx/krishaksh-backend/internal/observability"
	"github.com/hibiken/asynq"
)

const (
	TypeScanAnalyze  = "scan:analyze"
	TypePaymentEvent = "payment:event"
	TypeOrderNotify  = "order:notify"
	TypeRefundStart  = "refund:start"
	TypeAnalyticsEvt = "analytics:event"
)

type ScanAnalyzePayload struct {
	ScanID   string
	ImageURL string
	CropType string
}
type PaymentEventPayload struct {
	EventID    string
	EventType  string
	RawPayload []byte
}
type AnalyticsEventPayload struct {
	EventType string
	ActorUser string
	CompanyID string
	EntityID  string
	OrderID   string
	PaymentID string
	Metadata  map[string]interface{}
	Value     *float64
}

type QueueClient interface {
	EnqueueScanTask(scanID, imageURL, cropType string) error
	EnqueuePaymentEvent(eventID, eventType string, rawPayload []byte) error
	EnqueueOrderNotification(orderID, userID, status string) error
	EnqueueRefund(orderID, paymentID, reason string) error
	EnqueueAnalyticsEvent(payload AnalyticsEventPayload) error
	Close()
}

type queueClientImpl struct {
	client *asynq.Client
}

func NewQueueClient(redisAddr string) QueueClient {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	return &queueClientImpl{client: client}
}

func (q *queueClientImpl) EnqueueScanTask(scanID, imageURL, cropType string) error {
	payload, err := json.Marshal(ScanAnalyzePayload{
		ScanID:   scanID,
		ImageURL: imageURL,
		CropType: cropType,
	})
	if err != nil {
		return fmt.Errorf("queue marshal payload failed: %w", err)
	}

	task := asynq.NewTask(TypeScanAnalyze, payload, asynq.MaxRetry(3), asynq.Queue("scans"))
	info, err := q.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("queue enqueue failed: %w", err)
	}

	log.Printf("Successfully enqueued task: id=%s queue=%s", info.ID, info.Queue)
	observability.M.Inc("queue_enqueue_total:scan")
	return nil
}

func (q *queueClientImpl) Close() {
	if q.client != nil {
		q.client.Close()
	}
}

func (q *queueClientImpl) EnqueuePaymentEvent(eventID, eventType string, rawPayload []byte) error {
	payload, err := json.Marshal(PaymentEventPayload{EventID: eventID, EventType: eventType, RawPayload: rawPayload})
	if err != nil {
		return fmt.Errorf("queue marshal payment payload failed: %w", err)
	}
	_, err = q.client.Enqueue(asynq.NewTask(TypePaymentEvent, payload, asynq.MaxRetry(10), asynq.Queue("payments")))
	if err == nil {
		observability.M.Inc("queue_enqueue_total:payment")
	}
	return err
}

func (q *queueClientImpl) EnqueueOrderNotification(orderID, userID, status string) error {
	payload, err := json.Marshal(map[string]string{"order_id": orderID, "user_id": userID, "status": status})
	if err != nil {
		return err
	}
	_, err = q.client.Enqueue(asynq.NewTask(TypeOrderNotify, payload, asynq.MaxRetry(5), asynq.Queue("notifications")))
	if err == nil {
		observability.M.Inc("queue_enqueue_total:order_notification")
	}
	return err
}

func (q *queueClientImpl) EnqueueRefund(orderID, paymentID, reason string) error {
	payload, err := json.Marshal(map[string]string{"order_id": orderID, "payment_id": paymentID, "reason": reason})
	if err != nil {
		return err
	}
	_, err = q.client.Enqueue(asynq.NewTask(TypeRefundStart, payload, asynq.MaxRetry(7), asynq.Queue("refunds")))
	if err == nil {
		observability.M.Inc("queue_enqueue_total:refund")
	}
	return err
}

func (q *queueClientImpl) EnqueueAnalyticsEvent(payload AnalyticsEventPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = q.client.Enqueue(asynq.NewTask(TypeAnalyticsEvt, data, asynq.MaxRetry(5), asynq.Queue("analytics")))
	if err == nil {
		observability.M.Inc("queue_enqueue_total:analytics")
	}
	return err
}
