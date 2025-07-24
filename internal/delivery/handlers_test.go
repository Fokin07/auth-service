package delivery

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlexFox86/auth-service/internal/delivery/dto"
	"github.com/AlexFox86/auth-service/internal/models"
	"github.com/AlexFox86/auth-service/internal/pkg/crypto"
	mockrepo "github.com/AlexFox86/auth-service/internal/repository/mock"
	"github.com/AlexFox86/auth-service/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var errUserNotFound = errors.New("user not found")

func TestHandlerRegister(t *testing.T) {
	mockRepo := new(mockrepo.MockRepository)
	service := service.New(mockRepo, "secret", time.Hour)
	handler := NewHandler(service)

	tests := []struct {
		name           string
		requestBody    any
		mockSetup      func()
		expectedStatus int
	}{
		{
			name: "successful registration",
			requestBody: dto.RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).
					Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid request body",
			requestBody: map[string]any{
				"username": 123, // incorrect type
				"email":    "test@example.com",
				"password": "password123",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "validation error",
			requestBody: dto.RegisterRequest{
				Username: "te", // the name is too short
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "email already exists",
			requestBody: dto.RegisterRequest{
				Username: "testuser",
				Email:    "exists@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).
					Return(errEmailExists).Once()
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Register(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandlerLogin(t *testing.T) {
	mockRepo := new(mockrepo.MockRepository)
	service := service.New(mockRepo, "secret", time.Hour)
	handler := NewHandler(service)

	hashedPassword, _ := crypto.HashPassword("password123")

	tests := []struct {
		name           string
		requestBody    any
		mockSetup      func()
		expectedStatus int
	}{
		{
			name: "successful login",
			requestBody: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").
					Return(&models.User{
						ID:       uuid.New(),
						Username: "testuser",
						Email:    "test@example.com",
						Password: hashedPassword,
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid credentials",
			requestBody: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func() {
				mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").
					Return(&models.User{
						ID:       uuid.New(),
						Username: "testuser",
						Email:    "test@example.com",
						Password: hashedPassword,
					}, nil)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "user not found",
			requestBody: dto.LoginRequest{
				Email:    "notfound@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockRepo.On("GetUserByEmail", mock.Anything, "notfound@example.com").
					Return(nil, errUserNotFound)
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Login(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}
