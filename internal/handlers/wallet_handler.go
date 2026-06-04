package handlers

import (
	"BlockCertify/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WalletStatus(c *gin.Context) {
	arweaveAddress, arweaveBalance := config.Arweave.GetStatus()
	polygonAddress := config.Blockchain.GetAddress()

	c.JSON(http.StatusOK, gin.H{
		"arweave": gin.H{
			"address": arweaveAddress,
			"balance": arweaveBalance + " AR",
		},
		"polygon": gin.H{
			"address":         polygonAddress,
			"contractAddress": config.Params.GetString("polygon.contractAddress"),
			"network":         "Polygon Amoy Testnet",
		},
	})
}
