package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/whtvrr/Dental_Backend/internal/auth"
	"github.com/whtvrr/Dental_Backend/internal/domain/entities"
	"github.com/whtvrr/Dental_Backend/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthUseCase struct {
	userRepo   repositories.UserRepository
	jwtManager *auth.JWTManager
}

type SignUpRequest struct {
	Email     string            `json:"email" binding:"required,email"`
	Password  string            `json:"password" binding:"required,min=6"`
	FullName  string            `json:"full_name" binding:"required"`
	Role      entities.UserRole `json:"role" binding:"required"`
	PhoneNumber *string         `json:"phone_number,omitempty"`
}

type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type SignInResponse struct {
	Tokens *auth.TokenPair   `json:"tokens"`
	User   *UserInfo         `json:"user"`
}

type UserInfo struct {
	ID       string            `json:"id"`
	Email    string            `json:"email"`
	FullName string            `json:"full_name"`
	Role     entities.UserRole `json:"role"`
}

func NewAuthUseCase(userRepo repositories.UserRepository, jwtManager *auth.JWTManager) *AuthUseCase {
	return &AuthUseCase{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (uc *AuthUseCase) SignUp(ctx context.Context, req *SignUpRequest) (*auth.TokenPair, error) {
	// Validate role - only authenticatable roles can sign up
	if req.Role == entities.RoleClient {
		return nil, errors.New("clients cannot sign up through this endpoint")
	}

	if req.Role != entities.RoleAdmin && req.Role != entities.RoleDoctor && req.Role != entities.RoleReceptionist {
		return nil, errors.New("invalid role")
	}

	// Check if user already exists
	existingUser, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &entities.User{
		Email:        &req.Email,
		PasswordHash: &hashedPassword,
		Role:         req.Role,
		FullName:     req.FullName,
		PhoneNumber:  req.PhoneNumber,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate token pair
	return uc.jwtManager.GenerateTokenPair(user)
}

func (uc *AuthUseCase) SignIn(ctx context.Context, req *SignInRequest) (*auth.TokenPair, error) {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Check if user can authenticate
	if !user.CanAuthenticate() {
		return nil, errors.New("user cannot authenticate")
	}

	// Check if user has password
	if user.PasswordHash == nil {
		return nil, errors.New("user has no password set")
	}

	// Verify password
	if !auth.CheckPassword(*user.PasswordHash, req.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate token pair
	return uc.jwtManager.GenerateTokenPair(user)
}

func (uc *AuthUseCase) SignInWithUserInfo(ctx context.Context, req *SignInRequest) (*SignInResponse, error) {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Check if user can authenticate
	if !user.CanAuthenticate() {
		return nil, errors.New("user cannot authenticate")
	}

	// Check if user has password
	if user.PasswordHash == nil {
		return nil, errors.New("user has no password set")
	}

	// Verify password
	if !auth.CheckPassword(*user.PasswordHash, req.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate token pair
	tokenPair, err := uc.jwtManager.GenerateTokenPair(user)
	if err != nil {
		return nil, err
	}

	// Create user info response
	userInfo := &UserInfo{
		ID:       user.ID.Hex(),
		Email:    *user.Email,
		FullName: user.FullName,
		Role:     user.Role,
	}

	return &SignInResponse{
		Tokens: tokenPair,
		User:   userInfo,
	}, nil
}

func (uc *AuthUseCase) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*auth.TokenPair, error) {
	return uc.jwtManager.RefreshToken(req.RefreshToken)
}

func (uc *AuthUseCase) ValidateToken(tokenString string) (*auth.TokenClaims, error) {
	return uc.jwtManager.ValidateToken(tokenString)
}