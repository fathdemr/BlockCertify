package services

import (
	"BlockCertify/internal/dto"
	"BlockCertify/internal/repositories"
	"log/slog"
)

type UniversityService interface {
	GetUniversitiesFromDBRecord() ([]dto.UniversitiesResponse, error)
	GetUniversityByID(id string) (dto.UniversitiesResponse, error)
}

type universityService struct {
	repo repositories.UniversityRepository
}

func NewUniversityService(repo repositories.UniversityRepository) UniversityService {
	return &universityService{
		repo: repo,
	}
}

func (s *universityService) GetUniversitiesFromDBRecord() ([]dto.UniversitiesResponse, error) {

	universities, err := s.repo.GetUniversitiesFromDBRecord()
	if err != nil {
		slog.Error("Failed to get universities from DB: %v", err)
		return nil, err
	}
	return universities, nil
}

func (s *universityService) GetUniversityByID(id string) (dto.UniversitiesResponse, error) {

	university, err := s.repo.GetUniversityByID(id)
	if err != nil {
		slog.Error("Failed to get university by ID: %v", err)
		return dto.UniversitiesResponse{}, err
	}

	return dto.UniversitiesResponse{
		ID:   university.ID,
		Name: university.Name,
	}, nil
}
