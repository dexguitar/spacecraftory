package env

import (
	"github.com/caarlos0/env/v11"
)

type orderLoggerEnvConfig struct {
	Level  string `env:"ORDER_LOGGER_LEVEL,required"`
	AsJson bool   `env:"ORDER_LOGGER_AS_JSON,required"`
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
