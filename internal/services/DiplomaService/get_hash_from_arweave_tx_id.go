package DiplomaService

import "BlockCertify/internal/models"

func (s *DiplomaService) GetHashFromArweaveTxID(arweaveTxID string) (string, error) {
	var diploma models.Diploma
	err := s.DB.Where("arweave_tx_id = ?", arweaveTxID).First(&diploma).Error
	if err != nil {
		return "", err
	}
	hash := diploma.Hash
	return hash, nil
}
