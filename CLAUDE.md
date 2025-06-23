# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Development Commands

### Build & Run
- `make run` - Run the application (downloads deps automatically)
- `make build` - Build the application binary to `build/` directory  
- `make clean` - Clean build artifacts
- `make deps` - Download and tidy Go dependencies

### Testing
- `make test` - Run all tests (controllers, services, repositories)
- `go test -v ./internal/controllers` - Run controller integration tests
- `go test -v ./internal/services` - Run service unit tests
- `go test -v ./internal/repositories` - Run repository unit tests

### Code Quality
- `make fmt` - Format Go code
- `go mod tidy` - Clean up Go modules

### Docker
- `make docker-build` - Build Docker image
- `make docker-run` - Run Docker container on port 8080
- `make docker-stop` - Stop and remove Docker container
- `docker-compose up --build` - Run with Docker Compose

## Architecture Overview

This is a Go REST API using Clean Architecture with strict dependency injection:

### Layer Structure
```
Controllers (HTTP) → Services (Business Logic) → Repositories (Data Access)
```

### Key Dependencies
- **Framework**: Gin (HTTP router/middleware)
- **Logging**: Zap (structured logging throughout all layers)
- **Testing**: Testify (test suites and assertions)
- **Config**: godotenv (environment variables)
- **Storage**: In-memory (designed for easy database migration)

### Dependency Injection Flow
1. `cmd/server/main.go` - Application entry point and DI container
2. Repositories created first (data layer)
3. Services injected with repository interfaces
4. Controllers injected with service interfaces
5. Routes configured with controllers

### Core Models
- `Transaction` - Main financial transaction entity
- `MonthlyReport` - Aggregated financial reporting data
- Request/response DTOs with Gin validation tags

### Repository Pattern
All data access goes through `TransactionRepository` interface in `internal/repositories/interfaces.go`. Current implementation is in-memory but easily swappable for database persistence.

### Service Interfaces
- `TransactionService` - CRUD operations and business logic
- `ReportService` - Financial reporting and calculations

### Configuration
Environment variables loaded via `internal/config/config.go`:
- `PORT` (default: 8080)
- `ENVIRONMENT` (development/production)
- `DEFAULT_CURRENCY` (default: ARS)

### Logging Architecture
Structured logging with Zap across all layers:
- Controllers: HTTP request/response logging with performance metrics
- Services: Business logic operations and validation
- Repositories: Data access operations with timing
- Different log levels for development (debug) vs production (info)

### API Structure
- Health check: `GET /health`
- Transactions: `POST|GET|DELETE /api/v1/transactions`
- Reports: `GET /api/v1/reports/monthly/:year/:month`

### Testing Strategy
Comprehensive integration tests for all controllers using testify suites with isolated test environments and in-memory storage.