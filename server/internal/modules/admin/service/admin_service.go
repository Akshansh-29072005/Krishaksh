package service

import (
	"context"

	"github.com/aarcsx/krisho-backend/internal/models"
	"github.com/aarcsx/krisho-backend/internal/modules/admin/repository"
	"github.com/google/uuid"
)

type AdminService interface {
	ListUsers(ctx context.Context, limit int) ([]*models.User, error)
	ChangeRole(ctx context.Context, userID uuid.UUID, roleID int) error
	SuspendUser(ctx context.Context, userID uuid.UUID, suspended bool) error
	CreateDisease(ctx context.Context, d *models.Disease) error
	UpdateDisease(ctx context.Context, d *models.Disease) error
	DeleteDisease(ctx context.Context, id uuid.UUID) error
	CreateProduct(ctx context.Context, p *models.Product) error
	UpdateProduct(ctx context.Context, p *models.Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	ListSupport(ctx context.Context, limit int) ([]*models.SupportTicket, error)
	UpdateSupport(ctx context.Context, id uuid.UUID, status string, assignedTo *uuid.UUID) error
	ListCampaigns(ctx context.Context, limit int) ([]*models.AdvertisementCampaign, error)
	UpdateCampaign(ctx context.Context, id uuid.UUID, active bool) error
}

type adminServiceImpl struct{ repo repository.AdminRepository }

func NewAdminService(r repository.AdminRepository) AdminService { return &adminServiceImpl{repo: r} }
func (s *adminServiceImpl) ListUsers(ctx context.Context, limit int) ([]*models.User, error) {
	return s.repo.ListUsers(ctx, limit)
}
func (s *adminServiceImpl) ChangeRole(ctx context.Context, userID uuid.UUID, roleID int) error {
	return s.repo.ChangeUserRole(ctx, userID, roleID)
}
func (s *adminServiceImpl) SuspendUser(ctx context.Context, userID uuid.UUID, suspended bool) error {
	return s.repo.SuspendUser(ctx, userID, suspended)
}
func (s *adminServiceImpl) CreateDisease(ctx context.Context, d *models.Disease) error {
	return s.repo.CreateDisease(ctx, d)
}
func (s *adminServiceImpl) UpdateDisease(ctx context.Context, d *models.Disease) error {
	return s.repo.UpdateDisease(ctx, d)
}
func (s *adminServiceImpl) DeleteDisease(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteDisease(ctx, id)
}
func (s *adminServiceImpl) CreateProduct(ctx context.Context, p *models.Product) error {
	return s.repo.CreateProduct(ctx, p)
}
func (s *adminServiceImpl) UpdateProduct(ctx context.Context, p *models.Product) error {
	return s.repo.UpdateProduct(ctx, p)
}
func (s *adminServiceImpl) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteProduct(ctx, id)
}
func (s *adminServiceImpl) ListSupport(ctx context.Context, limit int) ([]*models.SupportTicket, error) {
	return s.repo.ListSupportTickets(ctx, limit)
}
func (s *adminServiceImpl) UpdateSupport(ctx context.Context, id uuid.UUID, status string, assignedTo *uuid.UUID) error {
	return s.repo.UpdateSupportTicket(ctx, id, status, assignedTo)
}
func (s *adminServiceImpl) ListCampaigns(ctx context.Context, limit int) ([]*models.AdvertisementCampaign, error) {
	return s.repo.ListCampaigns(ctx, limit)
}
func (s *adminServiceImpl) UpdateCampaign(ctx context.Context, id uuid.UUID, active bool) error {
	return s.repo.UpdateCampaignStatus(ctx, id, active)
}
