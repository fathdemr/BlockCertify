package repositories

import (
	"BlockCertify/internal/models"
	"log/slog"

	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type DepartmentRepository interface {
	GetDepartmentByID(facultyID uuid.UUID) ([]models.Department, error)
}

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{db: db}
}

func (r *departmentRepository) GetDepartmentByID(facultyID uuid.UUID) ([]models.Department, error) {
	var departments []models.Department

	err := r.db.Where("faculty_id = ?", facultyID).Find(&departments).Error
	if err != nil {
		slog.Error("Failed to get departments by faculty ID", "error", err)
		return nil, err
	}

	return departments, nil
}
