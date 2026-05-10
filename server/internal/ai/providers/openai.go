package providers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type OpenAIProvider struct {
	apiKey string
	hc     *http.Client
}

func NewOpenAIProvider(apiKey string) *OpenAIProvider {
	return &OpenAIProvider{apiKey: apiKey, hc: &http.Client{Timeout: 20 * time.Second}}
}

func (o *OpenAIProvider) Name() string { return "openai" }

func (o *OpenAIProvider) Infer(ctx context.Context, prompt, cropType string, imageBytes []byte, mimeType string) (string, error) {
	if o.apiKey == "" {
		return "", fmt.Errorf("openai key missing")
	}
	imgData := "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(imageBytes)
	payload := map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]interface{}{{
			"role": "user",
			"content": []map[string]interface{}{
				{"type": "text", "text": prompt + "\nCrop type: " + cropType},
				{"type": "image_url", "image_url": map[string]string{"url": imgData}},
			},
		}},
		"temperature": 0.1,
	}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)
	res, err := o.hc.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		return "", fmt.Errorf("openai status %d", res.StatusCode)
	}
	var out struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return "", err
	}
	if len(out.Choices) == 0 {
		return "", fmt.Errorf("empty openai response")
	}
	return out.Choices[0].Message.Content, nil
}
