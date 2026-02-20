package database

import (
	"BlockCertify/internal/config"
	"BlockCertify/internal/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := buildPostgresDSN(cfg)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Admin{},
		&models.Department{},
		&models.Diploma{},
		&models.DiplomaMetaData{},
		&models.Faculties{},
		&models.Student{},
		&models.Universities{},
		&models.User{},
	)
}
