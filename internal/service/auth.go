package service

import (
	"context"
	"fmt"
	"time"

	"github.com/AlexFox86/auth-service/internal/delivery/dto"
	"github.com/AlexFox86/auth-service/internal/models"
	"github.com/AlexFox86/auth-service/internal/pkg/crypto"
	"github.com/AlexFox86/auth-service/internal/pkg/token"
	"github.com/AlexFox86/auth-service/internal/repository/postgres"
)

// Service provides methods for authentication and registration
type Service struct {
	repo        postgres.Repository
	jwtSecret   []byte
	tokenExpiry time.Duration
}

// New creates a new authentication service
func New(repo postgres.Repository, jwtSecret string, tokenExpiry time.Duration) *Service {
	return &Service{
		repo:        repo,
		jwtSecret:   []byte(jwtSecret),
		tokenExpiry: tokenExpiry,
	}
}

// JwtSecret returns the 'jwtSecret' field
func (s *Service) JwtSecret() []byte {
	return s.jwtSecret
}

// TokenExpiry returns the 'tokenExpiry' field
func (s *Service) TokenExpiry() time.Duration {
	return s.tokenExpiry
}

// Register creates a new user
func (s *Service) Register(ctx context.Context, req *dto.RegisterRequest) (*models.User, error) {
	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

// Login performs user authentication
func (s *Service) Login(ctx context.Context, req *dto.LoginRequest) (*dto.Response, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, models.ErrInvalidCredentials
	}

	if err := crypto.CheckPassword(req.Password, user.Password); err != nil {
		return nil, models.ErrInvalidCredentials
	}

	token, err := token.GenerateToken(user, s.jwtSecret, s.tokenExpiry)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &dto.Response{
		Token: token,
		User:  *user,
	}, nil
}
