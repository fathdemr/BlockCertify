package UserService

import "BlockCertify/internal/models"

// Exists 1 -> User exists 0 -> User not exists
func (s *UserService) isExist(email string) (bool, error) {
	var count int64

	err := s.DB.Model(&models.User{}).
		Where("email = ?", email).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
