package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type iamClientGRPCEnvConfig struct {
	Host string `env:"IAM_CLIENT_GRPC_HOST,required"`
	Port string `env:"IAM_CLIENT_GRPC_PORT,required"`
}

type iamClientGRPCConfig struct {
	raw iamClientGRPCEnvConfig
}

func NewIAMClientGRPCConfig() (*iamClientGRPCConfig, error) {
	var raw iamClientGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &iamClientGRPCConfig{raw: raw}, nil
}

func (cfg *iamClientGRPCConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
