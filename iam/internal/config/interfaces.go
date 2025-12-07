package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type IAMGRPCConfig interface {
	Address() string
}

type PostgresConfig interface {
	Address() string
	MigrationDirectory() string
}

type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
	CacheTTL() time.Duration
}
