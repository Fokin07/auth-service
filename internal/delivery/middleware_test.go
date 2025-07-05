package delivery

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlexFox86/auth-service/internal/models"
	"github.com/AlexFox86/auth-service/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/AlexFox86/auth-service/internal/pkg/token"
	mockrepo "github.com/AlexFox86/auth-service/internal/repository/mock"
)

func TestAuthMiddleware(t *testing.T) {
	mockRepo := new(mockrepo.MockRepository)
	service := service.New(mockRepo, "secret", time.Hour)
	handler := NewHandler(service)

	// Создаем тестовый токен
	user := &models.User{
		ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Username: "testuser",
	}
	token, _ := token.GenerateToken(user, service.JwtSecret(), service.TokenExpiry())

	tests := []struct {
		name           string
		setupRequest   func(*http.Request)
		expectedStatus int
		expectedUserID string
	}{
		{
			name: "valid token",
			setupRequest: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer "+token)
			},
			expectedStatus: http.StatusOK,
			expectedUserID: user.ID.String(),
		},
		{
			name:           "missing authorization header",
			setupRequest:   func(r *http.Request) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "invalid token",
			setupRequest: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer invalid.token")
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/protected", nil)
			tt.setupRequest(req)
			w := httptest.NewRecorder()

			// Creating a test handler that checks the context
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				userID, ok := r.Context().Value(contextKeyUserID).(string)
				if tt.expectedUserID != "" {
					assert.True(t, ok)
					assert.Equal(t, tt.expectedUserID, userID)
				} else {
					assert.False(t, ok)
				}
				w.WriteHeader(http.StatusOK)
			})

			// Using middleware
			handler.AuthMiddleware(testHandler).ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
