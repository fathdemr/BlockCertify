package routes

import (
	"BlockCertify/internal/handlers"

	"github.com/gin-gonic/gin"
)

func UniversityRoutes(api *gin.RouterGroup) {
	api.GET("/universities", handlers.GetUniversities)
}
