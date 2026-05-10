package dto

// GoogleLoginRequest receives the OAuth ID Token from Android
type GoogleLoginRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

// TokenRefreshRequest rotates the expired Access Token
type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// AuthResponse maps the standardized JWT reply
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // seconds
}
