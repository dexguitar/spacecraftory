package env

import (
	"github.com/caarlos0/env/v11"
)

type orderKafkaEnvConfig struct {
	Brokers []string `env:"ORDER_KAFKA_BROKERS,required"`
}

type orderKafkaConfig struct {
	raw orderKafkaEnvConfig
}

func NewOrderKafkaConfig() (*orderKafkaConfig, error) {
	var raw orderKafkaEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderKafkaConfig{raw: raw}, nil
}

func (cfg *orderKafkaConfig) Brokers() []string {
	return cfg.raw.Brokers
}
