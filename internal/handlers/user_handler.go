package handlers

import (
	"BlockCertify/internal/config"
	"BlockCertify/internal/dto"
	"BlockCertify/internal/helper"
	apperrors "BlockCertify/internal/pkg/errors"
	"BlockCertify/internal/services/CacheService"
	"BlockCertify/internal/services/UniversityService"
	"BlockCertify/internal/services/UserService"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cacheService := CacheService.New(config.RedisClient)
	userService := UserService.New(config.DB, config.Params)
	userService.UseCacheService(cacheService)

	response, err := userService.Login(req)
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
			"error":   "Login failed",
			"details": err.Error(),
		})
		return
	}

	helper.SetCookie(c, "jwt", response.Token, time.Now().Add(time.Hour*1))

	c.JSON(http.StatusOK, response)
}
func RegisterAdmin(c *gin.Context) {

	var req dto.RegisterRequest
	userService := UserService.New(config.DB, config.Params)
	cacheService := CacheService.New(config.RedisClient)
	userService.UseCacheService(cacheService)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := userService.Register(req); err != nil {
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

func GetUniversities(c *gin.Context) {

	univerityService := UniversityService.New(config.DB)

	universities, err := univerityService.GetUniversitiesFromDBRecord()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, universities)
}

func Logout(c *gin.Context) {
	helper.ClearCookie(c, "jwt")
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}
