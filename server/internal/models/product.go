package models

import (
	"time"

	"github.com/google/uuid"
)

// Company represents a registered agrochemical/product company
type Company struct {
	ID         uuid.UUID `db:"id"`
	Name       string    `db:"name"`
	LogoURL    *string   `db:"logo_url"`
	Region     *string   `db:"region"`
	IsVerified bool      `db:"is_verified"`
	CreatedAt  time.Time `db:"created_at"`
}

// ProductCategory groups products (e.g., Fungicide, Fertilizer, Pesticide)
type ProductCategory struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
	Slug string    `db:"slug"`
}

// Product is the core commerce entity
type Product struct {
	ID                uuid.UUID `db:"id"`
	CompanyID         uuid.UUID `db:"company_id"`
	CategoryID        uuid.UUID `db:"category_id"`
	Name              string    `db:"name"`
	Description       string    `db:"description"`
	ImageURL          *string   `db:"image_url"`
	Price             float64   `db:"price"`
	Currency          string    `db:"currency"` // INR
	Unit              string    `db:"unit"`     // e.g., "500ml", "1kg"
	IsSponsored       bool      `db:"is_sponsored"`
	SponsoredPriority int       `db:"sponsored_priority"` // higher = ranked first in listing
	IsActive          bool      `db:"is_active"`
	StockAvailable    bool      `db:"stock_available"`
	StockQuantity     int       `db:"stock_quantity"`
	CropType          string    `db:"crop_type"` // primary crop compatibility
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}
