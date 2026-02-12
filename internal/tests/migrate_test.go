package tests

import (
	"BlockCertify/internal/config"
	"BlockCertify/internal/database"
	"BlockCertify/internal/models"
	"log/slog"
	"testing"
)

func TestAutoMigrate(t *testing.T) {

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "err", err)
	}

	db, err := database.Init(cfg.Db)
	if err != nil {
		slog.Error("failed to init database", "err", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.DiplomaMetaData{},
		&models.Diploma{},
	)

	if err != nil {
		slog.Error("failed to auto migrate", "err", err)
	}

}
