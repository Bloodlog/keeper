package web

import (
	"context"
	"keeper/internal/config"
	"keeper/internal/logger"
	"net/http"

	"go.uber.org/zap"
)

type DownloadFilePage interface {
	FileServerHandler(ctx context.Context) http.HandlerFunc
}

type DownloadFileServerHandler struct {
	log *logger.ZapLogger
	fs  http.Handler
}

func NewFileServerHandler(log *logger.ZapLogger, cfg *config.MainServerConfig) *DownloadFileServerHandler {
	return &DownloadFileServerHandler{
		log: log,
		fs: http.StripPrefix(
			cfg.BuildAgentsConfig.URLPrefix,
			http.FileServer(http.Dir(cfg.BuildAgentsConfig.DownloadDir)),
		),
	}
}

func (h *DownloadFileServerHandler) FileServerHandler(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-ctx.Done():
			http.Error(w, "server is shutting down", http.StatusServiceUnavailable)
			return
		default:
			h.log.InfoCtx(ctx, "serving download", zap.String("file", r.URL.Path))
			h.fs.ServeHTTP(w, r)
		}
	})
}
