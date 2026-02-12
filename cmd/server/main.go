package main

import (
	"BlockCertify/internal/config"
	"BlockCertify/internal/database"
	"BlockCertify/internal/handlers"
	"BlockCertify/internal/logger"
	"BlockCertify/internal/middleware"
	"BlockCertify/internal/repositories"
	"BlockCertify/internal/security"
	"BlockCertify/internal/services"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	logger.Init()

	//Load conf
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	r := gin.Default()

	// CORS config (allow frontend localhost:5173)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	db, err := database.Init(cfg.Db)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	contractRepo, err := repositories.NewContractRepository(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize contract repository: %v", err)
	}
	defer contractRepo.Close()

	userRepo := repositories.NewUserRepository(db)
	diplomaRepo := repositories.NewDiplomaRepository(db)

	tokenHelper := security.NewJWTHelper(
		cfg.JWTConfig.JWTSecret,
		cfg.JWTConfig.JWTExpireHours)

	//Mock Services
	//arweaveService := services.NewMockArweaveService("DEBUG_FAKE_ARWEAVE_TX")
	//blockchainService := services.NewMockBlockchainService()

	//Initialize services
	arweaveService := services.NewArweaveService(cfg)
	blockchainService := services.NewBlockChainService(cfg, contractRepo)
	diplomaService := services.NewDiplomaService(arweaveService, blockchainService, diplomaRepo)
	userService := services.NewUserService(userRepo, tokenHelper)
	AuthMiddleware := middleware.NewAuthMiddleware(tokenHelper, userRepo)

	//Initialize handlers
	diplomaHandler := handlers.NewDiplomaHandler(diplomaService)
	userHandler := handlers.NewUserHandler(userService)

	r.POST("/api/upload-diploma", diplomaHandler.Upload)
	r.POST("/api/verify-diploma", diplomaHandler.Verify)
	r.POST("/api/user/login", userHandler.Login)
	r.POST("/api/user/register", userHandler.Register)
	r.GET("/api/diploma-records", AuthMiddleware.Authorize(), diplomaHandler.GetDiplomaRecords)

	r.Static("/public", "./public")
	//Start server
	log.Printf("Server running on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
