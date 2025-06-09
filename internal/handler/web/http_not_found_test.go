package web

import (
	"keeper/internal/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/assert"
)

func TestNotFoundHandler(t *testing.T) {
	zl, _ := logger.NewZapLogger(zapcore.InfoLevel)
	handler := NewStaticPageHandler(zl)

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", http.NoBody)
	rec := httptest.NewRecorder()

	h := handler.NotFoundHandler(t.Context())
	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
