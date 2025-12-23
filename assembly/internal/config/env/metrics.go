package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type metricsEnvConfig struct {
	CollectorEndpoint string        `env:"OTEL_COLLECTOR_ENDPOINT" envDefault:""`
	CollectorInterval time.Duration `env:"METRICS_COLLECTOR_INTERVAL" envDefault:"10s"`
	ServiceName       string        `env:"SERVICE_NAME" envDefault:"assembly-service"`
	Environment       string        `env:"ENVIRONMENT" envDefault:"dev"`
}

type metricsConfig struct {
	raw metricsEnvConfig
}

func NewMetricsConfig() (*metricsConfig, error) {
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

