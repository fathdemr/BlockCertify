package UniversityService

import "BlockCertify/internal/models"

func (s *UniversityService) GetUniversityByID(id string) (models.Universities, error) {

	var university models.Universities

	err := s.DB.Where("id = ?", id).First(&university).Error
	if err != nil {
		return models.Universities{}, err
	}
	return university, nil
}
