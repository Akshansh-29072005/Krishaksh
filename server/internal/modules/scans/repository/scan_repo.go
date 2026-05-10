package repository

import (
	"context"
	"encoding/json"

	"github.com/aarcsx/krisho-backend/internal/database"
	"github.com/aarcsx/krisho-backend/internal/models"
	"github.com/google/uuid"
)

type ScanRepository interface {
	CreateScan(ctx context.Context, scan *models.Scan) error
	GetScanByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Scan, error)
	GetUserScans(ctx context.Context, userID uuid.UUID) ([]*models.Scan, error)
	UpdateScanStatus(ctx context.Context, id uuid.UUID, status string) error
	UpdateScanResult(ctx context.Context, id uuid.UUID, diseaseID *uuid.UUID, status string) error
	UpdateAIInference(ctx context.Context, id uuid.UUID, status string, provider *string, confidence *float64, processingErr *string, metadata map[string]interface{}) error
}

type scanRepoImpl struct {
	db *database.DB
}

func NewScanRepository(db *database.DB) ScanRepository {
	return &scanRepoImpl{db: db}
}

func (r *scanRepoImpl) CreateScan(ctx context.Context, scan *models.Scan) error {
	query := `INSERT INTO scans (id, user_id, image_url, crop_type, prediction_status, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Pool.Exec(ctx, query, scan.ID, scan.UserID, scan.ImageURL, scan.CropType, scan.PredictionStatus, scan.CreatedAt)
	return err
}

func (r *scanRepoImpl) GetScanByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Scan, error) {
	query := `SELECT id, user_id, image_url, crop_type, prediction_status, disease_id, ai_provider, confidence_score, processing_error, ai_metadata, completed_at, created_at 
	          FROM scans WHERE id = $1 AND user_id = $2 LIMIT 1`

	scan := &models.Scan{}
	err := r.db.Pool.QueryRow(ctx, query, id, userID).Scan(
		&scan.ID, &scan.UserID, &scan.ImageURL, &scan.CropType,
		&scan.PredictionStatus, &scan.DiseaseID, &scan.AIProvider, &scan.ConfidenceScore, &scan.ProcessingError, &scan.AIMetadata, &scan.CompletedAt, &scan.CreatedAt,
	)
	return scan, err
}

func (r *scanRepoImpl) GetUserScans(ctx context.Context, userID uuid.UUID) ([]*models.Scan, error) {
	query := `SELECT id, user_id, image_url, crop_type, prediction_status, disease_id, ai_provider, confidence_score, processing_error, ai_metadata, completed_at, created_at 
	          FROM scans WHERE user_id = $1 ORDER BY created_at DESC LIMIT 50`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scans []*models.Scan
	for rows.Next() {
		scan := &models.Scan{}
		err := rows.Scan(&scan.ID, &scan.UserID, &scan.ImageURL, &scan.CropType,
			&scan.PredictionStatus, &scan.DiseaseID, &scan.AIProvider, &scan.ConfidenceScore, &scan.ProcessingError, &scan.AIMetadata, &scan.CompletedAt, &scan.CreatedAt)
		if err != nil {
			return nil, err
		}
		scans = append(scans, scan)
	}
	return scans, nil
}

func (r *scanRepoImpl) UpdateScanStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE scans SET prediction_status = $1 WHERE id = $2`
	_, err := r.db.Pool.Exec(ctx, query, status, id)
	return err
}

func (r *scanRepoImpl) UpdateScanResult(ctx context.Context, id uuid.UUID, diseaseID *uuid.UUID, status string) error {
	query := `UPDATE scans SET prediction_status = $1, disease_id = $2 WHERE id = $3`
	_, err := r.db.Pool.Exec(ctx, query, status, diseaseID, id)
	return err
}

func (r *scanRepoImpl) UpdateAIInference(ctx context.Context, id uuid.UUID, status string, provider *string, confidence *float64, processingErr *string, metadata map[string]interface{}) error {
	meta := []byte(`{}`)
	if metadata != nil {
		b, _ := json.Marshal(metadata)
		meta = b
	}
	query := `UPDATE scans
	          SET prediction_status = $1, ai_provider = $2, confidence_score = $3, processing_error = $4, ai_metadata = $5, completed_at = NOW()
	          WHERE id = $6`
	_, err := r.db.Pool.Exec(ctx, query, status, provider, confidence, processingErr, meta, id)
	return err
}
