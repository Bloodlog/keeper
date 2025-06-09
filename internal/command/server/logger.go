package server

import (
	"context"
	"keeper/internal/logger"
	"log"

	"go.uber.org/zap"
)

// Init logger.
func initLogger(ctx context.Context) (*logger.ZapLogger, error) {
	l, err := logger.NewZapLogger(zap.InfoLevel)

	l.InfoCtx(
		ctx,
		"logging started...",
		zap.String("app", "logging"),
		zap.String("service", "main"),
	)

	if err != nil {
		log.Panic(err)
	}

	_ = l.WithContextFields(ctx,
		zap.String("app", "logging"),
		zap.String("service", "main"))

	defer l.Sync()

	return l, nil
}
