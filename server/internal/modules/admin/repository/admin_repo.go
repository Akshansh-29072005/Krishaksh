package repository

import (
	"context"

	"github.com/aarcsx/krishaksh-backend/internal/database"
	"github.com/aarcsx/krishaksh-backend/internal/models"
	"github.com/google/uuid"
)

type AdminRepository interface {
	ListUsers(ctx context.Context, limit int) ([]*models.User, error)
	ChangeUserRole(ctx context.Context, userID uuid.UUID, roleID int) error
	SuspendUser(ctx context.Context, userID uuid.UUID, suspended bool) error
	CreateDisease(ctx context.Context, d *models.Disease) error
	UpdateDisease(ctx context.Context, d *models.Disease) error
	DeleteDisease(ctx context.Context, id uuid.UUID) error
	CreateProduct(ctx context.Context, p *models.Product) error
	UpdateProduct(ctx context.Context, p *models.Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	ListSupportTickets(ctx context.Context, limit int) ([]*models.SupportTicket, error)
	UpdateSupportTicket(ctx context.Context, id uuid.UUID, status string, assignedTo *uuid.UUID) error
	ListCampaigns(ctx context.Context, limit int) ([]*models.AdvertisementCampaign, error)
	UpdateCampaignStatus(ctx context.Context, id uuid.UUID, active bool) error
}

type adminRepoImpl struct{ db *database.DB }

func NewAdminRepository(db *database.DB) AdminRepository { return &adminRepoImpl{db: db} }

func (r *adminRepoImpl) ListUsers(ctx context.Context, limit int) ([]*models.User, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.db.Pool.Query(ctx, `SELECT id, google_id, full_name, email, phone_number, device_token, role_id, village, language, created_at, updated_at FROM users ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.User
	for rows.Next() {
		u := &models.User{}
		if err := rows.Scan(&u.ID, &u.GoogleID, &u.FullName, &u.Email, &u.PhoneNumber, &u.DeviceToken, &u.RoleID, &u.Village, &u.Language, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, nil
}
func (r *adminRepoImpl) ChangeUserRole(ctx context.Context, userID uuid.UUID, roleID int) error {
	_, err := r.db.Pool.Exec(ctx, `UPDATE users SET role_id=$1, updated_at=NOW() WHERE id=$2`, roleID, userID)
	return err
}
func (r *adminRepoImpl) SuspendUser(ctx context.Context, userID uuid.UUID, suspended bool) error {
	_, err := r.db.Pool.Exec(ctx, `UPDATE users SET is_suspended=$1, updated_at=NOW() WHERE id=$2`, suspended, userID)
	return err
}
func (r *adminRepoImpl) CreateDisease(ctx context.Context, d *models.Disease) error {
	_, err := r.db.Pool.Exec(ctx, `INSERT INTO diseases (id,name,crop_type,severity,symptoms,prevention,description,image_url,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW())`, d.ID, d.Name, d.CropType, d.Severity, d.Symptoms, d.Prevention, d.Description, d.ImageURL)
	return err
}
func (r *adminRepoImpl) UpdateDisease(ctx context.Context, d *models.Disease) error {
	_, err := r.db.Pool.Exec(ctx, `UPDATE diseases SET name=$1,crop_type=$2,severity=$3,symptoms=$4,prevention=$5,description=$6,image_url=$7 WHERE id=$8`, d.Name, d.CropType, d.Severity, d.Symptoms, d.Prevention, d.Description, d.ImageURL, d.ID)
	return err
}
func (r *adminRepoImpl) DeleteDisease(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM diseases WHERE id=$1`, id)
	return err
}
func (r *adminRepoImpl) CreateProduct(ctx context.Context, p *models.Product) error {
	_, err := r.db.Pool.Exec(ctx, `INSERT INTO products (id,company_id,category_id,name,description,image_url,price,currency,unit,is_sponsored,sponsored_priority,is_active,stock_available,stock_quantity,crop_type,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,NOW(),NOW())`, p.ID, p.CompanyID, p.CategoryID, p.Name, p.Description, p.ImageURL, p.Price, p.Currency, p.Unit, p.IsSponsored, p.SponsoredPriority, p.IsActive, p.StockAvailable, p.StockQuantity, p.CropType)
	return err
}
func (r *adminRepoImpl) UpdateProduct(ctx context.Context, p *models.Product) error {
	_, err := r.db.Pool.Exec(ctx, `UPDATE products SET name=$1,description=$2,image_url=$3,price=$4,currency=$5,unit=$6,is_sponsored=$7,sponsored_priority=$8,is_active=$9,stock_available=$10,stock_quantity=$11,crop_type=$12,updated_at=NOW() WHERE id=$13`, p.Name, p.Description, p.ImageURL, p.Price, p.Currency, p.Unit, p.IsSponsored, p.SponsoredPriority, p.IsActive, p.StockAvailable, p.StockQuantity, p.CropType, p.ID)
	return err
}
func (r *adminRepoImpl) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM products WHERE id=$1`, id)
	return err
}
func (r *adminRepoImpl) ListSupportTickets(ctx context.Context, limit int) ([]*models.SupportTicket, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.db.Pool.Query(ctx, `SELECT id,user_id,title,description,status,priority,assigned_to,created_at,updated_at,resolved_at FROM support_tickets ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.SupportTicket
	for rows.Next() {
		t := &models.SupportTicket{}
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.AssignedTo, &t.CreatedAt, &t.UpdatedAt, &t.ResolvedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}
func (r *adminRepoImpl) UpdateSupportTicket(ctx context.Context, id uuid.UUID, status string, assignedTo *uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, `UPDATE support_tickets SET status=$1, assigned_to=$2, updated_at=NOW() WHERE id=$3`, status, assignedTo, id)
	return err
}
func (r *adminRepoImpl) ListCampaigns(ctx context.Context, limit int) ([]*models.AdvertisementCampaign, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.db.Pool.Query(ctx, `SELECT id, company_id, product_id, title, image_url, target_url, region, is_active, starts_at, ends_at, impressions, clicks, created_at FROM advertisement_campaigns ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.AdvertisementCampaign
	for rows.Next() {
		a := &models.AdvertisementCampaign{}
		if err := rows.Scan(&a.ID, &a.CompanyID, &a.ProductID, &a.Title, &a.ImageURL, &a.TargetURL, &a.Region, &a.IsActive, &a.StartsAt, &a.EndsAt, &a.Impressions, &a.Clicks, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}
func (r *adminRepoImpl) UpdateCampaignStatus(ctx context.Context, id uuid.UUID, active bool) error {
	_, err := r.db.Pool.Exec(ctx, `UPDATE advertisement_campaigns SET is_active=$1 WHERE id=$2`, active, id)
	return err
}
