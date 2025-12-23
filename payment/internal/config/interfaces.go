package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
	OtelEndpoint() string
	ServiceName() string
}

type PaymentGRPCConfig interface {
	Address() string
}

type TracingConfig interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
}
