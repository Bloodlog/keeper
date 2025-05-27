package logger

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/assert"
)

func TestDefaultSettings(t *testing.T) {
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	cfg := defaultSettings(level)

	assert.Equal(t, zapcore.InfoLevel, cfg.config.Level.Level())
	assert.Equal(t, false, cfg.config.Development)
	assert.Equal(t, "json", cfg.config.Encoding)

	assert.NotNil(t, cfg.config.Sampling)
	assert.Equal(t, 100, cfg.config.Sampling.Initial)
	assert.Equal(t, 100, cfg.config.Sampling.Thereafter)

	assert.Contains(t, cfg.config.OutputPaths, "stderr")
	assert.Contains(t, cfg.config.ErrorOutputPaths, "stderr")
}
