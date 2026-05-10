package handler

import (
	"net/http"
	"os"

	"github.com/aarcsx/krisho-backend/internal/core/response"
	"github.com/aarcsx/krisho-backend/internal/modules/auth/dto"
	"github.com/aarcsx/krisho-backend/internal/modules/auth/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req dto.GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Missing or invalid id_token")
		return
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID == "" {
		// Mock client id for missing env config
		clientID = "123456789-dummy.apps.googleusercontent.com"
	}

	res, err := h.service.LoginWithGoogle(c.Request.Context(), req, clientID)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login successful", res)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.TokenRefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	res, err := h.service.RefreshSession(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Session expired, please login again")
		return
	}

	response.Success(c, http.StatusOK, "Token refreshed successfully", res)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.TokenRefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "missing refresh_token")
		return
	}

	if err := h.service.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to logout")
		return
	}

	response.Success(c, http.StatusOK, "Logged out successfully", nil)
}
