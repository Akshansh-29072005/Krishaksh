package models

import (
	"time"

	"github.com/google/uuid"
)

type Scan struct {
	ID               uuid.UUID  `db:"id"`
	UserID           uuid.UUID  `db:"user_id"`
	ImageURL         string     `db:"image_url"`
	CropType         string     `db:"crop_type"`
	PredictionStatus string     `db:"prediction_status"` // UPLOADED, QUEUED, PROCESSING, COMPLETED, FAILED
	DiseaseID        *uuid.UUID `db:"disease_id"`        // Linked once AI analysis completes
	AIProvider       *string    `db:"ai_provider"`
	ConfidenceScore  *float64   `db:"confidence_score"`
	ProcessingError  *string    `db:"processing_error"`
	AIMetadata       []byte     `db:"ai_metadata"`
	CompletedAt      *time.Time `db:"completed_at"`
	CreatedAt        time.Time  `db:"created_at"`
}
