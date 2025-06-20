package controllers_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/maximicciullo/personal-finance-api/internal/models"
	"github.com/maximicciullo/personal-finance-api/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ReportControllerTestSuite struct {
	suite.Suite
	server *test.TestServer
}

func (suite *ReportControllerTestSuite) SetupTest() {
	suite.server = test.NewTestServer()
}

// Test GetMonthlyReport
func (suite *ReportControllerTestSuite) TestGetMonthlyReport_EmptyData() {
	// When
	w := suite.server.MakeRequest("GET", "/api/v1/reports/monthly/2024/6", nil)

	// Then
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	response := test.GetResponseJSON(suite.T(), w)
	assert.Equal(suite.T(), "June", response["month"])
	assert.Equal(suite.T(), float64(2024), response["year"])

	// Check empty totals using safe helpers
	totalIncome := test.SafeGetMap(suite.T(), response, "total_income")
	totalExpense := test.SafeGetMap(suite.T(), response, "total_expense")
	balance := test.SafeGetMap(suite.T(), response, "balance")

	assert.Empty(suite.T(), totalIncome)
	assert.Empty(suite.T(), totalExpense)
	assert.Empty(suite.T(), balance)

	// Check summary
	summary := test.SafeGetMap(suite.T(), response, "summary")
	assert.Equal(suite.T(), float64(0), summary["transaction_count"])
	assert.Equal(suite.T(), float64(0), summary["income_count"])
	assert.Equal(suite.T(), float64(0), summary["expense_count"])

	// Check transactions array
	transactions := test.SafeGetArray(suite.T(), response, "transactions")
	assert.Empty(suite.T(), transactions)
}

func (suite *ReportControllerTestSuite) TestGetMonthlyReport_WithData() {
	// Given - create transactions for June 2024
	transactions := []models.CreateTransactionRequest{
		{
			Type:        "income",
			Amount:      100000,
			Currency:    "ARS",
			Description: "Salary",
			Category:    "work",
			Date:        stringPtr("2024-06-15"),
		},
		{
			Type:        "expense",
			Amount:      25000,
			Currency:    "ARS",
			Description: "Groceries",
			Category:    "food",
			Date:        stringPtr("2024-06-16"),
		},
		{
			Type:        "expense",
			Amount:      200,
			Currency:    "USD",
			Description: "Netflix",
			Category:    "entertainment",
			Date:        stringPtr("2024-06-17"),
		},
		{
			Type:        "income",
			Amount:      500,
			Currency:    "USD",
			Description: "Freelance",
			Category:    "work",
			Date:        stringPtr("2024-06-18"),
		},
	}

	for _, req := range transactions {
		createResponse := suite.server.MakeRequest("POST", "/api/v1/transactions", req)
		assert.Equal(suite.T(), http.StatusCreated, createResponse.Code)
	}

	// When
	w := suite.server.MakeRequest("GET", "/api/v1/reports/monthly/2024/6", nil)

	// Then
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	response := test.GetResponseJSON(suite.T(), w)
	assert.Equal(suite.T(), "June", response["month"])
	assert.Equal(suite.T(), float64(2024), response["year"])

	// Check totals by currency using safe helpers
	totalIncome := test.SafeGetMap(suite.T(), response, "total_income")
	totalExpense := test.SafeGetMap(suite.T(), response, "total_expense")
	balance := test.SafeGetMap(suite.T(), response, "balance")

	assert.Equal(suite.T(), float64(100000), totalIncome["ARS"])
	assert.Equal(suite.T(), float64(500), totalIncome["USD"])

	assert.Equal(suite.T(), float64(25000), totalExpense["ARS"])
	assert.Equal(suite.T(), float64(200), totalExpense["USD"])

	assert.Equal(suite.T(), float64(75000), balance["ARS"]) // 100000 - 25000
	assert.Equal(suite.T(), float64(300), balance["USD"])   // 500 - 200

	// Check summary
	summary := test.SafeGetMap(suite.T(), response, "summary")
	assert.Equal(suite.T(), float64(4), summary["transaction_count"])
	assert.Equal(suite.T(), float64(2), summary["income_count"])
	assert.Equal(suite.T(), float64(2), summary["expense_count"])

	// Check category breakdown
	categoryBreakdown := test.SafeGetMap(suite.T(), summary, "category_breakdown")
	assert.Contains(suite.T(), categoryBreakdown, "work")
	assert.Contains(suite.T(), categoryBreakdown, "food")
	assert.Contains(suite.T(), categoryBreakdown, "entertainment")

	workCategory := test.SafeGetMap(suite.T(), categoryBreakdown, "work")
	assert.Equal(suite.T(), float64(2), workCategory["count"])

	workTotals := test.SafeGetMap(suite.T(), workCategory, "totals")
	assert.Equal(suite.T(), float64(100000), workTotals["ARS"])
	assert.Equal(suite.T(), float64(500), workTotals["USD"])

	// Check transactions array
	transactions_response := test.SafeGetArray(suite.T(), response, "transactions")
	assert.Len(suite.T(), transactions_response, 4)
}

func (suite *ReportControllerTestSuite) TestGetMonthlyReport_InvalidYear() {
	testCases := []struct {
		name string
		url  string
	}{
		{"invalid year format", "/api/v1/reports/monthly/invalid/6"},
		{"year too low", "/api/v1/reports/monthly/1800/6"},
		{"year too high", "/api/v1/reports/monthly/2040/6"},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			w := suite.server.MakeRequest("GET", tc.url, nil)
			assert.Equal(t, http.StatusBadRequest, w.Code)

			response := test.GetResponseJSON(t, w)
			assert.Equal(t, "Bad Request", response["error"])
		})
	}
}

func (suite *ReportControllerTestSuite) TestGetMonthlyReport_InvalidMonth() {
	testCases := []struct {
		name string
		url  string
	}{
		{"invalid month format", "/api/v1/reports/monthly/2024/invalid"},
		{"month too low", "/api/v1/reports/monthly/2024/0"},
		{"month too high", "/api/v1/reports/monthly/2024/13"},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			w := suite.server.MakeRequest("GET", tc.url, nil)
			assert.Equal(t, http.StatusBadRequest, w.Code)

			response := test.GetResponseJSON(t, w)
			assert.Equal(t, "Bad Request", response["error"])
		})
	}
}

func (suite *ReportControllerTestSuite) TestGetMonthlyReport_DifferentMonths() {
	// Given - create transactions in different months
	transactions := []models.CreateTransactionRequest{
		{
			Type:        "income",
			Amount:      1000,
			Currency:    "ARS",
			Description: "June income",
			Category:    "work",
			Date:        stringPtr("2024-06-15"),
		},
		{
			Type:        "expense",
			Amount:      500,
			Currency:    "ARS",
			Description: "July expense",
			Category:    "food",
			Date:        stringPtr("2024-07-15"),
		},
	}

	for _, req := range transactions {
		createResponse := suite.server.MakeRequest("POST", "/api/v1/transactions", req)
		assert.Equal(suite.T(), http.StatusCreated, createResponse.Code)
	}

	// When - get June report
	juneReport := suite.server.MakeRequest("GET", "/api/v1/reports/monthly/2024/6", nil)
	assert.Equal(suite.T(), http.StatusOK, juneReport.Code)

	juneResponse := test.GetResponseJSON(suite.T(), juneReport)
	juneSummary := test.SafeGetMap(suite.T(), juneResponse, "summary")
	assert.Equal(suite.T(), float64(1), juneSummary["transaction_count"])

	// When - get July report
	julyReport := suite.server.MakeRequest("GET", "/api/v1/reports/monthly/2024/7", nil)
	assert.Equal(suite.T(), http.StatusOK, julyReport.Code)

	julyResponse := test.GetResponseJSON(suite.T(), julyReport)
	julySummary := test.SafeGetMap(suite.T(), julyResponse, "summary")
	assert.Equal(suite.T(), float64(1), julySummary["transaction_count"])
}

// Test GetCurrentMonthReport
func (suite *ReportControllerTestSuite) TestGetCurrentMonthReport_EmptyData() {
	// When
	w := suite.server.MakeRequest("GET", "/api/v1/reports/current-month", nil)

	// Then
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	response := test.GetResponseJSON(suite.T(), w)
	
	// Should return current month and year
	now := time.Now()
	expectedMonth := now.Month().String()
	expectedYear := float64(now.Year())

	assert.Equal(suite.T(), expectedMonth, response["month"])
	assert.Equal(suite.T(), expectedYear, response["year"])

	// Check structure is correct
	assert.Contains(suite.T(), response, "total_income")
	assert.Contains(suite.T(), response, "total_expense")
	assert.Contains(suite.T(), response, "balance")
	assert.Contains(suite.T(), response, "transactions")
	assert.Contains(suite.T(), response, "summary")
}

func (suite *ReportControllerTestSuite) TestGetCurrentMonthReport_WithCurrentMonthData() {
	// Given - create transaction for current month
	now := time.Now()
	currentDate := now.Format("2006-01-02")

	request := models.CreateTransactionRequest{
		Type:        "income",
		Amount:      5000,
		Currency:    "ARS",
		Description: "Current month income",
		Category:    "work",
		Date:        &currentDate,
	}

	createResponse := suite.server.MakeRequest("POST", "/api/v1/transactions", request)
	assert.Equal(suite.T(), http.StatusCreated, createResponse.Code)

	// When
	w := suite.server.MakeRequest("GET", "/api/v1/reports/current-month", nil)

	// Then
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	response := test.GetResponseJSON(suite.T(), w)
	
	// Check that it includes the current month transaction
	summary := test.SafeGetMap(suite.T(), response, "summary")
	assert.Equal(suite.T(), float64(1), summary["transaction_count"])
	assert.Equal(suite.T(), float64(1), summary["income_count"])
	assert.Equal(suite.T(), float64(0), summary["expense_count"])

	totalIncome := test.SafeGetMap(suite.T(), response, "total_income")
	assert.Equal(suite.T(), float64(5000), totalIncome["ARS"])
}

func (suite *ReportControllerTestSuite) TestGetCurrentMonthReport_IgnoreOtherMonths() {
	// Given - create transactions for different months
	transactions := []models.CreateTransactionRequest{
		{
			Type:        "income",
			Amount:      1000,
			Currency:    "ARS",
			Description: "Last month income",
			Category:    "work",
			Date:        stringPtr("2024-05-15"), // Previous month
		},
		{
			Type:        "income",
			Amount:      2000,
			Currency:    "ARS",
			Description: "Current month income",
			Category:    "work",
			Date:        stringPtr(time.Now().Format("2006-01-02")), // Current month
		},
	}

	for _, req := range transactions {
		createResponse := suite.server.MakeRequest("POST", "/api/v1/transactions", req)
		assert.Equal(suite.T(), http.StatusCreated, createResponse.Code)
	}

	// When
	w := suite.server.MakeRequest("GET", "/api/v1/reports/current-month", nil)

	// Then
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	response := test.GetResponseJSON(suite.T(), w)
	
	// Should only include current month transaction
	summary := test.SafeGetMap(suite.T(), response, "summary")
	assert.Equal(suite.T(), float64(1), summary["transaction_count"])

	totalIncome := test.SafeGetMap(suite.T(), response, "total_income")
	assert.Equal(suite.T(), float64(2000), totalIncome["ARS"])
}

func (suite *ReportControllerTestSuite) TestGetCurrentMonthReport_ResponseFormat() {
	// When
	w := suite.server.MakeRequest("GET", "/api/v1/reports/current-month", nil)

	// Then
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	// Verify all required fields are present
	response := test.GetResponseJSON(suite.T(), w)
	requiredFields := []string{"month", "year", "total_income", "total_expense", "balance", "transactions", "summary"}
	
	for _, field := range requiredFields {
		assert.Contains(suite.T(), response, field, "Response should contain field: %s", field)
	}

	// Verify summary structure
	summary := test.SafeGetMap(suite.T(), response, "summary")
	summaryFields := []string{"transaction_count", "income_count", "expense_count", "category_breakdown"}
	
	for _, field := range summaryFields {
		assert.Contains(suite.T(), summary, field, "Summary should contain field: %s", field)
	}
}

// Test edge cases
func (suite *ReportControllerTestSuite) TestGetMonthlyReport_LeapYear() {
	// Test February in a leap year
	w := suite.server.MakeRequest("GET", "/api/v1/reports/monthly/2024/2", nil)
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	response := test.GetResponseJSON(suite.T(), w)
	assert.Equal(suite.T(), "February", response["month"])
	assert.Equal(suite.T(), float64(2024), response["year"])
}

func (suite *ReportControllerTestSuite) TestGetMonthlyReport_AllMonths() {
	// Test all months are valid
	for month := 1; month <= 12; month++ {
		url := fmt.Sprintf("/api/v1/reports/monthly/2024/%d", month)
		w := suite.server.MakeRequest("GET", url, nil)
		assert.Equal(suite.T(), http.StatusOK, w.Code, "Month %d should be valid", month)

		response := test.GetResponseJSON(suite.T(), w)
		assert.Equal(suite.T(), float64(2024), response["year"])
		assert.Contains(suite.T(), response, "month")
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}

func TestReportControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ReportControllerTestSuite))
}