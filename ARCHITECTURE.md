# 🏗️ Personal Finance API - Architecture Documentation

## 🎯 Overview

Clean Architecture REST API built with Go for personal finance management with transaction tracking, multi-currency support, and financial reporting.

**Tech Stack**: Go 1.21+ • Gin • Zap • Docker • In-memory storage

## 🏛️ Architecture

### Clean Architecture Layers
```
┌─────────────────────────────────────────────────────────┐
│                    External Interfaces                  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │
│  │    HTTP     │  │   Future    │  │   Future    │    │
│  │  REST API   │  │  GraphQL    │  │    gRPC     │    │
│  └─────────────┘  └─────────────┘  └─────────────┘    │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│                    Controllers Layer                    │
│           (HTTP handlers, request/response)            │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│                    Services Layer                       │
│              (Business logic, validation)              │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│                   Repositories Layer                    │
│              (Data access abstraction)                 │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│                    Data Storage                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │
│  │  In-Memory  │  │  PostgreSQL │  │   MongoDB   │    │
│  │  (Current)  │  │  (Future)   │  │  (Future)   │    │
│  └─────────────┘  └─────────────┘  └─────────────┘    │
└─────────────────────────────────────────────────────────┘
```

### System Architecture
```
┌─────────────────────────────────────────────────────────────────┐
│                         Client Layer                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐ │
│  │   Web App   │  │   Mobile    │  │   Future MCP Server     │ │
│  │             │  │     App     │  │   (Claude Integration)  │ │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                                │
                                │ HTTP/REST
                                │
┌─────────────────────────────────────────────────────────────────┐
│                      API Gateway Layer                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐ │
│  │   CORS      │  │   Logging   │  │     Authentication      │ │
│  │ Middleware  │  │ Middleware  │  │     (Future)            │ │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────────┐
│                     Personal Finance API                        │
│                                                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐ │
│  │   Health    │  │Transaction  │  │       Report            │ │
│  │ Controller  │  │ Controller  │  │     Controller          │ │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘ │
│                                │                                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐ │
│  │   Config    │  │Transaction  │  │       Report            │ │
│  │  Service    │  │  Service    │  │      Service            │ │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘ │
│                                │                                │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │              Transaction Repository                         │ │
│  │                   (Interface)                              │ │
│  └─────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────────┐
│                      Storage Layer                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐ │
│  │  In-Memory  │  │  PostgreSQL │  │      File System        │ │
│  │  (Current)  │  │  (Future)   │  │      (Backup)           │ │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

**Dependency Flow**:
```
main.go → Controllers → Services → Repositories → Storage
```

## 🔄 Core Components

### 1. Controllers (`internal/controllers/`)
- HTTP request/response handling
- Input validation and parsing
- Error formatting

### 2. Services (`internal/services/`)
- Business logic and validation
- Transaction orchestration
- Report generation

### 3. Repositories (`internal/repositories/`)
- Data access abstraction
- CRUD operations
- Thread-safe in-memory storage

### 4. Models (`internal/models/`)
```go
type Transaction struct {
    ID          int       `json:"id"`
    Type        string    `json:"type"`        // "expense" | "income"
    Amount      float64   `json:"amount"`
    Currency    string    `json:"currency"`    // "ARS", "USD", "EUR"
    Description string    `json:"description"`
    Category    string    `json:"category"`
    Date        time.Time `json:"date"`
}

type MonthlyReport struct {
    Year           int                 `json:"year"`
    Month          int                 `json:"month"`
    TotalIncome    float64            `json:"total_income"`
    TotalExpenses  float64            `json:"total_expenses"`
    NetAmount      float64            `json:"net_amount"`
    Categories     map[string]float64 `json:"categories"`
}
```

## 📊 Data Flow

### Request Processing
```
HTTP Request → Middleware (CORS, Logging) → Controller → Service → Repository → Storage
```

### Example: Create Transaction
```
POST /api/v1/transactions
└── TransactionController.CreateTransaction()
    ├── Parse & validate JSON
    └── TransactionService.CreateTransaction()
        ├── Apply business rules
        ├── Set defaults (currency, date)
        └── TransactionRepository.Create()
            └── Store in memory with ID generation
```

## 🔗 Dependencies

```go
"github.com/gin-gonic/gin"        // HTTP framework
"go.uber.org/zap"                 // Structured logging
"github.com/joho/godotenv"        // Environment config
"github.com/stretchr/testify"     // Testing
```

## 💾 Storage Strategy

**Current**: Thread-safe in-memory storage
- Fast operations, data persistence across requests
- Automatic ID generation and indexing

**Future**: Database migration ready
- Interface-based design allows easy swapping
- PostgreSQL, MongoDB, or SQLite options

## 🧪 Testing

- **Unit Tests**: Services and repositories
- **Integration Tests**: Full HTTP stack testing
- **Test Isolation**: Each test runs independently

## 🚀 Deployment

**Docker**: Multi-stage build with health checks
```yaml
services:
  personal-finance-api:
    build: .
    ports: ["8080:8080"]
    environment:
      - ENVIRONMENT=production
```

**Future**: Kubernetes deployment with auto-scaling

## 📝 API Endpoints

```
GET    /health                              # Health check
POST   /api/v1/transactions                 # Create transaction
GET    /api/v1/transactions                 # List with filters
DELETE /api/v1/transactions/:id             # Delete transaction
GET    /api/v1/reports/monthly/:year/:month # Monthly report
GET    /api/v1/reports/current-month        # Current month
```

## 🔮 Future Enhancements

### MCP Integration
```
Claude LLM ↔ MCP Server ↔ Personal Finance API
```

### Technical Roadmap
- **Authentication**: JWT/OAuth2
- **Database**: PostgreSQL migration
- **Caching**: Redis integration
- **Monitoring**: Metrics and observability
- **Multi-tenancy**: User account isolation

---

**Related**: [`README.md`](README.md) • [`CLAUDE.md`](CLAUDE.md) • [`TODO.md`](TODO.md)