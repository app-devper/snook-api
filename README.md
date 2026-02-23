# Snook API

REST API for a snooker/billiards table management and point-of-sale system, built with Go and the Gin framework.

## Tech Stack

- **Go** 1.26
- **Gin** — HTTP framework
- **MongoDB** — primary database
- **Redis** — session/cache store
- **JWT** — authentication via `golang-jwt/jwt/v5`
- **Logrus** — structured logging

## Project Structure

```
snook-api/
├── main.go                  # Entry point
├── db/                      # Database initialization (MongoDB + Redis)
├── middlewares/              # Authentication, authorization, CORS, recovery
└── app/
    ├── init.go              # Router setup and feature registration
    ├── core/
    │   ├── constant/        # Role constants (SUPER, ADMIN, etc.)
    │   └── errcode/         # Error code definitions
    ├── data/
    │   ├── entities/        # MongoDB document models
    │   └── repositories/    # Data access interfaces and implementations
    ├── domain/
    │   ├── init.go          # Repository dependency injection
    │   └── request/         # Request DTOs
    └── featues/             # Feature modules (route + usecase per feature)
        ├── booking/
        ├── creditor/
        ├── dashboard/
        ├── expense/
        ├── menu/
        ├── payment/
        ├── promotion/
        ├── report/
        ├── setting/
        ├── table/
        ├── table_order/
        └── table_session/
```

## Prerequisites

- Go 1.26+
- MongoDB instance
- Redis instance

## Environment Variables

Create a `.env` file in the project root with the following variables:

| Variable              | Description                          | Example            |
|-----------------------|--------------------------------------|--------------------|
| `PORT`                | Server port                          | `8587`             |
| `MONGO_HOST`          | MongoDB connection URI               | `mongodb+srv://...`|
| `MONGO_SNOOK_DB_NAME` | MongoDB database name               | `devper_snook`     |
| `REDIS_HOST`          | Redis connection URI                 | `redis://...`      |
| `SECRET_KEY`          | JWT signing secret                   | `your-secret-key`  |
| `CLIENT_ID`           | Client identifier for JWT validation | `000`              |
| `SYSTEM`              | System identifier for JWT validation | `SNOOK`            |

## Getting Started

```bash
# Clone the repository
git clone <repository-url>
cd snook-api

# Configure environment
cp .env.example .env   # then edit .env with your values

# Install dependencies
go mod download

# Run the server
go run main.go
```

The server starts on the port specified in `.env` (default `8587`).

## API

**Base path:** `/api/snook/v1`

### Feature Endpoints

| Module           | Route Prefix         | Description                  |
|------------------|----------------------|------------------------------|
| Table            | `/tables`            | Table CRUD                   |
| Table Session    | `/table-sessions`    | Table session management     |
| Booking          | `/bookings`          | Booking management           |
| Menu             | `/menus`             | Menu categories and items    |
| Table Order      | `/table-orders`      | Order management per table   |
| Payment          | `/payments`          | Payment processing           |
| Creditor         | `/creditors`         | Creditor management          |
| Promotion        | `/promotions`        | Promotion management         |
| Expense          | `/expenses`          | Expense tracking             |
| Setting          | `/settings`          | System settings              |
| Dashboard        | `/dashboards`        | Dashboard analytics          |
| Report           | `/reports`           | Report generation            |

### Authentication & Authorization

All endpoints require a **JWT Bearer token** in the `Authorization` header:

```
Authorization: Bearer <token>
```

The token is validated against `SECRET_KEY`, `CLIENT_ID`, and `SYSTEM` from the environment. Session validation is performed via Redis.

Write operations are restricted by role-based authorization (`SUPER`, `ADMIN`).
