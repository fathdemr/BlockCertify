package DiplomaService

import "BlockCertify/internal/models"

func (s *DiplomaService) GetAllDiplomaFromDatabase() <-chan models.Diploma {

	ch := make(chan models.Diploma)

	go func() {
		defer close(ch)

		var diplomas []models.Diploma
		if err := s.DB.Preload("MetaData").Find(&diplomas).Error; err != nil {
			return
		}

		for _, diploma := range diplomas {
			ch <- diploma
		}

	}()

	return ch
}
