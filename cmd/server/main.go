package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/maximicciullo/personal-finance-api/internal/config"
	"github.com/maximicciullo/personal-finance-api/internal/controllers"
	"github.com/maximicciullo/personal-finance-api/internal/middleware"
	"github.com/maximicciullo/personal-finance-api/internal/repositories"
	"github.com/maximicciullo/personal-finance-api/internal/services"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	if err := middleware.InitLogger(cfg.Environment); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer middleware.Logger.Sync()

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize repositories
	transactionRepo := repositories.NewMemoryTransactionRepository()

	// Initialize services
	transactionService := services.NewTransactionService(transactionRepo)
	reportService := services.NewReportService(transactionRepo)

	// Initialize controllers
	healthController := controllers.NewHealthController()
	transactionController := controllers.NewTransactionController(transactionService)
	reportController := controllers.NewReportController(reportService)

	// Setup routes
	router := setupRoutes(cfg, healthController, transactionController, reportController)

	// Start server
	printStartupInfo(cfg)
	middleware.Logger.Info("🚀 Server starting",
		zap.String("port", cfg.Port),
		zap.String("environment", cfg.Environment),
	)

	log.Fatal(router.Run(":" + cfg.Port))
}

func setupRoutes(
	cfg *config.Config,
	healthController *controllers.HealthController,
	transactionController *controllers.TransactionController,
	reportController *controllers.ReportController,
) *gin.Engine {
	router := gin.Default()

	// Global middleware
	router.Use(gin.Recovery())

	// Logging middleware based on environment
	if cfg.Environment == "production" {
		router.Use(middleware.ProductionLogger())
	} else {
		router.Use(middleware.DevelopmentLogger())
	}

	// CORS middleware based on environment
	if cfg.Environment == "production" {
		// Production CORS - restrict origins
		corsConfig := middleware.ProductionCORSConfig([]string{
			"https://your-frontend-domain.com",
			"https://api.your-domain.com",
		})
		router.Use(middleware.CORSWithConfig(corsConfig))
	} else {
		// Development CORS - permissive
		router.Use(middleware.DevelopmentCORS())
	}

	// Health check endpoint
	router.GET("/health", healthController.HealthCheck)

	// API routes group
	api := router.Group("/api/v1")
	{
		// Transaction routes
		transactions := api.Group("/transactions")
		{
			transactions.POST("", transactionController.CreateTransaction)
			transactions.GET("", transactionController.GetTransactions)
			transactions.GET("/:id", transactionController.GetTransaction)
			transactions.DELETE("/:id", transactionController.DeleteTransaction)
		}

		// Report routes
		reports := api.Group("/reports")
		{
			reports.GET("/monthly/:year/:month", reportController.GetMonthlyReport)
			reports.GET("/current-month", reportController.GetCurrentMonthReport)
		}
	}

	return router
}

func printStartupInfo(cfg *config.Config) {
	fmt.Printf("\n🚀 Personal Finance API\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("🌐 Server starting on port: %s\n", cfg.Port)
	fmt.Printf("🏗️  Environment: %s\n", cfg.Environment)
	fmt.Printf("💰 Default currency: %s\n", cfg.DefaultCurrency)

	baseURL := fmt.Sprintf("http://localhost:%s", cfg.Port)
	fmt.Printf("🔗 Base URL: %s\n", baseURL)

	fmt.Printf("\n📚 Available Endpoints:\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// Health endpoint
	fmt.Printf("🔍 Health Check:\n")
	fmt.Printf("  GET    %s/health\n", baseURL)

	// Transaction endpoints
	fmt.Printf("\n💳 Transactions:\n")
	fmt.Printf("  POST   %s/api/v1/transactions\n", baseURL)
	fmt.Printf("  GET    %s/api/v1/transactions\n", baseURL)
	fmt.Printf("  GET    %s/api/v1/transactions/:id\n", baseURL)
	fmt.Printf("  DELETE %s/api/v1/transactions/:id\n", baseURL)

	// Report endpoints
	fmt.Printf("\n📊 Reports:\n")
	fmt.Printf("  GET    %s/api/v1/reports/monthly/:year/:month\n", baseURL)
	fmt.Printf("  GET    %s/api/v1/reports/current-month\n", baseURL)

	// Quick test commands
	fmt.Printf("\n🧪 Quick Test Commands:\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("# Health check\n")
	fmt.Printf("curl %s/health\n\n", baseURL)

	fmt.Printf("# Create transaction\n")
	fmt.Printf("curl -X POST %s/api/v1/transactions \\\n", baseURL)
	fmt.Printf("  -H 'Content-Type: application/json' \\\n")
	fmt.Printf("  -d '{\"type\":\"expense\",\"amount\":1500,\"currency\":\"ARS\",\"description\":\"Coffee\",\"category\":\"food\"}'\n\n")

	fmt.Printf("# Get current month report\n")
	fmt.Printf("curl %s/api/v1/reports/current-month\n\n", baseURL)

	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("🎯 Ready to handle requests!\n\n")
}
