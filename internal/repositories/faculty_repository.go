package repositories

import (
	"BlockCertify/internal/models"

	"gorm.io/gorm"
)

type FacultyRepository interface {
	GetFaculties(universityID string) ([]models.Faculties, error)
}

type facultyRepository struct {
	db *gorm.DB
}

func NewFacultyRepository(db *gorm.DB) FacultyRepository {
	return &facultyRepository{
		db: db,
	}
}

func (r *facultyRepository) GetFaculties(universityID string) ([]models.Faculties, error) {

	var faculties []models.Faculties

	err := r.db.Where("university_id = ?", universityID).Find(&faculties).Error
	if err != nil {
		return nil, err
	}

	return faculties, nil
}
