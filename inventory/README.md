# Inventory Service

gRPC service with HTTP/REST gateway for managing spacecraft parts inventory. Provides operations to retrieve individual parts and list parts with flexible filtering.

## 🚀 Quick Start

```bash
# Start the service
go run cmd/grpc_server/main.go
```

The service will start on:

- **gRPC:** `localhost:50051`
- **HTTP Gateway:** `http://localhost:8081`

## 📡 API Endpoints

### 1. Get Part by UUID

Retrieves a single spacecraft part by its UUID.

#### HTTP/REST (via gRPC-Gateway)

```bash
curl http://localhost:8081/api/v1/parts/PART_UUID_HERE
```

#### gRPC

```bash
grpcurl -plaintext \
  -d '{"uuid": "PART_UUID_HERE"}' \
  localhost:50051 \
  inventory.v1.InventoryService/GetPart
```

**Response:**

```json
{
  "part": {
    "uuid": "123e4567-e89b-12d3-a456-426614174001",
    "name": "Quantum Drive Engine",
    "description": "High-efficiency quantum propulsion engine for interstellar travel",
    "price": 150000.0,
    "stock_quantity": 5,
    "category": "CATEGORY_ENGINE",
    "dimensions": {
      "length": 3.5,
      "width": 2.0,
      "height": 2.5,
      "weight": 500.0
    },
    "manufacturer": {
      "name": "SpaceTech Industries",
      "country": "USA",
      "website": "https://spacetech.example.com"
    },
    "tags": ["quantum", "propulsion", "interstellar"],
    "created_at": "2025-10-11T12:00:00Z",
    "updated_at": "2025-10-11T12:00:00Z"
  }
}
```

---

### 2. List Parts

Lists spacecraft parts with optional filtering by various criteria.

#### HTTP/REST: List All Parts

```bash
curl http://localhost:8081/api/v1/parts
```

#### HTTP/REST: Filter by Specific UUIDs

```bash
curl -X POST http://localhost:8081/api/v1/parts \
  -H "Content-Type: application/json" \
  -d '{
    "filter": {
      "uuids": [
        "123e4567-e89b-12d3-a456-426614174001",
        "123e4567-e89b-12d3-a456-426614174002"
      ]
    }
  }'
```

#### HTTP/REST: Filter by Category

```bash
curl -X POST http://localhost:8081/api/v1/parts \
  -H "Content-Type: application/json" \
  -d '{
    "filter": {
      "categories": ["CATEGORY_ENGINE", "CATEGORY_FUEL"]
    }
  }'
```

#### HTTP/REST: Filter by Multiple Criteria

```bash
curl -X POST http://localhost:8081/api/v1/parts \
  -H "Content-Type: application/json" \
  -d '{
    "filter": {
      "names": ["Quantum Drive Engine"],
      "categories": ["CATEGORY_ENGINE"],
      "manufacturer_countries": ["USA"],
      "tags": ["quantum"]
    }
  }'
```

#### gRPC: List Parts

```bash
# List all parts
grpcurl -plaintext \
  localhost:50051 \
  inventory.v1.InventoryService/ListParts

# Filter by categories
grpcurl -plaintext \
  -d '{
    "filter": {
      "categories": ["CATEGORY_ENGINE"]
    }
  }' \
  localhost:50051 \
  inventory.v1.InventoryService/ListParts
```

**Response:**

```json
{
  "parts": [
    {
      "uuid": "123e4567-e89b-12d3-a456-426614174001",
      "name": "Quantum Drive Engine",
      "description": "High-efficiency quantum propulsion engine for interstellar travel",
      "price": 150000.0,
      "stock_quantity": 5,
      "category": "CATEGORY_ENGINE",
      ...
    },
    ...
  ]
}
```

---

## 🏷️ Categories

Available spacecraft part categories:

| Category | Value                          | Description          |
| -------- | ------------------------------ | -------------------- |
| Unknown  | `CATEGORY_UNKNOWN_UNSPECIFIED` | Default/unspecified  |
| Engine   | `CATEGORY_ENGINE`              | Propulsion systems   |
| Fuel     | `CATEGORY_FUEL`                | Fuel cells and tanks |
| Porthole | `CATEGORY_PORTHOLE`            | Viewing windows      |
| Wing     | `CATEGORY_WING`                | Aerodynamic panels   |

---

## 🔍 Filtering Logic

The filtering system uses **AND** logic between different filter types and **OR** logic within each type:

### Examples:

**Single filter type (OR logic):**

```json
{
  "filter": {
    "categories": ["CATEGORY_ENGINE", "CATEGORY_FUEL"]
  }
}
```

Returns parts that are **ENGINE OR FUEL**

**Multiple filter types (AND logic):**

```json
{
  "filter": {
    "categories": ["CATEGORY_ENGINE"],
    "manufacturer_countries": ["USA", "Germany"]
  }
}
```

Returns parts that are:

- **ENGINE** (required)
- **AND** from **USA OR Germany** (required)

**Available filters:**

- `uuids` - List of specific part UUIDs (short-circuits other filters)
- `names` - Part names
- `categories` - Part categories
- `manufacturer_countries` - Manufacturer countries
- `tags` - Part tags

---

## 🧪 Example Workflows

### Find High-Tech Parts

```bash
# Find all quantum-tagged engine parts
curl -X POST http://localhost:8081/api/v1/parts \
  -H "Content-Type: application/json" \
  -d '{
    "filter": {
      "categories": ["CATEGORY_ENGINE"],
      "tags": ["quantum"]
    }
  }' | jq '.parts[] | {name: .name, price: .price}'
```

### Check Stock for Multiple Parts

```bash
# Get specific parts by UUID
curl -X POST http://localhost:8081/api/v1/parts \
  -H "Content-Type: application/json" \
  -d '{
    "filter": {
      "uuids": [
        "123e4567-e89b-12d3-a456-426614174001",
        "123e4567-e89b-12d3-a456-426614174002"
      ]
    }
  }' | jq '.parts[] | {name: .name, stock: .stock_quantity}'
```

### Find Parts by Manufacturer

```bash
# Find all parts made in Japan
curl -X POST http://localhost:8081/api/v1/parts \
  -H "Content-Type: application/json" \
  -d '{
    "filter": {
      "manufacturer_countries": ["Japan"]
    }
  }' | jq
```

---

## 🔧 Configuration

- **gRPC Port:** `50051`
- **HTTP Gateway Port:** `8081`
- **Read Header Timeout (HTTP):** `10s`
- **Shutdown Timeout:** `5s`

---

## 🗄️ Mock Data

The service initializes with 4 sample parts:

1. **Quantum Drive Engine** - $150,000 (ENGINE)
2. **Fusion Fuel Cell** - $75,000 (FUEL)
3. **Reinforced Porthole** - $25,000 (PORTHOLE)
4. **Aerodynamic Wing Panel** - $45,000 (WING)

---

## 🔍 Service Reflection

The service has gRPC reflection enabled for debugging:

```bash
# List all services
grpcurl -plaintext localhost:50051 list

# List all methods for InventoryService
grpcurl -plaintext localhost:50051 list inventory.v1.InventoryService

# Describe GetPart method
grpcurl -plaintext localhost:50051 describe inventory.v1.InventoryService.GetPart

# Describe ListParts method
grpcurl -plaintext localhost:50051 describe inventory.v1.InventoryService.ListParts
```

---

## 🌐 Integration Example

Using Inventory service with Order service:

```bash
# 1. Get available engine parts
PARTS=$(curl -s -X POST http://localhost:8081/api/v1/parts \
  -H "Content-Type: application/json" \
  -d '{
    "filter": {
      "categories": ["CATEGORY_ENGINE"]
    }
  }')

# Extract part UUIDs (requires jq)
PART_UUIDS=$(echo $PARTS | jq -r '.parts[].uuid')

echo "Available engines: $PART_UUIDS"

# 2. Create order with selected parts (Order Service)
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d "{
    \"user_uuid\": \"550e8400-e29b-41d4-a716-446655440000\",
    \"part_uuids\": [\"$(echo $PART_UUIDS | head -n1)\"]
  }"
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
grpcurl -plaintext localhost:50051 list

# Expected output:
# grpc.reflection.v1alpha.ServerReflection
# inventory.v1.InventoryService
```

---

## 🛡️ Error Handling

**gRPC Status Codes:**

- `OK` (0) - Success
- `INVALID_ARGUMENT` (3) - UUID is required or invalid
- `NOT_FOUND` (5) - Part not found
- `INTERNAL` (13) - Internal server error

**HTTP Status Codes (Gateway):**

- `200 OK` - Success
- `400 Bad Request` - Invalid parameters
- `404 Not Found` - Part not found
- `500 Internal Server Error` - Server error

---

## 📝 Proto Definition

Located at: `shared/proto/inventory/v1/inventory.proto`

```protobuf
service InventoryService {
    rpc GetPart(GetPartRequest) returns (GetPartResponse);
    rpc ListParts(ListPartsRequest) returns (ListPartsResponse);
}
```
