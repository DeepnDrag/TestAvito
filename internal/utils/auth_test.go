package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken("postgres", "password")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ValidateJWT(token, "password")
	assert.NoError(t, err)
	assert.Equal(t, "postgres", claims.UserName)
}

func TestValidateJWT_ValidToken(t *testing.T) {
	token, _ := GenerateToken("postgres", "password")
	claims, err := ValidateJWT(token, "password")
	assert.NoError(t, err)
	assert.Equal(t, "postgres", claims.UserName)
	assert.WithinDuration(t, time.Now(), claims.IssuedAt.Time, 10*time.Second)
}

func TestValidateJWT_InvalidSecret(t *testing.T) {
	token, _ := GenerateToken("postgres", "password")
	_, err := ValidateJWT(token, "wrongSecret")
	assert.Error(t, err)
}

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("password")
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	isValid := CheckPassword("password", hash)
	assert.True(t, isValid)
}

func TestCheckPassword_Valid(t *testing.T) {
	hash, _ := HashPassword("password")
	result := CheckPassword("password", hash)
	assert.True(t, result)
}

func TestCheckPassword_Invalid(t *testing.T) {
	hash, _ := HashPassword("password")
	result := CheckPassword("wrongPassword", hash)
	assert.False(t, result)
}
