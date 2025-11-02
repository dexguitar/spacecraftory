package env

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type orderHTTPEnvConfig struct {
	Host        string        `env:"ORDER_HTTP_HOST,required"`
	Port        string        `env:"ORDER_HTTP_PORT,required"`
	ReadTimeout time.Duration `env:"ORDER_HTTP_READ_TIMEOUT,required"`
}

type orderHTTPConfig struct {
	raw orderHTTPEnvConfig
}

func NewOrderHTTPConfig() (*orderHTTPConfig, error) {
	var raw orderHTTPEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderHTTPConfig{raw: raw}, nil
}

func (cfg *orderHTTPConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

func (cfg *orderHTTPConfig) ReadTimeout() time.Duration {
	return cfg.raw.ReadTimeout
}
