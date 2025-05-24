package logger

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

type ctxKeyZapFields struct{}

var zapFieldsKey = ctxKeyZapFields{}

type ZapFields map[string]zap.Field

// NewZapLogger returns a new ZapLogger configured with the provided options.
func NewZapLogger(level zapcore.Level) (*ZapLogger, error) {
	atomic := zap.NewAtomicLevelAt(level)
	settings := defaultSettings(atomic)

	l, err := settings.config.Build(settings.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return &ZapLogger{
		logger: l,
		level:  atomic,
	}, nil
}

func (z *ZapLogger) Sync() {
	_ = z.logger.Sync()
}

func (zf ZapFields) Append(fields ...zap.Field) ZapFields {
	zfCopy := make(ZapFields)
	for k, v := range zf {
		zfCopy[k] = v
	}

	for _, f := range fields {
		zfCopy[f.Key] = f
	}

	return zfCopy
}

type ZapLogger struct {
	logger *zap.Logger
	level  zap.AtomicLevel
}

func (z *ZapLogger) maskField(f zap.Field) zap.Field {
	const emailPartsCount = 2
	if f.Key == "password" {
		return zap.String(f.Key, "******")
	}

	if f.Key == "email" {
		email := f.String
		parts := strings.Split(email, "@")
		if len(parts) == emailPartsCount {
			return zap.String(f.Key, "***@"+parts[1])
		}
	}
	return f
}

func (z *ZapLogger) WithContextFields(ctx context.Context, fields ...zap.Field) context.Context {
	ctxFields, _ := ctx.Value(zapFieldsKey).(ZapFields)
	if ctxFields == nil {
		ctxFields = make(ZapFields)
	}

	merged := ctxFields.Append(fields...)
	return context.WithValue(ctx, zapFieldsKey, merged)
}

func (z *ZapLogger) withCtxFields(ctx context.Context, fields ...zap.Field) []zap.Field {
	fs := make(ZapFields)

	ctxFields, ok := ctx.Value(zapFieldsKey).(ZapFields)
	if ok {
		fs = ctxFields
	}

	fs = fs.Append(fields...)

	maskedFields := make([]zap.Field, 0, len(fs))
	for _, f := range fs {
		maskedFields = append(maskedFields, z.maskField(f))
	}

	return maskedFields
}

func (z *ZapLogger) InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	z.logger.Info(msg, z.withCtxFields(ctx, fields...)...)
}
