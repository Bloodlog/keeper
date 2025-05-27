package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewServerConfig(t *testing.T) {
	cfg := NewServerConfig()

	assert.Equal(t, "127.0.0.1", cfg.Server.Address)
	assert.Equal(t, 8080, cfg.Server.Port)

	assert.Equal(t, "postgres://keeper:password@localhost:5432/keeper?sslmode=disable", cfg.Database.DSN)

	assert.Equal(t, "secret", cfg.Security.EncryptionKey)
	assert.Equal(t, 24*time.Hour, cfg.Security.TokenTTL)

	assert.Equal(t, "database", cfg.Storage.StorageType)
}
