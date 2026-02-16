package routes

import (
	"BlockCertify/internal/handlers"

	"github.com/gin-gonic/gin"
)

func UniversityRoutes(v1 *gin.RouterGroup, h *handlers.UserHandler) {
	v1.GET("/universities", h.GetUniversities)
}
