package service

import (
    "context"

    appRepo "github.com/aarcsx/krisho-backend/internal/modules/app/repository"
    "github.com/aarcsx/krisho-backend/internal/models"
)

type AppService struct {
    repo *appRepo.AppRepository
}

func NewAppService(repo *appRepo.AppRepository) *AppService {
    return &AppService{repo: repo}
}

func (s *AppService) GetActiveCrops(ctx context.Context) ([]*models.Crop, error) {
    return s.repo.GetActiveCrops(ctx)
}

func (s *AppService) GetAppConfig(ctx context.Context) (*models.AppConfig, error) {
    return s.repo.GetAppConfig(ctx)
}
