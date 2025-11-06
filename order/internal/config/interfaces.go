package config

import "time"

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
