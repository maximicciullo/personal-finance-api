# Personal Finance API

A simple and clean REST API built in Go for managing personal finances. Track your income and expenses with support for multiple currencies and generate monthly reports.

## 🚀 Features

- **Transaction Management**: Create, read, and delete financial transactions
- **Multi-Currency Support**: Handle transactions in different currencies (ARS, USD, EUR, etc.)
- **Monthly Reports**: Generate detailed financial reports with category breakdowns
- **Clean Architecture**: Organized codebase following best practices
- **In-Memory Storage**: Simple storage solution (easily extensible to PostgreSQL/MongoDB)
- **RESTful API**: Standard HTTP methods and status codes
- **CORS Support**: Ready for frontend integration

## 🏗️ Architecture

The project follows a clean architecture pattern with clear separation of concerns:

```
├── cmd/server/           # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── controllers/     # HTTP handlers
│   ├── models/          # Data structures and DTOs
│   ├── repositories/    # Data access layer
│   ├── services/        # Business logic
│   └── utils/           # Helper functions
└── pkg/middleware/      # HTTP middleware
```

### Architecture Layers

- **Controllers**: Handle HTTP requests and responses
- **Services**: Contain business logic and validation
- **Repositories**: Abstract data access (currently in-memory, designed for easy DB migration)
- **Models**: Define data structures and request/response DTOs
- **Middleware**: Cross-cutting concerns (CORS, logging, etc.)

## 📁 Project Structure

```
personal-finance-api/
├── cmd/
│   └── server/
│       └── main.go                    # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go                  # Configuration management
│   ├── controllers/
│   │   ├── health_controller.go       # Health check endpoint
│   │   ├── transaction_controller.go  # Transaction CRUD operations
│   │   └── report_controller.go       # Report generation
│   ├── models/
│   │   ├── transaction.go             # Transaction data structures
│   │   └── report.go                  # Report data structures
│   ├── repositories/
│   │   ├── interfaces.go              # Repository contracts
│   │   └── memory_transaction_repository.go  # In-memory implementation
│   ├── services/
│   │   ├── interfaces.go              # Service contracts
│   │   ├── transaction_service.go     # Transaction business logic
│   │   └── report_service.go          # Report business logic
│   └── utils/
│       ├── response.go                # HTTP response helpers
│       └── validator.go               # Validation utilities
├── pkg/
│   └── middleware/
│       └── cors.go                    # CORS middleware
├── go.mod                             # Go module definition
├── go.sum                             # Go dependencies checksum
├── .gitignore                         # Git ignore file
├── README.md                          # Project documentation
└── Makefile                           # Build automation
```

## 🛠️ Prerequisites

- Go 1.21 or higher
- Git

## 📥 Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/maximicciullo/personal-finance-api.git
   cd personal-finance-api
   ```

2. **Download dependencies**
   ```bash
   make deps
   ```

3. **Run the application**
   ```bash
   make run
   ```

The API will be available at `http://localhost:8080`

## 🔧 Development

### Quick Start
```bash
# Install development dependencies
make dev-deps

# Run in development mode with auto-reload
make dev

# Run tests
make test

# Build the application
make build
```

### Available Make Commands
```bash
make help           # Show all available commands
make deps           # Download dependencies
make run            # Run the application
make dev            # Run with auto-reload (requires air)
make test           # Run tests
make test-coverage  # Run tests with coverage report
make build          # Build binary
make clean          # Clean build files
make fmt            # Format code
make lint           # Lint code (requires golangci-lint)
```

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Health Check
```http
GET /health
```

### Transactions

#### Create Transaction
```http
POST /api/v1/transactions
Content-Type: application/json

{
  "type": "expense",
  "amount": 15000,
  "currency": "ARS",
  "description": "Lunch: pizza, water, and flan dessert",
  "category": "food",
  "date": "2024-01-15"
}
```

#### Get All Transactions
```http
GET /api/v1/transactions
```

Query parameters:
- `type`: Filter by transaction type (`expense` or `income`)
- `category`: Filter by category
- `currency`: Filter by currency
- `from_date`: Filter from date (YYYY-MM-DD)
- `to_date`: Filter to date (YYYY-MM-DD)

#### Get Single Transaction
```http
GET /api/v1/transactions/{id}
```

#### Delete Transaction
```http
DELETE /api/v1/transactions/{id}
```

### Reports

#### Get Monthly Report
```http
GET /api/v1/reports/monthly/{year}/{month}
```

#### Get Current Month Report
```http
GET /api/v1/reports/current-month
```

## 💰 Example Usage

### Creating Transactions

**Expense Example:**
```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "type": "expense",
    "amount": 15000,
    "currency": "ARS",
    "description": "Lunch: pizza, water, and flan dessert",
    "category": "food"
  }'
```

**Income Example:**
```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "type": "income",
    "amount": 5000,
    "currency": "USD",
    "description": "Deel salary payment",
    "category": "salary"
  }'
```

### Getting Reports

```bash
# Current month report
curl http://localhost:8080/api/v1/reports/current-month

# Specific month report
curl http://localhost:8080/api/v1/reports/monthly/2024/1
```

## 🗄️ Data Models

### Transaction
```json
{
  "id": 1,
  "type": "expense",
  "amount": 15000,
  "currency": "ARS",
  "description": "Lunch: pizza, water, and flan dessert",
  "category": "food",
  "date": "2024-01-15T00:00:00Z",
  "created_at": "2024-01-15T14:30:00Z",
  "updated_at": "2024-01-15T14:30:00Z"
}
```

### Monthly Report
```json
{
  "month": "January",
  "year": 2024,
  "total_income": {
    "USD": 5000,
    "ARS": 770000
  },
  "total_expense": {
    "ARS": 45000
  },
  "balance": {
    "USD": 5000,
    "ARS": 725000
  },
  "transactions": [...],
  "summary": {
    "transaction_count": 3,
    "income_count": 2,
    "expense_count": 1,
    "category_breakdown": {
      "food": {
        "count": 1,
        "totals": {"ARS": 15000}
      }
    }
  }
}
```

## 🔮 Future Enhancements

### Database Integration
The current implementation uses in-memory storage for simplicity. The repository pattern makes it easy to add database support:

- **PostgreSQL**: For relational data with ACID compliance
- **MongoDB**: For document-based storage with flexible schemas

### MCP Integration
This API is designed to work with Model Context Protocol (MCP) for Claude integration:

```
"I had lunch: pizza, water, and flan dessert for 15000 pesos"
→ POST /api/v1/transactions with parsed data

"I received my Deel salary of 5000 USD"
→ POST /api/v1/transactions with income data
```

### Additional Features
- User authentication and authorization
- Transaction categories management
- Recurring transactions
- Budget tracking and alerts
- Exchange rate integration
- Data export (CSV, PDF)
- Transaction attachments (receipts)

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

## 🐳 Docker Support

```bash
# Build Docker image
make docker-build

# Run in Docker
make docker-run
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

Maxi Micciullo - [@maximicciullo](https://github.com/maximicciullo)

Project Link: [https://github.com/maximicciullo/personal-finance-api](https://github.com/maximicciullo/personal-finance-api)