# 🍔 Food Delivery Platform

A production-grade food delivery backend inspired by **Uber Eats**, built with Go and microservices architecture.

> Long-term learning project — progressively adopting technologies used at scale: gRPC, Kafka, Redis, OpenTelemetry, and Kubernetes.

---

## Architecture

```
┌─────────────┐     ┌──────────────────┐     ┌─────────────────┐
│   Client    │────▶│   API Gateway    │────▶│  user-service   │ :8081
└─────────────┘     │   (Phase 2)      │     ├─────────────────┤
                    └──────────────────┘     │restaurant-svc   │ :8082
                                             ├─────────────────┤
                                             │  order-service  │ :8083
                                             └─────────────────┘
                                                      │
                                             ┌────────▼────────┐
                                             │   PostgreSQL    │
                                             └─────────────────┘
```

---

## Services

| Service | Port | Status | Description |
|---|---|---|---|
| user-service | 8081 | ✅ Live | Registration, login, JWT auth |
| restaurant-service | 8082 | ✅ Live | Restaurant and menu management |
| order-service | 8083 | ✅ Live | Order lifecycle + state machine |

---

## Tech Stack

**Phase 1 (current)**
- **Go 1.24** — core language
- **Gin** — HTTP framework
- **GORM** — ORM
- **PostgreSQL** — primary database
- **JWT** — authentication
- **Docker Compose** — local dev environment

**Phase 2 (planned)**
- gRPC — inter-service communication
- Kafka — async messaging
- Redis — caching + geolocation
- API Gateway — rate limiting, auth middleware

**Phase 3 (planned)**
- OpenTelemetry + Jaeger — distributed tracing
- Prometheus + Grafana — metrics
- Kubernetes — container orchestration

---

## Getting Started

### Prerequisites
- Docker & Docker Compose
- Go 1.24+

### Run locally

```bash
git clone https://github.com/448267450/food-delivery-platform.git
cd food-delivery-platform/deploy
docker-compose up --build
```

Services will be available at:
- User Service: http://localhost:8081
- Restaurant Service: http://localhost:8082
- Order Service: http://localhost:8083

---

## API Reference

### User Service `:8081`

```bash
# Health check
GET /health

# Register
POST /api/v1/auth/register
{
  "name": "Ryan Ren",
  "email": "ryan@example.com",
  "password": "123456",
  "phone": "512-000-0000"
}

# Login
POST /api/v1/auth/login
{
  "email": "ryan@example.com",
  "password": "123456"
}

# Get profile
GET /api/v1/users/:id/profile
```

### Restaurant Service `:8082`

```bash
# Health check
GET /health

# Create restaurant
POST /api/v1/restaurants
{
  "owner_id": 1,
  "name": "Ryan's Burger",
  "description": "Best burgers in Austin",
  "address": "123 Main St, Austin TX",
  "phone": "512-000-1234"
}

# List all restaurants
GET /api/v1/restaurants

# Get restaurant detail (with menu items)
GET /api/v1/restaurants/:id

# Update restaurant
PUT /api/v1/restaurants/:id
{ "name": "New Name", "is_open": false }

# Delete restaurant
DELETE /api/v1/restaurants/:id

# Add menu item
POST /api/v1/restaurants/:id/menu
{
  "name": "Classic Cheeseburger",
  "description": "Beef patty with cheddar",
  "price": 12.99,
  "category": "burger"
}

# Update menu item
PUT /api/v1/restaurants/:id/menu/:itemId
{ "price": 14.99, "is_available": false }

# Delete menu item
DELETE /api/v1/restaurants/:id/menu/:itemId
```

### Order Service `:8083`

```bash
# Health check
GET /health

# Create order
POST /api/v1/orders
{
  "user_id": 1,
  "restaurant_id": 1,
  "delivery_address": "456 Oak Ave, Austin TX",
  "note": "No onions please",
  "items": [
    {
      "menu_item_id": 1,
      "name": "Classic Cheeseburger",
      "price": 12.99,
      "quantity": 2
    }
  ]
}

# Get order by ID (with items)
GET /api/v1/orders/:id

# Get all orders for a user
GET /api/v1/orders/user/:userId

# Get all orders for a restaurant
GET /api/v1/orders/restaurant/:restaurantId

# Update order status (state machine)
PUT /api/v1/orders/:id/status
{ "status": "PAID" }

# Cancel order
DELETE /api/v1/orders/:id
```

---

## Order State Machine

Orders follow a strict state transition flow — invalid transitions are rejected at the service layer:

```
PENDING ──▶ PAID ──▶ PREPARING ──▶ OUT_FOR_DELIVERY ──▶ DELIVERED
   │          │
   └──────────┴──▶ CANCELLED
```

| From | Allowed Transitions |
|---|---|
| PENDING | PAID, CANCELLED |
| PAID | PREPARING, CANCELLED |
| PREPARING | OUT_FOR_DELIVERY |
| OUT_FOR_DELIVERY | DELIVERED |
| DELIVERED | — (terminal) |
| CANCELLED | — (terminal) |

---

## Project Structure

```
food-delivery-platform/
├── services/
│   ├── user-service/
│   │   ├── cmd/main.go
│   │   ├── config/
│   │   └── internal/
│   │       ├── handler/      # HTTP layer
│   │       ├── service/      # Business logic
│   │       ├── repository/   # Database layer
│   │       └── model/        # Data models
│   ├── restaurant-service/   # Same structure
│   └── order-service/        # Same structure + state machine
├── deploy/                   # Docker Compose, Dockerfiles
└── go.mod
```

---

## Project Roadmap

- [x] Phase 1a: user-service — register, login, JWT auth
- [x] Phase 1b: restaurant-service — CRUD + menu management
- [x] Phase 1c: order-service — order lifecycle + state machine
- [ ] Phase 2a: gRPC inter-service communication
- [ ] Phase 2b: Kafka async messaging (order created events)
- [ ] Phase 3a: Redis geolocation for driver matching
- [ ] Phase 3b: Rate limiting + API Gateway
- [ ] Phase 4: OpenTelemetry + Jaeger distributed tracing
- [ ] Phase 5: Kubernetes deployment with HPA

---

## Contributing

PRs and issues are welcome. This is a learning project — feedback on Go idioms and architecture decisions is especially appreciated.

---

## License

MIT