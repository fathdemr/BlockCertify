package security

import (
	"fmt"
	"time"

	apperrors "BlockCertify/pkg/errors"

	"github.com/golang-jwt/jwt/v4"
)

type TokenHelper interface {
	Create(email string) (string, error)
	Verify(tokenString string) (jwt.MapClaims, error)
	ExpiresInSeconds() int64
}

type jwtHelper struct {
	secret []byte
	expire time.Duration
}

func NewJWTHelper(secret string, expireHours time.Duration) TokenHelper {
	return &jwtHelper{
		secret: []byte(secret),
		expire: time.Hour * expireHours,
	}
}

func (j *jwtHelper) Create(email string) (string, error) {

	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(j.expire).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(j.secret)
}

func (j *jwtHelper) Verify(tokenString string) (jwt.MapClaims, error) {

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

func (j *jwtHelper) ExpiresInSeconds() int64 {
	return int64(j.expire.Seconds())
}
