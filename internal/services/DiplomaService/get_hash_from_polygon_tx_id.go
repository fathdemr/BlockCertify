package DiplomaService

import "BlockCertify/internal/models"

func (s *DiplomaService) GetHashFromPolygonTxID(polygonTxID string) (string, error) {
	var diploma models.Diploma
	err := s.DB.Where("polygon_tx_id = ?", polygonTxID).First(&diploma).Error
	if err != nil {
		return "", err
	}
	hash := diploma.Hash
	return hash, nil
}
