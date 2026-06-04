package routes

import (
	"BlockCertify/internal/handlers"

	"github.com/gin-gonic/gin"
)

func WalletRoutes(api *gin.RouterGroup) {
	api.GET("/wallet/status", handlers.WalletStatus)
}
