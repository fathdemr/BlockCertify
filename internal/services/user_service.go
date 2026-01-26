package services

import (
	"BlockCertify/internal/dto"
	"BlockCertify/internal/helper"
	"BlockCertify/internal/models"
	"BlockCertify/internal/repositories"
	"BlockCertify/internal/security"
	apperrors "BlockCertify/pkg/errors"
	"fmt"

	"github.com/google/uuid"
)

type UserService interface {
	Register(req dto.RegisterRequest) error
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
}

type userService struct {
	repo        *repositories.UserRepository
	tokenHelper security.TokenHelper
}

func NewUserService(repo *repositories.UserRepository, tokenHelper security.TokenHelper) UserService {

	return &userService{
		repo:        repo,
		tokenHelper: tokenHelper,
	}
}

func (s *userService) Register(req dto.RegisterRequest) error {

	if err := helper.Validate.Struct(&req); err != nil {
		return err
	}

	existing, err := s.repo.Exists(req.Email)
	if err == nil && existing {
		return apperrors.New(apperrors.ErrUserExists, "User already exists", nil)
	}

	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("Failed to hash password: %w", err)
	}

	user := models.User{
		ID:          uuid.NewString(),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Password:    hashedPassword,
		Institution: req.Institution,
	}

	return s.repo.Create(&user)
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

	token, err := s.tokenHelper.Create(user.Email)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrTokenCreateFailed, "Token creation failed", err)
	}

	return &dto.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   s.tokenHelper.ExpiresInSeconds(),
	}, nil
}
