package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"github.com/aarcsx/krisho-backend/internal/config"
	"github.com/aarcsx/krisho-backend/internal/observability"
	"google.golang.org/protobuf/types/known/durationpb"
)

type CloudTasksClient struct {
	client    *cloudtasks.Client
	config    *config.Config
	queuePath string
}

const cloudTaskDispatchDeadline = 10 * time.Minute

func NewCloudTasksClient(ctx context.Context, cfg *config.Config) (QueueClient, error) {
	log.Printf("Initializing Cloud Tasks client...")
	log.Printf("Config params: ProjectID=%s, Location=%s, Queue=%s, WorkerURL=%s",
		cfg.GCPProjectID, cfg.GCPLocation, cfg.CloudTasksQueue, cfg.CloudTasksWorkerURL)

	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		log.Printf("ERROR: Failed to create cloud tasks client: %v", err)
		return nil, fmt.Errorf("failed to create cloud tasks client: %w", err)
	}

	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", cfg.GCPProjectID, cfg.GCPLocation, cfg.CloudTasksQueue)
	log.Printf("Using Cloud Tasks queue path: %s", queuePath)

	return &CloudTasksClient{
		client:    client,
		config:    cfg,
		queuePath: queuePath,
	}, nil
}

func (c *CloudTasksClient) EnqueueScanTask(scanID, imageURL, cropType string) error {
	log.Printf("Attempting to enqueue scan task: ScanID=%s, CropType=%s", scanID, cropType)

	payload, err := json.Marshal(ScanAnalyzePayload{
		ScanID:   scanID,
		ImageURL: imageURL,
		CropType: cropType,
	})
	if err != nil {
		log.Printf("ERROR: Queue marshal payload failed: %v", err)
		return fmt.Errorf("queue marshal payload failed: %w", err)
	}

	workerURL := strings.TrimSuffix(c.config.CloudTasksWorkerURL, "/") + "/tasks/scan"
	log.Printf("Worker endpoint URL: %s", workerURL)

	req := &cloudtaskspb.CreateTaskRequest{
		Parent: c.queuePath,
		Task: &cloudtaskspb.Task{
			MessageType: &cloudtaskspb.Task_HttpRequest{
				HttpRequest: &cloudtaskspb.HttpRequest{
					HttpMethod: cloudtaskspb.HttpMethod_POST,
					Url:        workerURL,
					Body:       payload,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
			},
			DispatchDeadline: durationpb.New(cloudTaskDispatchDeadline),
		},
	}

	log.Printf("Calling Cloud Tasks API to create task...")
	task, err := c.client.CreateTask(context.Background(), req)
	if err != nil {
		log.Printf("ERROR: Failed to create cloud task: %v", err)
		return fmt.Errorf("failed to create cloud task: %w", err)
	}

	log.Printf("SUCCESS: Enqueued Cloud Task: id=%s", task.Name)
	observability.M.Inc("queue_enqueue_total:scan")
	return nil
}

func (c *CloudTasksClient) EnqueuePaymentEvent(eventID, eventType string, rawPayload []byte) error {
	payload, err := json.Marshal(PaymentEventPayload{
		EventID:    eventID,
		EventType:  eventType,
		RawPayload: rawPayload,
	})
	if err != nil {
		return fmt.Errorf("queue marshal payment payload failed: %w", err)
	}

	workerURL := strings.TrimSuffix(c.config.CloudTasksWorkerURL, "/") + "/tasks/payment"

	req := &cloudtaskspb.CreateTaskRequest{
		Parent: c.queuePath,
		Task: &cloudtaskspb.Task{
			MessageType: &cloudtaskspb.Task_HttpRequest{
				HttpRequest: &cloudtaskspb.HttpRequest{
					HttpMethod: cloudtaskspb.HttpMethod_POST,
					Url:        workerURL,
					Body:       payload,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
			},
			DispatchDeadline: durationpb.New(cloudTaskDispatchDeadline),
		},
	}

	_, err = c.client.CreateTask(context.Background(), req)
	if err == nil {
		observability.M.Inc("queue_enqueue_total:payment")
	}
	return err
}

func (c *CloudTasksClient) EnqueueOrderNotification(orderID, userID, status string) error {
	payload, err := json.Marshal(map[string]string{"order_id": orderID, "user_id": userID, "status": status})
	if err != nil {
		return err
	}

	workerURL := strings.TrimSuffix(c.config.CloudTasksWorkerURL, "/") + "/tasks/notification"

	req := &cloudtaskspb.CreateTaskRequest{
		Parent: c.queuePath,
		Task: &cloudtaskspb.Task{
			MessageType: &cloudtaskspb.Task_HttpRequest{
				HttpRequest: &cloudtaskspb.HttpRequest{
					HttpMethod: cloudtaskspb.HttpMethod_POST,
					Url:        workerURL,
					Body:       payload,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
			},
			DispatchDeadline: durationpb.New(cloudTaskDispatchDeadline),
		},
	}

	_, err = c.client.CreateTask(context.Background(), req)
	if err == nil {
		observability.M.Inc("queue_enqueue_total:order_notification")
	}
	return err
}

func (c *CloudTasksClient) EnqueueRefund(orderID, paymentID, reason string) error {
	payload, err := json.Marshal(map[string]string{"order_id": orderID, "payment_id": paymentID, "reason": reason})
	if err != nil {
		return err
	}

	workerURL := strings.TrimSuffix(c.config.CloudTasksWorkerURL, "/") + "/tasks/refund"

	req := &cloudtaskspb.CreateTaskRequest{
		Parent: c.queuePath,
		Task: &cloudtaskspb.Task{
			MessageType: &cloudtaskspb.Task_HttpRequest{
				HttpRequest: &cloudtaskspb.HttpRequest{
					HttpMethod: cloudtaskspb.HttpMethod_POST,
					Url:        workerURL,
					Body:       payload,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
			},
			DispatchDeadline: durationpb.New(cloudTaskDispatchDeadline),
		},
	}

	_, err = c.client.CreateTask(context.Background(), req)
	if err == nil {
		observability.M.Inc("queue_enqueue_total:refund")
	}
	return err
}

func (c *CloudTasksClient) EnqueueAnalyticsEvent(payload AnalyticsEventPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	workerURL := strings.TrimSuffix(c.config.CloudTasksWorkerURL, "/") + "/tasks/analytics"

	req := &cloudtaskspb.CreateTaskRequest{
		Parent: c.queuePath,
		Task: &cloudtaskspb.Task{
			MessageType: &cloudtaskspb.Task_HttpRequest{
				HttpRequest: &cloudtaskspb.HttpRequest{
					HttpMethod: cloudtaskspb.HttpMethod_POST,
					Url:        workerURL,
					Body:       data,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
			},
			DispatchDeadline: durationpb.New(cloudTaskDispatchDeadline),
		},
	}

	_, err = c.client.CreateTask(context.Background(), req)
	if err == nil {
		observability.M.Inc("queue_enqueue_total:analytics")
	}
	return err
}

func (c *CloudTasksClient) Close() {
	if c.client != nil {
		_ = c.client.Close()
	}
}
