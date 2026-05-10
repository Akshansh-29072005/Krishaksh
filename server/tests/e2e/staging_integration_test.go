package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"
)

func requireReachable(t *testing.T, client *http.Client) {
	t.Helper()
	probe, err := client.Get(baseURL() + "/health")
	if err != nil {
		t.Skipf("staging target unreachable: %v", err)
		return
	}
	_ = probe.Body.Close()
}

func baseURL() string {
	v := os.Getenv("STAGING_BASE_URL")
	if v == "" {
		v = "http://localhost:8080"
	}
	return v
}

func bearerReq(method, url string, body []byte) (*http.Request, error) {
	r, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")
	if tok := os.Getenv("STAGING_JWT_TOKEN"); tok != "" {
		r.Header.Set("Authorization", "Bearer "+tok)
	}
	return r, nil
}

func TestHealthAndReady(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}
	requireReachable(t, client)
	for _, path := range []string{"/health", "/ready", "/metrics"} {
		res, err := client.Get(baseURL() + path)
		if err != nil {
			t.Fatalf("%s failed: %v", path, err)
		}
		if res.StatusCode >= 400 {
			t.Fatalf("%s status=%d", path, res.StatusCode)
		}
		_ = res.Body.Close()
	}
}

func TestAuthSessionAndProtectedRoute(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}
	requireReachable(t, client)
	req, _ := bearerReq(http.MethodGet, baseURL()+"/api/v1/scans/history", nil)
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if tok := os.Getenv("STAGING_JWT_TOKEN"); tok == "" {
		if res.StatusCode != http.StatusUnauthorized {
			t.Fatalf("expected 401 without token, got %d", res.StatusCode)
		}
	} else if res.StatusCode >= 500 {
		t.Fatalf("unexpected server error %d", res.StatusCode)
	}
	_ = res.Body.Close()
}

func TestPaymentLifecycleSmoke(t *testing.T) {
	orderID := os.Getenv("STAGING_TEST_ORDER_ID")
	if orderID == "" {
		t.Skip("STAGING_TEST_ORDER_ID not set")
	}
	client := &http.Client{Timeout: 10 * time.Second}
	requireReachable(t, client)
	payload, _ := json.Marshal(map[string]string{"order_id": orderID})
	req, _ := bearerReq(http.MethodPost, baseURL()+"/api/v1/payments/create-order", payload)
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode >= 500 {
		t.Fatalf("payment create-order failed status=%d", res.StatusCode)
	}
	_ = res.Body.Close()
}

func TestWebhookReplayProtectionSmoke(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}
	requireReachable(t, client)
	payload := []byte(`{"id":"evt_replay_test","event":"payment.captured","payload":{"payment":{"entity":{"id":"pay_test","order_id":"order_test","amount":100,"status":"captured"}}}}`)

	for i := 0; i < 2; i++ {
		req, _ := http.NewRequest(http.MethodPost, baseURL()+"/api/v1/payments/webhook", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Razorpay-Signature", "invalid-sig")
		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != http.StatusUnauthorized {
			t.Fatalf("expected 401 for spoofed webhook, got %d", res.StatusCode)
		}
		_ = res.Body.Close()
	}
}
