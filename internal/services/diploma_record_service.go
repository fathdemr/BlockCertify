package services

import (
	"BlockCertify/internal/dto"
	"BlockCertify/internal/repositories"
)

type HistoryService interface {
}

type historyService struct {
	repo repositories.DiplomaRepository
}

func NewHistoryService(repo repositories.DiplomaRepository) HistoryService {
	return &historyService{
		repo: repo,
	}
}

func (s *historyService) GetAllDiplomaFromDatabase() []dto.HistoryResponse {

	ch := s.repo.GetAllDiplomaFromDatabase()

	var response []dto.HistoryResponse

	for d := range ch {

		resp := dto.HistoryResponse{
			DiplomaID:  d.PublicID,
			UserName:   d.Owner,
			Department: d.MetaData.Department,
			CreateDate: d.CreatedAt,
			DiplomaPdf: d.ArweaveURL,
		}

		response = append(response, resp)
	}

	return response
}
