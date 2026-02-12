package middleware

import (
	"BlockCertify/internal/repositories"
	"BlockCertify/internal/security"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	Authorize() gin.HandlerFunc
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

func (s *authMiddleware) Authorize() gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "No Authorization header found",
			})
			return
		}

		tokenStr := strings.Replace(authHeader, "Bearer ", "", 10)
		tokenStr = strings.Replace(tokenStr, "bearer ", "", 10)

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
		}

		user, err := s.userRepo.FindByEmail(email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
		}

		c.Set("user", user)
		c.Next()

	}
}
