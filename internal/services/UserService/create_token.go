package UserService

import (
	"BlockCertify/internal/models"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var errProfileIDRequired = errors.New("email is required")

func (s *UserService) CreateToken(user *models.User) (string, error) {
	if user.Email == "" {
		return "", errProfileIDRequired
	}
	// Create the token 4 application
	token := jwt.New(jwt.SigningMethodHS256)
	Claims := make(jwt.MapClaims)
	Claims["id"] = user.ID
	Claims["email"] = user.Email
	Claims["Type"] = "Profile"
	Claims["Name"] = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	AppProfileTokenHour := s.params.GetUint64("crypto.token_expire_duration_hour")
	AppProfileTokenDuration := time.Duration(AppProfileTokenHour) * time.Hour
	// 86400 hours = 2 months
	Claims["exp"] = time.Now().Add(AppProfileTokenDuration).Unix()
	Claims["expd"] = time.Now().Add(AppProfileTokenDuration).Format("2006-01-02T15:04:05.000Z")
	Claims["iat"] = time.Now().Unix()
	token.Claims = Claims
	tokenString, err := token.SignedString([]byte(s.params.GetString("crypto.my_secret_key")))
	return tokenString, err
}
