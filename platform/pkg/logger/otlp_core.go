// OTLP Core Component
//
// What happens here:
// - record: a single log entry (level, message, time, field-attributes).
// - core: the "heart" of the logger. It decides "should I accept this entry" and "how to send it".
// - tee: a "splitter" that distributes one entry to multiple cores simultaneously.
//
// zapcore.Core interface (what any core must be able to do):
// - Enabled(level): decide whether to write an entry of this level.
// - With(fields): return a copy of core with additional fields (we account for them in Write).
// - Check(entry, ce): add itself to the list of entry recipients if the level matches.
// - Write(entry, fields): assemble the record and send it "where needed".
// - Sync(): flush buffers if any.
//
// OTLP flow architecture:
// zap.Logger -> zapcore.Tee -> SimpleOTLPCore -> OTLP Collector (gRPC) -> Elasticsearch
package logger

import (
	"context"
	"time"

	otelLog "go.opentelemetry.io/otel/log"
	"go.uber.org/zap/zapcore"
)

// emitTimeout - timeout for sending a single record, to avoid blocking the application
const emitTimeout = 500 * time.Millisecond

// SimpleOTLPCore converts zap entries to OpenTelemetry Records and sends them directly to OTLP
type SimpleOTLPCore struct {
	otlpLogger otelLog.Logger       // OTLP logger for sending records
	level      zapcore.LevelEnabler // minimum level for logging
}

// NewSimpleOTLPCore creates a new OTLP core that works directly with OTLP logger.
func NewSimpleOTLPCore(otlpLogger otelLog.Logger, level zapcore.LevelEnabler) *SimpleOTLPCore {
	return &SimpleOTLPCore{
		otlpLogger: otlpLogger,
		level:      level,
	}
}

// Enabled checks if a log of the given level should be recorded
func (c *SimpleOTLPCore) Enabled(level zapcore.Level) bool {
	return c.level.Enabled(level)
}

// With creates a new core with additional fields.
// In the current implementation, fields are processed in the Write method,
// so here a copy is created without changes.
func (c *SimpleOTLPCore) With(_ []zapcore.Field) zapcore.Core {
	return &SimpleOTLPCore{
		otlpLogger: c.otlpLogger,
		level:      c.level,
	}
}

// Check determines if the given log should be written by this core.
// Adds itself to CheckedEntry if the log level matches the settings.
func (c *SimpleOTLPCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}

	return ce
}

// Write converts zap Entry to OpenTelemetry Record and sends it to OTLP.
// Step by step:
//  1. Convert zap level to OTLP Severity (mapZapToOtelSeverity).
//  2. Build base Record: severity, body=message, timestamp (makeBaseRecord).
//  3. Encode zap fields into OTLP attributes (encodeFieldsToAttrs) and add them to Record.
//  4. Send the record via OTLP logger with a short timeout (emitWithTimeout),
//     to avoid blocking the application during network issues.
func (c *SimpleOTLPCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	severity := mapZapToOtelSeverity(entry.Level)
	record := makeBaseRecord(entry, severity)

	if len(fields) > 0 {
		attrs := encodeFieldsToAttrs(fields)
		if len(attrs) > 0 {
			record.AddAttributes(attrs...)
		}
	}

	c.emitWithTimeout(record)

	return nil
}

// Sync - synchronization not required: batching is done by OTLP SDK
func (c *SimpleOTLPCore) Sync() error {
	return nil
}

// mapZapToOtelSeverity - separate function for level conversion
func mapZapToOtelSeverity(level zapcore.Level) otelLog.Severity {
	switch level {
	case zapcore.DebugLevel:
		return otelLog.SeverityDebug
	case zapcore.InfoLevel:
		return otelLog.SeverityInfo
	case zapcore.WarnLevel:
		return otelLog.SeverityWarn
	case zapcore.ErrorLevel:
		return otelLog.SeverityError
	default:
		return otelLog.SeverityInfo
	}
}

// makeBaseRecord - builds a base record without attributes
func makeBaseRecord(entry zapcore.Entry, sev otelLog.Severity) otelLog.Record {
	r := otelLog.Record{}
	r.SetSeverity(sev)
	r.SetBody(otelLog.StringValue(entry.Message))
	r.SetTimestamp(entry.Time)

	return r
}

// encodeFieldsToAttrs - prepares attributes from zap fields.
// We use zapcore.NewMapObjectEncoder() to safely unfold []zapcore.Field
// into a key→value map. Then we transfer only basic types to OTLP KeyValue.
// Unsupported types are skipped (they will continue to live in stdout part via zap encoder).
func encodeFieldsToAttrs(fields []zapcore.Field) []otelLog.KeyValue {
	if len(fields) == 0 {
		return nil
	}

	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fields {
		f.AddTo(enc)
	}

	attrs := make([]otelLog.KeyValue, 0, len(enc.Fields))
	for k, v := range enc.Fields {
		switch val := v.(type) {
		case string:
			attrs = append(attrs, otelLog.String(k, val))
		case bool:
			attrs = append(attrs, otelLog.Bool(k, val))
		case int64:
			attrs = append(attrs, otelLog.Int64(k, val))
		case float64:
			attrs = append(attrs, otelLog.Float64(k, val))
		}
	}

	return attrs
}

// emitWithTimeout - sends to OTLP with a short timeout
func (c *SimpleOTLPCore) emitWithTimeout(record otelLog.Record) {
	if c.otlpLogger == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), emitTimeout)
	defer cancel()

	c.otlpLogger.Emit(ctx, record)
}
