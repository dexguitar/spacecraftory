module github.com/dexguitar/spacecraftory/inventory

go 1.25.2

require (
	github.com/dexguitar/spacecraftory/shared v0.0.0
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.38.0 // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251007200510-49b9836ed3ff // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251007200510-49b9836ed3ff // indirect
)

replace github.com/dexguitar/spacecraftory/shared => ../shared
