package delivery

import (
	"context"
	"net/http"

	"github.com/AlexFox86/auth-service/internal/pkg/token"
)

type contextKey string

const (
	contextKeyUserID contextKey = "userID"
)

// AuthMiddleware verifies the JWT token in the Authorization header
func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "authorization header required", http.StatusUnauthorized)
			return
		}

		claims, err := token.ValidateToken(authHeader, h.service.JwtSecret())
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
