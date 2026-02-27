package middleware
package middleware

import (
	"context"
	"log/slog"
	"net/http"






























}	})		handler.ServeHTTP(w, r.WithContext(ctx))		// Call the wrapped handler with enriched context		w.Header().Set("X-Trace-ID", traceID)		// Write trace ID to response header		ctx = context.WithValue(ctx, slog.KindKey, logger)		logger := slog.Default().With("trace_id", traceID)		// Add trace ID to structured logger		ctx := sdk.WithTraceID(r.Context(), traceID)		// Inject trace ID into context		}			traceID = uuid.New().String()		if traceID == "" {		// If not found, generate a new one		traceID := r.Header.Get("X-Trace-ID")		// Check for existing trace ID in header	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {func TraceIDMiddleware(handler http.Handler) http.Handler {// inject it into the request context, and return it in the response header.// TraceIDMiddleware wraps an http.Handler to extract or generate a trace ID,)	"github.com/google/uuid"	"github.com/ambientlabscomputing/mycelium_spine/sdk"