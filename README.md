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

    Client -->|HTTP/JSON| Gateway
    Gateway -->|gRPC| Auth
    Gateway -->|gRPC| Orders
    Gateway -.->|planned| Payments
    Gateway -.->|planned| Stock

    Auth --> AuthDB
    Orders -.->|not yet wired| OrdersDB
```

The Gateway exposes a REST/JSON interface and translates requests into gRPC calls to backend services.

Solid lines are implemented today; dashed lines are planned. For a more detailed architecture description, see [architecture.md](architecture.md).

## Implementation Status

| Service | Status | Notes |
|---------|--------|-------|
| Gateway | Working | HTTP → gRPC proxy for auth and orders, with request logging |
| Auth | Working | `CreateUser` and `Login` (JWT) backed by PostgreSQL |
| Orders | Stub | `CreateOrder` returns hardcoded data; no DB access yet |
| Payments | Scaffold | `go.mod` only, no Go source, empty Dockerfile |
| Stock | Scaffold | `go.mod` only, no Go source, empty Dockerfile |

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.26.5 |
| RPC Framework | gRPC (google.golang.org/grpc v1.83.0) |
| Serialization | Protocol Buffers (v1.36.11) |
| Code Generation | protoc, protoc-gen-go, protoc-gen-go-grpc |
| Database | PostgreSQL 16 (alpine) |
| DB Driver | pgx/v5 with connection pooling |
| Password Hashing | bcrypt (golang.org/x/crypto) |
| Tokens | JWT HS256 (github.com/golang-jwt/jwt/v5) |
| Config Loading | github.com/joho/godotenv (autoload) |
| Workspace | Go workspaces (`go.work`) |
| Containerization | Docker, Docker Compose |
| Hot Reload | [air](https://github.com/air-verse/air) |

## Project Structure

```
ecommerce-microservices/
├── api/                     # Shared API contracts (protobuf definitions + generated stubs)
│   ├── proto/               # .proto source files
│   │   ├── auth.proto       # AuthService: CreateUser, Login
│   │   ├── order.proto      # OrderService: CreateOrder, GetOrder
│   │   ├── payment.proto    # PaymentService: ProcessPayment
│   │   └── stock.proto      # StockService: CheckStock, ReserveStock
│   └── gen/                 # Generated Go gRPC code (gitignored, run `make gen`)
│       ├── auth/  order/  payment/  stock/
├── auth/                    # Auth microservice
│   ├── cmd/main.go          # gRPC server entry point
│   ├── internal/
│   │   ├── db/
│   │   │   ├── db.go        # pgxpool initialization with ping retry
│   │   │   └── users.go     # UserModel: CreateUser, GetUserByEmail
│   │   ├── handler/
│   │   │   ├── grpc.go      # AuthGRPCHandler setup + DB pool injection
│   │   │   ├── users.go     # CreateUser + SearchUsersByEmail
│   │   │   ├── token.go     # Login + VerifyToken (JWT issuance/verification)
│   │   │   └── types.go     # Claims struct
│   │   └── util/
│   │       └── bcrypt.go    # HashPassword, ComparePassword
│   ├── Dockerfile
│   ├── .air.toml
│   └── .env.example
├── gateway/                 # HTTP API Gateway
│   ├── cmd/main.go          # HTTP server + gRPC client dialing
│   ├── internal/
│   │   ├── handler/
│   │   │   ├── handler.go       # ping + order routes (order route auth-protected)
│   │   │   ├── auth_handler.go  # create_user + login + search_users routes
│   │   │   ├── users.go         # search_users route handler
│   │   │   └── types.go         # Request payload types
│   │   └── middleware/
│   │       └── auth.go           # RequireAuth JWT verification middleware
│   ├── Dockerfile           # Multi-stage build
│   ├── .air.toml
│   └── .env.example
├── orders/                  # Orders microservice
│   ├── cmd/main.go          # gRPC server entry point
│   ├── internal/handler/
│   │   └── grpc.go          # CreateOrder returns a stub response
│   ├── Dockerfile
│   ├── .air.toml
│   └── .env.example
├── payments/                # Payments microservice (scaffold: go.mod + empty Dockerfile)
├── stock/                   # Stock microservice (scaffold: go.mod + empty Dockerfile)
├── shared/                  # Common utilities module
│   ├── env.go               # GetEnvString(key, fallback)
│   ├── json.go              # WriteJSON, ReadJSON
│   ├── error.go             # WriteErrorBadRequest, WriteErrorServerError
│   ├── status_log.go        # LogOK, LogBadRequest, LogInternalServerError, LogNotFound
│   └── go.mod
├── scripts/                 # SQL init scripts and helper scripts
│   ├── auth_init.sql        # users table + email index
│   ├── orders_init.sql      # empty placeholder
│   ├── payments_init.sql    # empty placeholder
│   ├── stock_init.sql       # empty placeholder
│   ├── commit.sh            # stage, prompt for message, commit, push
│   ├── logintodb.sh         # psql into a running container
│   └── reset.sh             # docker compose down -v, optional restart
├── go.work                  # Go workspace file
├── Makefile                 # Protobuf generation targets
├── architecture.md          # Detailed architecture diagrams
└── docker-compose.yml       # Container orchestration
```

## Services

### API Contract (`api/`)

Canonical source of truth for all service interfaces. Contains `.proto` definitions and generated Go stubs.

| Proto | Service | RPCs |
|-------|---------|------|
| `auth.proto` | `AuthService` | `CreateUser`, `Login` |
| `order.proto` | `OrderService` | `CreateOrder`, `GetOrder` |
| `payment.proto` | `PaymentService` | `ProcessPayment` |
| `stock.proto` | `StockService` | `CheckStock`, `ReserveStock` |

> `api/gen/` is gitignored. Run `make gen` before building anything, or compilation will fail.

### Gateway (`gateway/`)

HTTP entrypoint that proxies requests to backend gRPC services. gRPC clients are created at startup with insecure (plaintext) transport credentials.

Listens on `:3000`.

| Method | Path | Description | Success |
|--------|------|-------------|---------|
| `GET` | `/api/v1/ping` | Health check, returns `"pong"` | 200 |
| `POST` | `/api/v1/create_user` | Creates a user via the Auth service | 201 |
| `POST` | `/api/v1/login` | Authenticates and returns a JWT string | 200 |
| `POST` | `/api/v1/orders` | Creates an order via the Orders service | 201 |

Every handler emits a request log line via the `shared` logging helpers in the form `METHOD URI STATUS DURATION`.

### Auth Service (`auth/`)

Handles user creation and authentication. Listens on `:5555`.

**gRPC methods:**
- `CreateUser` — hashes the password with bcrypt and inserts the user, returning email and name.
- `Login` — looks up the user by email, compares the bcrypt hash, and returns a signed JWT (HS256, 72 hour expiry) containing `email` and `name` claims.

The `GetUserByEmail` DB helper exists on `UserModel` and is used internally by `Login`. It is no longer exposed as an RPC.

Database: `auth_db` on `:6433`, accessed via a `pgx/v5` pool (max 10 / min 1 connections) that retries the initial ping up to 10 times.

### Orders Service (`orders/`)

Listens on `:4444`.

- `CreateOrder` — currently returns a hardcoded stub response and does not touch the database.
- `GetOrder` — defined in the proto but not implemented.

`orders-db` runs on `:5433` in docker-compose, but the service does not connect to it yet and `scripts/orders_init.sql` is empty.

### Payments Service (`payments/`)

Module scaffold only. Intended to handle `ProcessPayment(order_id, customer_id, amount)`.

- `go.mod` present, no Go source files, empty Dockerfile, not in docker-compose.

### Stock Service (`stock/`)

Module scaffold only. Intended to handle `CheckStock(product_id, quantity)` and `ReserveStock(order_id, items)`.

- `go.mod` present, no Go source files, empty Dockerfile, not in docker-compose.

### Shared (`shared/`)

Common utilities imported by the other modules as `ecommerce-shared`:

- `env.go` — `GetEnvString(key, fallback)`
- `json.go` — `WriteJSON`, `ReadJSON`
- `error.go` — `WriteErrorBadRequest`, `WriteErrorServerError`
- `status_log.go` — `LogOK`, `LogBadRequest`, `LogInternalServerError`, `LogNotFound`

## Prerequisites

- Go 1.26.5
- Protocol Buffers compiler (`protoc`)
- `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
- `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`
- Docker & Docker Compose (for running services with databases)
- [air](https://github.com/air-verse/air) (optional, for hot reload)

## Getting Started

### 1. Generate protobuf stubs

```bash
make gen
```

Generates Go gRPC stubs from `api/proto/*.proto` into `api/gen/{auth,order,payment,stock}`. This is required — the generated code is not committed.

### 2. Configure environment

Copy each `.env.example` to `.env`:

```bash
cp gateway/.env.example gateway/.env
cp auth/.env.example auth/.env
cp orders/.env.example orders/.env
```

> `docker-compose.yml` reads `./auth/.env` and `./gateway/.env` via `env_file`, so these must exist before running Docker Compose.

### 3. Run services locally

Run each module from its own directory in separate terminals:

```bash
# Terminal 1 — Auth service (needs auth-db reachable on :6433)
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

### 4. Run with Docker Compose

```bash
docker compose up --build
```

This starts:

| Container | Port | Description |
|-----------|------|-------------|
| `orders-db` | 5433 | PostgreSQL 16, healthchecked |
| `auth-db` | 6433 | PostgreSQL 16, healthchecked, runs `auth_init.sql` |
| `orders` | 4444 | Orders gRPC service |
| `auth` | 5555 | Auth gRPC service |
| `gateway` | 3000 | Public HTTP API |

Inside Compose, service discovery uses container names (`auth:5555`, `orders:4444`, `auth-db:5432`) via overridden environment variables.

> The payments and stock services are not yet part of docker-compose.

### 5. Try it out

```bash
# Health check
curl http://localhost:3000/api/v1/ping

# Create a user
curl -X POST http://localhost:3000/api/v1/create_user \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}'

# Login (returns a JWT string)
curl -X POST http://localhost:3000/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Create an order (stub response)
curl -X POST http://localhost:3000/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"customer_id":"1","items":[{"product_id":"3","quantity":3,"price":10.62}]}'
```

## Development

### Make targets

| Command | Description |
|---------|-------------|
| `make gen` | Regenerate gRPC stubs from protobuf into `api/gen/` |
| `make clean` | Remove `api/gen/*` and `api/proto/*.pb.go` |

### Helper scripts

| Script | Description |
|--------|-------------|
| `scripts/commit.sh` | Stage all changes, prompt for a message, commit and push |
| `scripts/logintodb.sh` | Prompt for container/user/db and open `psql` inside the container |
| `scripts/reset.sh` | `docker compose down -v`, then optionally restart the stack |

### Workspace

The project uses a Go workspace (`go.work`) covering `api`, `auth`, `gateway`, `orders`, `payments`, `shared`, and `stock`. Run `go work sync` after changing module dependencies.

Because modules are separate, build and test from within a module directory (for example `cd auth && go build ./...`) rather than from the repository root.

### Environment variables

| Variable | Default | Service |
|----------|---------|---------|
| `GATEWAY_PORT` | `3000` | gateway |
| `ORDER_SERVICE_URL` | `localhost:4444` | gateway |
| `AUTH_SERVICE_URL` | `localhost:5555` | gateway |
| `ORDERS_PORT` | `4444` | orders |
| `AUTH_PORT` | `5555` | auth |
| `AUTH_DB_CONN_STR` | `postgres://auth:auth@localhost:6433/auth_db` | auth |
| `jwt_secret` | `this is a basic secret change in production` | auth |

All values are loaded through `godotenv/autoload` from the service's `.env` file, with the defaults above as fallbacks.

> `jwt_secret` is intentionally lowercase in the code and is not present in `auth/.env.example`. Set it explicitly before deploying anywhere real.

### Hot reload with air

Each service has an `.air.toml` that builds `./cmd/main.go` into `./tmp/main`. Run `air` from inside the service directory:

```bash
cd auth && air
```

## Known Gaps

- `Orders.CreateOrder` returns hardcoded values and ignores the request payload.
- `Orders.GetOrder` is declared in the proto but not implemented.
- Orders, payments, and stock init SQL scripts are empty.
- The Gateway does not verify JWTs; there is no auth middleware on protected routes yet.
- All gRPC connections use insecure transport credentials.
- `LoginRequest` in `gateway/internal/handler/types.go` is unused — the login handler decodes into `authpb.LoginRequest` directly.
- There are no tests in the repository.

## Roadmap

- [x] Protobuf API contracts
- [x] Gateway with HTTP-to-gRPC translation
- [x] Auth service with DB connection pool and bcrypt hashing
- [x] Login with JWT issuance
- [x] Request logging middleware helpers
- [x] Docker / docker-compose setup
- [ ] JWT verification middleware in the gateway
- [ ] Orders DB integration and real `CreateOrder` / `GetOrder`
- [ ] Stock service implementation
- [ ] Payments service implementation
- [ ] Order saga: orders → stock reservation → payment
- [ ] Service discovery / centralized config
- [ ] TLS / production hardening
- [ ] Testing suite (unit + integration)
