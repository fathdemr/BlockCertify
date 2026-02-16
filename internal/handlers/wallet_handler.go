package handlers

import (
	"BlockCertify/internal/services"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	service services.WalletService
}

func NewWalletHandler(service services.WalletService) *WalletHandler {
	return &WalletHandler{
		service: service,
	}
}

func (h *WalletHandler) GetArweaveKeyFileJSON(c *gin.Context) {

	contentType := c.GetHeader("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{
			"error":   "Content-type header is not multipart/form-data",
			"details": "multipart/form-data required",
		})
		return
	}

	file, _, err := c.Request.FormFile("wallet")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"details": "wallet file is missing",
		})
		return
	}
	defer file.Close()

	keyBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"details": "failed to read wallet file",
		})
		return
	}

	wallet, err := h.service.ConnectWalletFromJSON(keyBytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"details": "failed to connect wallet",
		})
	}

	address := h.service.GetAddress(wallet)

	balance := h.service.GetBalance(wallet)

	c.JSON(http.StatusOK, gin.H{
		"address": address,
		"balance": balance,
		"status":  "wallet connected",
	})
}
