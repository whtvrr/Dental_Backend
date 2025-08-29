package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/whtvrr/Dental_Backend/internal/usecases"
)

type AuthHandler struct {
	authUseCase *usecases.AuthUseCase
}

func NewAuthHandler(authUseCase *usecases.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

// SignUp godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body usecases.SignUpRequest true "Sign up request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /auth/signup [post]
func (h *AuthHandler) SignUp(c *gin.Context) {
	var req usecases.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenPair, err := h.authUseCase.SignUp(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"tokens":  tokenPair,
	})
}

// SignIn godoc
// @Summary Sign in user
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body usecases.SignInRequest true "Sign in request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /auth/signin [post]
func (h *AuthHandler) SignIn(c *gin.Context) {
	var req usecases.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenPair, err := h.authUseCase.SignIn(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sign in successful",
		"tokens":  tokenPair,
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh the access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body usecases.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req usecases.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenPair, err := h.authUseCase.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"tokens":  tokenPair,
	})
}

// Me godoc
// @Summary Get current user information
// @Description Get the current authenticated user's information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	// Get user claims from context (set by auth middleware)
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": userClaims,
	})
}
