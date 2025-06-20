package services_test

import (
	"testing"
	"time"

	"github.com/maximicciullo/personal-finance-api/internal/middleware"
	"github.com/maximicciullo/personal-finance-api/internal/models"
	"github.com/maximicciullo/personal-finance-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// ReportServiceTestSuite is the test suite for ReportService
type ReportServiceTestSuite struct {
	suite.Suite
	mockRepo *MockTransactionRepository
	service  services.ReportService
}

func (suite *ReportServiceTestSuite) SetupTest() {
	// Initialize logger for testing
	middleware.InitLogger("test")
	
	suite.mockRepo = new(MockTransactionRepository)
	suite.service = services.NewReportService(suite.mockRepo)
}

func (suite *ReportServiceTestSuite) TearDownTest() {
	suite.mockRepo.AssertExpectations(suite.T())
}

// Test GetMonthlyReport
func (suite *ReportServiceTestSuite) TestGetMonthlyReport_Success() {
	// Given
	year := 2024
	month := 6
	startDate := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second)

	transactions := []models.Transaction{
		{
			ID:       1,
			Type:     "income",
			Amount:   50000,
			Currency: "ARS",
			Category: "salary",
			Date:     time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:       2,
			Type:     "expense",
			Amount:   15000,
			Currency: "ARS",
			Category: "food",
			Date:     time.Date(2024, 6, 16, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:       3,
			Type:     "expense",
			Amount:   200,
			Currency: "USD",
			Category: "entertainment",
			Date:     time.Date(2024, 6, 17, 0, 0, 0, 0, time.UTC),
		},
	}

	suite.mockRepo.On("GetByDateRange", mock.MatchedBy(func(start time.Time) bool {
		return start.Equal(startDate)
	}), mock.MatchedBy(func(end time.Time) bool {
		return end.Equal(endDate)
	})).Return(transactions, nil)

	// When
	result, err := suite.service.GetMonthlyReport(year, month)

	// Then
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	
	// Check basic info
	assert.Equal(suite.T(), "June", result.Month)
	assert.Equal(suite.T(), 2024, result.Year)
	
	// Check totals by currency
	assert.Equal(suite.T(), 50000.0, result.TotalIncome["ARS"])
	assert.Equal(suite.T(), 15000.0, result.TotalExpense["ARS"])
	assert.Equal(suite.T(), 200.0, result.TotalExpense["USD"])
	
	// Check balances
	assert.Equal(suite.T(), 35000.0, result.Balance["ARS"]) // 50000 - 15000
	assert.Equal(suite.T(), -200.0, result.Balance["USD"])  // 0 - 200
	
	// Check summary
	assert.Equal(suite.T(), 3, result.Summary.TransactionCount)
	assert.Equal(suite.T(), 1, result.Summary.IncomeCount)
	assert.Equal(suite.T(), 2, result.Summary.ExpenseCount)
	
	// Check category breakdown
	assert.Contains(suite.T(), result.Summary.CategoryBreakdown, "salary")
	assert.Contains(suite.T(), result.Summary.CategoryBreakdown, "food")
	assert.Contains(suite.T(), result.Summary.CategoryBreakdown, "entertainment")
	
	salaryCategory := result.Summary.CategoryBreakdown["salary"]
	assert.Equal(suite.T(), 1, salaryCategory.Count)
	assert.Equal(suite.T(), 50000.0, salaryCategory.Totals["ARS"])
	
	// Check transactions are included
	assert.Len(suite.T(), result.Transactions, 3)
}

func (suite *ReportServiceTestSuite) TestGetMonthlyReport_EmptyData() {
	// Given
	year := 2024
	month := 12
	startDate := time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second)

	emptyTransactions := []models.Transaction{}

	suite.mockRepo.On("GetByDateRange", mock.MatchedBy(func(start time.Time) bool {
		return start.Equal(startDate)
	}), mock.MatchedBy(func(end time.Time) bool {
		return end.Equal(endDate)
	})).Return(emptyTransactions, nil)

	// When
	result, err := suite.service.GetMonthlyReport(year, month)

	// Then
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	
	assert.Equal(suite.T(), "December", result.Month)
	assert.Equal(suite.T(), 2024, result.Year)
	
	// Should have empty maps
	assert.Empty(suite.T(), result.TotalIncome)
	assert.Empty(suite.T(), result.TotalExpense)
	assert.Empty(suite.T(), result.Balance)
	
	// Summary should be zero
	assert.Equal(suite.T(), 0, result.Summary.TransactionCount)
	assert.Equal(suite.T(), 0, result.Summary.IncomeCount)
	assert.Equal(suite.T(), 0, result.Summary.ExpenseCount)
	
	// Should have empty transactions
	assert.Empty(suite.T(), result.Transactions)
}

func (suite *ReportServiceTestSuite) TestGetMonthlyReport_InvalidYear() {
	testCases := []struct {
		name string
		year int
	}{
		{"year too low", 1800},
		{"year too high", 2040},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// When
			result, err := suite.service.GetMonthlyReport(tc.year, 6)

			// Then
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), "invalid year")
		})
	}
}

func (suite *ReportServiceTestSuite) TestGetMonthlyReport_InvalidMonth() {
	testCases := []struct {
		name  string
		month int
	}{
		{"month zero", 0},
		{"month negative", -1},
		{"month too high", 13},
		{"month way too high", 25},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// When
			result, err := suite.service.GetMonthlyReport(2024, tc.month)

			// Then
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), "month must be between 1 and 12")
		})
	}
}

func (suite *ReportServiceTestSuite) TestGetCurrentMonthReport_Success() {
	// Given
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())
	
	startDate := time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	transactions := []models.Transaction{
		{
			ID:       1,
			Type:     "income",
			Amount:   3000,
			Currency: "ARS",
			Category: "salary",
			Date:     now,
		},
	}

	suite.mockRepo.On("GetByDateRange", mock.MatchedBy(func(start time.Time) bool {
		return start.Equal(startDate)
	}), mock.MatchedBy(func(end time.Time) bool {
		return end.Equal(endDate)
	})).Return(transactions, nil)

	// When
	result, err := suite.service.GetCurrentMonthReport()

	// Then
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), now.Month().String(), result.Month)
	assert.Equal(suite.T(), currentYear, result.Year)
	assert.Equal(suite.T(), 1, result.Summary.TransactionCount)
	assert.Equal(suite.T(), 3000.0, result.TotalIncome["ARS"])
}

func TestReportServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ReportServiceTestSuite))
}