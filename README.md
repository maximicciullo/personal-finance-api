# 💰 Personal Finance API

A clean, RESTful API built with Go and Gin for managing personal finances with multi-currency support and reporting.

## ⚡ Features

- Transaction management (create, read, delete)
- Multi-currency support (ARS, USD, EUR, etc.)
- Monthly financial reports with category breakdowns
- Clean Architecture with Repository pattern
- Structured logging with Zap
- Comprehensive testing
- Docker ready

## 🚀 Quick Start

```bash
# Clone and setup
git clone https://github.com/maximicciullo/personal-finance-api.git
cd personal-finance-api
make deps

# Run the application
make run
```

API available at `http://localhost:8080`

## 📚 API Endpoints

```http
GET    /health                              # Health check
POST   /api/v1/transactions                 # Create transaction
GET    /api/v1/transactions                 # Get transactions (with filters)
DELETE /api/v1/transactions/:id             # Delete transaction
GET    /api/v1/reports/monthly/:year/:month # Monthly report
```

## 💡 Usage Example

```bash
# Create transaction
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "type": "expense",
    "amount": 15000,
    "currency": "ARS",
    "description": "Lunch",
    "category": "food",
    "date": "2024-06-19"
  }'

# Get monthly report
curl "http://localhost:8080/api/v1/reports/monthly/2024/6"
```

## ⚙️ Configuration

Create `.env` file for custom settings:

```env
PORT=8081                    # Server port (default: 8080)
ENVIRONMENT=development      # Environment mode
DEFAULT_CURRENCY=ARS         # Default transaction currency
```

## 🔧 Development Commands

```bash
make run          # Run application
make test         # Run all tests
make build        # Build binary
make fmt          # Format code
make docker-build # Build Docker image
make docker-run   # Run with Docker
```

## 🏗️ Architecture

```
Controllers (HTTP) → Services (Business Logic) → Repositories (Data Access)
```

- **Tech Stack**: Go 1.21+, Gin, Zap logging, Testify
- **Storage**: In-memory (easily extensible to database)
- **Pattern**: Clean Architecture with dependency injection

## 🐳 Docker

```bash
# With Docker Compose (recommended)
docker-compose up --build

# Direct Docker
make docker-build
make docker-run
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
  "date": "2024-06-19T00:00:00Z"
}
```

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

Maxi Micciullo - [@maximicciullo](https://github.com/maximicciullo)