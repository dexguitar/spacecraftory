package tracing

import "google.golang.org/grpc/metadata"

// metadataCarrier is an adapter between gRPC metadata and OpenTelemetry's TextMapCarrier.
// It implements the propagation.TextMapCarrier interface for trace context propagation.
type metadataCarrier metadata.MD

// Get returns the value for the specified key.
// Returns an empty string if the key doesn't exist.
func (mc metadataCarrier) Get(key string) string {
	values := metadata.MD(mc).Get(key)
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

// Set sets the value for the specified key.
func (mc metadataCarrier) Set(key, value string) {
	metadata.MD(mc).Set(key, value)
}

// Keys returns a list of all keys in the carrier.
func (mc metadataCarrier) Keys() []string {
	keys := make([]string, 0, len(mc))
	for k := range metadata.MD(mc) {
		keys = append(keys, k)
	}

	return keys
}
