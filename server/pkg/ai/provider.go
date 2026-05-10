package ai

import (
	"context"
	"fmt"
)

// AIAnalysisResult limits the returned data, ensuring the backend maps treatments later
type AIAnalysisResult struct {
	PredictedDisease string  `json:"predicted_disease"`
	Symptoms         string  `json:"symptoms"`
	ConfidenceScore  float64 `json:"confidence_score"`
}

// AIProvider interface guarantees loose coupling with internal providers or external ones (Gemini/OpenAI)
type AIProvider interface {
	AnalyzeCropImage(ctx context.Context, imageURL string, cropType string) (*AIAnalysisResult, error)
}

// GeminiProvider handles requests directly formatted for Google's Gemini Vision API
type GeminiProvider struct {
	apiKey string
}

func NewGeminiProvider(apiKey string) AIProvider {
	return &GeminiProvider{apiKey: apiKey}
}

func (p *GeminiProvider) AnalyzeCropImage(ctx context.Context, imageURL string, cropType string) (*AIAnalysisResult, error) {
	// MOCKED GEMINI CALL
	fmt.Printf("[Gemini API] Processing image from S3: %s for %s\n", imageURL, cropType)
	return &AIAnalysisResult{
		PredictedDisease: "Early Blight",
		Symptoms:         "Brown spots with concentric rings on lower leaves",
		ConfidenceScore:  0.94,
	}, nil
}

// OpenAIProvider handles requests targeting GPT-4o-vision
type OpenAIProvider struct {
	apiKey string
}

func NewOpenAIProvider(apiKey string) AIProvider {
	return &OpenAIProvider{apiKey: apiKey}
}

func (p *OpenAIProvider) AnalyzeCropImage(ctx context.Context, imageURL string, cropType string) (*AIAnalysisResult, error) {
	// MOCKED OPENAI CALL
	fmt.Printf("[OpenAI API] Processing image from S3: %s for %s\n", imageURL, cropType)
	return &AIAnalysisResult{
		PredictedDisease: "Healthy",
		Symptoms:         "None detected",
		ConfidenceScore:  0.98,
	}, nil
}
