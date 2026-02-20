package services

import (
	"BlockCertify/internal/dto"
	"BlockCertify/internal/helper"
	"BlockCertify/internal/models"
	apperrors "BlockCertify/internal/pkg/errors"
	"BlockCertify/internal/repositories"
	"BlockCertify/internal/security"
	"fmt"
	"log/slog"

	"github.com/gofrs/uuid/v5"
)

type UserService interface {
	Register(req dto.RegisterRequest) error
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
}

type userService struct {
	repo        repositories.UserRepository
	tokenHelper security.TokenHelper
	repoUni     repositories.UniversityRepository
}

func NewUserService(repo repositories.UserRepository, tokenHelper security.TokenHelper, repoUni repositories.UniversityRepository) UserService {

	return &userService{
		repo:        repo,
		tokenHelper: tokenHelper,
		repoUni:     repoUni,
	}
}

func (s *userService) Register(req dto.RegisterRequest) error {

	if err := helper.Validate.Struct(&req); err != nil {
		return err
	}

	existing, err := s.repo.Exists(req.Email)
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

	uni, err := s.repoUni.GetUniversityByID(req.UniversityID)
	if err != nil {
		slog.Error("Failed to get university by ID: %v", err)
		return apperrors.New(apperrors.ErrUniversityNotFound, "University not found", err)
	}

	tx := s.repo.CreateTransaction()
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

	admin := models.Admin{
		ID:           uuid.Must(uuid.NewV7()),
		UserID:       user.ID,
		UniversityID: uni.ID,
	}

	if err := tx.Create(&admin).Error; err != nil {
		slog.Error("Failed to create admin: %v", err)
		return apperrors.New(apperrors.ErrAdminCreationFailed, "Admin creation failed", err)
	}

	return tx.Commit().Error
}

func (s *userService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {

	if err := helper.Validate.Struct(&req); err != nil {
		return nil, err
	}

	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrInvalidCredentials, "Invalid email or password", nil)
	}

	ok := helper.VerifyPassword(user.Password, req.Password)
	if !ok {
		return nil, apperrors.New(apperrors.ErrInvalidCredentials, "Invalid email or password", nil)
	}

	token, err := s.tokenHelper.CreateToken(user.Email)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrTokenCreateFailed, "Token creation failed", err)
	}

	role := user.Role
	if role == "" {
		role = "admin" // Default for now
	}

	return &dto.LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: s.tokenHelper.ExpiresInSeconds(),
		Role:      "admin",
	}, nil
}
