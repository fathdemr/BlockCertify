package repositories

import (
	"BlockCertify/internal/models"

	"gorm.io/gorm"
)

type DiplomaRepository interface {
	CreateTransaction() *gorm.DB
	GetByDiplomaID(diplomaID string) (*models.Diploma, error)
	GetHashFromArweaveTxID(arweaveTxID string) (string, error)
	GetHashFromPolygonTxID(polygonTxID string) (string, error)
}

type diplomaRepository struct {
	db *gorm.DB
}

func NewDiplomaRepository(db *gorm.DB) DiplomaRepository {
	return &diplomaRepository{
		db: db,
	}
}

func (r *diplomaRepository) CreateTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *diplomaRepository) GetByDiplomaID(diplomaID string) (*models.Diploma, error) {
	var diploma models.Diploma
	err := r.db.Preload("MetaData").Where("public_id = ?", diplomaID).First(&diploma).Error
	if err != nil {
		return nil, err
	}
	return &diploma, nil
}

func (r *diplomaRepository) GetHashFromArweaveTxID(arweaveTxID string) (string, error) {
	var diploma models.Diploma
	err := r.db.Where("arweave_tx_id = ?", arweaveTxID).First(&diploma).Error
	if err != nil {
		return "", err
	}
	hash := diploma.Hash
	return hash, nil
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
