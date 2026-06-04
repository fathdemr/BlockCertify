package routes

import (
	"BlockCertify/internal/handlers"

	"github.com/gin-gonic/gin"
)

func DepartmentRoutes(api *gin.RouterGroup) {
	api.GET("/departments", handlers.GetDepartmentByID)
}
