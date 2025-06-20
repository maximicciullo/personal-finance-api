package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/maximicciullo/personal-finance-api/internal/config"
	"github.com/maximicciullo/personal-finance-api/internal/controllers"
	"github.com/maximicciullo/personal-finance-api/internal/repositories"
	"github.com/maximicciullo/personal-finance-api/internal/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

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
	router := setupRoutes(healthController, transactionController, reportController)

	// Start server
	fmt.Printf("🚀 Personal Finance API starting on port %s\n", cfg.Port)
	fmt.Println("📊 Available endpoints:")
	fmt.Println("  GET    /health")
	fmt.Println("  POST   /api/v1/transactions")
	fmt.Println("  GET    /api/v1/transactions")
	fmt.Println("  GET    /api/v1/transactions/:id")
	fmt.Println("  DELETE /api/v1/transactions/:id")
	fmt.Println("  GET    /api/v1/reports/monthly/:year/:month")
	fmt.Println("  GET    /api/v1/reports/current-month")

	log.Fatal(router.Run(":" + cfg.Port))
}

func setupRoutes(
	healthController *controllers.HealthController,
	transactionController *controllers.TransactionController,
	reportController *controllers.ReportController,
) *gin.Engine {
	router := gin.Default()

	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Health check
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

func corsMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}
