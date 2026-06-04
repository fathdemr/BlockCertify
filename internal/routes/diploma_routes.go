package routes

import (
	"BlockCertify/internal/handlers"

	"github.com/gin-gonic/gin"
)

func DiplomaRoutes(diploma *gin.RouterGroup) {

	diploma.POST("/upload", handlers.Upload)
	diploma.POST("/verify", handlers.Verify)
	diploma.GET("/records", handlers.GetDiplomaRecords)
	diploma.GET("/records/:diplomaId", handlers.GetDiplomaById)

}
