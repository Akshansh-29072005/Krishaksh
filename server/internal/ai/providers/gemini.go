package providers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GeminiProvider struct {
	apiKey string
	hc     *http.Client
}

func NewGeminiProvider(apiKey string) *GeminiProvider {
	return &GeminiProvider{apiKey: apiKey, hc: &http.Client{Timeout: 20 * time.Second}}
}

func (g *GeminiProvider) Name() string { return "gemini" }

func (g *GeminiProvider) Infer(ctx context.Context, prompt, cropType string, imageBytes []byte, mimeType string) (string, error) {
	if g.apiKey == "" {
		return "", fmt.Errorf("gemini key missing")
	}
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-flash-latest:generateContent"
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{{
			"parts": []map[string]interface{}{
				{"text": prompt + "\nCrop type: " + cropType},
				{"inline_data": map[string]string{"mime_type": mimeType, "data": base64.StdEncoding.EncodeToString(imageBytes)}},
			},
		}},
	}
	b, err := json.Marshal(payload)	
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", g.apiKey)
	res, err := g.hc.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		body, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("gemini status %d: %s", res.StatusCode, string(body))
	}
	var out struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return "", err
	}
	if len(out.Candidates) == 0 || len(out.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty gemini response")
	}
	return out.Candidates[0].Content.Parts[0].Text, nil
}
