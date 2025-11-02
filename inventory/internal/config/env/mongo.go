package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type mongoEnvConfig struct {
	Host     string `env:"INVENTORY_MONGO_HOST,required"`
	Port     string `env:"INVENTORY_MONGO_PORT,required"`
	Database string `env:"INVENTORY_MONGO_INITDB_DATABASE,required"`
	User     string `env:"INVENTORY_MONGO_INITDB_ROOT_USERNAME,required"`
	Password string `env:"INVENTORY_MONGO_INITDB_ROOT_PASSWORD,required"`
	AuthDB   string `env:"INVENTORY_MONGO_AUTH_DB,required"`
}

type mongoConfig struct {
	raw mongoEnvConfig
}

func NewMongoConfig() (*mongoConfig, error) {
	var raw mongoEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &mongoConfig{raw: raw}, nil
}

func (cfg *mongoConfig) URI() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=%s",
		cfg.raw.User,
		cfg.raw.Password,
		cfg.raw.Host,
		cfg.raw.Port,
		cfg.raw.Database,
		cfg.raw.AuthDB,
	)
}

func (cfg *mongoConfig) DatabaseName() string {
	return cfg.raw.Database
}