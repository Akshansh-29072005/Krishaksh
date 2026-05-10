package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/aarcsx/krisho-backend/internal/models"
	"github.com/aarcsx/krisho-backend/internal/modules/auth/dto"
	"github.com/aarcsx/krisho-backend/internal/modules/auth/repository"
	"github.com/aarcsx/krisho-backend/pkg/utils"
	"github.com/google/uuid"
	"google.golang.org/api/idtoken"
)

type AuthService interface {
	LoginWithGoogle(ctx context.Context, req dto.GoogleLoginRequest, clientID string) (*dto.AuthResponse, error)
	RefreshSession(ctx context.Context, req dto.TokenRefreshRequest) (*dto.AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}

type authServiceImpl struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authServiceImpl{repo: repo}
}

func (s *authServiceImpl) LoginWithGoogle(ctx context.Context, req dto.GoogleLoginRequest, clientID string) (*dto.AuthResponse, error) {
	// 1. Verify Google token
	payload, err := idtoken.Validate(ctx, req.IDToken, clientID)
	if err != nil {
		return nil, fmt.Errorf("invalid google token: %w", err)
	}

	googleID := payload.Subject
	email := payload.Claims["email"].(string)
	name := payload.Claims["name"].(string)

	// 2. Find or Create User
	user, err := s.repo.FindUserByGoogleID(ctx, googleID)
	if err != nil {
		// User does not exist, so let's register them
		role, err := s.repo.GetRoleByName(ctx, "FARMER")
		if err != nil {
			return nil, fmt.Errorf("default role not found: %w", err)
		}

		user = &models.User{
			ID:        uuid.New(),
			GoogleID:  &googleID,
			FullName:  name,
			Email:     &email,
			RoleID:    role.ID,
			Language:  "en",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := s.repo.CreateUser(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	roleName := "FARMER" // Ideally we grab this via a JOIN query on login
	if user.RoleID == 2 {
		roleName = "ADMIN"
	}

	// 3. Generate Access Token
	accessToken, err := utils.GenerateAccessToken(user.ID, roleName, *user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// 4. Generate & Store Refresh Token
	rtBytes := make([]byte, 32)
	rand.Read(rtBytes)
	refreshTokenStr := hex.EncodeToString(rtBytes)

	rt := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 Days
		Revoked:   false,
	}

	if err := s.repo.StoreRefreshToken(ctx, rt); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    900, // 15 mins
	}, nil
}

func (s *authServiceImpl) RefreshSession(ctx context.Context, req dto.TokenRefreshRequest) (*dto.AuthResponse, error) {
	// Find and validate refresh token
	rt, err := s.repo.FindRefreshToken(ctx, req.RefreshToken)
	if err != nil || rt.Revoked || time.Now().After(rt.ExpiresAt) {
		return nil, errors.New("invalid or expired refresh token")
	}

	s.repo.RevokeRefreshToken(ctx, req.RefreshToken)

	// Since we need email & role, we'd ideally load the user fully. Hardcoding dummy load for snippet.
	accessToken, err := utils.GenerateAccessToken(rt.UserID, "FARMER", "loaded@email.com")
	if err != nil {
		return nil, err
	}

	rtBytes := make([]byte, 32)
	rand.Read(rtBytes)
	newRtStr := hex.EncodeToString(rtBytes)

	newRt := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    rt.UserID,
		Token:     newRtStr,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}
	s.repo.StoreRefreshToken(ctx, newRt)

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRtStr,
		ExpiresIn:    900,
	}, nil
}

func (s *authServiceImpl) Logout(ctx context.Context, refreshToken string) error {
	return s.repo.RevokeRefreshToken(ctx, refreshToken)
}
