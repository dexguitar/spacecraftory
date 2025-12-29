package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const serviceName = "order-service"

// meter - factory for creating metric instruments
// Uses global MeterProvider initialized by platform/pkg/metrics
var meter = otel.Meter(serviceName)

var (
	// OrdersTotal - COUNTER for counting total number of created orders
	// Type: Int64Counter (monotonically increasing)
	// Usage: business metric for tracking order creation rate
	OrdersTotal metric.Int64Counter

	// OrdersRevenueTotal - COUNTER for tracking total revenue from orders
	// Type: Float64Counter (monotonically increasing)
	// Usage: business metric for tracking cumulative revenue
	OrdersRevenueTotal metric.Float64Counter
)

// InitMetrics initializes all order service metrics
// Should be called once at application startup after OpenTelemetry provider initialization
func InitMetrics() error {
	var err error

	// Create counter for total orders
	OrdersTotal, err = meter.Int64Counter(
		"orders_total",
		metric.WithDescription("Total number of orders created"),
	)
	if err != nil {
		return err
	}

	// Create counter for total revenue
	OrdersRevenueTotal, err = meter.Float64Counter(
		"orders_revenue_total",
		metric.WithDescription("Total revenue from orders"),
		metric.WithUnit("USD"),
	)
	if err != nil {
		return err
	}

	return nil
}
