package repositories

/*
type DiplomaRepository interface {
	CreateTransaction() *gorm.DB
	GetByDiplomaID(diplomaID string) (*models.Diploma, error)
	GetHashFromArweaveTxID(arweaveTxID string) (string, error)
	GetHashFromPolygonTxID(polygonTxID string) (string, error)
	GetAllDiplomaFromDatabase() <-chan models.Diploma
}

type diplomaRepository struct {
	DB *gorm.DB
}

func NewDiplomaRepository(db *gorm.DB) *DiplomaRepository {
	return &DiplomaRepository{
		DB: db,
	}
}

func (r *diplomaRepository) GetHashFromPolygonTxID(polygonTxID string) (string, error) {
	var diploma models.Diploma
	err := r.db.Where("polygon_tx_id = ?", polygonTxID).First(&diploma).Error
	if err != nil {
		return "", err
	}
	hash := diploma.Hash
	return hash, nil
}

func (r *diplomaRepository) GetAllDiplomaFromDatabase() <-chan models.Diploma {

	ch := make(chan models.Diploma)

	go func() {
		defer close(ch)

		var diplomas []models.Diploma
		if err := r.db.Preload("MetaData").Find(&diplomas).Error; err != nil {
			return
		}

		for _, diploma := range diplomas {
			ch <- diploma
		}

	}()

	return ch
}


*/
