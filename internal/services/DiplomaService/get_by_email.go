package DiplomaService

import "BlockCertify/internal/models"

func (s *DiplomaService) GetDiplomasByEmail(email string) <-chan models.Diploma {
	ch := make(chan models.Diploma)

	go func() {
		defer close(ch)

		var diplomas []models.Diploma
		if err := s.DB.
			Preload("MetaData").
			Joins("JOIN diploma_metadata ON diploma_metadata.diploma_id = diploma.id").
			Where("diploma_metadata.email = ?", email).
			Find(&diplomas).Error; err != nil {
			return
		}

		for _, diploma := range diplomas {
			ch <- diploma
		}
	}()

	return ch
}
