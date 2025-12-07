package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type orderGRPCClientEnvConfig struct {
	InventoryHost string `env:"ORDER_INVENTORY_GRPC_HOST,required"`
	InventoryPort string `env:"ORDER_INVENTORY_GRPC_PORT,required"`
	PaymentHost   string `env:"ORDER_PAYMENT_GRPC_HOST,required"`
	PaymentPort   string `env:"ORDER_PAYMENT_GRPC_PORT,required"`
	IAMHost       string `env:"ORDER_IAM_GRPC_HOST,required"`
	IAMPort       string `env:"ORDER_IAM_GRPC_PORT,required"`
}

type orderGRPCClientConfig struct {
	raw orderGRPCClientEnvConfig
}

func NewOrderGRPCClientConfig() (*orderGRPCClientConfig, error) {
	var raw orderGRPCClientEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &orderGRPCClientConfig{raw: raw}, nil
}

func (cfg *orderGRPCClientConfig) InventoryAddress() string {
	return net.JoinHostPort(cfg.raw.InventoryHost, cfg.raw.InventoryPort)
}

func (cfg *orderGRPCClientConfig) PaymentAddress() string {
	return net.JoinHostPort(cfg.raw.PaymentHost, cfg.raw.PaymentPort)
}

func (cfg *orderGRPCClientConfig) IAMAddress() string {
	return net.JoinHostPort(cfg.raw.IAMHost, cfg.raw.IAMPort)
}
