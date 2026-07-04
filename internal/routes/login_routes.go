package routes

import (
	"BlockCertify/internal/handlers"
	"BlockCertify/internal/middleware"
	"BlockCertify/internal/models"

	"github.com/gin-gonic/gin"
)

func UserRoutes(api *gin.RouterGroup, auth middleware.AuthMiddleware) {
	user := api.Group("/user")
	{
		user.POST("/login", handlers.Login)
		// Creating new admins requires an existing admin (bootstrap via data/seed.sql).
		user.POST("/register/admin", auth.Authorize(), auth.RequireRole(models.RoleAdmin), handlers.RegisterAdmin)
		user.POST("/logout", handlers.Logout)
		user.GET("/me", handlers.GetMe)
		user.PUT("/me", handlers.UpdateMe)
	}
}
