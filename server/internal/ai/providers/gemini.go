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
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", g.apiKey)
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{{
			"parts": []map[string]interface{}{
				{"text": prompt + "\nCrop type: " + cropType},
				{"inline_data": map[string]string{"mime_type": mimeType, "data": base64.StdEncoding.EncodeToString(imageBytes)}},
			},
		}},
		"generationConfig": map[string]interface{}{"temperature": 0.1, "responseMimeType": "application/json"},
	}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	res, err := g.hc.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		return "", fmt.Errorf("gemini status %d", res.StatusCode)
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
