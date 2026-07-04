package UserService

import (
	"BlockCertify/internal/models"
	"BlockCertify/internal/security"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var errProfileIDRequired = errors.New("email is required")

// CreateToken issues an RS256-signed JWT for the given user.
// The RSA key pair is read from config (crypto.rsa_private_key / crypto.rsa_public_key).
func (s *UserService) CreateToken(user *models.User) (string, error) {
	if user.Email == "" {
		return "", errProfileIDRequired
	}

	jwtHelper, err := security.NewJWTHelperFromParams(s.params)
	if err != nil {
		return "", err
	}

	AppProfileTokenHour := s.params.GetUint64("crypto.token_expire_duration_hour")
	AppProfileTokenDuration := time.Duration(AppProfileTokenHour) * time.Hour

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  string(user.Role),
		"Type":  "Profile",
		"Name":  fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		"exp":   time.Now().Add(AppProfileTokenDuration).Unix(),
		"expd":  time.Now().Add(AppProfileTokenDuration).Format("2006-01-02T15:04:05.000Z"),
		"iat":   time.Now().Unix(),
	}

	return jwtHelper.Sign(claims)
}
