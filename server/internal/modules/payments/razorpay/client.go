package razorpay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client interface {
	CreateOrder(ctx context.Context, amountPaise int64, currency, receipt string) (string, error)
}

type clientImpl struct {
	key, secret string
	hc          *http.Client
}

func NewClient(key, secret string) Client {
	return &clientImpl{key: key, secret: secret, hc: &http.Client{}}
}

func (c *clientImpl) CreateOrder(ctx context.Context, amountPaise int64, currency, receipt string) (string, error) {
	body, _ := json.Marshal(map[string]interface{}{"amount": amountPaise, "currency": currency, "receipt": receipt, "payment_capture": 1})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.razorpay.com/v1/orders", bytes.NewReader(body))
	req.SetBasicAuth(c.key, c.secret)
	req.Header.Set("Content-Type", "application/json")
	res, err := c.hc.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		return "", fmt.Errorf("razorpay order creation failed with status %d", res.StatusCode)
	}
	var out struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.ID, nil
}
