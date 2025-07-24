package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/AlexFox86/auth-service/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var errUserNotFound = errors.New("user not found")

// Repository interface for working with storage
type Repository interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

// PgRepository the structure for working with PostgreSQL database
type PgRepository struct {
	db *sqlx.DB
}

// NewPgRepository creates a new object of 'PgRepository' type
// and returns a pointer to it.
func NewPgRepository(db *sqlx.DB) *PgRepository {
	return &PgRepository{db: db}
}

// CreateUser creates a new user
func (r *PgRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt

	query := `
		INSERT INTO users (id, username, email, password, created_at, updated_at)
		VALUES (:id, :username, :email, :password, :created_at, :updated_at)`

	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

// GetUserByEmail gets the user by email
func (r *PgRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE email = $1`

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}
