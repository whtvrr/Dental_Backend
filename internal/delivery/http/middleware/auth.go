package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/whtvrr/Dental_Backend/internal/auth"
	"github.com/whtvrr/Dental_Backend/internal/domain/entities"
)

type AuthMiddleware struct {
	jwtManager *auth.JWTManager
}

func NewAuthMiddleware(jwtManager *auth.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

// RequireAuth validates JWT token and sets user claims in context
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := tokenParts[1]
		claims, err := m.jwtManager.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Ensure it's an access token
		if claims.Type != "access" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
			c.Abort()
			return
		}

		// Set user claims in context
		c.Set("user", claims)
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("user_email", claims.Email)
		c.Set("user_full_name", claims.FullName)

		c.Next()
	}
}

// RequireRoles checks if user has one of the required roles
func (m *AuthMiddleware) RequireRoles(roles ...entities.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		role := userRole.(entities.UserRole)
		
		// Admin can access everything
		if role == entities.RoleAdmin {
			c.Next()
			return
		}

		// Check if user has required role
		for _, requiredRole := range roles {
			if role == requiredRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}

// RequireAdmin checks if user is admin
func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return m.RequireRoles(entities.RoleAdmin)
}

// RequireDoctor checks if user is doctor or admin
func (m *AuthMiddleware) RequireDoctor() gin.HandlerFunc {
	return m.RequireRoles(entities.RoleDoctor, entities.RoleAdmin)
}

// RequireReceptionist checks if user is receptionist or admin
func (m *AuthMiddleware) RequireReceptionist() gin.HandlerFunc {
	return m.RequireRoles(entities.RoleReceptionist)
}

// RequireDoctorOrReceptionist checks if user is doctor, receptionist, or admin
func (m *AuthMiddleware) RequireDoctorOrReceptionist() gin.HandlerFunc {
	return m.RequireRoles(entities.RoleDoctor, entities.RoleReceptionist)
}

// RequireStaff checks if user is staff (doctor, receptionist, or admin)
func (m *AuthMiddleware) RequireStaff() gin.HandlerFunc {
	return m.RequireRoles(entities.RoleDoctor, entities.RoleReceptionist, entities.RoleAdmin)
}