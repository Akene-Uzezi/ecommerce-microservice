# Architecture

## System Overview

```mermaid
graph TB
    Client[Client]
    Gateway[Gateway :3000]
    Auth[Auth Service :5555]
    Orders[Orders Service :4444]
    Payments[Payments Service :9000]
    Stock[Stock Service :8888]
    Products[Products Service :7777]
    AuthDB[(Auth DB :6433)]
    OrdersDB[(Orders DB :5433)]
    ProductsDB[(Products DB :7433)]
    PaymentsDB[(Payments DB :9433)]
    StockDB[(Stock DB :8433)]

    Client -->|HTTP/JSON| Gateway
    Gateway -->|gRPC| Auth
    Gateway -->|gRPC| Orders
    Gateway -->|gRPC| Products
    Gateway -->|gRPC| Payments
    Gateway -->|gRPC| Stock

    Auth --> AuthDB
    Orders --> OrdersDB
    Products --> ProductsDB
    Payments --> PaymentsDB
    Stock --> StockDB
```

All services are implemented and connected. The Gateway exposes HTTP/JSON endpoints and proxies requests to all backend gRPC services. Each service has its own dedicated PostgreSQL database.

## Request Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant G as Gateway
    participant A as Auth Service
    participant O as Orders Service
    participant P as Products Service
    participant PM as Payments Service
    participant S as Stock Service
    participant ADB as Auth DB
    participant ODB as Orders DB
    participant PDB as Products DB

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
    O->>ODB: INSERT INTO orders + order_items
    ODB-->>O: Order created
    O-->>G: Order response
    G-->>C: HTTP 201 JSON

    C->>G: POST /api/v1/payments (Bearer token)
    G->>PM: gRPC ProcessPayment
    PM-->>G: Mock payment response
    G-->>C: HTTP 200 JSON

    C->>G: GET /api/v1/stock?product_name=X (Bearer token)
    G->>S: gRPC CheckStock
    S->>PDB: SELECT quantity FROM products
    PDB-->>S: Quantity
    S-->>G: Stock response
    G-->>C: HTTP 200 JSON
```

## Data Layer

```mermaid
graph LR
    AuthDB[(auth_db :6433)]
    OrdersDB[(orders_db :5433)]
    ProductsDB[(products_db :7433)]
    PaymentsDB[(payments_db :9433)]
    StockDB[(stock_db :8433)]

    AuthDB -->|users table| AuthSchema[id, email, password, name]
    OrdersDB -->|orders + order_items tables| OrdersSchema[order + items]
    ProductsDB -->|products table| ProductsSchema[id, name, price, quantity]
    PaymentsDB -->|payments table| PaymentsSchema[id, order_id, payment_id, status]
    StockDB -->|products table| StockSchema[id, name, price, quantity]
```

- `auth_db` (`:6433`) — created by `scripts/auth_init.sql`, used by the Auth service
- `orders_db` (`:5433`) — created by `scripts/orders_init.sql`, used by the Orders service
- `products_db` (`:7433`) — created by `scripts/products_init.sql`, used by the Products service
- `payments_db` (`:9433`) — created by `scripts/payments_init.sql`, used by the Payments service
- `stock_db` (`:8433`) — created by `scripts/stock_init.sql`, used by the Stock service

## Services

| Service | Port | Protocol | Database | Purpose |
|---------|------|----------|----------|---------|
| Gateway | 3000 | HTTP/JSON | — | Public API entrypoint; gRPC client to all backend services; JWT auth middleware on protected routes |
| Auth | 5555 | gRPC | :6433 | User creation, authentication, token verification, user lookup |
| Orders | 4444 | gRPC | :5433 | Order creation with DB persistence, order retrieval, product stock check |
| Products | 7777 | gRPC | :7433 | Product catalog management with DB-backed handlers |
| Payments | 9000 | gRPC | :9433 | Mock payment processing with success/failure logic |
| Stock | 8888 | gRPC | :8433 | Inventory management with stock check and reservation |
