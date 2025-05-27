package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAgentConfig(t *testing.T) {
	cfg := NewAgentConfig()

	assert.Equal(t, "127.0.0.1", cfg.RemoteServer.Address)
	assert.Equal(t, 8081, cfg.RemoteServer.Port)

	assert.Equal(t, "", cfg.RemoteServer.CACert)
	assert.Equal(t, 5*time.Second, cfg.RemoteServer.Timeout)

	assert.Equal(t, false, cfg.RemoteServer.EnableTLS)
}
