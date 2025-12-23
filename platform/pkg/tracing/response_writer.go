package tracing

import (
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

// traceResponseWriter wraps http.ResponseWriter to track response status
// and add trace ID to response headers.
type traceResponseWriter struct {
	http.ResponseWriter
	statusCode  int
	span        trace.Span
	headerAdded bool
}

// addTraceIDHeader adds trace ID to response headers if not already added.
func (w *traceResponseWriter) addTraceIDHeader() {
	if !w.headerAdded {
		traceID := w.span.SpanContext().TraceID().String()
		if traceID != "" {
			w.ResponseWriter.Header().Set(HTTPTraceIDHeader, traceID)
		}
		w.headerAdded = true
	}
}

// WriteHeader intercepts the status code and adds trace ID header.
func (w *traceResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.addTraceIDHeader()
	w.ResponseWriter.WriteHeader(code)
}

// Write intercepts body writes and ensures headers are set.
func (w *traceResponseWriter) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}

	w.addTraceIDHeader()

	return w.ResponseWriter.Write(b)
}
