// =============================================================================
// PLATFORM METRICS - OPENTELEMETRY METRICS PROVIDER INITIALIZATION
// =============================================================================
//
// This package implements metrics collection initialization based on OpenTelemetry.
// It is a platform module that configures the infrastructure for sending metrics
// from the application to the OpenTelemetry Collector.
//
// METRICS COLLECTION ARCHITECTURE:
//
// # Application → MeterProvider → Reader → Exporter → OTLP Collector → Prometheus
//
// Components:
// 1. MeterProvider - central object for managing metrics
// 2. Reader - reads metrics from the application and aggregates them
// 3. Exporter - sends metrics to external systems
// 4. OTLP Collector - receives metrics and forwards them to Prometheus
//
// METRICS DELIVERY MODEL:
//
// Push Model (used here):
// - Application actively sends metrics to the collector
// - Periodic delivery (e.g., every 10 seconds)
// - Does not require opening ports on the application
//
// Pull Model (alternative):
// - Prometheus scraping endpoints (/metrics)
// - Collector polls the application
package metrics

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

const (
	// defaultTimeout - timeout for sending a batch of metrics to the collector
	// Important: should not be too large to avoid blocking the application
	// Usually 5-10 seconds is enough for most cases
	defaultTimeout = 5 * time.Second

	// defaultInterval - default interval for sending metrics
	// Trade-off between metric freshness and network load
	// Usually: 10-60 seconds depending on requirements
	defaultInterval = 10 * time.Second
)

// =============================================================================
// GLOBAL VARIABLES - SINGLETON PATTERN
// =============================================================================
//
// Using global variables for MeterProvider and Exporter because:
// 1. OpenTelemetry recommends one MeterProvider per application
// 2. All parts of the application should use the same exporter
// 3. Simplifies dependency injection in large applications
// 4. Allows proper resource cleanup during shutdown
//
// Alternatives:
// - Dependency Injection Container (wire, dig)
// - Context-based dependency passing
// - Service Locator pattern
var (
	// exporter - component for sending metrics via OTLP (OpenTelemetry Protocol)
	// OTLP - standard protocol for telemetry transfer between components
	// Supports HTTP and gRPC transports (using gRPC for better performance)
	exporter *otlpmetricgrpc.Exporter

	// meterProvider - central OpenTelemetry object for managing metrics
	// Responsible for:
	// - Creating Meters for different application components
	// - Aggregating metrics (summing, calculating percentiles)
	// - Periodic metric sending via Reader
	// - Managing metric lifecycle
	meterProvider *metric.MeterProvider
)

// =============================================================================
// CONFIGURATION INTERFACE - DEPENDENCY INVERSION PRINCIPLE
// =============================================================================
//
// Config defines an abstraction for metrics configuration.
// This applies the dependency inversion principle:
// - Platform module does not depend on specific configuration implementation
// - Allows using different configuration sources (files, environment variables)
// - Simplifies testing (can create mock configuration)

// Config defines the configuration interface for metrics initialization
type Config interface {
	// CollectorEndpoint returns the OTLP collector address
	// Usually: "localhost:4317" (gRPC) or "localhost:4318" (HTTP)
	CollectorEndpoint() string

	// CollectorInterval returns the metrics sending interval
	// Trade-off between metric freshness and network load
	// Usually: 10-60 seconds depending on requirements
	CollectorInterval() time.Duration

	// ServiceName returns the name of the service for resource identification
	ServiceName() string

	// Environment returns the deployment environment (e.g., "dev", "staging", "production")
	Environment() string
}

// =============================================================================
// METRICS INITIALIZATION - BOOTSTRAP PROCESS
// =============================================================================

// InitProvider initializes the global OpenTelemetry metrics provider
//
// This method performs the following steps:
// 1. Creates an OTLP exporter for sending metrics to the collector
// 2. Configures MeterProvider with periodic metric reading
// 3. Sets the global metrics provider for the entire application
//
// IMPORTANT: Must be called once at application startup, before creating metrics
//
// Parameters:
// - ctx: context for cancellation and timeouts
// - cfg: configuration with collector address and sending interval
//
// Returns an error if connection to the collector fails
func InitProvider(ctx context.Context, cfg Config) error {
	var err error

	// =========================================================================
	// STEP 1: CREATE OTLP EXPORTER
	// =========================================================================
	//
	// OTLP (OpenTelemetry Protocol) - standard protocol for telemetry
	// OTLP advantages:
	// - Efficient binary serialization (protobuf)
	// - Batching support to reduce network load
	// - Standardized format for compatibility
	// - Built-in compression and retry logic support
	//
	exporter, err = otlpmetricgrpc.New(
		ctx,
		// Collector address (usually localhost:4317 for gRPC)
		otlpmetricgrpc.WithEndpoint(cfg.CollectorEndpoint()),

		// Disable TLS for local development
		// In production use TLS and authentication
		otlpmetricgrpc.WithInsecure(),

		// Timeout for sending each batch of metrics
		// Prevents application blocking during network issues
		otlpmetricgrpc.WithTimeout(defaultTimeout),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create metrics exporter")
	}

	// =========================================================================
	// STEP 2: DETERMINE EXPORT INTERVAL
	// =========================================================================
	interval := cfg.CollectorInterval()
	if interval == 0 {
		interval = defaultInterval
	}

	// =========================================================================
	// STEP 3: CREATE RESOURCE WITH SERVICE ATTRIBUTES
	// =========================================================================
	//
	// Resource describes the entity producing telemetry
	// Standard attributes:
	// - service.name - name of the service
	// - deployment.environment - deployment environment (dev/staging/prod)
	//
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName()),
			semconv.DeploymentEnvironment(cfg.Environment()),
		),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create resource")
	}

	// =========================================================================
	// STEP 4: CREATE METER PROVIDER
	// =========================================================================
	//
	// MeterProvider - central object for managing metrics
	// Uses Builder pattern for configuration:
	//
	// Reader is responsible for:
	// - Periodic reading of metrics from the application
	// - Data aggregation (summing counters, calculating percentiles)
	// - Sending metric batches via exporter
	// - Managing metric state between sends
	//
	meterProvider = metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(
			// PeriodicReader automatically reads and sends metrics
			metric.NewPeriodicReader(
				exporter, // Use the created OTLP exporter
				// Interval for reading and sending metrics
				// Frequent sends = more current metrics, but more load
				metric.WithInterval(interval),
			),
		),
	)

	// =========================================================================
	// STEP 5: SET GLOBAL PROVIDER
	// =========================================================================
	//
	// Set the created provider as global
	// After this, all otel.Meter() calls will use our provider
	//
	// Global provider vs local:
	// + Easy to use in any part of the application
	// + No need to pass dependencies via dependency injection
	// - Complicates testing (need to reset global state)
	// - May create issues with modular testing
	//
	otel.SetMeterProvider(meterProvider)

	return nil
}

// =============================================================================
// GETTER METHODS - ACCESS TO INTERNAL COMPONENTS
// =============================================================================

// GetMeterProvider returns the current metrics provider
//
// Used for:
// - Creating additional Readers (e.g., for testing)
// - Getting metric statistics
// - Low-level access to provider functionality
//
// In most cases, it's better to use otel.Meter() to create metrics
func GetMeterProvider() *metric.MeterProvider {
	return meterProvider
}

// =============================================================================
// GRACEFUL SHUTDOWN - PROPER RESOURCE CLEANUP
// =============================================================================

// Shutdown properly closes the metrics provider and exporter
//
// GRACEFUL SHUTDOWN IMPORTANCE:
// - Sends remaining metrics before closing
// - Releases network connections and goroutines
// - Prevents data loss during application restart
// - Honors the collector contract (proper gRPC connection closure)
//
// # Should be called in defer or signal handler when application terminates
//
// Parameters:
// - ctx: context with timeout for graceful shutdown (usually 5-10 seconds)
//
// Returns an error if resources cannot be properly closed
func Shutdown(ctx context.Context) error {
	// Check if there's anything to close
	if meterProvider == nil && exporter == nil {
		return nil
	}

	var errs []error

	// First close the metrics provider
	// This will stop sending new metrics and send remaining ones
	if meterProvider != nil {
		if err := meterProvider.Shutdown(ctx); err != nil {
			errs = append(errs, errors.Wrap(err, "failed to shutdown meter provider"))
		}
	}

	// Then close the exporter
	// This will close network connections and release resources
	if exporter != nil {
		if err := exporter.Shutdown(ctx); err != nil {
			errs = append(errs, errors.Wrap(err, "failed to shutdown exporter"))
		}
	}

	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}
