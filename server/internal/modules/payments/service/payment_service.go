package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aarcsx/krishaksh-backend/internal/models"
	orderRepo "github.com/aarcsx/krishaksh-backend/internal/modules/orders/repository"
	"github.com/aarcsx/krishaksh-backend/internal/modules/payments/razorpay"
	paymentRepo "github.com/aarcsx/krishaksh-backend/internal/modules/payments/repository"
	"github.com/aarcsx/krishaksh-backend/pkg/queue"
	"github.com/google/uuid"
)

type PaymentService interface {
	CreateRazorpayOrder(ctx context.Context, userID, orderID uuid.UUID) (map[string]interface{}, error)
	HandleWebhook(ctx context.Context, signature string, payload []byte) error
}

type paymentServiceImpl struct {
	orders        orderRepo.OrderRepository
	payments      paymentRepo.PaymentRepository
	razor         razorpay.Client
	queueClient   queue.QueueClient
	webhookSecret string
	rpKeyID       string
	rpKeySecret   string
}

func NewPaymentService(o orderRepo.OrderRepository, p paymentRepo.PaymentRepository, r razorpay.Client, q queue.QueueClient, webhookSecret, keyID, keySecret string) PaymentService {
	return &paymentServiceImpl{orders: o, payments: p, razor: r, queueClient: q, webhookSecret: webhookSecret, rpKeyID: keyID, rpKeySecret: keySecret}
}

func (s *paymentServiceImpl) CreateRazorpayOrder(ctx context.Context, userID, orderID uuid.UUID) (map[string]interface{}, error) {
	order, err := s.orders.GetOrderByID(ctx, orderID, userID)
	if err != nil {
		return nil, errors.New("order not found")
	}
	amountPaise := int64(order.GrandTotal * 100)
	providerOrderID, err := s.razor.CreateOrder(ctx, amountPaise, order.Currency, order.ID.String())
	if err != nil {
		return nil, err
	}
	pay := &models.Payment{ID: uuid.New(), OrderID: order.ID, Provider: "razorpay", ProviderOrderID: providerOrderID, Status: "created", Amount: order.GrandTotal, Currency: order.Currency, Metadata: json.RawMessage(`{}`)}
	if err := s.payments.CreatePayment(ctx, pay); err != nil {
		return nil, err
	}
	_ = s.orders.LinkPayment(ctx, order.ID, pay.ID)
	return map[string]interface{}{"key_id": s.rpKeyID, "order_id": providerOrderID, "amount": amountPaise, "currency": order.Currency, "receipt": order.ID.String()}, nil
}

func (s *paymentServiceImpl) HandleWebhook(ctx context.Context, signature string, payload []byte) error {
	if !validateHMAC(signature, payload, s.webhookSecret) {
		return errors.New("invalid webhook signature")
	}
	var env struct {
		ID      string `json:"id"`
		Event   string `json:"event"`
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
	if err := json.Unmarshal(payload, &env); err != nil {
		return err
	}
	we := &models.WebhookEvent{ID: uuid.New(), Provider: "razorpay", EventID: env.ID, EventType: env.Event, Signature: &signature, Payload: payload}
	if err := s.payments.UpsertWebhookEvent(ctx, we); err != nil {
		return err
	}
	if err := s.queueClient.EnqueuePaymentEvent(env.ID, env.Event, payload); err != nil {
		return err
	}
	return nil
}

func ValidateCheckoutSignature(orderID, paymentID, signature, secret string) bool {
	msg := fmt.Sprintf("%s|%s", orderID, paymentID)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(msg))
	return hex.EncodeToString(h.Sum(nil)) == signature
}

func validateHMAC(signature string, payload []byte, secret string) bool {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil)) == signature
}
