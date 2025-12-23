package config

import (
	"time"

	"github.com/IBM/sarama"
)

type LoggerConfig interface {
	Level() string
	AsJson() bool
	OtelEndpoint() string
	ServiceName() string
}

type MetricsConfig interface {
	CollectorEndpoint() string
	CollectorInterval() time.Duration
	ServiceName() string
	Environment() string
}

type TracingConfig interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
}

type HTTPConfig interface {
	Address() string
	ReadTimeout() time.Duration
}

type GRPCClientConfig interface {
	InventoryAddress() string
	PaymentAddress() string
	IAMAddress() string
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
