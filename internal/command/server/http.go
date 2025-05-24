package server

import (
	"context"
	"errors"
	"fmt"
	"keeper/internal/config"
	"keeper/internal/logger"
	"log"
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// Init http server.
func initHTTPServer(
	ctx context.Context,
	g *errgroup.Group,
	cfg *config.MainServerConfig,
	router http.Handler,
	l *logger.ZapLogger,
) {
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.Port),
		Handler: router,
	}
	g.Go(func() (err error) {
		l.InfoCtx(ctx, "Starting HTTP server", zap.String("addr", httpServer.Addr), zap.Int("port", cfg.Server.Port))

		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("listen and serve failed: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		defer log.Print("server has been shutdown")
		<-ctx.Done()

		shutdownTimeoutCtx, cancelShutdownTimeoutCtx := context.WithTimeout(context.Background(), timeoutServerShutdown)
		defer cancelShutdownTimeoutCtx()

		if httpServer != nil {
			if err := httpServer.Shutdown(shutdownTimeoutCtx); err != nil {
				l.InfoCtx(ctx, "HTTP server Shutdown: %v", zap.Error(err))
			}
		}

		return nil
	})
}
