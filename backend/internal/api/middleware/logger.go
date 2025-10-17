package middleware

import (
	"log/slog"
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/leak-streaming/leak-streaming/backend/internal/platform/telemetry"
)

func RequestLogger(log *slog.Logger) func(http.Handler) http.Handler {
	logger := log.With("component", "http")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &responseRecorder{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(rec, r)

			logEntry := logger.With(
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration_ms", time.Since(start).Milliseconds(),
				"request_id", chimiddleware.GetReqID(r.Context()),
			)

			if correlationID := telemetry.CorrelationIDFromContext(r.Context()); correlationID != "" {
				logEntry = logEntry.With("correlation_id", correlationID)
			}

			logEntry.Info("request completed")
		})
	}
}

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (r *responseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
