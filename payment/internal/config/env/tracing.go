package env

import (
	"github.com/caarlos0/env/v11"
)

type tracingEnvConfig struct {
	CollectorEndpoint string `env:"PAYMENT_OTEL_COLLECTOR_ENDPOINT" envDefault:""`
	ServiceName       string `env:"PAYMENT_SERVICE_NAME" envDefault:"payment-service"`
	Environment       string `env:"PAYMENT_ENVIRONMENT" envDefault:"dev"`
	ServiceVersion    string `env:"PAYMENT_SERVICE_VERSION" envDefault:"1.0.0"`
}

type tracingConfig struct {
	raw tracingEnvConfig
}

func NewPaymentTracingConfig() (*tracingConfig, error) {
	var raw tracingEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &tracingConfig{raw: raw}, nil
}

func (cfg *tracingConfig) CollectorEndpoint() string {
	return cfg.raw.CollectorEndpoint
}

func (cfg *tracingConfig) ServiceName() string {
	return cfg.raw.ServiceName
}

func (cfg *tracingConfig) Environment() string {
	return cfg.raw.Environment
}

func (cfg *tracingConfig) ServiceVersion() string {
	return cfg.raw.ServiceVersion
}
