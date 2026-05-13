package repository

import (
	"context"

	"github.com/aarcsx/krisho-backend/internal/database"
	"github.com/aarcsx/krisho-backend/internal/models"
	"github.com/google/uuid"
)

type AuthRepository interface {
	FindUserByGoogleID(ctx context.Context, googleID string) (*models.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	GetRoleByName(ctx context.Context, name string) (*models.Role, error)
	StoreRefreshToken(ctx context.Context, token *models.RefreshToken) error
	RevokeRefreshToken(ctx context.Context, tokenStr string) error
	FindRefreshToken(ctx context.Context, tokenStr string) (*models.RefreshToken, error)
	UpdateDeviceToken(ctx context.Context, userID uuid.UUID, deviceToken string) error
}

type authRepoImpl struct {
	db *database.DB
}

func NewAuthRepository(db *database.DB) AuthRepository {
	return &authRepoImpl{db: db}
}

func (r *authRepoImpl) FindUserByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	query := `SELECT id, google_id, full_name, email, phone_number, device_token, role_id, village, language, created_at, updated_at 
	          FROM users WHERE google_id = $1 LIMIT 1`

	user := &models.User{}
	err := r.db.Pool.QueryRow(ctx, query, googleID).Scan(
		&user.ID, &user.GoogleID, &user.FullName, &user.Email,
		&user.PhoneNumber, &user.DeviceToken, &user.RoleID,
		&user.Village, &user.Language, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err // ErrNoRows handled by service layer usually
	}
	return user, nil
}

func (r *authRepoImpl) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	query := `SELECT id, google_id, full_name, email, phone_number, device_token, role_id, village, language, created_at, updated_at 
	          FROM users WHERE id = $1 LIMIT 1`

	user := &models.User{}
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&user.ID, &user.GoogleID, &user.FullName, &user.Email,
		&user.PhoneNumber, &user.DeviceToken, &user.RoleID,
		&user.Village, &user.Language, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *authRepoImpl) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, google_id, full_name, email, language, role_id, device_token, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Pool.Exec(ctx, query,
		user.ID, user.GoogleID, user.FullName, user.Email, user.Language,
		user.RoleID, user.DeviceToken, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *authRepoImpl) UpdateDeviceToken(ctx context.Context, userID uuid.UUID, deviceToken string) error {
	query := `UPDATE users SET device_token = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Pool.Exec(ctx, query, deviceToken, userID)
	return err
}

func (r *authRepoImpl) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	query := `SELECT id, name, permissions FROM roles WHERE name = $1 LIMIT 1`
	role := &models.Role{}
	err := r.db.Pool.QueryRow(ctx, query, name).Scan(&role.ID, &role.Name, &role.Permissions)
	return role, err
}

func (r *authRepoImpl) StoreRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (id, user_id, token, expires_at, revoked) 
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Pool.Exec(ctx, query, token.ID, token.UserID, token.Token, token.ExpiresAt, token.Revoked)
	return err
}

func (r *authRepoImpl) RevokeRefreshToken(ctx context.Context, tokenStr string) error {
	query := `UPDATE refresh_tokens SET revoked = TRUE WHERE token = $1`
	_, err := r.db.Pool.Exec(ctx, query, tokenStr)
	return err
}

func (r *authRepoImpl) FindRefreshToken(ctx context.Context, tokenStr string) (*models.RefreshToken, error) {
	query := `SELECT id, user_id, token, expires_at, revoked FROM refresh_tokens WHERE token = $1 LIMIT 1`
	token := &models.RefreshToken{}
	err := r.db.Pool.QueryRow(ctx, query, tokenStr).Scan(&token.ID, &token.UserID, &token.Token, &token.ExpiresAt, &token.Revoked)
	return token, err
}
