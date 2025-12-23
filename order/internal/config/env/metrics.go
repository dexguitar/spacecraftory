package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type metricsEnvConfig struct {
	CollectorEndpoint string        `env:"ORDER_OTEL_COLLECTOR_ENDPOINT" envDefault:""`
	CollectorInterval time.Duration `env:"ORDER_METRICS_COLLECTOR_INTERVAL" envDefault:"10s"`
	ServiceName       string        `env:"ORDER_SERVICE_NAME" envDefault:"order-service"`
	Environment       string        `env:"ORDER_ENVIRONMENT" envDefault:"dev"`
}

type metricsConfig struct {
	raw metricsEnvConfig
}

func NewOrderMetricsConfig() (*metricsConfig, error) {
	var raw metricsEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &metricsConfig{raw: raw}, nil
}

func (cfg *metricsConfig) CollectorEndpoint() string {
	return cfg.raw.CollectorEndpoint
}

func (cfg *metricsConfig) CollectorInterval() time.Duration {
	return cfg.raw.CollectorInterval
}

func (cfg *metricsConfig) ServiceName() string {
	return cfg.raw.ServiceName
}

func (cfg *metricsConfig) Environment() string {
	return cfg.raw.Environment
}
