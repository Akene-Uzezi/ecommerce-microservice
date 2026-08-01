# Ecommerce Microservices

A Go-based ecommerce microservices platform using gRPC for inter-service communication and an HTTP API Gateway as the public entrypoint. The project is organized as a Go workspace with independent modules per service.

## Architecture

```
Client (HTTP/JSON)
    │
    ▼
Gateway (:3000) ──gRPC──▶ Orders (:4444)
                              │
                    (future) ▼
                    Stock / Payments
```

The Gateway exposes a REST/JSON interface and translates requests into gRPC calls to backend services. Services communicate directly via gRPC.

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.26.5 |
| RPC Framework | gRPC |
| Serialization | Protocol Buffers |
| Code Generation | protoc, protoc-gen-go, protoc-gen-go-grpc |
| Workspace | Go workspaces (`go.work`) |

## Project Structure

```
ecommerce-microservices/
├── api/                     # Shared API contracts (protobuf definitions + generated stubs)
│   ├── proto/               # .proto source files
│   └── gen/                 # Generated Go gRPC code
├── gateway/                 # HTTP API Gateway
├── orders/                  # Orders microservice
├── payments/                # Payments microservice (scaffold)
├── stock/                   # Stock microservice (scaffold)
├── shared/                  # Common utilities
├── scripts/                 # SQL init scripts, commit helper
├── go.work                  # Go workspace file
├── go.mod                   # Root module (ecommerce-api)
├── Makefile                 # Protobuf generation targets
└── docker-compose.yml       # Container orchestration (to be configured)
```

## Services

### API Contract (`api/`)

Canonical source of truth for all service interfaces. Contains `.proto` definitions and generated Go stubs.

- **order.proto** — `OrderService`: `CreateOrder`, `GetOrder`
- **stock.proto** — `StockService`: `CheckStock`, `ReserveStock`
- **payment.proto** — `PaymentService`: `ProcessPayment`

### Gateway (`gateway/`)

HTTP entrypoint that proxies requests to backend gRPC services.

- Listens on `:3000`
- `GET /api/v1/ping` — health check
- `POST /api/v1/orders` — creates an order via the Orders service

### Orders Service (`orders/`)

Manages order creation and retrieval.

- Listens on `:4444`
- `CreateOrder` — returns a stub response (not yet connected to stock/payment flow)
- `GetOrder` — defined in proto but not yet implemented

### Payments Service (`payments/`)

Module scaffold only. Intended to handle `ProcessPayment(order_id, customer_id, amount)`.

### Stock Service (`stock/`)

Module scaffold only. Intended to handle `CheckStock(product_id, quantity)` and `ReserveStock(order_id, items)`.

### Shared (`shared/`)

Common utilities:

- `env.go` — `GetEnvString(key, fallback)`
- `json.go` — HTTP JSON helpers: `WriteJSON`, `ReadJSON`

## Prerequisites

- Go 1.26.5
- Protocol Buffers compiler (`protoc`)
- protoc-gen-go (`go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`)
- protoc-gen-go-grpc (`go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`)

## Getting Started

### 1. Generate Protobuf Stubs

```bash
make gen
```

This generates Go gRPC stubs from `api/proto/*.proto` into `api/gen/`.

### 2. Run Services

Run each module from its directory:

```bash
# Terminal 1 — Orders service
cd orders && go run ./cmd

# Terminal 2 — Gateway
cd gateway && go run ./cmd
```

### 3. Test

```bash
# Health check
curl http://localhost:3000/api/v1/ping

# Create an order
curl -X POST http://localhost:3000/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"customer_id":"1","items":[{"product_id":"3","quantity":3,"price":10.62}]}'
```

## Development

### Make Targets

| Command | Description |
|---------|-------------|
| `make gen` | Regenerate gRPC stubs from protobuf |
| `make clean` | Remove generated files |

### Workspace

The project uses a Go workspace (`go.work`) to allow cross-module imports (e.g., the gateway imports `ecommerce-api/gen/order`). Run `go work sync` if modules change.

## Roadmap

- [x] Protobuf API contracts
- [x] Orders service (skeleton)
- [x] Gateway with HTTP-to-gRPC translation
- [ ] Stock service implementation
- [ ] Payments service implementation
- [ ] Database integration (orders, stock)
- [ ] Docker / docker-compose setup
- [ ] Service discovery / config
- [ ] TLS / production hardening
- [ ] Testing suite (unit + integration)
