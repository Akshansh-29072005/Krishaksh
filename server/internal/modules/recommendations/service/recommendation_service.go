package service

import (
	"context"

	"fmt"
	"github.com/aarcsx/krishaksh-backend/internal/models"
	diseaseRepo "github.com/aarcsx/krishaksh-backend/internal/modules/diseases/repository"
	productRepo "github.com/aarcsx/krishaksh-backend/internal/modules/products/repository"
	scanRepo "github.com/aarcsx/krishaksh-backend/internal/modules/scans/repository"
	"github.com/google/uuid"
)

// RecommendationResponse bundles all relevant data for a scan result screen
type RecommendationResponse struct {
	Scan     *models.Scan      `json:"scan"`
	Disease  *models.Disease   `json:"disease,omitempty"`
	Products []*models.Product `json:"recommended_products"`
}

type RecommendationService interface {
	GetForScan(ctx context.Context, userID uuid.UUID, scanID uuid.UUID) (*RecommendationResponse, error)
}

type recommendationServiceImpl struct {
	scanRepo    scanRepo.ScanRepository
	diseaseRepo diseaseRepo.DiseaseRepository
	productRepo productRepo.ProductRepository
}

func NewRecommendationService(
	sr scanRepo.ScanRepository,
	dr diseaseRepo.DiseaseRepository,
	pr productRepo.ProductRepository,
) RecommendationService {
	return &recommendationServiceImpl{
		scanRepo:    sr,
		diseaseRepo: dr,
		productRepo: pr,
	}
}

func (s *recommendationServiceImpl) GetForScan(ctx context.Context, userID uuid.UUID, scanID uuid.UUID) (*RecommendationResponse, error) {
	// 1. Load the scan, ensuring user owns it
	scan, err := s.scanRepo.GetScanByID(ctx, scanID, userID)
	if err != nil {
		return nil, fmt.Errorf("scan not found: %w", err)
	}

	resp := &RecommendationResponse{Scan: scan}

	// 2. If scan is complete and has a mapped disease, load intelligence
	if scan.DiseaseID != nil {
		disease, err := s.diseaseRepo.GetByID(ctx, *scan.DiseaseID)
		if err == nil {
			resp.Disease = disease
		}

		// 3. Ranked product fetching — ORDER is controlled in SQL via disease_product_mappings
		// Sponsored products naturally float to top via `is_sponsored DESC, priority ASC`
		products, err := s.productRepo.GetByDiseaseID(ctx, *scan.DiseaseID)
		if err == nil {
			resp.Products = products
		}
	}

	return resp, nil
}
