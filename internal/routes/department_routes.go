package routes

import (
	"BlockCertify/internal/handlers"

	"github.com/gin-gonic/gin"
)

func DepartmentRoutes(api *gin.RouterGroup, h *handlers.DepartmentHandler) {
	api.POST("/departments", h.GetDepartmentByID)
}
