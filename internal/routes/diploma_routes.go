package routes

import (
	"BlockCertify/internal/handlers"
	"BlockCertify/internal/middleware"
	"BlockCertify/internal/models"

	"github.com/gin-gonic/gin"
)

func DiplomaRoutes(diploma *gin.RouterGroup, auth middleware.AuthMiddleware) {

	// Issuance endpoints spend platform funds (Arweave + Polygon) — admin only.
	diploma.POST("/upload", auth.Authorize(), auth.RequireRole(models.RoleAdmin), handlers.Upload)
	diploma.POST("/prepare-upload", auth.Authorize(), auth.RequireRole(models.RoleAdmin), handlers.PrepareUpload)
	diploma.POST("/confirm-upload", auth.Authorize(), auth.RequireRole(models.RoleAdmin), handlers.ConfirmUpload)

	// Verification is public by design: anyone can check a diploma.
	diploma.POST("/verify", handlers.Verify)
	diploma.GET("/records", handlers.GetDiplomaRecords)
	diploma.GET("/records/:diplomaId", handlers.GetDiplomaById)
}
