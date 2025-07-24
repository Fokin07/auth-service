package mockrepo

import (
	"context"
	"time"

	"github.com/AlexFox86/auth-service/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockRepository mock repository for testing
type MockRepository struct {
	mock.Mock
}

// CreateUser creates a new user
func (m *MockRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	args := m.Called(ctx, &user)

	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt

	return user, args.Error(0)
}

// GetUserByEmail gets the user by email
func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
