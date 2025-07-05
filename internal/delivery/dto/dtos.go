package dto

import "github.com/AlexFox86/auth-service/internal/models"

// RegisterRequest registration request
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginRequest login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Response response with a token
type Response struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}
