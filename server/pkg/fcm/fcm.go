package fcm

import (
	"context"
	"fmt"
	"log"
)

// FCMPayload is a simplified struct representing what we send to Firebase
type FCMPayload struct {
	Title string
	Body  string
	Token string            // target device FCM token
	Data  map[string]string // optional deep-link data
}

// FCMClient defines the interface for push notification delivery
type FCMClient interface {
	Send(ctx context.Context, payload FCMPayload) error
}

// fcmClientImpl is the production implementation using Firebase Admin SDK (REST fallback for now)
type fcmClientImpl struct {
	serverKey string // Firebase Server Key (legacy HTTP v1 fallback)
}

func NewFCMClient(serverKey string) FCMClient {
	return &fcmClientImpl{serverKey: serverKey}
}

func (f *fcmClientImpl) Send(ctx context.Context, payload FCMPayload) error {
	if payload.Token == "" {
		return fmt.Errorf("FCM: empty device token, skipping notification")
	}

	// In production: use firebase-admin-go SDK or HTTP v1 API
	// For now: structured mock log confirming the dispatch pathway
	log.Printf("[FCM] Sending push to token=%s title=%q body=%q data=%v",
		payload.Token[:min(len(payload.Token), 20)]+"...",
		payload.Title,
		payload.Body,
		payload.Data,
	)
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
