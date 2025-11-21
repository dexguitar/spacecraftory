package config

import (
	"time"

	"github.com/IBM/sarama"
)

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type HTTPConfig interface {
	Address() string
	ReadTimeout() time.Duration
}

type GRPCClientConfig interface {
	InventoryAddress() string
	PaymentAddress() string
}

type PostgresConfig interface {
	Address() string
	ExternalPort() string
	MigrationDirectory() string
}
type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type OrderAssembledConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}
