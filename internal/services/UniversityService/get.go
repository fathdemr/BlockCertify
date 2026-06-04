package UniversityService

import (
	"BlockCertify/internal/dto"
	"BlockCertify/internal/models"
	"log/slog"
)

func (s *UniversityService) GetUniversitiesFromDBRecord() ([]dto.UniversitiesResponse, error) {
	var response []dto.UniversitiesResponse

	err := s.DB.
		Model(&models.Universities{}).
		Select("id, name").
		Scan(&response).Error

	if err != nil {
		slog.Error("failed to fetch universities", "err", err)
		return nil, err
	}

	return response, nil
}
