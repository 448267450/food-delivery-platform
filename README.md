# 🍔 Food Delivery Platform

A production-grade food delivery backend inspired by **Uber Eats**, built with Go and microservices architecture.

> This project is a long-term learning journey — progressively adopting technologies used at scale: gRPC, Kafka, Redis, OpenTelemetry, and Kubernetes.

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

| Service | Port | Description |
|---|---|---|
| user-service | 8081 | Registration, login, JWT auth |
| restaurant-service | 8082 | Restaurant and menu management |
| order-service | 8083 | Order lifecycle + state machine |

---

## Tech Stack

**Phase 1 (current)**
- **Go 1.22** — core language
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
- Go 1.22+

### Run locally

```bash
# Clone the repo
git clone https://github.com/448267450/food-delivery-platform.git
cd food-delivery-platform

# Start all services + PostgreSQL
cd deploy && docker-compose up --build
```

Services will be available at:
- User Service: http://localhost:8081
- Restaurant Service: http://localhost:8082
- Order Service: http://localhost:8083

---

## API Reference

### User Service

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

---

## Order State Machine

Orders follow a strict state transition flow:

```
PENDING ──▶ PAID ──▶ PREPARING ──▶ OUT_FOR_DELIVERY ──▶ DELIVERED
   │          │
   └──────────┴──▶ CANCELLED
```

Invalid transitions are rejected at the service layer.

---

## Project Roadmap

- [x] Phase 1: Single services with REST API
- [ ] Phase 2: Microservices with gRPC + Kafka
- [ ] Phase 3: Redis geolocation for driver matching
- [ ] Phase 4: Observability with OpenTelemetry
- [ ] Phase 5: K8s deployment with HPA

---

## Project Structure

```
food-delivery-platform/
├── services/
│   ├── user-service/
│   │   ├── cmd/main.go
│   │   ├── config/
│   │   └── internal/
│   │       ├── handler/    # HTTP layer
│   │       ├── service/    # Business logic
│   │       ├── repository/ # Database layer
│   │       └── model/      # Data models
│   ├── restaurant-service/
│   └── order-service/
├── pkg/                    # Shared utilities
├── deploy/                 # Docker, K8s configs
└── go.mod
```

---

## Contributing

PRs and issues are welcome. This is a learning project — feedback on Go idioms and architecture decisions is especially appreciated.

---

## License

MIT
