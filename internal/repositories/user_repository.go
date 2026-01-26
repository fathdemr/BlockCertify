package repositories

import (
	"BlockCertify/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Exists 1 -> User exists 0 -> User not exists
func (r *UserRepository) Exists(email string) (bool, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}
