package handlers

import (
	"BlockCertify/internal/config"
	"BlockCertify/internal/dto"
	"BlockCertify/internal/helper"
	apperrors "BlockCertify/internal/pkg/errors"
	"BlockCertify/internal/security"
	"BlockCertify/internal/services/CacheService"
	"BlockCertify/internal/services/UniversityService"
	"BlockCertify/internal/services/UserService"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func extractEmailFromToken(c *gin.Context) (string, error) {
	tokenStr := ""

	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		tokenStr = strings.TrimPrefix(strings.TrimPrefix(authHeader, "Bearer "), "bearer ")
	} else if cookie, err := c.Cookie("jwt"); err == nil && cookie != "" {
		tokenStr = cookie
	}

	if tokenStr == "" {
		return "", apperrors.New(apperrors.ErrInvalidToken, "No Authorization header", nil)
	}
	jwtHelper := security.NewJWTHelper(
		config.Params.GetString("crypto.my_secret_key"),
		time.Duration(config.Params.GetInt64("crypto.token_expire_duration_hour")),
	)
	claims, err := jwtHelper.Verify(tokenStr)
	if err != nil {
		return "", err
	}
	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return "", apperrors.New(apperrors.ErrInvalidToken, "Invalid token claims", nil)
	}
	return email, nil
}

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

func GetMe(c *gin.Context) {
	email, err := extractEmailFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userService := UserService.New(config.DB, config.Params)
	user, err := userService.FindByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, dto.MeResponse{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      string(user.Role),
	})
}

func UpdateMe(c *gin.Context) {
	email, err := extractEmailFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req dto.UpdateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userService := UserService.New(config.DB, config.Params)
	user, err := userService.FindByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashed)
	}

	if err := userService.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, dto.MeResponse{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      string(user.Role),
	})
}
