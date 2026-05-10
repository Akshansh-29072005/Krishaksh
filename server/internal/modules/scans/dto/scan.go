package dto

import "time"

type CreateScanRequest struct {
	CropType string `json:"crop_type" binding:"required"`
	ImageKey string `json:"image_key" binding:"required"` // The path where it was uploaded to S3
}

type ScanResponse struct {
	ID               string    `json:"id"`
	ImageURL         string    `json:"image_url"`
	CropType         string    `json:"crop_type"`
	PredictionStatus string    `json:"prediction_status"`
	DiseaseName      *string   `json:"disease_name,omitempty"`
	ConfidenceScore  *float64  `json:"confidence_score,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}
