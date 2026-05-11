package models

import (
	"time"

	"github.com/google/uuid"
)

type Scan struct {
	ID               uuid.UUID  `db:"id" json:"id"`
	UserID           uuid.UUID  `db:"user_id" json:"user_id"`
	ImageURL         string     `db:"image_url" json:"image_url"`
	CropType         string     `db:"crop_type" json:"crop_type"`
	PredictionStatus string     `db:"prediction_status" json:"prediction_status"`
	DiseaseID        *uuid.UUID `db:"disease_id" json:"disease_id"`
	DiseaseName      *string    `db:"disease_name" json:"disease_name"`     // AI predicted name
	AISymptoms       []string   `db:"ai_symptoms" json:"ai_symptoms"`       // AI extracted symptoms
	AIProvider       *string    `db:"ai_provider" json:"ai_provider"`
	ConfidenceScore  *float64   `db:"confidence_score" json:"confidence_score"`
	ProcessingError  *string    `db:"processing_error" json:"processing_error"`
	AIMetadata       []byte     `db:"ai_metadata" json:"ai_metadata"`
	CompletedAt      *time.Time `db:"completed_at" json:"completed_at"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
}
