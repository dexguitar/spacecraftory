# ============================
# Stage 1: Build stage
# ============================

# Use official Go image based on Alpine — lightweight, fast and secure
FROM golang:1.24.2-alpine AS builder

# Install git — required for downloading dependencies from private and public repositories
RUN apk add --no-cache git

# Set working directory inside the container
WORKDIR /app

# Copy go.work and its checksum — this allows Go to properly link modules
COPY go.work go.work.sum ./

# Copy go.mod and go.sum of all dependent modules — needed for caching go mod download
COPY platform/go.mod platform/go.sum ./platform/
COPY shared/go.mod shared/go.sum ./shared/
COPY inventory/go.mod inventory/go.sum ./inventory/

# Download all dependencies specified in module files
RUN go mod download

# Copy source code of all modules — including dependencies, as they are needed during compilation
COPY platform ./platform
COPY shared ./shared
COPY inventory ./inventory

# Download grpc-health-probe — a utility for checking gRPC service health
ADD https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.41/grpc_health_probe-darwin-arm64 ./grpc-health-probe

# Make the downloaded file executable
RUN chmod +x grpc-health-probe

# Build Inventory service binary for Darwin architecture without CGO
RUN CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o app-inventory ./inventory/cmd/grpc_server/main.go


# ============================
# Stage 2: Final image
# ============================

# Use clean Alpine image as final container — it will be lightweight and secure
FROM alpine:3.21.3

# Create system user without root privileges — best practice for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory in container
WORKDIR /app

# Run application under non-root user
USER appuser

# Copy compiled binary from builder stage
COPY --from=builder /app/app-inventory .

# Copy grpc-health-probe binary to /bin
COPY --from=builder /app/grpc-health-probe /bin/grpc-health-probe

# Copy .env file for application configuration
COPY deploy/compose/inventory/.env .env

# Copy .env file for application configuration
COPY deploy/compose/inventory/.env /deploy/compose/inventory/.env

# Expose Inventory gRPC service port
EXPOSE 50051

# Set startup command — run our binary
ENTRYPOINT ["./app-inventory"]