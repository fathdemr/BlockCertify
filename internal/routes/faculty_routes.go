package routes

import (
	"BlockCertify/internal/handlers"

	"github.com/gin-gonic/gin"
)

func FacultyRoutes(api *gin.RouterGroup) {
	api.GET("/faculties", handlers.GetFaculties)
}
