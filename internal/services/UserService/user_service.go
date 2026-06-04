package UserService

import (
	"BlockCertify/internal/models"
	"BlockCertify/internal/services/CacheService"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type UserService struct {
	DB           *gorm.DB
	CacheService *CacheService.CacheService
	params       *viper.Viper
}

func New(db *gorm.DB, params *viper.Viper) *UserService {
	return &UserService{
		DB:     db,
		params: params,
	}
}

func (s *UserService) UseCacheService(cache *CacheService.CacheService) {
	s.CacheService = cache
}

func (s *UserService) Create(user *models.User) error {
	return s.DB.Create(user).Error
}

func (s *UserService) CreateAdmin(admin *models.Admin) error {
	return s.DB.Create(admin).Error
}

func (s *UserService) CreateTransaction() *gorm.DB {
	return s.DB.Begin()
}

func (s *UserService) ExpiresInSeconds() int64 {
	seconds := time.Duration(s.params.GetUint64("crypto.contact_token_expire_duration_minute")) * time.Second
	return int64(seconds.Seconds())
}
