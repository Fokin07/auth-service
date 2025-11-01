package token

import (
	"fmt"
	"strings"
	"time"

	"github.com/AlexFox86/auth-service/internal/models"
	"github.com/golang-jwt/jwt"
)

// GenerateToken creates a JWT token
func GenerateToken(user *models.User, jwtSecret []byte, tokenExpiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":      user.ID.String(),
		"username": user.Username,
		"exp":      time.Now().Add(tokenExpiry).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken checks the JWT token
func ValidateToken(tokenString string, jwtSecret []byte) (jwt.MapClaims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	return claims, nil
}
