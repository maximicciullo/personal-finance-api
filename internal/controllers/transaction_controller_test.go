package controllers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/maximicciullo/personal-finance-api/internal/models"
	"github.com/maximicciullo/personal-finance-api/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TransactionControllerTestSuite struct {
	suite.Suite
	server *test.TestServer
}

func (suite *TransactionControllerTestSuite) SetupTest() {
	suite.server = test.NewTestServer()
}

// Test CreateTransaction
func (suite *TransactionControllerTestSuite) TestCreateTransaction_Success() {
	// Given
	request := models.CreateTransactionRequest{
		Type:        "expense",
		Amount:      15000,
		Currency:    "ARS",
		Description: "Lunch at restaurant",
		Category:    "food",
	}

	// When
	w := suite.server.MakeRequest("POST", "/api/v1/transactions", request)

	// Then
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	response := test.GetResponseJSON(suite.T(), w)
	assert.Equal(suite.T(), float64(1), response["id"])
	assert.Equal(suite.T(), "expense", response["type"])
	assert.Equal(suite.T(), float64(15000), response["amount"])
	assert.Equal(suite.T(), "ARS", response["currency"])
	assert.Equal(suite.T(), "Lunch at restaurant", response["description"])
	assert.Equal(suite.T(), "food", response["category"])
	assert.Contains(suite.T(), response, "created_at")
	assert.Contains(suite.T(), response, "updated_at")
}

func (suite *TransactionControllerTestSuite) TestCreateTransaction_WithDate() {
	// Given
	date := "2024-06-19"
	request := models.CreateTransactionRequest{
		Type:        "income",
		Amount:      50000,
		Currency:    "ARS",
		Description: "Freelance payment",
		Category:    "work",
		Date:        &date,
	}

	// When
	w := suite.server.MakeRequest("POST", "/api/v1/transactions", request)

	// Then
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	response := test.GetResponseJSON(suite.T(), w)
	assert.Contains(suite.T(), response["date"], "2024-06-19")
}

func (suite *TransactionControllerTestSuite) TestCreateTransaction_DefaultCurrency() {
	// Given
	request := models.CreateTransactionRequest{
		Type:        "expense",
		Amount:      100,
		Description: "Coffee",
		Category:    "food",
		// Currency not provided
	}

	// When
	w := suite.server.MakeRequest("POST", "/api/v1/transactions", request)

	// Then
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	response := test.GetResponseJSON(suite.T(), w)
	assert.Equal(suite.T(), "ARS", response["currency"])
}

func (suite *TransactionControllerTestSuite) TestCreateTransaction_ValidationErrors() {
	testCases := []struct {
		name    string
		request models.CreateTransactionRequest
		status  int
	}{
		{
			name: "missing type",
			request: models.CreateTransactionRequest{
				Amount:      100,
				Description: "Test",
				Category:    "test",
			},
			status: http.StatusBadRequest,
		},
		{
			name: "invalid type",
			request: models.CreateTransactionRequest{
				Type:        "invalid",
				Amount:      100,
				Description: "Test",
				Category:    "test",
			},
			status: http.StatusBadRequest,
		},
		{
			name: "zero amount",
			request: models.CreateTransactionRequest{
				Type:        "expense",
				Amount:      0,
				Description: "Test",
				Category:    "test",
			},
			status: http.StatusBadRequest,
		},
		{
			name: "negative amount",
			request: models.CreateTransactionRequest{
				Type:        "expense",
				Amount:      -100,
				Description: "Test",
				Category:    "test",
			},
			status: http.StatusBadRequest,
		},
		{
			name: "missing description",
			request: models.CreateTransactionRequest{
				Type:     "expense",
				Amount:   100,
				Category: "test",
			},
			status: http.StatusBadRequest,
		},
		{
			name: "missing category",
			request: models.CreateTransactionRequest{
				Type:        "expense",
				Amount:      100,
				Description: "Test",
			},
			status: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			w := suite.server.MakeRequest("POST", "/api/v1/transactions", tc.request)
			assert.Equal(t, tc.status, w.Code)

			response := test.GetResponseJSON(t, w)
			assert.Contains(t, response, "error")
			assert.Contains(t, response, "message")
		})
	}
}

func (suite *TransactionControllerTestSuite) TestCreateTransaction_InvalidDate() {
	// Given
	invalidDate := "invalid-date"
	request := models.CreateTransactionRequest{
		Type:        "expense",
		Amount:      100,
		Description: "Test",
		Category:    "test",
		Date:        &invalidDate,
	}

	// When
	w := suite.server.MakeRequest("POST", "/api/v1/transactions", request)

	// Then
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	response := test.GetResponseJSON(suite.T(), w)
	assert.Contains(suite.T(), response["message"], "date format")
}

// Test GetTransactions
func (suite *TransactionControllerTestSuite) TestGetTransactions_EmptyList() {
	// When
	w := suite.server.MakeRequest("GET", "/api/v1/transactions", nil)

	// Then
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var transactions []interface{}
	err := json.Unmarshal(w.Body.Bytes(), &transactions)
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), transactions)
}

func (suite *TransactionControllerTestSuite) TestGetTransactions_WithData() {
	// Given - create some transactions
	transactions := []models.CreateTransactionRequest{
		{
			Type:        "expense",
			Amount:      100,
			Currency:    "ARS",
			Description: "Coffee",
			Category:    "food",
		},
		{
			Type:        "income",
			Amount:      1000,
			Currency:    "USD",
			Description: "Freelance",
			Category:    "work",
		},
	}

	for _, req := range transactions {
		suite.server.MakeRequest("POST", "/api/v1/transactions", req)
	}

	// When
	w := suite.server.MakeRequest("GET", "/api/v1/transactions", nil)

	// Then
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), response, 2)
}

func (suite *TransactionControllerTestSuite) TestGetTransactions_WithFilters() {
	// Given - create transactions with different types
	transactions := []models.CreateTransactionRequest{
		{Type: "expense", Amount: 100, Description: "Coffee", Category: "food", Currency: "ARS"},
		{Type: "income", Amount: 1000, Description: "Salary", Category: "work", Currency: "ARS"},
		{Type: "expense", Amount: 50, Description: "Transport", Category: "transport", Currency: "USD"},
	}

	for _, req := range transactions {
		suite.server.MakeRequest("POST", "/api/v1/transactions", req)
	}

	testCases := []struct {
		name           string
		query          string
		expectedCount  int
		expectedType   string
		expectedCurrency string
	}{
		{
			name:          "filter by type expense",
			query:         "?type=expense",
			expectedCount: 2,
			expectedType:  "expense",
		},
		{
			name:          "filter by type income",
			query:         "?type=income",
			expectedCount: 1,
			expectedType:  "income",
		},
		{
			name:          "filter by category",
			query:         "?category=food",
			expectedCount: 1,
		},
		{
			name:             "filter by currency",
			query:            "?currency=USD",
			expectedCount:    1,
			expectedCurrency: "USD",
		},
		{
			name:          "multiple filters",
			query:         "?type=expense&currency=ARS",
			expectedCount: 1,
			expectedType:  "expense",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			w := suite.server.MakeRequest("GET", "/api/v1/transactions"+tc.query, nil)
			assert.Equal(t, http.StatusOK, w.Code)

			var response []map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response, tc.expectedCount)

			if tc.expectedType != "" && len(response) > 0 {
				assert.Equal(t, tc.expectedType, response[0]["type"])
			}

			if tc.expectedCurrency != "" && len(response) > 0 {
				assert.Equal(t, tc.expectedCurrency, response[0]["currency"])
			}
		})
	}
}

// Test GetTransaction
func (suite *TransactionControllerTestSuite) TestGetTransaction_Success() {
	// Given - create a transaction
	request := models.CreateTransactionRequest{
		Type:        "expense",
		Amount:      100,
		Description: "Test",
		Category:    "test",
	}

	createResponse := suite.server.MakeRequest("POST", "/api/v1/transactions", request)
	assert.Equal(suite.T(), http.StatusCreated, createResponse.Code)

	// When
	w := suite.server.MakeRequest("GET", "/api/v1/transactions/1", nil)

	// Then
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	response := test.GetResponseJSON(suite.T(), w)
	assert.Equal(suite.T(), float64(1), response["id"])
	assert.Equal(suite.T(), "expense", response["type"])
	assert.Equal(suite.T(), float64(100), response["amount"])
}

func (suite *TransactionControllerTestSuite) TestGetTransaction_NotFound() {
	// When
	w := suite.server.MakeRequest("GET", "/api/v1/transactions/999", nil)

	// Then
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	response := test.GetResponseJSON(suite.T(), w)
	assert.Equal(suite.T(), "Not Found", response["error"])
	assert.Contains(suite.T(), response["message"], "not found")
}

func (suite *TransactionControllerTestSuite) TestGetTransaction_InvalidID() {
	// When
	w := suite.server.MakeRequest("GET", "/api/v1/transactions/invalid", nil)

	// Then
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	response := test.GetResponseJSON(suite.T(), w)
	assert.Equal(suite.T(), "Bad Request", response["error"])
	assert.Contains(suite.T(), response["message"], "Invalid transaction ID")
}

// Test DeleteTransaction
func (suite *TransactionControllerTestSuite) TestDeleteTransaction_Success() {
	// Given - create a transaction
	request := models.CreateTransactionRequest{
		Type:        "expense",
		Amount:      100,
		Description: "Test",
		Category:    "test",
	}

	createResponse := suite.server.MakeRequest("POST", "/api/v1/transactions", request)
	assert.Equal(suite.T(), http.StatusCreated, createResponse.Code)

	// When
	w := suite.server.MakeRequest("DELETE", "/api/v1/transactions/1", nil)

	// Then
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	response := test.GetResponseJSON(suite.T(), w)
	assert.Contains(suite.T(), response["message"], "deleted successfully")

	// Verify transaction is deleted
	getResponse := suite.server.MakeRequest("GET", "/api/v1/transactions/1", nil)
	assert.Equal(suite.T(), http.StatusNotFound, getResponse.Code)
}

func (suite *TransactionControllerTestSuite) TestDeleteTransaction_NotFound() {
	// When
	w := suite.server.MakeRequest("DELETE", "/api/v1/transactions/999", nil)

	// Then
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	response := test.GetResponseJSON(suite.T(), w)
	assert.Equal(suite.T(), "Not Found", response["error"])
}

func (suite *TransactionControllerTestSuite) TestDeleteTransaction_InvalidID() {
	// When
	w := suite.server.MakeRequest("DELETE", "/api/v1/transactions/invalid", nil)

	// Then
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	response := test.GetResponseJSON(suite.T(), w)
	assert.Equal(suite.T(), "Bad Request", response["error"])
}

func TestTransactionControllerTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionControllerTestSuite))
}