# 💰 Personal Finance API

A clean, RESTful API built with Go and Gin for managing personal finances with multi-currency support and comprehensive reporting capabilities.

## 🚀 Features

- ✨ **Transaction Management**: Create, read, and delete financial transactions
- 💱 **Multi-Currency Support**: Handle transactions in different currencies (ARS, USD, EUR, etc.)
- 📊 **Monthly Reports**: Detailed financial reports with category breakdowns
- 🏗️ **Clean Architecture**: Separation of concerns with controllers, services, and repositories
- 🔄 **Repository Pattern**: Easy migration from in-memory to database storage
- 🌐 **RESTful API**: Clean, intuitive endpoints following REST conventions
- ⚡ **High Performance**: Built with Gin framework for exceptional speed
- 🔧 **Development Tools**: Comprehensive Makefile with useful commands
- 🛡️ **Input Validation**: Automatic JSON binding and validation
- 🌍 **CORS Support**: Cross-origin resource sharing enabled
- 📝 **Structured Logging**: Zap-powered logging across all layers with performance metrics
- 🔧 **Environment Configuration**: Flexible configuration via environment variables
- 🧪 **Comprehensive Testing**: Integration tests for all controllers with test helpers
- 🐳 **Docker Ready**: Complete Docker setup with multi-stage builds

## ⚡ Tech Stack

- **Backend**: Go 1.21+
- **Framework**: [Gin](https://gin-gonic.com/) - High-performance HTTP web framework
- **Testing**: [Testify](https://github.com/stretchr/testify) - Testing toolkit with assertions and test suites
- **Logging**: [Zap](https://github.com/uber-go/zap) - Blazing fast, structured, leveled logging
- **Architecture**: Clean Architecture with Repository pattern
- **Storage**: In-memory (easily extensible to PostgreSQL/MongoDB)
- **Configuration**: Environment variables with godotenv
- **Validation**: Built-in JSON binding and validation with Gin
- **CORS**: Custom middleware for cross-origin requests

## 📋 Prerequisites

- Go 1.21 or higher
- Git

## 🚀 Quick Start

### 1. Clone the repository
```bash
git clone https://github.com/maximicciullo/personal-finance-api.git
cd personal-finance-api
```

### 2. Install dependencies
```bash
make deps
```

### 3. Configure environment (optional)
```bash
# Create .env file for custom configuration
echo "PORT=8081
ENVIRONMENT=development
DEFAULT_CURRENCY=ARS" > .env
```

### 4. Run the application
```bash
make run
```

The API will be available at `http://localhost:8081` (or port 8080 if not configured)

## 📚 API Endpoints

### Health Check
```http
GET /health
```

### Transactions
```http
POST   /api/v1/transactions              # Create transaction
GET    /api/v1/transactions              # Get all transactions (with filters)
GET    /api/v1/transactions/:id          # Get transaction by ID
DELETE /api/v1/transactions/:id          # Delete transaction
```

### Reports
```http
GET /api/v1/reports/monthly/:year/:month # Get monthly report
GET /api/v1/reports/current-month        # Get current month report
```

## 💡 Usage Examples

### Create a Transaction
```bash
curl -X POST http://localhost:8081/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "type": "expense",
    "amount": 15000,
    "currency": "ARS",
    "description": "Lunch at restaurant",
    "category": "food",
    "date": "2024-06-19"
  }'
```

### Get Transactions with Filters
```bash
# Get all expense transactions
curl "http://localhost:8081/api/v1/transactions?type=expense"

# Get transactions by category
curl "http://localhost:8081/api/v1/transactions?category=food"

# Get transactions by date range
curl "http://localhost:8081/api/v1/transactions?from_date=2024-06-01&to_date=2024-06-30"
```

### Get Monthly Report
```bash
curl "http://localhost:8081/api/v1/reports/monthly/2024/6"
```

## 🏗️ Project Structure

```
personal-finance-api/
├── cmd/server/              # Application entrypoint
│   └── main.go
├── internal/                # Private application code
│   ├── config/             # Configuration management
│   ├── controllers/        # HTTP handlers (with Zap logging)
│   ├── middleware/         # HTTP middleware (CORS, Zap logging, etc.)
│   ├── models/             # Data models and business entities
│   ├── repositories/       # Data access layer (with Zap logging)
│   ├── services/           # Business logic (with Zap logging)
│   ├── test/               # Test utilities and helpers
│   └── utils/              # Utility functions
├── build/                  # Build artifacts (created by make build)
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
├── Dockerfile              # Docker image definition
├── docker-compose.yml      # Docker Compose configuration
├── .dockerignore          # Docker ignore file
├── .env                   # Environment variables (create manually)
├── Makefile               # Build and development commands
└── README.md              # Project documentation
```

## 🔧 Development Commands

```bash
# Run the application
make run

# Development mode with auto-reload
make dev

# Run all tests
make test

# Run integration tests
make test-integration

# Run tests with coverage
make test-coverage

# Run tests in watch mode
make test-watch

# Build the application
make build

# Format code
make fmt

# Lint code
make lint

# Security check
make security

# Build for multiple platforms
make build-all

# Clean build files
make clean

# Show all available commands
make help
```

## 📝 Transaction Model

```json
{
  "id": 1,
  "type": "expense",
  "amount": 15000,
  "currency": "ARS",
  "description": "Lunch at restaurant",
  "category": "food",
  "date": "2024-06-19T00:00:00Z",
  "created_at": "2024-06-19T10:30:00Z",
  "updated_at": "2024-06-19T10:30:00Z"
}
```

## 📊 Monthly Report Model

```json
{
  "month": "June",
  "year": 2024,
  "total_income": {
    "ARS": 100000,
    "USD": 500
  },
  "total_expense": {
    "ARS": 75000,
    "USD": 200
  },
  "balance": {
    "ARS": 25000,
    "USD": 300
  },
  "transactions": [...],
  "summary": {
    "transaction_count": 15,
    "income_count": 3,
    "expense_count": 12,
    "category_breakdown": {
      "food": {
        "count": 5,
        "totals": {
          "ARS": 25000
        }
      }
    }
  }
}
```

## ⚙️ Configuration

Create a `.env` file in the root directory to customize settings:

```env
# Server Configuration
PORT=8081                    # API server port (default: 8080)
ENVIRONMENT=development      # Environment: development|production
DEFAULT_CURRENCY=ARS        # Default currency for transactions

# Future Database Configuration (when migrating from in-memory)
# DB_HOST=localhost
# DB_PORT=5432
# DB_NAME=personal_finance
# DB_USER=your_user
# DB_PASSWORD=your_password
```

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Server port | `8080` | No |
| `ENVIRONMENT` | Runtime environment | `development` | No |
| `DEFAULT_CURRENCY` | Default currency for transactions | `ARS` | No |

### Port Configuration

If you encounter "address already in use" error:

```bash
# Option 1: Use environment variable
PORT=8081 make run

# Option 2: Create .env file
echo "PORT=8081" > .env
make run

# Option 3: Find and kill process using port 8080
lsof -i :8080                # Find process ID
kill -9 <PID>               # Kill the process
```

## 📝 Logging System

This project implements comprehensive structured logging using **Zap** across all architectural layers.

### Logging Architecture

The logging system follows the clean architecture pattern:

- **🎛️ Controllers**: HTTP request/response logging with client info and performance metrics
- **⚙️ Services**: Business logic operations, validation, and processing times
- **🗄️ Repository**: Data access operations, query performance, and transaction counts
- **📊 Middleware**: HTTP middleware logging with request/response bodies and headers

### Log Levels by Environment

#### Development Mode
- **Level**: Debug
- **Format**: Colorized console output
- **Features**: Request/response bodies, headers, debug information
- **Performance**: Detailed timing for all operations

#### Production Mode  
- **Level**: Info/Warn
- **Format**: JSON structured logs
- **Features**: Essential information only, security-focused
- **Performance**: Optimized for minimal overhead

### Log Categories

#### 📥 HTTP Requests
```bash
📥 HTTP Request | method=POST path=/api/v1/transactions client_ip=127.0.0.1
```

#### 📤 HTTP Responses  
```bash
📤 HTTP Response | ✅ 201 /api/v1/transactions latency=15ms size=256bytes
```

#### 🎛️ Controller Layer
```bash
🎛️ CreateTransaction started | client_ip=127.0.0.1 user_agent=curl/7.68.0
```

#### ⚙️ Service Layer
```bash
⚙️ CreateTransaction - validation passed | validation_duration=1ms
```

#### 🗄️ Repository Layer
```bash
🗄️ Create transaction completed successfully | transaction_id=1 duration=2ms
```

#### ⚡ Performance Metrics
```bash
⚡ CreateTransaction service call | duration=15ms transaction_id=1
```

#### ❌ Error Logging
```bash
❌ CreateTransaction - validation failed | layer=service error="amount must be positive"
```

### Using the Logger

The logger is automatically injected into all layers:

```go
// In Controllers
c.logger.Controller("Operation started", zap.String("param", value))

// In Services  
s.logger.Service("Processing data", zap.Int("count", len(data)))

// In Repositories
r.logger.Repository("Query executed", zap.Duration("duration", time.Since(start)))

// Performance logging
logger.Performance("Database query", duration, zap.Int("records", count))

// Error logging
logger.Error("controller", "Validation failed", err, zap.Any("request", req))
```

### Log Output Examples

#### Development Console Output
```bash
🚀 Server starting | port=8081 environment=development
📥 HTTP Request | method=POST path=/api/v1/transactions client_ip=127.0.0.1
🎛️ CreateTransaction started | client_ip=127.0.0.1
⚙️ CreateTransaction - validation passed | validation_duration=1ms  
🗄️ Create transaction completed | transaction_id=1 duration=2ms
⚡ CreateTransaction service call | duration=15ms success=true
📤 HTTP Response | ✅ 201 /api/v1/transactions latency=18ms size=256bytes
```

#### Production JSON Logs
```json
{
  "level": "info",
  "timestamp": "2024-06-19T10:30:15.123Z",
  "caller": "controllers/transaction_controller.go:45",
  "message": "🎛️ CreateTransaction started",
  "layer": "controller",
  "client_ip": "192.168.1.100",
  "user_agent": "PostmanRuntime/7.32.3"
}
```

## 🏗️ Architecture & Design Patterns

### Clean Architecture Layers

This project strictly follows Clean Architecture principles with proper dependency flow:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Controllers   │───▶│    Services     │───▶│  Repositories   │
│   (HTTP Layer)  │    │ (Business Logic)│    │  (Data Layer)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│    Gin HTTP     │    │     Domain      │    │   In-Memory     │
│   (Framework)   │    │    Models       │    │   Storage       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Dependency Injection

- **Controllers** ← depend on → **Service Interfaces**
- **Services** ← depend on → **Repository Interfaces**  
- **Repositories** ← implement → **Repository Interfaces**
- **Models** ← are used by → **All Layers**

### Layer Responsibilities

#### 🎛️ Controllers (`internal/controllers/`)
- HTTP request/response handling
- Input validation and binding
- Error response formatting
- Request logging and metrics

#### ⚙️ Services (`internal/services/`)
- Business logic implementation
- Data validation and transformation
- Orchestration between repositories
- Business rule enforcement

#### 🗄️ Repositories (`internal/repositories/`)
- Data access and persistence
- Query implementation
- Data mapping and conversion
- Storage-specific operations

#### 📊 Models (`internal/models/`)
- Domain entities and data structures
- Business value objects
- Request/response DTOs
- Constants and enums

## 🧪 Testing

This project includes comprehensive integration tests for all controllers, ensuring API reliability and correctness.

### Test Structure

```
internal/
├── controllers/
│   ├── health_controller_test.go        # Health endpoint tests
│   ├── transaction_controller_test.go   # Transaction CRUD tests  
│   └── report_controller_test.go        # Report generation tests
└── test/
    └── helper.go                        # Test utilities and helpers
```

### Running Tests

#### All Tests
```bash
# Run all tests
make test

# Run with verbose output
go test -v ./...
```

#### Test Categories
```bash
# Integration tests only (controllers)
make test-integration

# Unit tests only (services, repositories, models)
make test-unit

# Tests with coverage report
make test-coverage
open coverage.html

# Benchmark tests
make test-bench

# Watch mode (requires entr: brew install entr)
make test-watch
```

#### Test Examples
```bash
# Test specific controller
go test -v ./internal/controllers/ -run TestHealthController

# Test with race detection
go test -race ./...

# Test with timeout
go test -timeout 30s ./...
```

### Test Coverage

The tests cover:

- ✅ **Health Controller**: Status endpoint validation
- ✅ **Transaction Controller**: 
  - CRUD operations (Create, Read, Delete)
  - Input validation and error handling
  - Query parameter filtering
  - JSON response formatting
- ✅ **Report Controller**:
  - Monthly report generation
  - Current month reports
  - Multi-currency calculations
  - Category breakdowns
  - Edge cases (invalid dates, empty data)

### Manual Testing
```bash
# Test health endpoint
curl http://localhost:8081/health

# Create a test transaction
curl -X POST http://localhost:8081/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "type": "income",
    "amount": 50000,
    "currency": "ARS",
    "description": "Freelance payment",
    "category": "work"
  }'

# Get all transactions
curl http://localhost:8081/api/v1/transactions

# Get current month report
curl http://localhost:8081/api/v1/reports/current-month
```

### Test Features

- **Test Suites**: Organized using testify/suite for setup and teardown
- **Test Server**: Isolated test environment with in-memory storage
- **Helper Functions**: Utilities for JSON assertions and HTTP requests
- **Comprehensive Coverage**: All endpoints, error cases, and edge conditions
- **Integration Testing**: End-to-end API testing with real HTTP requests
- **Validation Testing**: Input validation, type checking, and constraint testing

## 🔄 Database Migration

The application currently uses in-memory storage, but it's designed to easily migrate to a database:

1. Implement the repository interfaces for your chosen database
2. Update the dependency injection in `main.go`
3. Add database configuration to the config package

## 🐳 Docker Support

### Using Docker Compose (Recommended)
```bash
# Build and run with docker-compose
docker-compose up --build

# Run in background
docker-compose up -d --build

# Stop services
docker-compose down
```

### Using Docker directly
```bash
# Build Docker image
make docker-build
# or manually:
docker build -t personal-finance-api:latest .

# Run Docker container
make docker-run
# or manually:
docker run -p 8080:8080 --name personal-finance-api personal-finance-api:latest

# Stop Docker container
make docker-stop
# or manually:
docker stop personal-finance-api && docker rm personal-finance-api
```

### Environment Variables in Docker
```bash
# Run with custom environment
docker run -p 8081:8081 \
  -e PORT=8081 \
  -e ENVIRONMENT=production \
  -e DEFAULT_CURRENCY=USD \
  --name personal-finance-api \
  personal-finance-api:latest
```

### Health Check
```bash
# Check if container is healthy
docker ps
# Look for "healthy" status

# Manual health check
curl http://localhost:8080/health
```

## 🛠️ Troubleshooting

### Common Issues

#### Port Already in Use
```bash
# Error: listen tcp :8080: bind: address already in use
# Solution 1: Use different port
PORT=8081 make run

# Solution 2: Kill process using port 8080
lsof -i :8080
kill -9 <PID>

# Solution 3: Use .env file
echo "PORT=8081" > .env
```

#### Module Issues
```bash
# Error: module not found
go mod download
go mod tidy

# Clean and rebuild
make clean
make deps
make build
```

#### Docker Build Issues
```bash
# Error: build context too large
# Make sure .dockerignore is properly configured

# Error: module not found during build
docker build --no-cache -t personal-finance-api .

# Check container logs
docker logs personal-finance-api
```

#### Docker Runtime Issues
```bash
# Container exits immediately
docker logs personal-finance-api

# Port binding issues
docker run -p 8081:8080 personal-finance-api  # Map to different host port

# Environment variable issues
docker run -e PORT=8080 personal-finance-api
```

### Development Tips

1. **Use development mode for auto-reload**:
   ```bash
   make dev  # Requires air: go install github.com/cosmtrek/air@latest
   ```

2. **Check logs for debugging**:
   ```bash
   # Gin provides detailed request logs in development mode
   # Look for [GIN] prefixed log entries
   ```

3. **Validate JSON payloads**:
   ```bash
   # Use tools like jq to validate JSON
   echo '{"type":"expense","amount":100}' | jq .
   ```

## 🚀 Deployment

### Local Build
```bash
# Build for current platform
make build

# Build for multiple platforms
make build-all
```

### Docker Deployment
```bash
# Development
docker-compose up --build

# Production
docker-compose -f docker-compose.yml up -d --build
```

### Cloud Deployment

#### Using Docker
```bash
# Build and tag for registry
docker build -t your-registry/personal-finance-api:v1.0.0 .
docker push your-registry/personal-finance-api:v1.0.0

# Deploy to your cloud platform
# (Example commands for different platforms)

# AWS ECS/Fargate
aws ecs update-service --cluster your-cluster --service personal-finance-api

# Google Cloud Run
gcloud run deploy personal-finance-api --image your-registry/personal-finance-api:v1.0.0

# Azure Container Instances
az container create --resource-group myResourceGroup --name personal-finance-api
```

### Environment variables for production
```env
PORT=8080
ENVIRONMENT=production
DEFAULT_CURRENCY=ARS

# Production optimizations
GIN_MODE=release
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

Maxi Micciullo - [@maximicciullo](https://github.com/maximicciullo)

Project Link: [https://github.com/maximicciullo/personal-finance-api](https://github.com/maximicciullo/personal-finance-api)

## 🙏 Acknowledgments

- Built with [Gin](https://gin-gonic.com/) - The fastest full-featured web framework for Go
- Inspired by clean architecture principles
- Following Go best practices and conventions

---

**Note**: This API uses port 8081 by default in examples to avoid conflicts. You can configure any port using the `PORT` environment variable.