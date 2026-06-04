package DepartmentService

import (
	"BlockCertify/internal/dto"
	"BlockCertify/internal/models"
	"log/slog"

	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type DepartmentService struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *DepartmentService {
	return &DepartmentService{DB: db}
}

func (s *DepartmentService) GetDepartmentByID(facultyID uuid.UUID) ([]dto.DepartmentResponse, error) {
	var departments []models.Department

	err := s.DB.Where("faculty_id = ?", facultyID).Find(&departments).Error
	if err != nil {
		slog.Error("Failed to get departments by faculty ID", "error", err)
		return nil, err
	}

	var response []dto.DepartmentResponse
	for _, department := range departments {
		response = append(response, dto.DepartmentResponse{
			ID:   department.ID,
			Name: department.DepartmentName,
		})
	}
	return response, nil
}
