# Ecommerce Microservices

A Go-based ecommerce microservices platform using gRPC for inter-service communication and an HTTP API Gateway as the public entrypoint. The project is organized as a Go workspace with independent modules per service.

## Architecture

```mermaid
graph TB
    Client[Client]
    Gateway[Gateway :3000]
    Auth[Auth Service :5555]
    Orders[Orders Service :4444]
    Payments[Payments Service]
    Stock[Stock Service]
    AuthDB[(Auth DB :6433)]
    OrdersDB[(Orders DB :5433)]
    PaymentDB[(Payments DB)]
    StockDB[(Stock DB)]

    Client -->|HTTP/JSON| Gateway
    Gateway -->|gRPC| Auth
    Gateway -->|gRPC| Orders
    Gateway -->|gRPC| Payments
    Gateway -->|gRPC| Stock

    Auth --> AuthDB
    Orders --> OrdersDB
    Payments --> PaymentDB
    Stock --> StockDB
```

The Gateway exposes a REST/JSON interface and translates requests into gRPC calls to backend services. Services communicate directly via gRPC.

For a more detailed architecture description, see [architecture.md](architecture.md).

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.26.5 |
| RPC Framework | gRPC |
| Serialization | Protocol Buffers |
| Code Generation | protoc, protoc-gen-go, protoc-gen-go-grpc |
| Database | PostgreSQL 16 |
| DB Driver | pgx/v5 with connection pooling |
| Password Hashing | bcrypt (golang.org/x/crypto) |
| Workspace | Go workspaces (`go.work`) |
| Containerization | Docker, Docker Compose |
| Hot Reload | [air](https://github.com/cosmtrek/air) |

## Project Structure

```
ecommerce-microservices/
├── api/                     # Shared API contracts (protobuf definitions + generated stubs)
│   ├── proto/               # .proto source files
│   │   ├── auth.proto       # AuthService: CreateUser, Login
│   │   ├── order.proto      # OrderService: CreateOrder, GetOrder
│   │   ├── payment.proto    # PaymentService: ProcessPayment
│   │   └── stock.proto      # StockService: CheckStock, ReserveStock
│   └── gen/                 # Generated Go gRPC code (empty until `make gen`)
├── auth/                    # Auth microservice
│   ├── cmd/
│   ├── internal/
│   │   ├── db/              # PostgreSQL connection pool (pgxpool)
│   │   │   ├── db.go        # Pool initialization with retry logic
│   │   │   └── users.go     # User model (CreateUser stub)
│   │   └── handler/         # gRPC handler
│   │       └── grpc.go      # CreateUser stub, Login not implemented
│   ├── Dockerfile
│   └── .air.toml
├── gateway/                 # HTTP API Gateway
│   ├── cmd/
│   ├── internal/
│   │   └── handler/         # HTTP handlers proxying to gRPC services
│   │       ├── handler.go   # Order routes
│   │       ├── auth_handler.go # Auth routes
│   │       └── types.go     # Request payload types
│   ├── Dockerfile           # Multi-stage build
│   └── .air.toml
├── orders/                  # Orders microservice
│   ├── cmd/
│   ├── internal/
│   │   └── handler/         # gRPC handler
│   │       └── grpc.go      # CreateOrder stub response, GetOrder not implemented
│   ├── Dockerfile
│   └── .air.toml
├── payments/                # Payments microservice (scaffold)
│   ├── go.mod
│   ├── internal/
│   └── Dockerfile           # Empty scaffold
├── stock/                   # Stock microservice (scaffold)
│   ├── go.mod
│   ├── internal/
│   └── Dockerfile           # Empty scaffold
├── shared/                  # Common utilities
│   ├── env.go               # GetEnvString(key, fallback)
│   ├── json.go              # HTTP JSON helpers: WriteJSON, ReadJSON
│   ├── error.go             # HTTP error helpers: WriteErrorBadRequest, WriteErrorServerError
│   ├── bcrypt.go            # Password hashing: HashPassword, ComparePassword
│   └── go.mod
├── scripts/                 # SQL init scripts and helper scripts
│   ├── auth_init.sql
│   ├── orders_init.sql
│   ├── payments_init.sql
│   ├── stock_init.sql
│   ├── commit.sh
│   ├── logintodb.sh
│   └── reset.sh
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

> Note: `api/gen/` is empty until protobuf stubs are generated with `make gen`.

### Gateway (`gateway/`)

HTTP entrypoint that proxies requests to backend gRPC services.

- Listens on `:3000`
- `GET /api/v1/ping` — health check
- `POST /api/v1/orders` — creates an order via the Orders service
- `POST /api/v1/create_user` — creates a user via the Auth service
- Environment: `gateway/.env.example` → copy to `gateway/.env`

### Auth Service (`auth/`)

Handles user creation and authentication.

- Listens on `:5555`
- `CreateUser` — handler stub (returns nil; DB pool initialized with retry logic)
- `Login` — not yet implemented
- Database: `auth_db` on `:6433` with `pgx/v5` connection pooling
- Environment: `auth/.env.example` → copy to `auth/.env`
- Uses `os.Getenv` for configuration (no godotenv autoload)

### Orders Service (`orders/`)

Manages order creation and retrieval.

- Listens on `:4444`
- `CreateOrder` — returns a stub response (not yet connected to stock/payment flow)
- `GetOrder` — defined in proto but not yet implemented
- Database: `orders_db` on `:5433` (init script exists, service not yet integrated)
- Environment: `orders/.env.example` → copy to `orders/.env`

### Payments Service (`payments/`)

Module scaffold only. Intended to handle `ProcessPayment(order_id, customer_id, amount)`.

- `go.mod` present
- No Go source files yet
- Dockerfile is empty

### Stock Service (`stock/`)

Module scaffold only. Intended to handle `CheckStock(product_id, quantity)` and `ReserveStock(order_id, items)`.

- `go.mod` present
- No Go source files yet
- Dockerfile is empty

### Shared (`shared/`)

Common utilities:

- `env.go` — `GetEnvString(key, fallback)`
- `json.go` — HTTP JSON helpers: `WriteJSON`, `ReadJSON`
- `error.go` — HTTP error helpers: `WriteErrorBadRequest`, `WriteErrorServerError`
- `bcrypt.go` — Password hashing: `HashPassword`, `ComparePassword`

## Prerequisites

- Go 1.26.5
- Protocol Buffers compiler (`protoc`)
- protoc-gen-go (`go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`)
- protoc-gen-go-grpc (`go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`)
- Docker & Docker Compose (for running services with databases)
- [air](https://github.com/cosmtrek/air) (optional, for hot reload during development)

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

Or use `air` for hot reloading:

```bash
cd auth && air
cd orders && air
cd gateway && air
```

### 3. Configure Environment

Copy `.env.example` to `.env` in each service directory:

```bash
cp gateway/.env.example gateway/.env
cp auth/.env.example auth/.env
cp orders/.env.example orders/.env
```

### 4. Run with Docker Compose

```bash
docker compose up --build
```

This starts:
- `orders-db` on `:5433`
- `auth-db` on `:6433`
- `orders` service on `:4444`
- `auth` service on `:5555`
- `gateway` service on `:3000`

> Note: The payments and stock services are not yet included in docker-compose.

### 5. Test

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

Each service reads configuration from environment variables. Copy the `.env.example` file in each service directory to `.env` and adjust values as needed.

| Variable | Default | Service | Loaded via |
|----------|---------|---------|------------|
| `GATEWAY_PORT` | `3000` | gateway | godotenv |
| `ORDERS_PORT` | `4444` | orders | godotenv |
| `AUTH_PORT` | `5555` | auth | os.Getenv |
| `ORDER_SERVICE_URL` | `localhost:4444` | gateway | godotenv |
| `AUTH_SERVICE_URL` | `localhost:5555` | gateway | os.Getenv |
| `AUTH_DB_CONN_STR` | `postgres://auth:auth@localhost:6433/auth_db` | auth | os.Getenv |

## Roadmap

- [x] Protobuf API contracts
- [x] Orders service (skeleton)
- [x] Auth service (skeleton) with DB connection pool
- [x] Gateway with HTTP-to-gRPC translation
- [x] Docker / docker-compose setup
- [ ] Stock service implementation
- [ ] Payments service implementation
- [ ] Orders DB integration
- [ ] Auth service full implementation (CreateUser, Login)
- [ ] Service discovery / config
- [ ] TLS / production hardening
- [ ] Testing suite (unit + integration)
