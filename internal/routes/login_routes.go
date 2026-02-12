package routes

import (
	"BlockCertify/internal/handlers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(api *gin.RouterGroup, h *handlers.UserHandler) {
	user := api.Group("/user")
	{
		user.POST("/login", h.Login)
		user.POST("/register", h.Register)
	}
}
