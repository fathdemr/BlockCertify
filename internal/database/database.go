package database

import (
	"BlockCertify/internal/config"
	"BlockCertify/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := buildPostgresDSN(cfg)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&models.Diploma{}); err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&models.DiplomaMetaData{}); err != nil {
		return nil, err
	}

	return db, nil
}
