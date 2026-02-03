package repositories

import "gorm.io/gorm"

type DiplomaRepository interface {
	CreateTransaction() *gorm.DB
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
