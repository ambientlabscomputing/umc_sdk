package middleware

import (
	"log/slog"
	"net/http"

	"github.com/ambientlabscomputing/mycelium_spine/sdk"
	"github.com/ambientlabscomputing/umc_sdk/logging"
	"github.com/google/uuid"
)

// TraceIDMiddleware wraps an http.Handler to extract or generate a trace ID,
// inject it into the request context, and return it in the response header.
func TraceIDMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for existing trace ID in header
		traceID := r.Header.Get("X-Trace-ID")

		// If not found, generate a new one
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// Inject trace ID into context
		ctx := sdk.WithTraceID(r.Context(), traceID)

		// Add trace ID to structured logger
		logger := slog.Default().With("trace_id", traceID)
		ctx = logging.WithLogger(ctx, logger)

		// Write trace ID to response header
		w.Header().Set("X-Trace-ID", traceID)

		// Call the wrapped handler with enriched context
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}