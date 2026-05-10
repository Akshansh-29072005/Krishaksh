package models

import (
	"time"

	"github.com/google/uuid"
)

// Disease captures the complete metadata for a detected plant illness
type Disease struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	CropType    string    `db:"crop_type"`
	Severity    string    `db:"severity"` // LOW, MEDIUM, HIGH, CRITICAL
	Symptoms    string    `db:"symptoms"`
	Prevention  string    `db:"prevention"`
	Description string    `db:"description"`
	ImageURL    *string   `db:"image_url"`
	CreatedAt   time.Time `db:"created_at"`
}

// Treatment maps a specific action to cure/control a disease
type Treatment struct {
	ID          uuid.UUID `db:"id"`
	DiseaseID   uuid.UUID `db:"disease_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	StepOrder   int       `db:"step_order"`
	CreatedAt   time.Time `db:"created_at"`
}

// DiseaseProductMapping is the relational bridge controlled ONLY by the backend
type DiseaseProductMapping struct {
	ID        uuid.UUID `db:"id"`
	DiseaseID uuid.UUID `db:"disease_id"`
	ProductID uuid.UUID `db:"product_id"`
	Priority  int       `db:"priority"` // 1 = highest recommended
	CreatedAt time.Time `db:"created_at"`
}
