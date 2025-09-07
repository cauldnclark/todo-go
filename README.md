# Todo API Server

A Go-based REST API server for todo management with modern backend technologies.

## Tech Stack

- **Backend**: Go
- **Database**: PostgreSQL
- **Cache**: Redis
- **Message Queue**: Kafka
- **Authentication**: Google OAuth2
- **Frontend**: React + TypeScript

## Prerequisites

- Go 1.19+
- PostgreSQL
- Redis
- Node.js (for frontend)
- [Goose](https://github.com/pressly/goose) for database migrations

## Quick Start

### Database Setup

1. Install Goose globally:
   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```

2. Run migrations:
   ```bash
   goose -dir migrations postgres "host=localhost port=5432 user=jc password=pass1234 dbname=todo sslmode=disable" up
   ```

### Development

1. Clone the repository
2. Copy `.env.sample` to `.env` and configure your environment variables
3. Start the API server:
   ```bash
   go run cmd/api/main.go
   ```
4. Start the frontend development server:
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

## Project Structure

```
├── cmd/api/          # Application entry point
├── internal/         # Private application code
│   ├── config/       # Configuration management
│   ├── handlers/     # HTTP handlers
│   ├── middleware/   # HTTP middleware
│   ├── models/       # Data models
│   ├── repository/   # Data access layer
│   └── service/      # Business logic
├── migrations/       # Database migrations
└── frontend/         # React frontend application
```

## Creating New Migrations

```bash
goose create migration_name sql -dir migrations
```
