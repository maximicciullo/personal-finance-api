package services_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/maximicciullo/personal-finance-api/internal/middleware"
	"github.com/maximicciullo/personal-finance-api/internal/models"
	"github.com/maximicciullo/personal-finance-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockTransactionRepository is a mock implementation of TransactionRepository
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(transaction *models.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetByID(id int) (*models.Transaction, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetAll() ([]models.Transaction, error) {
	args := m.Called()
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByFilters(filters models.TransactionFilters) ([]models.Transaction, error) {
	args := m.Called(filters)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByDateRange(startDate, endDate time.Time) ([]models.Transaction, error) {
	args := m.Called(startDate, endDate)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTransactionRepository) Update(transaction *models.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

// TransactionServiceTestSuite is the test suite for TransactionService
type TransactionServiceTestSuite struct {
	suite.Suite
	mockRepo *MockTransactionRepository
	service  services.TransactionService
}

func (suite *TransactionServiceTestSuite) SetupTest() {
	// Initialize logger for testing
	middleware.InitLogger("test")
	
	suite.mockRepo = new(MockTransactionRepository)
	suite.service = services.NewTransactionService(suite.mockRepo)
}

func (suite *TransactionServiceTestSuite) TearDownTest() {
	suite.mockRepo.AssertExpectations(suite.T())
}

// Test CreateTransaction
func (suite *TransactionServiceTestSuite) TestCreateTransaction_Success() {
	// Given
	request := &models.CreateTransactionRequest{
		Type:        "expense",
		Amount:      1500.50,
		Currency:    "ARS",
		Description: "Test expense",
		Category:    "food",
	}

	expectedTransaction := &models.Transaction{
		Type:        "expense",
		Amount:      1500.50,
		Currency:    "ARS",
		Description: "Test expense",
		Category:    "food",
	}

	// Set expectation: repository Create should be called and should succeed
	suite.mockRepo.On("Create", mock.MatchedBy(func(t *models.Transaction) bool {
		return t.Type == expectedTransaction.Type &&
			t.Amount == expectedTransaction.Amount &&
			t.Currency == expectedTransaction.Currency &&
			t.Description == expectedTransaction.Description &&
			t.Category == expectedTransaction.Category
	})).Return(nil).Run(func(args mock.Arguments) {
		// Simulate what repository does - set ID and timestamps
		transaction := args.Get(0).(*models.Transaction)
		transaction.ID = 1
		transaction.CreatedAt = time.Now()
		transaction.UpdatedAt = time.Now()
	})

	// When
	result, err := suite.service.CreateTransaction(request)

	// Then
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), 1, result.ID)
	assert.Equal(suite.T(), "expense", result.Type)
	assert.Equal(suite.T(), 1500.50, result.Amount)
	assert.Equal(suite.T(), "ARS", result.Currency)
	assert.Equal(suite.T(), "Test expense", result.Description)
	assert.Equal(suite.T(), "food", result.Category)
}

func (suite *TransactionServiceTestSuite) TestCreateTransaction_WithCustomDate() {
	// Given
	customDate := "2024-06-15"
	request := &models.CreateTransactionRequest{
		Type:        "income",
		Amount:      5000,
		Currency:    "USD",
		Description: "Freelance payment",
		Category:    "work",
		Date:        &customDate,
	}

	suite.mockRepo.On("Create", mock.AnythingOfType("*models.Transaction")).Return(nil).Run(func(args mock.Arguments) {
		transaction := args.Get(0).(*models.Transaction)
		transaction.ID = 2
	})

	// When
	result, err := suite.service.CreateTransaction(request)

	// Then
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	expectedDate, _ := time.Parse("2006-01-02", "2024-06-15")
	assert.Equal(suite.T(), expectedDate.Format("2006-01-02"), result.Date.Format("2006-01-02"))
}

func (suite *TransactionServiceTestSuite) TestCreateTransaction_DefaultCurrency() {
	// Given
	request := &models.CreateTransactionRequest{
		Type:        "expense",
		Amount:      100,
		Description: "Test without currency",
		Category:    "test",
		// Currency not specified
	}

	suite.mockRepo.On("Create", mock.MatchedBy(func(t *models.Transaction) bool {
		return t.Currency == "ARS" // Should default to ARS
	})).Return(nil).Run(func(args mock.Arguments) {
		transaction := args.Get(0).(*models.Transaction)
		transaction.ID = 3
	})

	// When
	result, err := suite.service.CreateTransaction(request)

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "ARS", result.Currency)
}

func (suite *TransactionServiceTestSuite) TestCreateTransaction_ValidationErrors() {
	testCases := []struct {
		name    string
		request *models.CreateTransactionRequest
		error   string
	}{
		{
			name: "invalid type",
			request: &models.CreateTransactionRequest{
				Type:        "invalid",
				Amount:      100,
				Description: "Test",
				Category:    "test",
			},
			error: "type must be 'expense' or 'income'",
		},
		{
			name: "zero amount",
			request: &models.CreateTransactionRequest{
				Type:        "expense",
				Amount:      0,
				Description: "Test",
				Category:    "test",
			},
			error: "amount must be positive",
		},
		{
			name: "negative amount",
			request: &models.CreateTransactionRequest{
				Type:        "expense",
				Amount:      -100,
				Description: "Test",
				Category:    "test",
			},
			error: "amount must be positive",
		},
		{
			name: "empty description",
			request: &models.CreateTransactionRequest{
				Type:     "expense",
				Amount:   100,
				Category: "test",
			},
			error: "description is required",
		},
		{
			name: "empty category",
			request: &models.CreateTransactionRequest{
				Type:        "expense",
				Amount:      100,
				Description: "Test",
			},
			error: "category is required",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// When
			result, err := suite.service.CreateTransaction(tc.request)

			// Then
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), tc.error)
		})
	}
}

func (suite *TransactionServiceTestSuite) TestCreateTransaction_InvalidDate() {
	// Given
	invalidDate := "invalid-date"
	request := &models.CreateTransactionRequest{
		Type:        "expense",
		Amount:      100,
		Description: "Test",
		Category:    "test",
		Date:        &invalidDate,
	}

	// When
	result, err := suite.service.CreateTransaction(request)

	// Then
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "invalid date format")
}

func (suite *TransactionServiceTestSuite) TestCreateTransaction_RepositoryError() {
	// Given
	request := &models.CreateTransactionRequest{
		Type:        "expense",
		Amount:      100,
		Description: "Test",
		Category:    "test",
	}

	suite.mockRepo.On("Create", mock.AnythingOfType("*models.Transaction")).Return(errors.New("database error"))

	// When
	result, err := suite.service.CreateTransaction(request)

	// Then
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "database error")
}

// Test GetTransaction
func (suite *TransactionServiceTestSuite) TestGetTransaction_Success() {
	// Given
	expectedTransaction := &models.Transaction{
		ID:          1,
		Type:        "expense",
		Amount:      100,
		Currency:    "ARS",
		Description: "Test",
		Category:    "test",
	}

	suite.mockRepo.On("GetByID", 1).Return(expectedTransaction, nil)

	// When
	result, err := suite.service.GetTransaction(1)

	// Then
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), expectedTransaction, result)
}

func (suite *TransactionServiceTestSuite) TestGetTransaction_InvalidID() {
	testCases := []int{0, -1, -10}

	for _, id := range testCases {
		suite.T().Run(fmt.Sprintf("ID_%d", id), func(t *testing.T) {
			// When
			result, err := suite.service.GetTransaction(id)

			// Then
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), "invalid transaction ID")
		})
	}
}

func (suite *TransactionServiceTestSuite) TestGetTransaction_NotFound() {
	// Given
	suite.mockRepo.On("GetByID", 999).Return(nil, errors.New("transaction not found"))

	// When
	result, err := suite.service.GetTransaction(999)

	// Then
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "transaction not found")
}

// Test GetTransactions
func (suite *TransactionServiceTestSuite) TestGetTransactions_Success() {
	// Given
	filters := models.TransactionFilters{
		Type:     "expense",
		Category: "food",
	}

	expectedTransactions := []models.Transaction{
		{ID: 1, Type: "expense", Amount: 100, Category: "food"},
		{ID: 2, Type: "expense", Amount: 200, Category: "food"},
	}

	suite.mockRepo.On("GetByFilters", filters).Return(expectedTransactions, nil)

	// When
	result, err := suite.service.GetTransactions(filters)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), expectedTransactions, result)
}

func (suite *TransactionServiceTestSuite) TestGetTransactions_EmptyResult() {
	// Given
	filters := models.TransactionFilters{Type: "income"}
	emptyResult := []models.Transaction{}

	suite.mockRepo.On("GetByFilters", filters).Return(emptyResult, nil)

	// When
	result, err := suite.service.GetTransactions(filters)

	// Then
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
}

func (suite *TransactionServiceTestSuite) TestGetTransactions_RepositoryError() {
	// Given
	filters := models.TransactionFilters{}
	suite.mockRepo.On("GetByFilters", filters).Return([]models.Transaction{}, errors.New("database error"))

	// When
	result, err := suite.service.GetTransactions(filters)

	// Then
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

// Test DeleteTransaction
func (suite *TransactionServiceTestSuite) TestDeleteTransaction_Success() {
	// Given
	suite.mockRepo.On("Delete", 1).Return(nil)

	// When
	err := suite.service.DeleteTransaction(1)

	// Then
	assert.NoError(suite.T(), err)
}

func (suite *TransactionServiceTestSuite) TestDeleteTransaction_InvalidID() {
	testCases := []int{0, -1, -5}

	for _, id := range testCases {
		suite.T().Run(fmt.Sprintf("ID_%d", id), func(t *testing.T) {
			// When
			err := suite.service.DeleteTransaction(id)

			// Then
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid transaction ID")
		})
	}
}

func (suite *TransactionServiceTestSuite) TestDeleteTransaction_NotFound() {
	// Given
	suite.mockRepo.On("Delete", 999).Return(errors.New("transaction not found"))

	// When
	err := suite.service.DeleteTransaction(999)

	// Then
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "transaction not found")
}

func TestTransactionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionServiceTestSuite))
}