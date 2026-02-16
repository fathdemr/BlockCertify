package routes

import (
	"BlockCertify/internal/handlers"

	"github.com/gin-gonic/gin"
)

func WalletRoutes(api *gin.RouterGroup, h *handlers.WalletHandler) {
	api.POST("/upload-key-file", h.GetArweaveKeyFileJSON)
}
