package routes

import (
	"BlockCertify/internal/handlers"

	"github.com/gin-gonic/gin"
)

func DiplomaRoutes(diploma *gin.RouterGroup, d *handlers.DiplomaHandler) {

	diploma.POST("/prepare", d.PrepareUpload)
	diploma.POST("/confirm", d.ConfirmUpload)
	diploma.POST("/verify", d.Verify)
	diploma.GET("/records", d.GetDiplomaRecords)
	diploma.GET("/records/:diplomaId", d.GetDiplomaById)

}
