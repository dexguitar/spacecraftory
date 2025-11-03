package logger

import (
	"context"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Key string

const (
	traceIDKey Key = "trace_id"
	userIDKey  Key = "user_id"
)

// Global singleton logger
var (
	globalLogger *logger
	initOnce     sync.Once
	dynamicLevel zap.AtomicLevel
)

// logger is a wrapper over zap.Logger with context enrichment support
type logger struct {
	zapLogger *zap.Logger
}

// Init initializes the global logger.
func Init(levelStr string, asJSON bool) error {
	initOnce.Do(func() {
		dynamicLevel = zap.NewAtomicLevelAt(parseLevel(levelStr))

		encoderCfg := buildProductionEncoderConfig()

		var encoder zapcore.Encoder
		if asJSON {
			encoder = zapcore.NewJSONEncoder(encoderCfg)
		} else {
			encoder = zapcore.NewConsoleEncoder(encoderCfg)
		}

		core := zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			dynamicLevel,
		)

		zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))

		globalLogger = &logger{
			zapLogger: zapLogger,
		}
	})

	return nil
}

func buildProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",                 // timestamp
		LevelKey:       "level",                     // log level
		NameKey:        "logger",                    // logger name, if used
		CallerKey:      "caller",                    // where the log was called from
		MessageKey:     "message",                   // message text
		StacktraceKey:  "stacktrace",                // stacktrace for errors
		LineEnding:     zapcore.DefaultLineEnding,   // line break
		EncodeLevel:    zapcore.CapitalLevelEncoder, // INFO, ERROR
		EncodeTime:     zapcore.ISO8601TimeEncoder,  // readable ISO 8601 format
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // short caller
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
