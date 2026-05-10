package service

import (
	"context"
	"time"

	"github.com/aarcsx/krishaksh-backend/internal/modules/analytics/repository"
	"github.com/google/uuid"
)

type AnalyticsService interface {
	Track(ctx context.Context, e repository.Event) error
	Dashboard(ctx context.Context, days int) (map[string]int64, error)
	CompanyDashboard(ctx context.Context, companyID uuid.UUID, days int) (map[string]int64, error)
}

type analyticsServiceImpl struct {
	repo repository.AnalyticsRepository
}

func NewAnalyticsService(r repository.AnalyticsRepository) AnalyticsService {
	return &analyticsServiceImpl{repo: r}
}

func (s *analyticsServiceImpl) Track(ctx context.Context, e repository.Event) error {
	if e.OccurredAt.IsZero() {
		e.OccurredAt = time.Now()
	}
	return s.repo.InsertEvent(ctx, e)
}
func (s *analyticsServiceImpl) Dashboard(ctx context.Context, days int) (map[string]int64, error) {
	if days <= 0 {
		days = 30
	}
	return s.repo.DashboardMetrics(ctx, time.Now().AddDate(0, 0, -days))
}
func (s *analyticsServiceImpl) CompanyDashboard(ctx context.Context, companyID uuid.UUID, days int) (map[string]int64, error) {
	if days <= 0 {
		days = 30
	}
	return s.repo.CompanyMetrics(ctx, companyID, time.Now().AddDate(0, 0, -days))
}
