package service

import (
	"context"
	"time"

	analyticsRepo "github.com/aarcsx/krisho-backend/internal/modules/analytics/repository"
	"github.com/google/uuid"
)

type CompanyService interface {
	Dashboard(ctx context.Context, companyID uuid.UUID, days int) (map[string]int64, error)
}

type companyServiceImpl struct {
	analytics analyticsRepo.AnalyticsRepository
}

func NewCompanyService(a analyticsRepo.AnalyticsRepository) CompanyService {
	return &companyServiceImpl{analytics: a}
}

func (s *companyServiceImpl) Dashboard(ctx context.Context, companyID uuid.UUID, days int) (map[string]int64, error) {
	if days <= 0 {
		days = 30
	}
	return s.analytics.CompanyMetrics(ctx, companyID, time.Now().AddDate(0, 0, -days))
}
