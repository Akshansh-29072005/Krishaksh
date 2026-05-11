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
	ID                uuid.UUID `db:"id" json:"id"`
	CompanyID         uuid.UUID `db:"company_id" json:"company_id"`
	CategoryID        uuid.UUID `db:"category_id" json:"category_id"`
	Name              string    `db:"name" json:"name"`
	Description       string    `db:"description" json:"description"`
	ImageURL          *string   `db:"image_url" json:"image_url"`
	Price             float64   `db:"price" json:"price"`
	Currency          string    `db:"currency" json:"currency"`
	Unit              string    `db:"unit" json:"unit"`
	IsSponsored       bool      `db:"is_sponsored" json:"is_sponsored"`
	SponsoredPriority int       `db:"sponsored_priority" json:"sponsored_priority"`
	IsActive          bool      `db:"is_active" json:"is_active"`
	StockAvailable    bool      `db:"stock_available" json:"stock_available"`
	StockQuantity     int       `db:"stock_quantity" json:"stock_quantity"`
	CropType          string    `db:"crop_type" json:"crop_type"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}
