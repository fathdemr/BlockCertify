package UserService

import (
	"BlockCertify/internal/dto"
	"BlockCertify/internal/helper"
	"BlockCertify/internal/models"
	apperrors "BlockCertify/internal/pkg/errors"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/gofrs/uuid/v5"
)

var ErrToManyAttempts = errors.New("to many attempts")

func (s *UserService) Register(req dto.RegisterRequest) error {

	redisKey := fmt.Sprintf("%s:register:attempt", req.Email)
	registerAttemptDuration := time.Duration(1) * time.Minute
	registerAttempt := 5
	registerAttemptPenalty := time.Duration(1) * time.Minute

	if value, err := s.CacheService.Get(redisKey); err != nil {
		_ = s.CacheService.SetWithExpireDuration(redisKey, "1", registerAttemptDuration)
	} else {
		if attempts, _ := strconv.Atoi(value); attempts >= registerAttempt {
			attempts++
			_ = s.CacheService.UpdateKeyValue(redisKey, strconv.Itoa(attempts))
			_ = s.CacheService.ExtendTTL(redisKey, registerAttemptPenalty)

			return ErrToManyAttempts
		} else {
			attempts++
			_ = s.CacheService.UpdateKeyValue(redisKey, strconv.Itoa(attempts))
		}
	}

	if err := helper.Validate.Struct(&req); err != nil {
		return err
	}

	existing, err := s.isExist(req.Email)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	if existing {
		return apperrors.New(apperrors.ErrUserExists, "User with this email already exists", nil)
	}

	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("Failed to hash password: %w", err)
	}

	tx := s.CreateTransaction()
	defer tx.Rollback()

	user := models.User{
		ID:        uuid.Must(uuid.NewV7()),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      models.RoleAdmin,
	}

	if err := tx.Create(&user).Error; err != nil {
		slog.Error("Failed to create user: %v", err)
		return apperrors.New(apperrors.ErrUserCreationFailed, "User creation failed", err)
	}

	//If Everything is OK then delete the caches
	_ = s.CacheService.Delete(redisKey)

	return tx.Commit().Error
}

func (s *UserService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {

	// Bruteforce check with attempt limit 5 times for 60 seconds
	redisKey := fmt.Sprintf("%s:login:attempt", req.Email)
	loginAttemptDuration := time.Duration(60) * time.Second
	loginAttempt := 5
	loginAttemptPenalty := time.Duration(1) * time.Minute

	// If redis is not set, set duration 1 to 60 seconds to log in
	if value, errCache := s.CacheService.Get(redisKey); errCache != nil {
		_ = s.CacheService.SetWithExpireDuration(redisKey, "1", loginAttemptDuration)
	} else {
		// If current attempt bigger than penalty count or equal then set the extend the TTL
		if attempt, _ := strconv.Atoi(value); attempt >= loginAttempt {
			attempt++
			_ = s.CacheService.UpdateKeyValue(redisKey, strconv.Itoa(attempt))
			_ = s.CacheService.ExtendTTL(redisKey, loginAttemptPenalty)

			return nil, ErrToManyAttempts
		} else {
			attempt++
			_ = s.CacheService.UpdateKeyValue(redisKey, strconv.Itoa(attempt))
		}
	}

	if err := helper.Validate.Struct(&req); err != nil {
		return nil, err
	}

	user, err := s.FindByEmail(req.Email)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrInvalidCredentials, "Invalid email or password", nil)
	}

	ok := helper.VerifyPassword(user.Password, req.Password)
	if !ok {
		return nil, apperrors.New(apperrors.ErrInvalidCredentials, "Invalid email or password", nil)
	}

	token, err := s.CreateToken(user)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrTokenCreateFailed, "Token creation failed", err)
	}

	role := user.Role
	if role == "" {
		role = "admin" // Default for now
	}

	//If login is OK delete key
	_ = s.CacheService.Delete(redisKey)

	return &dto.LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: s.ExpiresInSeconds(),
		Role:      "admin",
	}, nil
}
