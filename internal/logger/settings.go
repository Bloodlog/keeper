package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type settings struct {
	config *zap.Config
	opts   []zap.Option
}

const (
	samplingInitial    = 100
	samplingThereafter = 100
)

func defaultSettings(level zap.AtomicLevel) *settings {
	config := &zap.Config{
		Level:       level,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    samplingInitial,
			Thereafter: samplingThereafter,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "@timestamp",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return &settings{
		config: config,
		opts: []zap.Option{
			zap.AddCallerSkip(1),
		},
	}
}
