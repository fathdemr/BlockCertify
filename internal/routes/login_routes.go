package routes

import (
	"BlockCertify/internal/handlers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(api *gin.RouterGroup) {
	user := api.Group("/user")
	{
		user.POST("/login", handlers.Login)
		user.POST("/register/admin", handlers.RegisterAdmin)
		user.POST("/logout", handlers.Logout)
	}
}
