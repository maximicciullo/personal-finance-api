package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/maximicciullo/personal-finance-api/internal/controllers"
	"github.com/maximicciullo/personal-finance-api/internal/repositories"
	"github.com/maximicciullo/personal-finance-api/internal/services"
	"github.com/stretchr/testify/assert"
)

// TestServer wraps the test dependencies
type TestServer struct {
	Router                *gin.Engine
	TransactionRepo       *repositories.MemoryTransactionRepository
	TransactionService    services.TransactionService
	TransactionController *controllers.TransactionController
	ReportService         services.ReportService
	ReportController      *controllers.ReportController
	HealthController      *controllers.HealthController
}

// NewTestServer creates a new test server with all dependencies
func NewTestServer() *TestServer {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Initialize repositories
	transactionRepo := repositories.NewMemoryTransactionRepository()

	// Initialize services
	transactionService := services.NewTransactionService(transactionRepo)
	reportService := services.NewReportService(transactionRepo)

	// Initialize controllers
	healthController := controllers.NewHealthController()
	transactionController := controllers.NewTransactionController(transactionService)
	reportController := controllers.NewReportController(reportService)

	// Setup router
	router := setupTestRoutes(healthController, transactionController, reportController)

	return &TestServer{
		Router:                router,
		TransactionRepo:       transactionRepo,
		TransactionService:    transactionService,
		TransactionController: transactionController,
		ReportService:         reportService,
		ReportController:      reportController,
		HealthController:      healthController,
	}
}

// setupTestRoutes configures routes for testing
func setupTestRoutes(
	healthController *controllers.HealthController,
	transactionController *controllers.TransactionController,
	reportController *controllers.ReportController,
) *gin.Engine {
	router := gin.New()

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

// MakeRequest performs an HTTP request and returns the response
func (ts *TestServer) MakeRequest(method, url string, body interface{}) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, url, reqBody)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	w := httptest.NewRecorder()
	ts.Router.ServeHTTP(w, req)

	return w
}

// AssertJSON checks if the response body contains expected JSON
func AssertJSON(t *testing.T, w *httptest.ResponseRecorder, expected interface{}) {
	var actual interface{}
	err := json.Unmarshal(w.Body.Bytes(), &actual)
	assert.NoError(t, err, "Response should be valid JSON")

	expectedJSON, _ := json.Marshal(expected)
	var expectedInterface interface{}
	json.Unmarshal(expectedJSON, &expectedInterface)

	assert.Equal(t, expectedInterface, actual)
}

// AssertJSONContains checks if response contains specific fields
func AssertJSONContains(t *testing.T, w *httptest.ResponseRecorder, expectedFields map[string]interface{}) {
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON")

	for key, expectedValue := range expectedFields {
		assert.Contains(t, response, key, "Response should contain field: %s", key)
		if expectedValue != nil {
			assert.Equal(t, expectedValue, response[key], "Field %s should have expected value", key)
		}
	}
}

// GetResponseJSON unmarshals response body to interface
func GetResponseJSON(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON")
	return response
}

// SafeGetMap safely extracts a map from the response, returning empty map if nil or wrong type
func SafeGetMap(t *testing.T, response map[string]interface{}, key string) map[string]interface{} {
	value, exists := response[key]
	if !exists {
		t.Logf("Key %s not found in response", key)
		return make(map[string]interface{})
	}
	
	if value == nil {
		return make(map[string]interface{})
	}
	
	mapValue, ok := value.(map[string]interface{})
	if !ok {
		t.Logf("Key %s is not a map, got type: %T", key, value)
		return make(map[string]interface{})
	}
	
	return mapValue
}

// SafeGetArray safely extracts an array from the response, returning empty array if nil or wrong type
func SafeGetArray(t *testing.T, response map[string]interface{}, key string) []interface{} {
	value, exists := response[key]
	if !exists {
		t.Logf("Key %s not found in response", key)
		return make([]interface{}, 0)
	}
	
	if value == nil {
		return make([]interface{}, 0)
	}
	
	arrayValue, ok := value.([]interface{})
	if !ok {
		t.Logf("Key %s is not an array, got type: %T", key, value)
		return make([]interface{}, 0)
	}
	
	return arrayValue
}