# LQ Studio Booking System

A photography studio booking backend built with Go, following Clean Architecture principles. Supports three package types with intelligent theme assignment and conflict prevention.

## Tech Stack

- **Language**: Go 1.24
- **Web Framework**: Echo v4
- **Database**: PostgreSQL 16
- **Migration Tool**: Goose
- **Query Generator**: sqlc (type-safe SQL)
- **Logging**: Zap
- **Documentation**: Swagger/OpenAPI

## Features

- ✅ **Three Package Types**
  - Package A: 30 minutes, single theme (user selects)
  - Package B: 60 minutes, sequential themes (auto-assigned)
  - Package C: 120 minutes, exclusive studio
- ✅ **Intelligent Availability Checking**
- ✅ **Conflict Prevention** with database transactions
- ✅ **Clean Architecture** with proper layer separation
- ✅ **Type-Safe Database Access** via sqlc
- ✅ **Swagger API Documentation**
- ✅ **Malaysia Timezone Support** (UTC+8)

## Project Structure

```
lqstudio-backend/
├── cmd/api/              # Application entry point
├── internal/
│   ├── domain/           # Business logic (entities, interfaces, services)
│   ├── application/      # Use cases & DTOs
│   ├── infrastructure/   # Database, HTTP, logging
│   └── config/           # Configuration
├── migrations/           # Database migrations
├── sql/                  # SQL queries for sqlc
└── docs/                 # Generated Swagger docs
```

## Quick Start

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- Make

### 1. Clone and Install Dependencies

```bash
cd lqstudio-backend
go mod download
```

### 2. Start Database

```bash
make docker-run
```

Wait for PostgreSQL to be healthy (check with `docker ps`).

### 3. Run Migrations

```bash
make migrate-up
```

This creates all tables and seeds themes/packages.

### 4. Generate sqlc Code (if needed)

```bash
make sqlc-generate
```

### 5. Run the Server

```bash
make run
```

The API will be available at `http://localhost:8080`.

## API Documentation

Once the server is running, access Swagger documentation at:

```
http://localhost:8080/swagger/index.html
```

## API Endpoints

### Public Endpoints

**Health Check**
```
GET /health
```

**Get Availability**
```
GET /api/v1/packages/{packageId}/availability?date=2026-01-25
```

**Create Booking**
```
POST /api/v1/bookings
Content-Type: application/json

{
  "user_name": "John Doe",
  "user_email": "john@example.com",
  "package_id": 1,
  "theme_id": 1,
  "booking_date": "2026-01-25",
  "start_time": "09:00"
}
```

**Get User Bookings**
```
GET /api/v1/users/{userId}/bookings
```

**Cancel Booking**
```
PATCH /api/v1/bookings/{bookingId}/cancel
```

## Makefile Commands

```bash
# Development
make run              # Run the application
make build            # Build binary
make watch            # Live reload (requires air)

# Database
make docker-run       # Start PostgreSQL container
make docker-down      # Stop PostgreSQL container
make migrate-up       # Run migrations
make migrate-down     # Rollback last migration
make migrate-create   # Create new migration

# Code Generation
make sqlc-generate    # Generate sqlc code
make swagger          # Generate Swagger docs

# Testing
make test             # Run all tests
make itest            # Run integration tests

# Utilities
make clean            # Remove binary
make setup            # Complete setup (docker + migrations + sqlc)
```

## Database Schema

### Core Tables

- **users** - Customer information
- **themes** - Photography themes (A, B)
- **packages** - Package definitions (A, B, C)
- **bookings** - Booking records
- **booking_theme_slots** - Theme slot assignments (source of truth for availability)

### Key Relationships

```
users 1 → * bookings
packages 1 → * bookings
bookings 1 → * booking_theme_slots
themes 1 → * booking_theme_slots
```

## Business Logic Highlights

### Package B Sequential Assignment

When booking Package B (60 minutes):
1. Checks if Theme A is free for first 30 minutes AND Theme B is free for second 30 minutes
2. If yes → assigns A then B
3. Otherwise, checks reversed order (B then A)
4. If neither works → booking rejected

### Conflict Detection

Uses overlap detection rule:
```
existing.start < requested.end AND existing.end > requested.start
```

All bookings are created with **SERIALIZABLE** transaction isolation to prevent double-booking.

## Configuration

Environment variables are loaded from `.env`:

```env
PORT=8080
APP_ENV=local
DB_HOST=psql_bp
DB_PORT=5432
DB_DATABASE=lqstudio
DB_USERNAME=lqstudio_user
DB_PASSWORD=lqstudio_password_123
TZ=Asia/Kuala_Lumpur
CORS_ALLOWED_ORIGINS=http://localhost:3000
```

## Development Workflow

1. **Make code changes**
2. **Run** `make sqlc-generate` if SQL queries changed
3. **Run** `make swagger` if API annotations changed
4. **Test** with `make test`
5. **Run** with `make watch` for live reload

## Example Usage

### 1. Check availability for Package A
```bash
curl "http://localhost:8080/api/v1/packages/1/availability?date=2026-01-25"
```

### 2. Book a session
```bash
curl -X POST http://localhost:8080/api/v1/bookings \
  -H "Content-Type: application/json" \
  -d '{
    "user_name": "Jane Smith",
    "user_email": "jane@example.com",
    "package_id": 2,
    "booking_date": "2026-01-25",
    "start_time": "10:00"
  }'
```

### 3. View user's bookings
```bash
curl "http://localhost:8080/api/v1/users/1/bookings"
```

## Testing

Run tests:
```bash
make test
```

Run integration tests (requires database):
```bash
make itest
```

## Architecture

This project follows **Clean Architecture** principles:

- **Domain Layer**: Pure business logic, no dependencies
- **Application Layer**: Use cases orchestrating business workflows
- **Infrastructure Layer**: External concerns (database, HTTP, logging)

Dependencies point inward: Infrastructure → Application → Domain

## License

MIT

## Support

For questions or issues, please contact support@lqstudio.com
