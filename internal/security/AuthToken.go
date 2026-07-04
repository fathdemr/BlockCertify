package security

import (
	"BlockCertify/internal/models"
	apperrors "BlockCertify/internal/pkg/errors"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// TokenHelper signs and verifies JWT tokens.
type TokenHelper interface {
	Sign(claims jwt.MapClaims) (string, error)
	Verify(tokenString string) (jwt.MapClaims, error)
}

// JwtHelper implements TokenHelper using RS256 (RSA) keys.
// The private key is required only for signing; verification uses the public key.
type JwtHelper struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	expire     time.Duration
}

// NewRSAJWTHelper builds a JwtHelper from PEM-encoded RSA keys.
// privatePEM may be empty for verify-only usage.
func NewRSAJWTHelper(privatePEM, publicPEM string, expire time.Duration) (*JwtHelper, error) {

	if publicPEM == "" {
		return nil, fmt.Errorf("rsa public key is empty")
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicPEM))
	if err != nil {
		return nil, fmt.Errorf("failed to parse rsa public key: %w", err)
	}

	helper := &JwtHelper{
		publicKey: publicKey,
		expire:    expire,
	}

	if privatePEM != "" {
		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privatePEM))
		if err != nil {
			return nil, fmt.Errorf("failed to parse rsa private key: %w", err)
		}
		helper.privateKey = privateKey
	}

	return helper, nil
}

// Sign creates an RS256-signed token from the given claims.
// If "exp" is not set, it is derived from the helper's expire duration.
func (j *JwtHelper) Sign(claims jwt.MapClaims) (string, error) {

	if j.privateKey == nil {
		return "", apperrors.New(apperrors.ErrTokenCreateFailed, "RSA private key is not configured", nil)
	}

	if _, ok := claims["exp"]; !ok {
		claims["exp"] = time.Now().Add(j.expire).Unix()
	}
	if _, ok := claims["iat"]; !ok {
		claims["iat"] = time.Now().Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(j.privateKey)
}

func (j *JwtHelper) Verify(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.publicKey, nil
	})

	if err != nil {
		if verr, ok := err.(*jwt.ValidationError); ok && verr.Errors&jwt.ValidationErrorExpired != 0 {
			return nil, apperrors.New(apperrors.ErrTokenExpired, "Token is Expired", err)
		}
		return nil, apperrors.New(apperrors.ErrInvalidToken, "Invalid token", err)
	}

	if !token.Valid {
		return nil, apperrors.New(apperrors.ErrInvalidToken, "Invalid token", nil)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, apperrors.New(apperrors.ErrInvalidToken, "Invalid token claims", nil)
	}

	return claims, nil
}

func (j *JwtHelper) GetActor(c *gin.Context) *models.User {
	var actor models.User
	if m, isExist := c.Get("actor"); isExist {
		actor = m.(models.User)
	}
	return &actor
}
