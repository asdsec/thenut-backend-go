package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("../")

	// todo: implement test env

	assert.NoError(t, err)
	assert.Equal(t, "postgres", config.DBDriver)
	assert.Equal(t, "postgresql://root:secret@localhost:5432/thenut?sslmode=disable", config.DBSource)
	assert.Equal(t, "0.0.0.0:8080", config.ServerAddress)
	assert.Equal(t, "12345678901234567890123456789012", config.TokenSymmetricKey)
	assert.Equal(t, 15*time.Minute, config.AccessTokenDuration)
}
