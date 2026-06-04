package UniversityService

import (
	"gorm.io/gorm"
)

type UniversityService struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *UniversityService {
	return &UniversityService{
		DB: db,
	}
}
