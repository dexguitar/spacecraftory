# Order Service

HTTP/REST API service for managing spacecraft orders. Handles order creation, payment, retrieval, and cancellation.

## 🚀 Quick Start

```bash
# Start the service
go run cmd/http_server/main.go
```

The service will start on **http://localhost:8080**

## 📡 API Endpoints

### 1. Create Order

Creates a new order for spacecraft parts.

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_uuid": "550e8400-e29b-41d4-a716-446655440000",
    "part_uuids": [
      "123e4567-e89b-12d3-a456-426614174001",
      "123e4567-e89b-12d3-a456-426614174002"
    ]
  }'
```

**Response:**

```json
{
  "order_uuid": "123e4567-e89b-12d3-a456-426614174000",
  "total_price": 200.0
}
```

---

### 2. Get Order by UUID

Retrieves order details by UUID.

```bash
curl http://localhost:8080/api/v1/orders/123e4567-e89b-12d3-a456-426614174000
```

**Response:**

```json
{
  "order_uuid": "123e4567-e89b-12d3-a456-426614174000",
  "user_uuid": "550e8400-e29b-41d4-a716-446655440000",
  "part_uuids": [
    "123e4567-e89b-12d3-a456-426614174001",
    "123e4567-e89b-12d3-a456-426614174002"
  ],
  "total_price": 200.0,
  "status": "PENDING_PAYMENT"
}
```

---

### 3. Pay for Order

Processes payment for an order.

```bash
curl -X POST http://localhost:8080/api/v1/orders/123e4567-e89b-12d3-a456-426614174000/pay \
  -H "Content-Type: application/json" \
  -d '{
    "payment_method": "CARD"
  }'
```

**Response:**

```json
{
  "transaction_uuid": "789e4567-e89b-12d3-a456-426614174999"
}
```

**Available Payment Methods:**

- `CARD` - Bank card
- `SBP` - Fast Payment System
- `CREDIT_CARD` - Credit card
- `INVESTOR_MONEY` - Investor funds

---

### 4. Cancel Order

Cancels a pending order (only available for unpaid orders).

```bash
curl -X POST http://localhost:8080/api/v1/orders/123e4567-e89b-12d3-a456-426614174000/cancel
```

**Response:** `204 No Content` (on success)

**Error Responses:**

- `404 Not Found` - Order not found
- `409 Conflict` - Order is already paid and cannot be cancelled

---

## 📊 Order Statuses

- `PENDING_PAYMENT` - Order created, awaiting payment
- `PAID` - Order has been paid
- `CANCELLED` - Order has been cancelled

---

## 🧪 Full Example Flow

```bash
# 1. Create an order
ORDER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_uuid": "550e8400-e29b-41d4-a716-446655440000",
    "part_uuids": ["123e4567-e89b-12d3-a456-426614174001"]
  }')

# Extract order UUID (requires jq)
ORDER_UUID=$(echo $ORDER_RESPONSE | jq -r '.order_uuid')

echo "Created order: $ORDER_UUID"

# 2. Check order details
curl http://localhost:8080/api/v1/orders/$ORDER_UUID

# 3. Pay for the order
curl -X POST http://localhost:8080/api/v1/orders/$ORDER_UUID/pay \
  -H "Content-Type: application/json" \
  -d '{"payment_method": "CARD"}'

# 4. Verify payment went through
curl http://localhost:8080/api/v1/orders/$ORDER_UUID
```

---

## 🔧 Configuration

- **HTTP Port:** `8080`
- **Read Header Timeout:** `5s` (protection against Slowloris attacks)
- **Shutdown Timeout:** `10s`

---

## 🛡️ Error Handling

All errors follow a standard format:

```json
{
  "code": 404,
  "message": "Order 123e4567-e89b-12d3-a456-426614174000 not found"
}
```

**Common HTTP Status Codes:**

- `200 OK` - Success
- `201 Created` - Resource created
- `204 No Content` - Success with no body
- `400 Bad Request` - Invalid request data
- `404 Not Found` - Resource not found
- `409 Conflict` - Operation not allowed (e.g., cancelling paid order)
- `500 Internal Server Error` - Server error
