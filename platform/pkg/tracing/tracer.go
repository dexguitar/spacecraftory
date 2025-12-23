package tracing

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	// DefaultCompressor - default compression algorithm
	DefaultCompressor = "gzip"
	// DefaultRetryEnabled - enable retries by default
	DefaultRetryEnabled = true
	// DefaultRetryInitialInterval - initial interval between retries
	DefaultRetryInitialInterval = 500 * time.Millisecond
	// DefaultRetryMaxInterval - maximum interval between retries
	DefaultRetryMaxInterval = 5 * time.Second
	// DefaultRetryMaxElapsedTime - maximum time for all retries
	DefaultRetryMaxElapsedTime = 30 * time.Second
	// DefaultTimeout - default timeout for operations
	DefaultTimeout = 5 * time.Second
)

// serviceName - service name for tracing
var serviceName string

type Config interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
}

// InitTracer initializes the global OpenTelemetry tracer.
// The function returns an error if initialization fails.
func InitTracer(ctx context.Context, cfg Config) error {
	// Store service name for use in spans
	serviceName = cfg.ServiceName()

	// Create exporter for sending traces to OpenTelemetry Collector via gRPC
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(cfg.CollectorEndpoint()), // Collector address
		otlptracegrpc.WithInsecure(),                        // Disable TLS for local development
		otlptracegrpc.WithTimeout(DefaultTimeout),
		otlptracegrpc.WithCompressor(DefaultCompressor),
		otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
			Enabled:         DefaultRetryEnabled,
			InitialInterval: DefaultRetryInitialInterval,
			MaxInterval:     DefaultRetryMaxInterval,
			MaxElapsedTime:  DefaultRetryMaxElapsedTime,
		}),
	)
	if err != nil {
		return err
	}

	// Create resource with service metadata
	// Resource adds attributes to each trace, helping to identify the source
	attributeResource, err := resource.New(ctx,
		resource.WithAttributes(
			// Use standard OpenTelemetry attributes
			semconv.ServiceName(cfg.ServiceName()),
			semconv.ServiceVersion(cfg.ServiceVersion()),
			attribute.String("environment", cfg.Environment()),
		),
		// Automatically detect host, OS, and other system attributes
		resource.WithHost(),
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithContainer(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		return err
	}

	// Create trace provider with configured exporter and resource
	// BatchSpanProcessor collects spans in batches for efficient sending
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(attributeResource),
		// Configure trace sampling:
		// 1. ParentBased - respect parent span's sampling decision
		// 2. TraceIDRatioBased(1.0) - keep 100% of traces (1.0 = 100%)
		// In production, it's recommended to use a lower percentage (0.1 = 10%)
		// to reduce load on the tracing system
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(1.0))),
	)

	// Set global trace provider
	otel.SetTracerProvider(tracerProvider)

	// Configure context propagation for cross-service communication:
	// 1. TraceContext - W3C standard for passing trace ID and parent span ID via HTTP headers
	//    Allows linking requests between services into a single trace
	// 2. Baggage - mechanism for passing additional metadata between services
	//    For example: user_id, tenant_id, request_id and other business contexts
	// Propagation is a mechanism for passing tracing context between services
	// When a request passes through multiple services, propagation allows:
	// - Maintaining the connection between all spans in the call chain
	// - Passing additional context between services
	// - Ensuring end-to-end tracing of the entire request
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return nil
}

// ShutdownTracer closes the global OpenTelemetry tracer.
// This function should be called when the application terminates.
func ShutdownTracer(ctx context.Context) error {
	provider := otel.GetTracerProvider()
	if provider == nil {
		return nil
	}

	// Cast to concrete type to call Shutdown
	tracerProvider, ok := provider.(*sdktrace.TracerProvider)
	if !ok {
		return nil
	}

	// Close trace provider on exit
	err := tracerProvider.Shutdown(ctx)
	if err != nil {
		// Errors during shutdown are not critical, but may lead to loss of last traces
		return err
	}

	return nil
}

// StartSpan creates a new span and returns it with a new context.
// This is a convenient wrapper over trace.Tracer.Start that uses the global tracer.
//
// Difference between Tracer and TracerProvider:
// 1. TracerProvider - is a tracer factory that:
//   - Manages tracer lifecycle
//   - Configures trace export
//   - Controls sampling
//   - Stores global settings
//
// 2. Tracer - is a specific tool for creating spans:
//   - Creates spans for a specific service/component
//   - Manages relationships between spans
//   - Adds attributes to spans
//   - Tracks execution context
//
// Tracer name (serviceName):
// - Used to identify the source of spans in the tracing system
// - Allows grouping spans by service in UI (e.g., in Jaeger)
// - If a tracer with this name already exists - it is returned
// - If not - a new tracer with this name is created
// - In our case, we use the service name from configuration
//
// Span creation:
// - When creating the first (root) span, a new trace ID is generated
// - All subsequent spans in the chain inherit this trace ID
// - Trace ID links all spans of a single request/operation
// - If trace ID already exists in context, a new one is not generated
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	// Get tracer from global provider
	// Use service name from configuration for better identification in Jaeger
	return otel.Tracer(serviceName).Start(ctx, name, opts...)
}

// SpanFromContext returns the current active span from context.
// If no span exists, NoopSpan is returned.
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// TraceIDFromContext extracts trace ID from context.
// Returns a string with trace ID or empty string if trace is not found.
func TraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ""
	}

	return span.SpanContext().TraceID().String()
}
