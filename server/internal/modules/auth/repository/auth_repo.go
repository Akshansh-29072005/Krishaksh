package repository

import (
	"context"

	"github.com/aarcsx/krisho-backend/internal/database"
	"github.com/aarcsx/krisho-backend/internal/models"
)

type AuthRepository interface {
	FindUserByGoogleID(ctx context.Context, googleID string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	GetRoleByName(ctx context.Context, name string) (*models.Role, error)
	StoreRefreshToken(ctx context.Context, token *models.RefreshToken) error
	RevokeRefreshToken(ctx context.Context, tokenStr string) error
	FindRefreshToken(ctx context.Context, tokenStr string) (*models.RefreshToken, error)
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

func (r *authRepoImpl) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, google_id, full_name, email, language, role_id, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Pool.Exec(ctx, query,
		user.ID, user.GoogleID, user.FullName, user.Email, user.Language,
		user.RoleID, user.CreatedAt, user.UpdatedAt,
	)
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
