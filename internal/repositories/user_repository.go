package repositories

import (
	"BlockCertify/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
	Exists(email string) (bool, error)
	Create(user *models.User) error
	CreateAdmin(admin *models.Admin) error
	CreateTransaction() *gorm.DB
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Exists 1 -> User exists 0 -> User not exists
func (r *userRepository) Exists(email string) (bool, error) {
	var count int64

	err := r.db.Model(&models.User{}).
		Where("email = ?", email).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) CreateAdmin(admin *models.Admin) error {
	return r.db.Create(admin).Error
}

func (r *userRepository) CreateTransaction() *gorm.DB {
	return r.db.Begin()
}
