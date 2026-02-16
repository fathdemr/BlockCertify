package repositories

import (
	"BlockCertify/internal/dto"
	"BlockCertify/internal/models"
	"log/slog"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
	Exists(email string) (bool, error)
	Create(user *models.User) error
	GetUniversitiesFromDBRecord() ([]dto.UniversitiesResponse, error)
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

func (r *userRepository) GetUniversitiesFromDBRecord() ([]dto.UniversitiesResponse, error) {

	var response []dto.UniversitiesResponse

	err := r.db.
		Model(&models.Universities{}).
		Select("id, name").
		Scan(&response).Error

	if err != nil {
		slog.Error("failed to fetch universities", "err", err)
		return nil, err
	}

	return response, nil
}
