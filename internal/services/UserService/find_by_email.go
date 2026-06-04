package UserService

import "BlockCertify/internal/models"

func (s *UserService) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
