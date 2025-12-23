package tracing

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

// HTTP header constants for trace ID propagation
const (
	// HTTPTraceIDHeader is the HTTP header key for trace ID
	HTTPTraceIDHeader = "X-Trace-ID"
)

// HTTPHandlerMiddleware creates an HTTP middleware for tracing incoming requests.
// The middleware extracts trace context from incoming headers and creates a new span for each request.
// It also adds the trace ID to response headers for client correlation.
func HTTPHandlerMiddleware(serviceName string) func(http.Handler) http.Handler {
	tracer := otel.GetTracerProvider().Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract trace context from incoming request headers
			ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

			// Create span name from HTTP method and path
			spanName := fmt.Sprintf("%s %s", r.Method, r.URL.Path)

			// Create span with HTTP attributes
			ctx, span := tracer.Start(
				ctx,
				spanName,
				createHTTPSpanAttributes(r)...,
			)
			defer span.End()

			// Wrap response writer to capture status code and add trace ID
			wrw := &traceResponseWriter{
				ResponseWriter: w,
				span:           span,
				headerAdded:    false,
			}

			// Call the actual handler with enriched context
			next.ServeHTTP(wrw, r.WithContext(ctx))

			// Add response status to span attributes
			span.SetAttributes(semconv.HTTPResponseStatusCode(wrw.statusCode))
		})
	}
}

// createHTTPSpanAttributes creates standard HTTP span attributes.
func createHTTPSpanAttributes(r *http.Request) []trace.SpanStartOption {
	return []trace.SpanStartOption{
		trace.WithAttributes(
			semconv.HTTPRequestMethodKey.String(r.Method),
			semconv.URLPath(r.URL.Path),
			semconv.ServerAddress(r.Host),
			semconv.UserAgentOriginal(r.UserAgent()),
		),
		trace.WithSpanKind(trace.SpanKindServer),
	}
}
