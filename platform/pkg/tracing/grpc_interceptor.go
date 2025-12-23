package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Header constants for trace ID propagation
const (
	// TraceIDHeader is the gRPC metadata key for trace ID
	TraceIDHeader = "x-trace-id"
)

// UnaryServerInterceptor creates a gRPC unary server interceptor for tracing incoming requests.
// The interceptor extracts trace context from incoming metadata and creates a new span for each request.
//
// If the incoming request contains trace context (from a calling service), the new span will be
// a child of that trace. Otherwise, a new root span will be created.
func UnaryServerInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	tracer := otel.GetTracerProvider().Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Extract metadata from incoming context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		// Extract trace context from metadata
		ctx = propagator.Extract(ctx, metadataCarrier(md))

		// Create new span for this request
		ctx, span := tracer.Start(
			ctx,
			info.FullMethod,
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		// Add trace ID to response headers for client correlation
		ctx = AddTraceIDToResponse(ctx)

		// Call the actual handler
		resp, err := handler(ctx, req)
		if err != nil {
			span.RecordError(err)
		}

		return resp, err
	}
}

// UnaryClientInterceptor creates a gRPC unary client interceptor for tracing outgoing requests.
// The interceptor injects trace context into outgoing metadata to propagate the trace to called services.
func UnaryClientInterceptor(serviceName string) grpc.UnaryClientInterceptor {
	tracer := otel.GetTracerProvider().Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// Determine span name based on existing trace context
		spanName := formatSpanName(ctx, method)

		// Create span for client call
		ctx, span := tracer.Start(
			ctx,
			spanName,
			trace.WithSpanKind(trace.SpanKindClient),
		)
		defer span.End()

		// Create carrier for trace context propagation
		carrier := metadataCarrier(extractOutgoingMetadata(ctx))

		// Inject trace context into metadata
		propagator.Inject(ctx, carrier)

		// Update context with injected metadata
		ctx = metadata.NewOutgoingContext(ctx, metadata.MD(carrier))

		// Call the actual service
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			span.RecordError(err)
		}

		return err
	}
}

// formatSpanName formats the span name based on whether a trace context exists.
func formatSpanName(ctx context.Context, method string) string {
	if !trace.SpanContextFromContext(ctx).IsValid() {
		return "client." + method
	}

	return method
}

// extractOutgoingMetadata extracts outgoing metadata from context and creates a copy.
func extractOutgoingMetadata(ctx context.Context) metadata.MD {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return metadata.New(nil)
	}

	return md.Copy()
}

// GetTraceIDFromContext extracts trace ID from context.
// Useful for logging and returning trace ID to clients.
func GetTraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ""
	}

	return span.SpanContext().TraceID().String()
}

// AddTraceIDToResponse adds trace ID to outgoing gRPC response metadata.
// This allows clients to retrieve the trace ID for debugging.
func AddTraceIDToResponse(ctx context.Context) context.Context {
	traceID := GetTraceIDFromContext(ctx)
	if traceID == "" {
		return ctx
	}

	md := extractOutgoingMetadata(ctx)
	md.Set(TraceIDHeader, traceID)

	return metadata.NewOutgoingContext(ctx, md)
}
