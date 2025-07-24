package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AlexFox86/auth-service/internal/delivery/dto"
	"github.com/AlexFox86/auth-service/internal/models"
	"github.com/AlexFox86/auth-service/internal/pkg/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mockrepo "github.com/AlexFox86/auth-service/internal/repository/mock"
)

var (
	errUserNotFound = errors.New("user not found")
	errEmailExists  = errors.New("email already exists")
)

func TestNew(t *testing.T) {
	t.Run("test New()", func(t *testing.T) {
		mockRepo := new(mockrepo.MockRepository)

		expService := &Service{
			repo:        mockRepo,
			jwtSecret:   []byte("secret"),
			tokenExpiry: time.Hour,
		}

		service := New(mockRepo, "secret", time.Hour)
		assert.Equal(t, expService, service)
	})
}

func TestServiceRegister(t *testing.T) {
	tests := []struct {
		name        string
		req         *dto.RegisterRequest
		mockSetup   func(*mockrepo.MockRepository)
		expected    models.User
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
			expected: models.User{
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
					Return(errEmailExists)
			},
			expected:    models.User{},
			expectedErr: errEmailExists,
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
				assert.Equal(t, tt.expected, user)
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
					Return(nil, errUserNotFound)
			},
			expected:    nil,
			expectedErr: ErrInvalidCredentials,
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
			expectedErr: ErrInvalidCredentials,
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
