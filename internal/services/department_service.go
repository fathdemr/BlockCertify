package services

import (
	"BlockCertify/internal/dto"
	"BlockCertify/internal/repositories"

	"github.com/gofrs/uuid/v5"
)

type DepartmentService interface {
	GetDepartmentByID(facultyID uuid.UUID) ([]dto.DepartmentResponse, error)
}

type departmentService struct {
	repo repositories.DepartmentRepository
}

func NewDepartmentService(repo repositories.DepartmentRepository) DepartmentService {
	return &departmentService{repo: repo}
}

func (s *departmentService) GetDepartmentByID(facultyID uuid.UUID) ([]dto.DepartmentResponse, error) {

	departments, err := s.repo.GetDepartmentByID(facultyID)
	if err != nil {
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
