package env

import (
	"github.com/caarlos0/env/v11"
)

type orderLoggerEnvConfig struct {
	Level        string `env:"ORDER_LOGGER_LEVEL,required"`
	AsJson       bool   `env:"ORDER_LOGGER_AS_JSON,required"`
	OtelEndpoint string `env:"ORDER_OTEL_COLLECTOR_ENDPOINT" envDefault:""`
	ServiceName  string `env:"ORDER_SERVICE_NAME" envDefault:"order-service"`
}

type orderLoggerConfig struct {
	raw orderLoggerEnvConfig
}

func NewOrderLoggerConfig() (*orderLoggerConfig, error) {
	var raw orderLoggerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderLoggerConfig{raw: raw}, nil
}

func (cfg *orderLoggerConfig) Level() string {
	return cfg.raw.Level
}

func (cfg *orderLoggerConfig) AsJson() bool {
	return cfg.raw.AsJson
}

func (cfg *orderLoggerConfig) OtelEndpoint() string {
	return cfg.raw.OtelEndpoint
}

func (cfg *orderLoggerConfig) ServiceName() string {
	return cfg.raw.ServiceName
}
