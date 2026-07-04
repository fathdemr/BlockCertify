package middleware

import (
	"BlockCertify/internal/models"
	"BlockCertify/internal/repositories"
	"BlockCertify/internal/security"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	Authorize() gin.HandlerFunc
	RequireRole(role models.UserRole) gin.HandlerFunc
}

type authMiddleware struct {
	jwtHelper security.TokenHelper
	userRepo  repositories.UserRepository
}

func NewAuthMiddleware(jwtHelper security.TokenHelper, userRepo repositories.UserRepository) AuthMiddleware {
	return &authMiddleware{
		jwtHelper: jwtHelper,
		userRepo:  userRepo,
	}
}

// extractToken reads the JWT from the Authorization header, falling back to the "jwt" cookie.
func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		tokenStr := strings.TrimSpace(authHeader)
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		tokenStr = strings.TrimPrefix(tokenStr, "bearer ")
		return strings.TrimSpace(tokenStr)
	}
	if cookie, err := c.Cookie("jwt"); err == nil {
		return strings.TrimSpace(cookie)
	}
	return ""
}

func (s *authMiddleware) Authorize() gin.HandlerFunc {

	return func(c *gin.Context) {

		tokenStr := extractToken(c)
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "No Authorization header found",
			})
			return
		}

		claims, err := s.jwtHelper.Verify(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		email, ok := claims["email"].(string)
		if !ok || email == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}

		user, err := s.userRepo.FindByEmail(email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			return
		}

		c.Set("user", user)
		c.Set("email", email)
		c.Next()
	}
}

// RequireRole must run after Authorize; it rejects users that don't have the given role.
func (s *authMiddleware) RequireRole(role models.UserRole) gin.HandlerFunc {

	return func(c *gin.Context) {

		value, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		user, ok := value.(*models.User)
		if !ok || user.Role != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			return
		}

		c.Next()
	}
}
