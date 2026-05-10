package repository

import (
	"context"

	"github.com/aarcsx/krisho-backend/internal/database"
	"github.com/aarcsx/krisho-backend/internal/models"
	"github.com/google/uuid"
)

type DiseaseRepository interface {
	GetAll(ctx context.Context) ([]*models.Disease, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Disease, error)
	GetTreatments(ctx context.Context, diseaseID uuid.UUID) ([]*models.Treatment, error)
}

type diseaseRepoImpl struct {
	db *database.DB
}

func NewDiseaseRepository(db *database.DB) DiseaseRepository {
	return &diseaseRepoImpl{db: db}
}

func (r *diseaseRepoImpl) GetAll(ctx context.Context) ([]*models.Disease, error) {
	query := `SELECT id, name, crop_type, severity, symptoms, prevention, description, image_url, created_at 
	          FROM diseases ORDER BY name ASC`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var diseases []*models.Disease
	for rows.Next() {
		d := &models.Disease{}
		if err := rows.Scan(&d.ID, &d.Name, &d.CropType, &d.Severity, &d.Symptoms, &d.Prevention, &d.Description, &d.ImageURL, &d.CreatedAt); err != nil {
			return nil, err
		}
		diseases = append(diseases, d)
	}
	return diseases, nil
}

func (r *diseaseRepoImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Disease, error) {
	query := `SELECT id, name, crop_type, severity, symptoms, prevention, description, image_url, created_at 
	          FROM diseases WHERE id = $1`
	d := &models.Disease{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&d.ID, &d.Name, &d.CropType, &d.Severity, &d.Symptoms, &d.Prevention, &d.Description, &d.ImageURL, &d.CreatedAt,
	)
	return d, err
}

func (r *diseaseRepoImpl) GetTreatments(ctx context.Context, diseaseID uuid.UUID) ([]*models.Treatment, error) {
	query := `SELECT id, disease_id, title, description, step_order, created_at 
	          FROM treatments WHERE disease_id = $1 ORDER BY step_order ASC`
	rows, err := r.db.Pool.Query(ctx, query, diseaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var treatments []*models.Treatment
	for rows.Next() {
		t := &models.Treatment{}
		if err := rows.Scan(&t.ID, &t.DiseaseID, &t.Title, &t.Description, &t.StepOrder, &t.CreatedAt); err != nil {
			return nil, err
		}
		treatments = append(treatments, t)
	}
	return treatments, nil
}
