# Architecture

## System Overview

```mermaid
graph TB
    Client[Client]
    Gateway[Gateway :3000]
    Auth[Auth Service :5555]
    Orders[Orders Service :4444]
    Payments[Payments Service]
    Stock[Stock Service]
    Products[Products Service]
    AuthDB[(Auth DB :6433)]
    OrdersDB[(Orders DB :5433)]

    Client -->|HTTP/JSON| Gateway
    Gateway -->|gRPC| Auth
    Gateway -->|gRPC| Orders
    Gateway -.->|planned| Payments
    Gateway -.->|planned| Stock
    Gateway -.->|planned| Products

    Auth --> AuthDB
    Orders -.->|not yet wired| OrdersDB
```

Solid lines are implemented today; dashed lines are planned. The Payments and Stock services are scaffolds (no Go source, no containers). The Products service has DB models and a gRPC handler stub but no server entrypoint, Dockerfile, or DB wiring yet, and no container. Their databases do not exist yet. The Orders service runs but does not yet connect to `orders-db`.

## Request Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant G as Gateway
    participant A as Auth Service
    participant O as Orders Service
    participant ADB as Auth DB

    C->>G: POST /api/v1/create_user
    G->>A: gRPC CreateUser
    A->>ADB: INSERT INTO users
    ADB-->>A: User created
    A-->>G: User response
    G-->>C: HTTP 201 JSON

    C->>G: POST /api/v1/login
    G->>A: gRPC Login
    A->>ADB: SELECT user by email
    ADB-->>A: User row
    A-->>G: JWT token
    G-->>C: HTTP 200 {token, user}

    C->>G: POST /api/v1/orders (Bearer token)
    G->>A: gRPC VerifyToken
    A-->>G: valid = true
    G->>O: gRPC CreateOrder
    O-->>G: Stub order (no DB)
    G-->>C: HTTP 201 JSON
```

## Data Layer

```mermaid
graph LR
    Auth[(auth_db)]
    Orders[(orders_db)]

    Auth -->|users table| AuthSchema[id, email, password, name]
    Orders -->|orders table - unused| OrdersSchema[declared in orders_init.sql, not yet created by the service]
```

- `auth_db` (`:6433`) is created by `scripts/auth_init.sql` and used by the Auth service.
- `orders_db` (`:5433`) has a running, healthchecked container, but `scripts/orders_init.sql` is empty and the Orders service does not connect to it yet.
- `payments_db` and `stock_db` are planned and do not exist.

## Services

| Service | Port | Protocol | Database | Purpose |
|---------|------|----------|----------|---------|
| Gateway | 3000 | HTTP/JSON | — | Public API entrypoint; gRPC client to Auth/Orders; JWT auth middleware on protected routes |
| Auth | 5555 | gRPC | :6433 | User creation, authentication, token verification, user lookup |
| Orders | 4444 | gRPC | :5433 (unused) | Order creation (stub) and retrieval (unimplemented) |
| Payments | — | gRPC (planned) | — (planned) | Payment processing (scaffold) |
| Stock | — | gRPC (planned) | — (planned) | Inventory management (scaffold) |
| Products | — | gRPC (planned) | — (planned) | Product catalog management (db + handler stubs, no server/Dockerfile) |
