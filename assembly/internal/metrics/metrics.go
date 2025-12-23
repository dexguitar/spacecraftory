package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const serviceName = "assembly-service"

// meter - factory for creating metric instruments
// Uses global MeterProvider initialized by platform/pkg/metrics
var meter = otel.Meter(serviceName)

var (
	// AssemblyDuration - HISTOGRAM for measuring assembly time
	// Type: Float64Histogram (distribution of values)
	// Usage: tracking assembly operation performance
	// Automatically creates metrics: _count, _sum, _bucket for percentiles
	AssemblyDuration metric.Float64Histogram
)

// InitMetrics initializes all assembly service metrics
// Should be called once at application startup after OpenTelemetry provider initialization
func InitMetrics() error {
	var err error

	// Create histogram for assembly duration with appropriate buckets
	// Buckets are optimized for assembly operations (can take seconds to minutes)
	AssemblyDuration, err = meter.Float64Histogram(
		"assembly_duration_seconds",
		metric.WithDescription("Duration of spacecraft assembly operations"),
		metric.WithUnit("s"),
		// Bucket boundaries for assembly operations
		// 1s, 2s, 5s, 10s, 15s, 20s, 30s, 45s, 60s, 120s
		metric.WithExplicitBucketBoundaries(
			1.0, 2.0, 5.0, 10.0, 15.0, 20.0, 30.0, 45.0, 60.0, 120.0,
		),
	)
	if err != nil {
		return err
	}

	return nil
}

