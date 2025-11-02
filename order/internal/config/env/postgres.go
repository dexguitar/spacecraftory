package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type orderPostgresEnvConfig struct {
	Host               string `env:"ORDER_POSTGRES_HOST,required"`
	Port               string `env:"ORDER_POSTGRES_PORT,required"`
	ExternalPort       string `env:"ORDER_EXTERNAL_POSTGRES_PORT,required"`
	User               string `env:"ORDER_POSTGRES_USER,required"`
	Password           string `env:"ORDER_POSTGRES_PASSWORD,required"`
	Database           string `env:"ORDER_POSTGRES_DB,required"`
	SSLMode            string `env:"ORDER_POSTGRES_SSL_MODE,required"`
	MigrationDirectory string `env:"ORDER_POSTGRES_MIGRATION_DIRECTORY,required"`
}

type orderPostgresConfig struct {
	raw orderPostgresEnvConfig
}

func NewOrderPostgresConfig() (*orderPostgresConfig, error) {
	var raw orderPostgresEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderPostgresConfig{raw: raw}, nil
}

func (cfg *orderPostgresConfig) Address() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.raw.Host, cfg.raw.Port, cfg.raw.User, cfg.raw.Password, cfg.raw.Database, cfg.raw.SSLMode)
}

func (cfg *orderPostgresConfig) ExternalPort() string {
	return cfg.raw.ExternalPort
}

func (cfg *orderPostgresConfig) MigrationDirectory() string {
	return cfg.raw.MigrationDirectory
}
