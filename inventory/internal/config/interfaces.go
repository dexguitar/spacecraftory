package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
	OtelEndpoint() string
	ServiceName() string
}

type InventoryGRPCConfig interface {
	Address() string
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}

type IAMClientGRPCConfig interface {
	Address() string
}
