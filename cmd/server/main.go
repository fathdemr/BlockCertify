package main

import (
	"BlockCertify/internal/config"
	"BlockCertify/internal/database"
	"BlockCertify/internal/handlers"
	"BlockCertify/internal/repositories"
	"BlockCertify/internal/security"
	"BlockCertify/internal/services"
	"log"
	"net/http"
	"time"

	"github.com/rs/cors"
)

func main() {
	//Load conf
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

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

	//Initialize handlers
	diplomaHandler := handlers.NewDiplomaHandler(diplomaService)
	userHandler := handlers.NewUserHandler(userService)

	//Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/upload-diploma", diplomaHandler.Upload)
	mux.HandleFunc("/api/verify-diploma", diplomaHandler.Verify)
	mux.HandleFunc("/api/user/login", userHandler.Login)
	mux.HandleFunc("/api/user/register", userHandler.Register)
	mux.Handle("/", http.FileServer(http.Dir("public")))

	//Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           int((12 * time.Hour).Seconds()),
	})

	handler := c.Handler(mux)
	//Start server
	log.Printf("Server running on port %s", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
