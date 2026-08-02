# Ecommerce Microservices

A Go-based ecommerce microservices platform using gRPC for inter-service communication and an HTTP API Gateway as the public entrypoint. The project is organized as a Go workspace with independent modules per service.

## Architecture

```
Client (HTTP/JSON)
     │
     ▼
Gateway (:3000)
     │
     ├──gRPC──▶ Auth (:5555)
     │              │
     │              ▼
     │         Auth DB (PostgreSQL :6433)
     │
     └──gRPC──▶ Orders (:4444)
                    │
                    ▼
               Orders DB (PostgreSQL :5433)
```

The Gateway exposes a REST/JSON interface and translates requests into gRPC calls to backend services. Services communicate directly via gRPC.

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.26.5 |
| RPC Framework | gRPC |
| Serialization | Protocol Buffers |
| Code Generation | protoc, protoc-gen-go, protoc-gen-go-grpc |
| Database | PostgreSQL 16 |
| Workspace | Go workspaces (`go.work`) |
| Containerization | Docker, Docker Compose |

## Project Structure

```
ecommerce-microservices/
├── api/                     # Shared API contracts (protobuf definitions + generated stubs)
│   ├── proto/               # .proto source files
│   │   ├── auth.proto       # AuthService: CreateUser, Login
│   │   ├── order.proto      # OrderService: CreateOrder, GetOrder
│   │   ├── payment.proto    # PaymentService: ProcessPayment
│   │   └── stock.proto      # StockService: CheckStock, ReserveStock
│   └── gen/                 # Generated Go gRPC code
├── auth/                    # Auth microservice
│   ├── cmd/
│   ├── internal/
│   │   ├── db/              # PostgreSQL connection pool
│   │   └── handler/         # gRPC handler
│   └── Dockerfile
├── gateway/                 # HTTP API Gateway
│   ├── cmd/
│   ├── internal/
│   │   └── handler/         # HTTP handlers proxying to gRPC services
│   └── Dockerfile
├── orders/                  # Orders microservice
│   ├── cmd/
│   ├── internal/
│   │   └── handler/         # gRPC handler
│   └── Dockerfile
├── payments/                # Payments microservice (scaffold)
│   └── Dockerfile
├── stock/                   # Stock microservice (scaffold)
│   └── Dockerfile
├── shared/                  # Common utilities
│   ├── env.go               # GetEnvString(key, fallback)
│   └── json.go              # HTTP JSON helpers: WriteJSON, ReadJSON
├── scripts/                 # SQL init scripts, commit helper
│   ├── auth_init.sql
│   ├── orders_init.sql
│   ├── payments_init.sql
│   └── stock_init.sql
├── go.work                  # Go workspace file
├── go.mod                   # Root module (ecommerce-api)
├── Makefile                 # Protobuf generation targets
└── docker-compose.yml       # Container orchestration
```

## Services

### API Contract (`api/`)

Canonical source of truth for all service interfaces. Contains `.proto` definitions and generated Go stubs.

- **auth.proto** — `AuthService`: `CreateUser`, `Login`
- **order.proto** — `OrderService`: `CreateOrder`, `GetOrder`
- **payment.proto** — `PaymentService`: `ProcessPayment`
- **stock.proto** — `StockService`: `CheckStock`, `ReserveStock`

### Gateway (`gateway/`)

HTTP entrypoint that proxies requests to backend gRPC services.

- Listens on `:3000`
- `GET /api/v1/ping` — health check
- `POST /api/v1/orders` — creates an order via the Orders service

### Auth Service (`auth/`)

Handles user creation and authentication.

- Listens on `:5555`
- `CreateUser` — stub response with hardcoded values
- `Login` — defined in proto, handler not yet implemented
- Database: `auth_db` on `:6433`

### Orders Service (`orders/`)

Manages order creation and retrieval.

- Listens on `:4444`
- `CreateOrder` — returns a stub response (not yet connected to stock/payment flow)
- `GetOrder` — defined in proto but not yet implemented
- Database: `orders_db` on `:5433`

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
- Docker & Docker Compose (for running services with databases)

## Getting Started

### 1. Generate Protobuf Stubs

```bash
make gen
```

This generates Go gRPC stubs from `api/proto/*.proto` into `api/gen/`.

### 2. Run Services Locally

Run each module from its directory in separate terminals:

```bash
# Terminal 1 — Auth service
cd auth && go run ./cmd

# Terminal 2 — Orders service
cd orders && go run ./cmd

# Terminal 3 — Gateway
cd gateway && go run ./cmd
```

### 3. Run with Docker Compose

```bash
docker compose up --build
```

This starts:
- `orders-db` on `:5433`
- `auth-db` on `:6433`
- `orders` service on `:4444`
- `auth` service on `:5555`
- `gateway` on `:3000`

### 4. Test

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

The project uses a Go workspace (`go.work`) to allow cross-module imports. Run `go work sync` if modules change.

### Environment Variables

Services use `.env` files via `github.com/joho/godotenv/autoload`. Key variables:

| Variable | Default | Service |
|----------|---------|---------|
| `GATEWAY_PORT` | `3000` | gateway |
| `ORDERS_PORT` | `4444` | orders |
| `AUTH_PORT` | `5555` | auth |
| `ORDER_SERVICE_URL` | `localhost:4444` | gateway |
| `AUTH_SERVICE_URL` | `localhost:5555` | gateway |
| `AUTH_DB_CONN_STR` | — | auth |
| `ORDERS_DB_CONN_STR` | — | orders |

## Roadmap

- [x] Protobuf API contracts
- [x] Orders service (skeleton)
- [x] Auth service (skeleton)
- [x] Gateway with HTTP-to-gRPC translation
- [x] Docker / docker-compose setup
- [ ] Stock service implementation
- [ ] Payments service implementation
- [ ] Database integration (orders, auth)
- [ ] Service discovery / config
- [ ] TLS / production hardening
- [ ] Testing suite (unit + integration)
