package telemetry

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const CorrelationIDKey contextKey = "correlation_id"
const correlationHeader = "X-Correlation-ID"

func CorrelationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := r.Header.Get(correlationHeader)
		if correlationID == "" {
			correlationID = uuid.NewString()
		}

		ctx := context.WithValue(r.Context(), CorrelationIDKey, correlationID)
		w.Header().Set(correlationHeader, correlationID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CorrelationIDFromContext(ctx context.Context) string {
	if value, ok := ctx.Value(CorrelationIDKey).(string); ok {
		return value
	}
	return ""
}
