package UserService

import "BlockCertify/internal/models"

func (s *UserService) Update(user *models.User) error {
	return s.DB.Save(user).Error
}
