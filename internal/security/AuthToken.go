package security

import (
	"BlockCertify/internal/models"
	apperrors "BlockCertify/internal/pkg/errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type JwtHelper struct {
	secret []byte
	expire time.Duration
}

func NewJWTHelper(secret string, expireHours time.Duration) *JwtHelper {
	return &JwtHelper{
		secret: []byte(secret),
		expire: time.Hour * expireHours,
	}
}

func (j *JwtHelper) Verify(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, apperrors.New(apperrors.ErrInvalidToken, "Invalid token", nil)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, apperrors.New(apperrors.ErrInvalidToken, "Invalid token claims", nil)
	}

	if claims["exp"].(float64) < float64(time.Now().Unix()) {
		return nil, apperrors.New(apperrors.ErrTokenExpired, "Token is Expired", nil)
	}

	return claims, nil
}

func (j *JwtHelper) GetJwtSecretKey() []byte {
	return j.secret
}

func (j *JwtHelper) GetActor(c *gin.Context) *models.User {
	var actor models.User
	if m, isExist := c.Get("actor"); isExist {
		actor = m.(models.User)
	}
	return &actor
}
