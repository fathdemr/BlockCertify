package FacultyService

import (
	"BlockCertify/internal/dto"
	"BlockCertify/internal/models"

	"github.com/gofrs/uuid/v5"
)

func (s *FacultyService) GetAllFacultiesByID(id uuid.UUID) ([]dto.FacultiesResponse, error) {
	if id == uuid.Nil {
		return nil, nil
	}

	var faculties []models.Faculties

	err := s.DB.Where("university_id = ?", id).Find(&faculties).Error
	if err != nil {
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
