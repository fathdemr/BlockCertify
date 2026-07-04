package main

import (
	"BlockCertify/internal/config"
	"BlockCertify/internal/logger"
	"BlockCertify/internal/middleware"
	"BlockCertify/internal/models"
	"BlockCertify/internal/repositories"
	"BlockCertify/internal/routes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	logger.Init()

	if err := config.InitConfigFile("./internal"); err != nil {
		panic(err)
	}
	if err := config.InitJWT(); err != nil {
		panic(err)
	}
	if err := config.InitDB(); err != nil {
		panic(err)
	}
	if err := config.RedisClient.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	if err := config.DB.AutoMigrate(
		&models.User{},
		&models.Universities{},
		&models.Faculties{},
		&models.Department{},
		&models.Admin{},
		&models.Student{},
		&models.Diploma{},
		&models.DiplomaMetaData{},
	); err != nil {
		panic(err)
	}
	slog.Info("✅ migration tamamlandı")

	app := gin.New()

	app.ForwardedByClientIP = true
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"https://blockcertify.uk",
		"https://www.blockcertify.uk",
		"http://localhost:3000",
		"http://localhost:5173",
	}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{
		"Origin",
		"Authorization",
		"Content-Type",
		"Bilet",
		"ApiKey",
		"ApiSecret",
		"X-Forwarded-For",
		"X-Real-Ip",
		"User-Agent",
		"Referer",
		"Accept-Language",
		"Accept-Encoding",
		"Cache-Control",
		"Connection",
		"DNT",
		"X-Requested-With",
		"Sec-Fetch-Site",
		"Sec-Fetch-Mode",
		"Sec-Fetch-Dest",
		"X-Device-Id",
		"X-Device-Model",
		"X-OS-Version",
		"X-Client-Version",
		"X-Platform",
		"X-Timezone",
		"X-Session-Id",
		"X-App-Id",
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}
	app.Use(cors.New(corsConfig))

	routes.HealthRoutes(app)

	authMw := middleware.NewAuthMiddleware(config.JWT, repositories.NewUserRepository(config.DB))

	exapi := app.Group("/exapi")
	routes.UserRoutes(exapi, authMw)
	routes.UniversityRoutes(exapi)
	routes.PingRoutes(exapi)
	routes.FacultyRoutes(exapi)
	routes.DepartmentRoutes(exapi)

	api := app.Group("/api")
	api.Use(middleware.RateLimit(60, time.Minute))
	routes.DiplomaRoutes(api, authMw)
	routes.WalletRoutes(api)

	// ── 4. HTTP Server ─────────────────────────────────────────────────────────
	port := config.Params.GetString("app.port")
	if port == "" {
		port = "5075"
	}
	if !strings.HasPrefix(port, ":") {
		port = "0.0.0.0:" + port
	}

	srv := &http.Server{
		Addr:              port,
		Handler:           app,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		slog.Info("✅ sunucu başlatıldı", "addr", "http://"+port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("❌ sunucu hatası", "error", err)
			os.Exit(1)
		}
	}()

	// ── 5. Graceful Shutdown ───────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("🔄 sunucu kapatılıyor...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("❌ graceful shutdown başarısız", "error", err)
	}

	if sqlDB, err := config.DB.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			slog.Error("❌ postgres kapatılamadı", "error", err)
		}
	}
	if config.RedisClient != nil {
		if err := config.RedisClient.Close(); err != nil {
			slog.Error("❌ redis kapatılamadı", "error", err)
		}
	}

	slog.Info("✅ sunucu düzgünce kapatıldı")
}
