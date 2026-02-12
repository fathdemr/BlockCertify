package database

import (
	"BlockCertify/internal/config"
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
