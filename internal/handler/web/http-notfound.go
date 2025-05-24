package web

import (
	"context"
	"keeper/internal/logger"
	"net/http"

	"go.uber.org/zap"
)

type StaticPage interface {
	NotFoundHandler(ctx context.Context) http.HandlerFunc
}

type StaticPageHandler struct {
	log *logger.ZapLogger
}

func NewStaticPageHandler(log *logger.ZapLogger) *StaticPageHandler {
	return &StaticPageHandler{log: log}
}

func (h *StaticPageHandler) NotFoundHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.log.InfoCtx(ctx, "Route not found",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
		)
		w.WriteHeader(http.StatusNotFound)
	}
}
