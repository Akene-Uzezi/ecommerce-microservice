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

## Request Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant G as Gateway
    participant A as Auth Service
    participant O as Orders Service
    participant ADB as Auth DB
    participant ODB as Orders DB

    C->>G: POST /api/v1/create_user
    G->>A: gRPC CreateUser
    A->>ADB: INSERT INTO users
    ADB-->>A: User created
    A-->>G: User response
    G-->>C: HTTP 201 JSON

    C->>G: POST /api/v1/orders
    G->>O: gRPC CreateOrder
    O->>ODB: INSERT INTO orders
    ODB-->>O: Order created
    O-->>G: Order response
    G-->>C: HTTP 201 JSON
```

## Data Layer

```mermaid
graph LR
    Auth[(auth_db)]
    Orders[(orders_db)]
    Payments[(payments_db)]
    Stock[(stock_db)]

    Auth -->|users table| AuthSchema[id, email, password, name]
    Orders -->|orders table| OrdersSchema[id, customer_id, items, status, total, created_at]
    Payments -->|payments table| PaymentSchema[id, order_id, customer_id, amount, status]
    Stock -->|stock table| StockSchema[product_id, quantity, reserved]
```

## Services

| Service | Port | Protocol | Database | Purpose |
|---------|------|----------|----------|---------|
| Gateway | 3000 | HTTP/JSON | — | Public API entrypoint |
| Auth | 5555 | gRPC | :6433 | User creation and authentication |
| Orders | 4444 | gRPC | :5433 | Order creation and retrieval |
| Payments | — | gRPC | — | Payment processing |
| Stock | — | gRPC | — | Inventory management |
