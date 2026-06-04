package DiplomaService

import "BlockCertify/internal/models"

func (s *DiplomaService) GetByDiplomaID(diplomaID string) (*models.Diploma, error) {
	var diploma models.Diploma
	err := s.DB.Preload("MetaData").Where("public_id = ?", diplomaID).First(&diploma).Error
	if err != nil {
		return nil, err
	}
	return &diploma, nil
}
