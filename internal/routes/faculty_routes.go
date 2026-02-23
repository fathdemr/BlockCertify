package routes

import (
	"BlockCertify/internal/handlers"

	"github.com/gin-gonic/gin"
)

func FacultyRoutes(api *gin.RouterGroup, h *handlers.FacultyHandler) {
	api.POST("/faculties", h.GetFaculties)
}
