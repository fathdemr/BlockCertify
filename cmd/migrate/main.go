package main

import (
	"BlockCertify/internal/models"
	"fmt"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		mustEnv("DB_LIVE_HOST"),
		getEnv("DB_LIVE_PORT", "5432"),
		mustEnv("DB_LIVE_USER_NAME"),
		mustEnv("DB_LIVE_PASSWORD"),
		mustEnv("DB_LIVE_DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		slog.Error("DB bağlantısı başarısız", "error", err)
		os.Exit(1)
	}

	slog.Info("DB bağlandı, migration başlıyor...")

	err = db.AutoMigrate(
		&models.User{},
		&models.Universities{},
		&models.Faculties{},
		&models.Department{},
		&models.Admin{},
		&models.Student{},
		&models.Diploma{},
		&models.DiplomaMetaData{},
	)
	if err != nil {
		slog.Error("Migration başarısız", "error", err)
		os.Exit(1)
	}

	slog.Info("✅ Migration tamamlandı")
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		slog.Error("zorunlu env var eksik", "key", key)
		os.Exit(1)
	}
	return v
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
