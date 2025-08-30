package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/whtvrr/Dental_Backend/internal/delivery/http/response"
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
// @Success 201 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Router /auth/signup [post]
func (h *AuthHandler) SignUp(c *gin.Context) {
	var req usecases.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	tokenPair, err := h.authUseCase.SignUp(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Created("User created successfully", tokenPair))
}

// SignIn godoc
// @Summary Sign in user
// @Description Authenticate user and return tokens with user role
// @Tags auth
// @Accept json
// @Produce json
// @Param request body usecases.SignInRequest true "Sign in request"
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 401 {object} response.StandardResponse "Unauthorized"
// @Router /auth/signin [post]
func (h *AuthHandler) SignIn(c *gin.Context) {
	var req usecases.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	result, err := h.authUseCase.SignInWithUserInfo(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Unauthorized(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Sign in successful", result))
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh the access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body usecases.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 401 {object} response.StandardResponse "Unauthorized"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req usecases.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	tokenPair, err := h.authUseCase.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Unauthorized(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Token refreshed successfully", tokenPair))
}

// Me godoc
// @Summary Get current user information
// @Description Get the current authenticated user's information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.StandardResponse
// @Failure 401 {object} response.StandardResponse "Unauthorized"
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	// Get user claims from context (set by auth middleware)
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not authenticated"))
		return
	}

	c.JSON(http.StatusOK, response.OK("User information retrieved successfully", userClaims))
}
