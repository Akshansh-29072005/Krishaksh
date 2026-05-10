package repository

import (
	"context"

	"github.com/aarcsx/krishaksh-backend/internal/database"
	"github.com/aarcsx/krishaksh-backend/internal/models"
)

type AdvertisementRepository interface {
	GetFeaturedAds(ctx context.Context, region *string) ([]*models.AdvertisementCampaign, error)
	IncrementImpressions(ctx context.Context, adID string) error
	IncrementClicks(ctx context.Context, adID string) error
}

type adRepoImpl struct {
	db *database.DB
}

func NewAdvertisementRepository(db *database.DB) AdvertisementRepository {
	return &adRepoImpl{db: db}
}

func (r *adRepoImpl) GetFeaturedAds(ctx context.Context, region *string) ([]*models.AdvertisementCampaign, error) {
	// Region-aware: null region = nationwide campaign visible to everyone
	query := `SELECT id, company_id, product_id, title, image_url, target_url, region, is_active, starts_at, ends_at, impressions, clicks, created_at
	          FROM advertisement_campaigns
	          WHERE is_active = true AND starts_at <= NOW() AND ends_at >= NOW()
	            AND (region IS NULL OR region = $1)
	          ORDER BY created_at DESC LIMIT 5`

	rows, err := r.db.Pool.Query(ctx, query, region)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ads []*models.AdvertisementCampaign
	for rows.Next() {
		a := &models.AdvertisementCampaign{}
		if err := rows.Scan(
			&a.ID, &a.CompanyID, &a.ProductID, &a.Title, &a.ImageURL, &a.TargetURL,
			&a.Region, &a.IsActive, &a.StartsAt, &a.EndsAt, &a.Impressions, &a.Clicks, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		ads = append(ads, a)
	}
	return ads, nil
}

func (r *adRepoImpl) IncrementImpressions(ctx context.Context, adID string) error {
	_, err := r.db.Pool.Exec(ctx, `UPDATE advertisement_campaigns SET impressions = impressions + 1 WHERE id = $1`, adID)
	return err
}

func (r *adRepoImpl) IncrementClicks(ctx context.Context, adID string) error {
	_, err := r.db.Pool.Exec(ctx, `UPDATE advertisement_campaigns SET clicks = clicks + 1 WHERE id = $1`, adID)
	return err
}
