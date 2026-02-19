package repositories

import (
	"BlockCertify/internal/dto"
	"BlockCertify/internal/models"
	"log/slog"

	"gorm.io/gorm"
)

type UniversityRepository interface {
	GetUniversityByID(id string) (models.Universities, error)
	GetUniversitiesFromDBRecord() ([]dto.UniversitiesResponse, error)
}

type universityRepository struct {
	db *gorm.DB
}

func NewUniversityRepository(db *gorm.DB) UniversityRepository {
	return &universityRepository{
		db: db,
	}
}

func (r *universityRepository) GetUniversityByID(id string) (models.Universities, error) {

	var university models.Universities

	err := r.db.Where("id = ?", id).First(&university).Error
	if err != nil {
		return models.Universities{}, err
	}
	return university, nil
}

func (r *universityRepository) GetUniversitiesFromDBRecord() ([]dto.UniversitiesResponse, error) {

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
