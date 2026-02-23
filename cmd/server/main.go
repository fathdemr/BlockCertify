package main

import (
	"BlockCertify/internal/config"
	"BlockCertify/internal/database"
	"BlockCertify/internal/handlers"
	"BlockCertify/internal/logger"
	"BlockCertify/internal/middleware"
	"BlockCertify/internal/repositories"
	"BlockCertify/internal/routes"
	"BlockCertify/internal/security"
	"BlockCertify/internal/services"
	"log"
	"log/slog"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	logger.Init()

	//Load conf
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config: %v", err)
	}

	r := gin.Default()

	// CORS config
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost", // Docker frontend (Nginx port 80)
			"http://localhost:80",
			"http://localhost:5173", // Vite dev server
			"http://localhost:8080",
			"http://localhost:3000",
			"http://185.252.234.84",
			"http://185.252.234.84:80",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	db, err := database.Init(cfg.Db)
	if err != nil {
		slog.Error("Failed to initialize database: %v", err)
	}

	err = database.Migrate(db)
	if err != nil {
		slog.Error("Failed to migrate database: %v", err)
	}

	contractRepo, err := repositories.NewContractRepository(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize contract repository: %v", err)
	}
	defer contractRepo.Close()

	//Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	diplomaRepo := repositories.NewDiplomaRepository(db)
	uniRepo := repositories.NewUniversityRepository(db)
	facultyRepo := repositories.NewFacultyRepository(db)
	departmentRepo := repositories.NewDepartmentRepository(db)

	tokenHelper := security.NewJWTHelper(
		cfg.JWTConfig.JWTSecret,
		cfg.JWTConfig.JWTExpireHours,
	)

	//Mock Services
	//arweaveService := services.NewMockArweaveService("DEBUG_FAKE_ARWEAVE_TX")
	//blockchainService := services.NewMockBlockchainService()

	//Initialize services
	arweaveService := services.NewArweaveService(cfg)
	blockchainService := services.NewBlockChainService(cfg, contractRepo)
	diplomaService := services.NewDiplomaService(arweaveService, blockchainService, diplomaRepo)
	userService := services.NewUserService(userRepo, tokenHelper, uniRepo)
	AuthMiddleware := middleware.NewAuthMiddleware(tokenHelper, userRepo)
	uniService := services.NewUniversityService(uniRepo)
	walletService := services.NewWalletService()
	facultyService := services.NewFacultyService(facultyRepo)
	departmentService := services.NewDepartmentService(departmentRepo)

	//Initialize handlers
	diplomaHandler := handlers.NewDiplomaHandler(diplomaService)
	userHandler := handlers.NewUserHandler(userService, uniService)
	walletHandler := handlers.NewWalletHandler(walletService)
	facultyHandler := handlers.NewFacultyHandler(facultyService)
	departmentHandler := handlers.NewDepartmentHandler(departmentService)

	api := r.Group("/api/v1")
	auth := api.Group("/auth")
	diploma := api.Group("/diploma")
	wallet := api.Group("/wallet")

	//Public routes
	routes.UserRoutes(auth, userHandler)
	routes.UniversityRoutes(api, userHandler)
	routes.PingRoutes(api)
	routes.FacultyRoutes(api, facultyHandler)
	routes.DepartmentRoutes(api, departmentHandler)

	//Protected routes
	diploma.Use(AuthMiddleware.Authorize())
	routes.DiplomaRoutes(diploma, diplomaHandler)
	routes.WalletRoutes(wallet, walletHandler)

	r.Static("/public", "./public")
	//Start server
	log.Printf("Server running on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
