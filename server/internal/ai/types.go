package ai

import "time"

type Prediction struct {
	Disease    string   `json:"disease"`
	Confidence float64  `json:"confidence"`
	Symptoms   []string `json:"symptoms"`
}

type VisionResponse struct {
	CropType    string       `json:"crop_type"`
	Predictions []Prediction `json:"predictions"`
	Notes       string       `json:"notes"`
}

type InferenceResult struct {
	Provider          string
	Response          VisionResponse
	PrimaryPrediction Prediction
	Uncertain         bool
	Latency           time.Duration
	PromptVersion     string
	RawResponse       string
	InputBytes        int
	OutputBytes       int
}

type InferenceRequest struct {
	ImageURL  string
	CropType  string
	Timeout   time.Duration
	TraceID   string
	RequestID string
}
