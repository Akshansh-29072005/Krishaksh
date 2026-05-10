package models

import (
	"time"

	"github.com/google/uuid"
)

// AdvertisementCampaign tracks sponsored promotions from companies
type AdvertisementCampaign struct {
	ID          uuid.UUID  `db:"id"`
	CompanyID   uuid.UUID  `db:"company_id"`
	ProductID   *uuid.UUID `db:"product_id"` // optional product tie-in
	Title       string     `db:"title"`
	ImageURL    string     `db:"image_url"`
	TargetURL   *string    `db:"target_url"`
	Region      *string    `db:"region"` // nil = nationwide
	IsActive    bool       `db:"is_active"`
	StartsAt    time.Time  `db:"starts_at"`
	EndsAt      time.Time  `db:"ends_at"`
	Impressions int64      `db:"impressions"` // updated by background job
	Clicks      int64      `db:"clicks"`
	CreatedAt   time.Time  `db:"created_at"`
}

// RecommendationRule drives dynamic backend rankings — no hardcoded logic
type RecommendationRule struct {
	ID          uuid.UUID `db:"id"`
	DiseaseID   uuid.UUID `db:"disease_id"`
	ProductID   uuid.UUID `db:"product_id"`
	Priority    int       `db:"priority"`     // 1 = highest
	IsSponsored bool      `db:"is_sponsored"` // sponsored products float to top
	Region      *string   `db:"region"`
	IsActive    bool      `db:"is_active"`
	CreatedAt   time.Time `db:"created_at"`
}
