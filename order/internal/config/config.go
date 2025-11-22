package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/dexguitar/spacecraftory/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger                 LoggerConfig
	HTTP                   HTTPConfig
	GRPCClient             GRPCClientConfig
	Postgres               PostgresConfig
	Kafka                  KafkaConfig
	OrderPaidProducer      OrderPaidProducerConfig
	OrderAssembledConsumer OrderAssembledConsumerConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewOrderLoggerConfig()
	if err != nil {
		return err
	}

	httpCfg, err := env.NewOrderHTTPConfig()
	if err != nil {
		return err
	}

	grpcClientCfg, err := env.NewOrderGRPCClientConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewOrderPostgresConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewOrderKafkaConfig()
	if err != nil {
		return err
	}

	orderPaidProducerCfg, err := env.NewOrderPaidProducerConfig()
	if err != nil {
		return err
	}

	orderAssembledConsumerCfg, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                 loggerCfg,
		HTTP:                   httpCfg,
		GRPCClient:             grpcClientCfg,
		Postgres:               postgresCfg,
		Kafka:                  kafkaCfg,
		OrderPaidProducer:      orderPaidProducerCfg,
		OrderAssembledConsumer: orderAssembledConsumerCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
