package services

import (
	"BlockCertify/internal/dto"
	"BlockCertify/internal/repositories"
	"log/slog"

	"github.com/gofrs/uuid/v5"
)

type FacultyService interface {
	GetFaculties(universityID uuid.UUID) ([]dto.FacultiesResponse, error)
}

type facultyService struct {
	repo repositories.FacultyRepository
}

func NewFacultyService(repo repositories.FacultyRepository) FacultyService {
	return &facultyService{
		repo: repo,
	}
}

func (s *facultyService) GetFaculties(universityID uuid.UUID) ([]dto.FacultiesResponse, error) {
	faculties, err := s.repo.GetFaculties(universityID.String())
	if err != nil {
		slog.Error("Failed to get faculties from DB: %v", err)
		return nil, err
	}

	var response []dto.FacultiesResponse
	for _, faculty := range faculties {
		response = append(response, dto.FacultiesResponse{
			ID:   faculty.ID,
			Name: faculty.FacultyName,
		})
	}

	return response, nil
}
