package logger

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func init() {
	// Initialize logger with a "no-op" writer to avoid cluttering console and slowing down benchmarks
	InitForBenchmark()
}

func BenchmarkGlobalLogger(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info(ctx, "test message")
	}
}

func BenchmarkWithLogger(b *testing.B) {
	log := With(zap.String("static_field", "static_value"))
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info(ctx, "test message")
	}
}

func BenchmarkWithContextLogger(b *testing.B) {
	ctx := context.WithValue(context.Background(), traceIDKey, "trace-123")
	ctx = context.WithValue(ctx, userIDKey, "user-456")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WithContext(ctx).Info(ctx, "test message")
	}
}

func BenchmarkChainLogger(b *testing.B) {
	ctx := context.WithValue(context.Background(), traceIDKey, "trace-123")
	ctx = context.WithValue(ctx, userIDKey, "user-456")

	log := With(zap.String("static_field", "static_value"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info(ctx, "test message")
	}
}
