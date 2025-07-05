package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hashed, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, password, hashed)
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword123"
	wrongPassword := "wrongpassword"

	hashed, _ := HashPassword(password)

	t.Run("correct password", func(t *testing.T) {
		err := CheckPassword(password, hashed)
		assert.NoError(t, err)
	})

	t.Run("wrong password", func(t *testing.T) {
		err := CheckPassword(wrongPassword, hashed)
		assert.Error(t, err)
	})

	t.Run("invalid hash", func(t *testing.T) {
		err := CheckPassword(password, "invalidhash")
		assert.Error(t, err)
	})
}

func TestGenerateRandomString(t *testing.T) {
	str, err := GenerateRandomString(32)
	assert.NoError(t, err)
	assert.Len(t, str, 43)
}
