package handlers

import (
	"BlockCertify/internal/dto"
	"BlockCertify/internal/services"
	"net/http"

	apperrors "BlockCertify/pkg/errors"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Login(c *gin.Context) {

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Login(req)
	if err != nil {
		appErr, ok := err.(*apperrors.AppError)
		if ok {
			details := ""
			if appErr.Err != nil {
				details = appErr.Err.Error()
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   appErr.Message,
				"details": details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Login failed",
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) Register(c *gin.Context) {

	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := h.service.Register(req); err != nil {
		appErr, ok := err.(*apperrors.AppError)
		if ok {
			details := ""
			if appErr.Err != nil {
				details = appErr.Err.Error()
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   appErr.Message,
				"details": details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Success",
	})

}
