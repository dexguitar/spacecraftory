package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/dexguitar/spacecraftory/iam/internal/config/env"
)

var appConfig *config

type config struct {
	Logger   LoggerConfig
	IAMGRPC  IAMGRPCConfig
	Postgres PostgresConfig
	Redis    RedisConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	iamGRPCCfg, err := env.NewIAMGRPCConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	redisCfg, err := env.NewRedisConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:   loggerCfg,
		IAMGRPC:  iamGRPCCfg,
		Postgres: postgresCfg,
		Redis:    redisCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
