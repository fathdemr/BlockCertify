package handlers

import (
	"BlockCertify/internal/dto"
	"BlockCertify/internal/helper"
	"BlockCertify/internal/services"
	"net/http"
	"time"

	apperrors "BlockCertify/pkg/errors"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service    services.UserService
	uniService services.UniversityService
}

func NewUserHandler(service services.UserService, uniService services.UniversityService) *UserHandler {
	return &UserHandler{
		service:    service,
		uniService: uniService,
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

	helper.SetCookie(c, "jwt", response.Token, time.Now().Add(time.Hour*1))

	c.JSON(http.StatusOK, response)
}
func (h *UserHandler) RegisterAdmin(c *gin.Context) {

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

func (h *UserHandler) GetUniversities(c *gin.Context) {

	universities, err := h.uniService.GetUniversitiesFromDBRecord()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, universities)
}

func (h *UserHandler) Logout(c *gin.Context) {
	helper.ClearCookie(c, "jwt")
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}
