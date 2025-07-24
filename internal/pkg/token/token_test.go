package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/AlexFox86/auth-service/internal/models"
)

func TestGenerateToken(t *testing.T) {
	t.Run("generate token", func(t *testing.T) {
		user := &models.User{
			ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Username: "testuser",
		}
		_, err := GenerateToken(user, []byte("secret"), time.Hour)
		assert.NoError(t, err)
	})
}

func TestValidateToken(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		user := &models.User{
			ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Username: "testuser",
		}
		tokenString, err := GenerateToken(user, []byte("secret"), time.Hour)
		assert.NoError(t, err)

		claims, err := ValidateToken(tokenString, []byte("secret"))
		assert.NoError(t, err)
		assert.Equal(t, user.ID.String(), claims["sub"])
		assert.Equal(t, user.Username, claims["username"])
	})

	t.Run("invalid token", func(t *testing.T) {
		claims, err := ValidateToken("invalid.token.string", []byte("secret"))
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("wrong signing method", func(t *testing.T) {
		// Creating a token with an incorrect signature method
		testToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"sub": "123",
		})
		tokenString, _ := testToken.SignedString([]byte("key"))

		claims, err := ValidateToken(tokenString, []byte("secret"))
		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}
