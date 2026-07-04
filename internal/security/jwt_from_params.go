package security

import (
	"time"

	"github.com/spf13/viper"
)

// NewJWTHelperFromParams builds an RSA JwtHelper from config values:
// crypto.rsa_private_key, crypto.rsa_public_key, crypto.token_expire_duration_hour.
func NewJWTHelperFromParams(params *viper.Viper) (*JwtHelper, error) {
	expire := time.Duration(params.GetUint64("crypto.token_expire_duration_hour")) * time.Hour
	return NewRSAJWTHelper(
		params.GetString("crypto.rsa_private_key"),
		params.GetString("crypto.rsa_public_key"),
		expire,
	)
}
