// Package logger provides a dual-write logger with stdout and OpenTelemetry Collector output.
//
// ARCHITECTURE:
//
// The logger uses zapcore.NewTee for parallel writes to two destinations:
// 1. Stdout (for Kubernetes/container environments)
// 2. OpenTelemetry Collector (for centralized log collection)
//
// FLOW:
//
//	Application
//	    ↓ (logger.Info/Error)
//	zap.Logger
//	    ↓
//	zapcore.Tee
//	   ↙        ↘
//	StdoutCore   SimpleOTLPCore
//	   ↓             ↓
//	os.Stdout   OTLP Collector → Elasticsearch
//
// USAGE:
//
//	// Initialize logger at application startup
//	err := logger.Init("info", true)
//	// or with OTLP:
//	err := logger.InitWithOTLP("info", true, "otel-collector:4317", "my-service", "dev")
//	defer logger.Close()
//
//	// Use the logger
//	logger.Info(ctx, "message", zap.String("key", "value"))
package logger

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	otelLog "go.opentelemetry.io/otel/log"
	otelLogSdk "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Context keys for enrichment
type Key string

const (
	traceIDKey Key = "trace_id"
	userIDKey  Key = "user_id"
)

// shutdownTimeout - timeout for graceful shutdown OTLP provider
const shutdownTimeout = 2 * time.Second

// Global singleton logger
var (
	globalLogger *logger
	initOnce     sync.Once
	dynamicLevel zap.AtomicLevel
	otelProvider *otelLogSdk.LoggerProvider // OTLP provider for graceful shutdown
)

// logger is a wrapper over zap.Logger with context enrichment support
type logger struct {
	zapLogger *zap.Logger
}

// Init initializes the global logger with stdout output only.
func Init(levelStr string, asJSON bool) error {
	initOnce.Do(func() {
		dynamicLevel = zap.NewAtomicLevelAt(parseLevel(levelStr))

		core := createStdoutCore(asJSON, dynamicLevel)
		zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))

		globalLogger = &logger{
			zapLogger: zapLogger,
		}
	})

	return nil
}

// InitWithOTLP initializes the global logger with dual-write to stdout and OTLP.
// Parameters:
//   - levelStr: log level ("debug", "info", "warn", "error")
//   - asJSON: whether to format stdout output as JSON
//   - otlpEndpoint: OpenTelemetry Collector gRPC endpoint (e.g., "otel-collector:4317")
//   - serviceName: name of the service for telemetry
//   - environment: deployment environment (e.g., "dev", "staging", "production")
func InitWithOTLP(levelStr string, asJSON bool, otlpEndpoint, serviceName, environment string) error {
	initOnce.Do(func() {
		dynamicLevel = zap.NewAtomicLevelAt(parseLevel(levelStr))

		cores := buildCores(asJSON, dynamicLevel, otlpEndpoint, serviceName, environment)
		zapLogger := zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(2))

		globalLogger = &logger{
			zapLogger: zapLogger,
		}
	})

	return nil
}

// buildCores creates the slice of cores for zapcore.Tee.
// Always includes stdout core, optionally adds OTLP core.
func buildCores(asJSON bool, level zap.AtomicLevel, otlpEndpoint, serviceName, environment string) []zapcore.Core {
	cores := []zapcore.Core{
		createStdoutCore(asJSON, level),
	}

	if otlpEndpoint != "" {
		if otlpCore := createOTLPCore(otlpEndpoint, serviceName, environment, level); otlpCore != nil {
			cores = append(cores, otlpCore)
		}
	}

	return cores
}

// createStdoutCore creates a core for writing to stdout.
func createStdoutCore(asJSON bool, level zap.AtomicLevel) zapcore.Core {
	encoderCfg := buildProductionEncoderConfig()

	var encoder zapcore.Encoder
	if asJSON {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	return zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		level,
	)
}

// createOTLPCore creates a core for sending to OpenTelemetry Collector.
// Returns nil on connection error (graceful degradation).
func createOTLPCore(endpoint, serviceName, environment string, level zapcore.LevelEnabler) *SimpleOTLPCore {
	otlpLogger, err := createOTLPLogger(endpoint, serviceName, environment)
	if err != nil {
		// Graceful degradation: OTLP not available, but stdout continues working
		return nil
	}

	return NewSimpleOTLPCore(otlpLogger, level)
}

// createOTLPLogger creates an OTLP logger with configured exporter and resources.
// Uses BatchProcessor for efficient log sending.
func createOTLPLogger(endpoint, serviceName, environment string) (otelLog.Logger, error) {
	ctx := context.Background()

	exporter, err := createOTLPExporter(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	rs, err := createOTLPResource(ctx, serviceName, environment)
	if err != nil {
		return nil, err
	}

	provider := otelLogSdk.NewLoggerProvider(
		otelLogSdk.WithResource(rs),
		otelLogSdk.WithProcessor(otelLogSdk.NewBatchProcessor(exporter)),
	)
	otelProvider = provider // save for shutdown

	return provider.Logger("app"), nil
}

// createOTLPExporter creates a gRPC exporter for OTLP Collector
func createOTLPExporter(ctx context.Context, endpoint string) (*otlploggrpc.Exporter, error) {
	return otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(endpoint),
		otlploggrpc.WithInsecure(), // For local development; use TLS in production
	)
}

// createOTLPResource creates service metadata for telemetry
func createOTLPResource(ctx context.Context, serviceName, environment string) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			attribute.String("deployment.environment", environment),
		),
	)
}

func buildProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}

// SetLevel dynamically changes the log level
func SetLevel(levelStr string) {
	if dynamicLevel == (zap.AtomicLevel{}) {
		return
	}

	dynamicLevel.SetLevel(parseLevel(levelStr))
}

// InitForBenchmark initializes a no-op logger for benchmarks.
func InitForBenchmark() {
	core := zapcore.NewNopCore()

	globalLogger = &logger{
		zapLogger: zap.New(core),
	}
}

// Logger returns the global enrichment-aware logger
func Logger() *logger {
	return globalLogger
}

// SetNopLogger sets the global logger to no-op mode.
// Perfect for unit tests.
func SetNopLogger() {
	globalLogger = &logger{
		zapLogger: zap.NewNop(),
	}
}

// Sync flushes logger buffers
func Sync() error {
	if globalLogger != nil {
		return globalLogger.zapLogger.Sync()
	}

	return nil
}

// Close gracefully shuts down the logger and OTLP provider.
// Stops OTLP provider with timeout to send remaining logs.
func Close() error {
	err := Sync()
	if err != nil {
		return err
	}

	if otelProvider != nil {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		return otelProvider.Shutdown(ctx)
	}

	return nil
}

// With creates a new enrichment-aware logger with additional fields
func With(fields ...zap.Field) *logger {
	if globalLogger == nil {
		return &logger{zapLogger: zap.NewNop()}
	}

	return &logger{
		zapLogger: globalLogger.zapLogger.With(fields...),
	}
}

// WithContext creates an enrichment-aware logger with context
func WithContext(ctx context.Context) *logger {
	if globalLogger == nil {
		return &logger{zapLogger: zap.NewNop()}
	}

	return &logger{
		zapLogger: globalLogger.zapLogger.With(fieldsFromContext(ctx)...),
	}
}

// Debug enrichment-aware debug log
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Debug(ctx, msg, fields...)
}

// Info enrichment-aware info log
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Info(ctx, msg, fields...)
}

// Warn enrichment-aware warn log
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Warn(ctx, msg, fields...)
}

// Error enrichment-aware error log
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Error(ctx, msg, fields...)
}

// Fatal enrichment-aware fatal log
func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Fatal(ctx, msg, fields...)
}

// Instance methods for enrichment loggers (logger)

func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Debug(msg, allFields...)
}

func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Info(msg, allFields...)
}

func (l *logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Warn(msg, allFields...)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Error(msg, allFields...)
}

func (l *logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Fatal(msg, allFields...)
}

// parseLevel converts string level to zapcore.Level
func parseLevel(levelStr string) zapcore.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// fieldsFromContext extracts enrichment fields from context
func fieldsFromContext(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)

	if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
		fields = append(fields, zap.String(string(traceIDKey), traceID))
	}

	if userID, ok := ctx.Value(userIDKey).(string); ok && userID != "" {
		fields = append(fields, zap.String(string(userIDKey), userID))
	}

	return fields
}
