package repository

import (
	"context"

	"github.com/aarcsx/krishaksh-backend/internal/database"
	"github.com/aarcsx/krishaksh-backend/internal/models"
	"github.com/google/uuid"
)

type ProductRepository interface {
	GetAll(ctx context.Context, cropType string, sponsored bool) ([]*models.Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	GetByDiseaseID(ctx context.Context, diseaseID uuid.UUID) ([]*models.Product, error)
}

type productRepoImpl struct {
	db *database.DB
}

func NewProductRepository(db *database.DB) ProductRepository {
	return &productRepoImpl{db: db}
}

func (r *productRepoImpl) GetAll(ctx context.Context, cropType string, sponsored bool) ([]*models.Product, error) {
	// Sponsored products float to top; secondary sort by priority then name
	query := `SELECT id, company_id, category_id, name, description, image_url, price, currency, unit,
	                 is_sponsored, sponsored_priority, is_active, stock_available, stock_quantity, crop_type, created_at, updated_at
	          FROM products WHERE is_active = true
	          ORDER BY is_sponsored DESC, sponsored_priority DESC, name ASC
	          LIMIT 50`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		p := &models.Product{}
		if err := rows.Scan(
			&p.ID, &p.CompanyID, &p.CategoryID, &p.Name, &p.Description, &p.ImageURL,
			&p.Price, &p.Currency, &p.Unit, &p.IsSponsored, &p.SponsoredPriority,
			&p.IsActive, &p.StockAvailable, &p.StockQuantity, &p.CropType, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *productRepoImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	query := `SELECT id, company_id, category_id, name, description, image_url, price, currency, unit,
	                 is_sponsored, sponsored_priority, is_active, stock_available, stock_quantity, crop_type, created_at, updated_at
	          FROM products WHERE id = $1`
	p := &models.Product{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.CompanyID, &p.CategoryID, &p.Name, &p.Description, &p.ImageURL,
		&p.Price, &p.Currency, &p.Unit, &p.IsSponsored, &p.SponsoredPriority,
		&p.IsActive, &p.StockAvailable, &p.StockQuantity, &p.CropType, &p.CreatedAt, &p.UpdatedAt,
	)
	return p, err
}

func (r *productRepoImpl) GetByDiseaseID(ctx context.Context, diseaseID uuid.UUID) ([]*models.Product, error) {
	// JOIN through disease_product_mappings for db-driven recommendation logic
	query := `SELECT p.id, p.company_id, p.category_id, p.name, p.description, p.image_url, p.price, p.currency, p.unit,
	                 p.is_sponsored, p.sponsored_priority, p.is_active, p.stock_available, p.stock_quantity, p.crop_type, p.created_at, p.updated_at
	          FROM products p
	          JOIN disease_product_mappings m ON p.id = m.product_id
	          WHERE m.disease_id = $1 AND p.is_active = true
	          ORDER BY p.is_sponsored DESC, m.priority ASC, p.sponsored_priority DESC`

	rows, err := r.db.Pool.Query(ctx, query, diseaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		p := &models.Product{}
		if err := rows.Scan(
			&p.ID, &p.CompanyID, &p.CategoryID, &p.Name, &p.Description, &p.ImageURL,
			&p.Price, &p.Currency, &p.Unit, &p.IsSponsored, &p.SponsoredPriority,
			&p.IsActive, &p.StockAvailable, &p.StockQuantity, &p.CropType, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
