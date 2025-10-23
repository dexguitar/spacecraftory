![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/dexguitar/facd6b4c2386f87586848b74c73136e4/raw/coverage.json)
# Payment Service

gRPC service with HTTP/REST gateway for processing spacecraft order payments. Generates transaction UUIDs and handles payment method validation.

## 🚀 Quick Start

```bash
# Start the service
go run cmd/grpc_server/main.go
```

The service will start on:

- **gRPC:** `localhost:50052`
- **HTTP Gateway:** `http://localhost:8082`

## 📡 API Endpoints

### Pay for Order

Processes payment for an order and returns a transaction UUID.

#### HTTP/REST (via gRPC-Gateway)

```bash
curl -X POST http://localhost:8082/api/v1/payments \
  -H "Content-Type: application/json" \
  -d '{
    "order_uuid": "123e4567-e89b-12d3-a456-426614174000",
    "user_uuid": "550e8400-e29b-41d4-a716-446655440000",
    "payment_method": "PAYMENT_METHOD_CARD"
  }'
```

**Response:**

```json
{
  "transaction_uuid": "789e4567-e89b-12d3-a456-426614174999"
}
```

#### gRPC

```bash
# Using grpcurl (must be installed)
grpcurl -plaintext \
  -d '{
    "order_uuid": "123e4567-e89b-12d3-a456-426614174000",
    "user_uuid": "550e8400-e29b-41d4-a716-446655440000",
    "payment_method": "PAYMENT_METHOD_CARD"
  }' \
  localhost:50052 \
  payment.v1.PaymentService/PayOrder
```

**Response:**

```json
{
  "transactionUuid": "789e4567-e89b-12d3-a456-426614174999"
}
```

---

## 💳 Payment Methods

The following payment methods are supported:

| Method         | Value                                | Description         |
| -------------- | ------------------------------------ | ------------------- |
| Unknown        | `PAYMENT_METHOD_UNKNOWN_UNSPECIFIED` | Default/unspecified |
| Card           | `PAYMENT_METHOD_CARD`                | Bank card           |
| SBP            | `PAYMENT_METHOD_SBP`                 | Fast Payment System |
| Credit Card    | `PAYMENT_METHOD_CREDIT_CARD`         | Credit card         |
| Investor Money | `PAYMENT_METHOD_INVESTOR_MONEY`      | Investor funds      |

---

## 🧪 Example: Full Payment Flow

```bash
# Using HTTP Gateway
curl -X POST http://localhost:8082/api/v1/payments \
  -H "Content-Type: application/json" \
  -d '{
    "order_uuid": "123e4567-e89b-12d3-a456-426614174000",
    "user_uuid": "550e8400-e29b-41d4-a716-446655440000",
    "payment_method": "PAYMENT_METHOD_CARD"
  }' | jq

# Using gRPC
grpcurl -plaintext \
  -d '{
    "order_uuid": "123e4567-e89b-12d3-a456-426614174000",
    "user_uuid": "550e8400-e29b-41d4-a716-446655440000",
    "payment_method": "PAYMENT_METHOD_CARD"
  }' \
  localhost:50052 \
  payment.v1.PaymentService/PayOrder | jq
```

---

## 🔧 Configuration

- **gRPC Port:** `50052`
- **HTTP Gateway Port:** `8082`
- **Read Header Timeout (HTTP):** `10s`
- **Shutdown Timeout:** `5s`

---

## 🔍 Service Reflection

The service has gRPC reflection enabled for debugging:

```bash
# List all services
grpcurl -plaintext localhost:50052 list

# List all methods for PaymentService
grpcurl -plaintext localhost:50052 list payment.v1.PaymentService

# Describe PayOrder method
grpcurl -plaintext localhost:50052 describe payment.v1.PaymentService.PayOrder
```

**Output:**

```
payment.v1.PaymentService.PayOrder is a method:
rpc PayOrder ( .payment.v1.PayOrderRequest ) returns ( .payment.v1.PayOrderResponse );
```

---

## 🌐 Integration with Order Service

The Payment service is designed to work with the Order service:

```bash
# 1. Create an order (Order Service)
ORDER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/payments \
  -H "Content-Type: application/json" \
  -d '{
    "user_uuid": "550e8400-e29b-41d4-a716-446655440000",
    "part_uuids": ["123e4567-e89b-12d3-a456-426614174001"]
  }')

ORDER_UUID=$(echo $ORDER_RESPONSE | jq -r '.order_uuid')

# 2. Process payment (Payment Service)
PAYMENT_RESPONSE=$(curl -s -X POST http://localhost:8082/api/v1/payments \
  -H "Content-Type: application/json" \
  -d "{
    \"order_uuid\": \"$ORDER_UUID\",
    \"user_uuid\": \"550e8400-e29b-41d4-a716-446655440000\",
    \"payment_method\": \"PAYMENT_METHOD_CARD\"
  }")

TRANSACTION_UUID=$(echo $PAYMENT_RESPONSE | jq -r '.transaction_uuid')

echo "Payment processed: $TRANSACTION_UUID"

# 3. Update order with payment (Order Service)
curl -X POST http://localhost:8080/api/v1/payments/$ORDER_UUID/pay \
  -H "Content-Type: application/json" \
  -d '{"payment_method": "CARD"}'
```

---

## 🛠️ Development Tools

### Installing grpcurl

```bash
# macOS
brew install grpcurl

# Linux
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Or via task
task grpcurl:install
```

### Testing gRPC Connection

```bash
# Check if service is running
grpcurl -plaintext localhost:50052 list

# Expected output:
# grpc.reflection.v1alpha.ServerReflection
# payment.v1.PaymentService
```

---

## 🛡️ Error Handling

gRPC errors follow standard gRPC status codes:

- `OK` (0) - Success
- `INVALID_ARGUMENT` (3) - Invalid request parameters
- `NOT_FOUND` (5) - Resource not found
- `INTERNAL` (13) - Internal server error

HTTP Gateway errors are translated to standard HTTP status codes:

- `200 OK` - Success
- `400 Bad Request` - Invalid parameters
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

---

## 📝 Proto Definition

Located at: `shared/proto/payment/v1/payment.proto`

```protobuf
service PaymentService {
    rpc PayOrder (PayOrderRequest) returns (PayOrderResponse) {
        option (google.api.http) = {
            post: "/api/v1/payments"
            body: "*"
        };
    };
}
```
