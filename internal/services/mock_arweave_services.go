package services

type MockArweaveService struct {
	TxID string
	Err  error
}

func NewMockArweaveService(txID string) *MockArweaveService {
	return &MockArweaveService{
		TxID: txID,
	}
}

func (m *MockArweaveService) Upload(filePath, fileHash string) (string, error) {
	if m.Err != nil {
		return "", m.Err
	}

	return m.TxID, nil
}
