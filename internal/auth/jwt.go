package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/whtvrr/Dental_Backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenClaims struct {
	UserID   primitive.ObjectID    `json:"user_id"`
	Email    string               `json:"email"`
	FullName string               `json:"full_name"`
	Role     entities.UserRole    `json:"role"`
	Type     string               `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type JWTManager struct {
	secretKey        string
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
}

func NewJWTManager(secretKey string, accessTokenTTL, refreshTokenTTL time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:       secretKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (j *JWTManager) GenerateTokenPair(user *entities.User) (*TokenPair, error) {
	if user.Email == nil {
		return nil, errors.New("user email is required for authentication")
	}

	accessToken, err := j.generateToken(user, "access", j.accessTokenTTL)
	if err != nil {
		return nil, err
	}

	refreshToken, err := j.generateToken(user, "refresh", j.refreshTokenTTL)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(j.accessTokenTTL.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

func (j *JWTManager) generateToken(user *entities.User, tokenType string, ttl time.Duration) (string, error) {
	claims := TokenClaims{
		UserID:   user.ID,
		Email:    *user.Email,
		FullName: user.FullName,
		Role:     user.Role,
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTManager) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (j *JWTManager) RefreshToken(refreshTokenString string) (*TokenPair, error) {
	claims, err := j.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	if claims.Type != "refresh" {
		return nil, errors.New("invalid token type")
	}

	// Create a user object from claims for token generation
	user := &entities.User{
		ID:       claims.UserID,
		Email:    &claims.Email,
		FullName: claims.FullName,
		Role:     claims.Role,
	}

	return j.GenerateTokenPair(user)
}