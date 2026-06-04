package FacultyService

import (
	"gorm.io/gorm"
)

type FacultyService struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *FacultyService {
	return &FacultyService{
		DB: db,
	}
}
