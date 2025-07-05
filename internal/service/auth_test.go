package service

import (
	"context"
	"testing"
	"time"

	"github.com/AlexFox86/auth-service/internal/delivery/dto"
	"github.com/AlexFox86/auth-service/internal/models"
	"github.com/AlexFox86/auth-service/internal/pkg/crypto"
	"github.com/AlexFox86/auth-service/internal/pkg/token"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mockrepo "github.com/AlexFox86/auth-service/internal/repository/mock"
)

func TestServiceRegister(t *testing.T) {
	tests := []struct {
		name        string
		req         *dto.RegisterRequest
		mockSetup   func(*mockrepo.MockRepository)
		expected    *models.User
		expectedErr error
	}{
		{
			name: "successful registration",
			req: &dto.RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(mr *mockrepo.MockRepository) {
				mr.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).
					Return(nil).
					Run(func(args mock.Arguments) {
						user := args.Get(1).(*models.User)
						user.ID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
					})
			},
			expected: &models.User{
				ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Username: "testuser",
				Email:    "test@example.com",
			},
			expectedErr: nil,
		},
		{
			name: "email already exists",
			req: &dto.RegisterRequest{
				Username: "testuser",
				Email:    "exists@example.com",
				Password: "password123",
			},
			mockSetup: func(mr *mockrepo.MockRepository) {
				mr.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).
					Return(models.ErrEmailExists)
			},
			expected:    nil,
			expectedErr: models.ErrEmailExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockrepo.MockRepository)
			tt.mockSetup(mockRepo)

			service := New(mockRepo, "secret", time.Hour)
			user, err := service.Register(context.Background(), tt.req)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.ID, user.ID)
				assert.Equal(t, tt.expected.Username, user.Username)
				assert.Equal(t, tt.expected.Email, user.Email)
				assert.NotEmpty(t, user.Password)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestServiceLogin(t *testing.T) {
	hashedPassword, _ := crypto.HashPassword("password123")

	tests := []struct {
		name        string
		req         *dto.LoginRequest
		mockSetup   func(*mockrepo.MockRepository)
		expected    *dto.Response
		expectedErr error
	}{
		{
			name: "successful login",
			req: &dto.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(mr *mockrepo.MockRepository) {
				mr.On("GetUserByEmail", mock.Anything, "test@example.com").
					Return(&models.User{
						ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Username: "testuser",
						Email:    "test@example.com",
						Password: hashedPassword,
					}, nil)
			},
			expected: &dto.Response{
				User: models.User{
					ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Username: "testuser",
					Email:    "test@example.com",
					Password: hashedPassword,
				},
			},
			expectedErr: nil,
		},
		{
			name: "user not found",
			req: &dto.LoginRequest{
				Email:    "notfound@example.com",
				Password: "password123",
			},
			mockSetup: func(mr *mockrepo.MockRepository) {
				mr.On("GetUserByEmail", mock.Anything, "notfound@example.com").
					Return(nil, models.ErrUserNotFound)
			},
			expected:    nil,
			expectedErr: models.ErrInvalidCredentials,
		},
		{
			name: "wrong password",
			req: &dto.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(mr *mockrepo.MockRepository) {
				mr.On("GetUserByEmail", mock.Anything, "test@example.com").
					Return(&models.User{
						ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Username: "testuser",
						Email:    "test@example.com",
						Password: hashedPassword,
					}, nil)
			},
			expected:    nil,
			expectedErr: models.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockrepo.MockRepository)
			tt.mockSetup(mockRepo)

			service := New(mockRepo, "secret", time.Hour)
			resp, err := service.Login(context.Background(), tt.req)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.User.ID, resp.User.ID)
				assert.Equal(t, tt.expected.User.Username, resp.User.Username)
				assert.Equal(t, tt.expected.User.Email, resp.User.Email)
				assert.NotEmpty(t, resp.Token)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestServiceValidateToken(t *testing.T) {
	service := New(nil, "secret", time.Hour)

	t.Run("valid token", func(t *testing.T) {
		user := &models.User{
			ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Username: "testuser",
		}
		tokenString, err := token.GenerateToken(user, service.jwtSecret, service.tokenExpiry)
		assert.NoError(t, err)

		claims, err := token.ValidateToken(tokenString, service.jwtSecret)
		assert.NoError(t, err)
		assert.Equal(t, user.ID.String(), claims["sub"])
		assert.Equal(t, user.Username, claims["username"])
	})

	t.Run("invalid token", func(t *testing.T) {
		claims, err := token.ValidateToken("invalid.token.string", service.jwtSecret)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("wrong signing method", func(t *testing.T) {
		// Creating a token with an incorrect signature method
		testToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"sub": "123",
		})
		tokenString, _ := testToken.SignedString([]byte("key"))

		claims, err := token.ValidateToken(tokenString, service.jwtSecret)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}
