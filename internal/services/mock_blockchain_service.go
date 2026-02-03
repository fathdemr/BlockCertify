package services

import (
	"BlockCertify/internal/dto"
)

type MockBlockchainService struct{}

func NewMockBlockchainService() *MockBlockchainService {
	return &MockBlockchainService{}
}

func (m *MockBlockchainService) VerifyDiploma(string) (bool, string, error) {
	return false, "", nil
}

func (m *MockBlockchainService) StoreDiploma(_, _ string) (*dto.BlockchainResult, error) {
	return &dto.BlockchainResult{
		TransactionHash: "DEBUG_FAKE_POLYGON_TX",
		BlockNumber:     0,
	}, nil
}
