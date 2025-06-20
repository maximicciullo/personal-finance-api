package repositories_test

import (
	"testing"
	"time"

	"github.com/maximicciullo/personal-finance-api/internal/middleware"
	"github.com/maximicciullo/personal-finance-api/internal/models"
	"github.com/maximicciullo/personal-finance-api/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// MemoryTransactionRepositoryTestSuite is the test suite for MemoryTransactionRepository
type MemoryTransactionRepositoryTestSuite struct {
	suite.Suite
	repo *repositories.MemoryTransactionRepository
}

func (suite *MemoryTransactionRepositoryTestSuite) SetupTest() {
	// Initialize logger for testing
	middleware.InitLogger("test")

	suite.repo = repositories.NewMemoryTransactionRepository()
}

// Test Create
func (suite *MemoryTransactionRepositoryTestSuite) TestCreate_Success() {
	// Given
	transaction := &models.Transaction{
		Type:        "expense",
		Amount:      1500.50,
		Currency:    "ARS",
		Description: "Test expense",
		Category:    "food",
		Date:        time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC),
	}

	// When
	beforeCreate := time.Now()
	err := suite.repo.Create(transaction)
	afterCreate := time.Now()

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, transaction.ID)
	assert.NotZero(suite.T(), transaction.CreatedAt)
	assert.NotZero(suite.T(), transaction.UpdatedAt)

	// Verify timestamps are within reasonable range
	assert.True(suite.T(), transaction.CreatedAt.After(beforeCreate) || transaction.CreatedAt.Equal(beforeCreate))
	assert.True(suite.T(), transaction.CreatedAt.Before(afterCreate) || transaction.CreatedAt.Equal(afterCreate))
	assert.True(suite.T(), transaction.UpdatedAt.After(beforeCreate) || transaction.UpdatedAt.Equal(beforeCreate))
	assert.True(suite.T(), transaction.UpdatedAt.Before(afterCreate) || transaction.UpdatedAt.Equal(afterCreate))

	// CreatedAt and UpdatedAt should be very close (within 1 second)
	timeDiff := transaction.UpdatedAt.Sub(transaction.CreatedAt)
	assert.True(suite.T(), timeDiff >= 0 && timeDiff < time.Second)
}

func (suite *MemoryTransactionRepositoryTestSuite) TestCreate_MultipleTransactions() {
	// Given
	transaction1 := &models.Transaction{
		Type:        "expense",
		Amount:      100,
		Currency:    "ARS",
		Description: "First transaction",
		Category:    "food",
		Date:        time.Now(),
	}

	transaction2 := &models.Transaction{
		Type:        "income",
		Amount:      200,
		Currency:    "USD",
		Description: "Second transaction",
		Category:    "work",
		Date:        time.Now(),
	}

	// When
	err1 := suite.repo.Create(transaction1)
	err2 := suite.repo.Create(transaction2)

	// Then
	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	assert.Equal(suite.T(), 1, transaction1.ID)
	assert.Equal(suite.T(), 2, transaction2.ID)
}

func (suite *MemoryTransactionRepositoryTestSuite) TestCreate_ConcurrentAccess() {
	// Test thread safety by creating transactions concurrently
	done := make(chan bool)
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			transaction := &models.Transaction{
				Type:        "expense",
				Amount:      float64(index * 100),
				Currency:    "ARS",
				Description: "Concurrent transaction",
				Category:    "test",
				Date:        time.Now(),
			}
			err := suite.repo.Create(transaction)
			assert.NoError(suite.T(), err)
			assert.NotZero(suite.T(), transaction.ID)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify all transactions were created
	transactions, err := suite.repo.GetAll()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), transactions, numGoroutines)
}

// Test GetByID
func (suite *MemoryTransactionRepositoryTestSuite) TestGetByID_Success() {
	// Given
	transaction := &models.Transaction{
		Type:        "income",
		Amount:      5000,
		Currency:    "USD",
		Description: "Test income",
		Category:    "work",
		Date:        time.Now(),
	}
	suite.repo.Create(transaction)

	// When
	result, err := suite.repo.GetByID(1)

	// Then
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), 1, result.ID)
	assert.Equal(suite.T(), "income", result.Type)
	assert.Equal(suite.T(), 5000.0, result.Amount)
	assert.Equal(suite.T(), "USD", result.Currency)
	assert.Equal(suite.T(), "Test income", result.Description)
	assert.Equal(suite.T(), "work", result.Category)
}

func (suite *MemoryTransactionRepositoryTestSuite) TestGetByID_NotFound() {
	// When
	result, err := suite.repo.GetByID(999)

	// Then
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "transaction not found")
}

func (suite *MemoryTransactionRepositoryTestSuite) TestGetByID_MultipleTransactions() {
	// Given - create multiple transactions
	transactions := []*models.Transaction{
		{Type: "expense", Amount: 100, Currency: "ARS", Description: "First", Category: "food", Date: time.Now()},
		{Type: "income", Amount: 200, Currency: "USD", Description: "Second", Category: "work", Date: time.Now()},
		{Type: "expense", Amount: 300, Currency: "EUR", Description: "Third", Category: "transport", Date: time.Now()},
	}

	for _, tx := range transactions {
		suite.repo.Create(tx)
	}

	// When & Then - verify each transaction can be retrieved
	for i, expectedTx := range transactions {
		result, err := suite.repo.GetByID(i + 1)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), expectedTx.Description, result.Description)
		assert.Equal(suite.T(), expectedTx.Amount, result.Amount)
	}
}

// Test GetAll
func (suite *MemoryTransactionRepositoryTestSuite) TestGetAll_EmptyRepository() {
	// When
	result, err := suite.repo.GetAll()

	// Then
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Empty(suite.T(), result)
}

func (suite *MemoryTransactionRepositoryTestSuite) TestGetAll_WithTransactions() {
	// Given
	transactions := []*models.Transaction{
		{Type: "expense", Amount: 100, Currency: "ARS", Description: "First", Category: "food", Date: time.Now()},
		{Type: "income", Amount: 200, Currency: "USD", Description: "Second", Category: "work", Date: time.Now()},
	}

	for _, tx := range transactions {
		suite.repo.Create(tx)
	}

	// When
	result, err := suite.repo.GetAll()

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), "First", result[0].Description)
	assert.Equal(suite.T(), "Second", result[1].Description)
}

func (suite *MemoryTransactionRepositoryTestSuite) TestGetAll_ReturnsCopy() {
	// Given
	transaction := &models.Transaction{
		Type: "expense", Amount: 100, Currency: "ARS",
		Description: "Test", Category: "food", Date: time.Now(),
	}
	suite.repo.Create(transaction)

	// When
	result1, _ := suite.repo.GetAll()
	result2, _ := suite.repo.GetAll()

	// Then - verify we get different slices (copies)
	assert.NotSame(suite.T(), &result1, &result2)

	// Modify one result and verify the other is unchanged
	result1[0].Description = "Modified"
	assert.Equal(suite.T(), "Modified", result1[0].Description)
	assert.Equal(suite.T(), "Test", result2[0].Description)
}

// Test GetByFilters
func (suite *MemoryTransactionRepositoryTestSuite) TestGetByFilters_NoFilters() {
	// Given
	transactions := []*models.Transaction{
		{Type: "expense", Amount: 100, Currency: "ARS", Description: "Expense", Category: "food", Date: time.Now()},
		{Type: "income", Amount: 200, Currency: "USD", Description: "Income", Category: "work", Date: time.Now()},
	}

	for _, tx := range transactions {
		suite.repo.Create(tx)
	}

	// When
	filters := models.TransactionFilters{}
	result, err := suite.repo.GetByFilters(filters)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
}

func (suite *MemoryTransactionRepositoryTestSuite) TestGetByFilters_TypeFilter() {
	// Given
	transactions := []*models.Transaction{
		{Type: "expense", Amount: 100, Currency: "ARS", Description: "Expense1", Category: "food", Date: time.Now()},
		{Type: "income", Amount: 200, Currency: "USD", Description: "Income1", Category: "work", Date: time.Now()},
		{Type: "expense", Amount: 300, Currency: "EUR", Description: "Expense2", Category: "transport", Date: time.Now()},
	}

	for _, tx := range transactions {
		suite.repo.Create(tx)
	}

	// When
	filters := models.TransactionFilters{Type: "expense"}
	result, err := suite.repo.GetByFilters(filters)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	for _, tx := range result {
		assert.Equal(suite.T(), "expense", tx.Type)
	}
}

func (suite *MemoryTransactionRepositoryTestSuite) TestGetByFilters_CategoryFilter() {
	// Given
	transactions := []*models.Transaction{
		{Type: "expense", Amount: 100, Currency: "ARS", Description: "Food1", Category: "food", Date: time.Now()},
		{Type: "expense", Amount: 200, Currency: "USD", Description: "Food2", Category: "food", Date: time.Now()},
		{Type: "expense", Amount: 300, Currency: "EUR", Description: "Transport", Category: "transport", Date: time.Now()},
	}

	for _, tx := range transactions {
		suite.repo.Create(tx)
	}

	// When
	filters := models.TransactionFilters{Category: "food"}
	result, err := suite.repo.GetByFilters(filters)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	for _, tx := range result {
		assert.Equal(suite.T(), "food", tx.Category)
	}
}

func (suite *MemoryTransactionRepositoryTestSuite) TestGetByFilters_CurrencyFilter() {
	// Given
	transactions := []*models.Transaction{
		{Type: "expense", Amount: 100, Currency: "ARS", Description: "ARS1", Category: "food", Date: time.Now()},
		{Type: "income", Amount: 200, Currency: "USD", Description: "USD1", Category: "work", Date: time.Now()},
		{Type: "expense", Amount: 300, Currency: "ARS", Description: "ARS2", Category: "transport", Date: time.Now()},
	}

	for _, tx := range transactions {
		suite.repo.Create(tx)
	}

	// When
	filters := models.TransactionFilters{Currency: "ARS"}
	result, err := suite.repo.GetByFilters(filters)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	for _, tx := range result {
		assert.Equal(suite.T(), "ARS", tx.Currency)
	}
}

func (suite *MemoryTransactionRepositoryTestSuite) TestGetByFilters_DateRangeFilter() {
	// Given
	baseDate := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	transactions := []*models.Transaction{
		{Type: "expense", Amount: 100, Currency: "ARS", Description: "Before", Category: "food", Date: baseDate.AddDate(0, 0, -5)},
		{Type: "expense", Amount: 200, Currency: "ARS", Description: "During", Category: "food", Date: baseDate},
		{Type: "expense", Amount: 300, Currency: "ARS", Description: "After", Category: "food", Date: baseDate.AddDate(0, 0, 5)},
	}

	for _, tx := range transactions {
		suite.repo.Create(tx)
	}

	// When
	fromDate := baseDate.AddDate(0, 0, -1)
	toDate := baseDate.AddDate(0, 0, 1)
	filters := models.TransactionFilters{
		FromDate: &fromDate,
		ToDate:   &toDate,
	}
	result, err := suite.repo.GetByFilters(filters)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), "During", result[0].Description)
}

func (suite *MemoryTransactionRepositoryTestSuite) TestGetByFilters_MultipleFilters() {
	// Given
	transactions := []*models.Transaction{
		{Type: "expense", Amount: 100, Currency: "ARS", Description: "Match", Category: "food", Date: time.Now()},
		{Type: "income", Amount: 200, Currency: "ARS", Description: "WrongType", Category: "food", Date: time.Now()},
		{Type: "expense", Amount: 300, Currency: "USD", Description: "WrongCurrency", Category: "food", Date: time.Now()},
		{Type: "expense", Amount: 400, Currency: "ARS", Description: "WrongCategory", Category: "transport", Date: time.Now()},
	}

	for _, tx := range transactions {
		suite.repo.Create(tx)
	}

	// When
	filters := models.TransactionFilters{
		Type:     "expense",
		Currency: "ARS",
		Category: "food",
	}
	result, err := suite.repo.GetByFilters(filters)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), "Match", result[0].Description)
}

func (suite *MemoryTransactionRepositoryTestSuite) TestGetByFilters_NoMatches() {
	// Given
	transaction := &models.Transaction{
		Type: "expense", Amount: 100, Currency: "ARS",
		Description: "Test", Category: "food", Date: time.Now(),
	}
	suite.repo.Create(transaction)

	// When
	filters := models.TransactionFilters{Type: "nonexistent"}
	result, err := suite.repo.GetByFilters(filters)

	// Then
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
}

// Test GetByDateRange
func (suite *MemoryTransactionRepositoryTestSuite) TestGetByDateRange_Success() {
	// Given
	baseDate := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	transactions := []*models.Transaction{
		{Type: "expense", Amount: 100, Currency: "ARS", Description: "Before", Category: "food", Date: baseDate.AddDate(0, 0, -2)},
		{Type: "expense", Amount: 200, Currency: "ARS", Description: "Start", Category: "food", Date: baseDate},
		{Type: "expense", Amount: 300, Currency: "ARS", Description: "Middle", Category: "food", Date: baseDate.AddDate(0, 0, 1)},
		{Type: "expense", Amount: 400, Currency: "ARS", Description: "End", Category: "food", Date: baseDate.AddDate(0, 0, 2)},
		{Type: "expense", Amount: 500, Currency: "ARS", Description: "After", Category: "food", Date: baseDate.AddDate(0, 0, 4)},
	}

	for _, tx := range transactions {
		suite.repo.Create(tx)
	}

	// When
	startDate := baseDate
	endDate := baseDate.AddDate(0, 0, 2)
	result, err := suite.repo.GetByDateRange(startDate, endDate)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 3)

	descriptions := make([]string, len(result))
	for i, tx := range result {
		descriptions[i] = tx.Description
	}
	assert.Contains(suite.T(), descriptions, "Start")
	assert.Contains(suite.T(), descriptions, "Middle")
	assert.Contains(suite.T(), descriptions, "End")
}

func (suite *MemoryTransactionRepositoryTestSuite) TestGetByDateRange_EmptyResult() {
	// Given
	transaction := &models.Transaction{
		Type: "expense", Amount: 100, Currency: "ARS",
		Description: "Test", Category: "food",
		Date: time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
	}
	suite.repo.Create(transaction)

	// When - search in a different date range
	startDate := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 7, 31, 0, 0, 0, 0, time.UTC)
	result, err := suite.repo.GetByDateRange(startDate, endDate)

	// Then
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
}

// Test Delete
func (suite *MemoryTransactionRepositoryTestSuite) TestDelete_Success() {
	// Given
	transaction := &models.Transaction{
		Type: "expense", Amount: 100, Currency: "ARS",
		Description: "Test", Category: "food", Date: time.Now(),
	}
	suite.repo.Create(transaction)

	// When
	err := suite.repo.Delete(1)

	// Then
	assert.NoError(suite.T(), err)

	// Verify transaction is deleted
	result, err := suite.repo.GetByID(1)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)

	// Verify repository is empty
	all, err := suite.repo.GetAll()
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), all)
}

func (suite *MemoryTransactionRepositoryTestSuite) TestDelete_NotFound() {
	// When
	err := suite.repo.Delete(999)

	// Then
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "transaction not found")
}

func (suite *MemoryTransactionRepositoryTestSuite) TestDelete_MiddleTransaction() {
	// Given - create multiple transactions
	transactions := []*models.Transaction{
		{Type: "expense", Amount: 100, Currency: "ARS", Description: "First", Category: "food", Date: time.Now()},
		{Type: "income", Amount: 200, Currency: "USD", Description: "Second", Category: "work", Date: time.Now()},
		{Type: "expense", Amount: 300, Currency: "EUR", Description: "Third", Category: "transport", Date: time.Now()},
	}

	for _, tx := range transactions {
		suite.repo.Create(tx)
	}

	// When - delete middle transaction
	err := suite.repo.Delete(2)

	// Then
	assert.NoError(suite.T(), err)

	// Verify only 2 transactions remain
	all, err := suite.repo.GetAll()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), all, 2)

	// Verify correct transactions remain
	assert.Equal(suite.T(), "First", all[0].Description)
	assert.Equal(suite.T(), "Third", all[1].Description)

	// Verify deleted transaction cannot be found
	result, err := suite.repo.GetByID(2)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

// Test Update
func (suite *MemoryTransactionRepositoryTestSuite) TestUpdate_Success() {
	// Given
	original := &models.Transaction{
		Type: "expense", Amount: 100, Currency: "ARS",
		Description: "Original", Category: "food", Date: time.Now(),
	}
	suite.repo.Create(original)

	// Capture original timestamps and wait a bit
	originalCreatedAt := original.CreatedAt
	time.Sleep(10 * time.Millisecond) // Ensure timestamp difference

	// When
	updated := &models.Transaction{
		ID:          1,
		Type:        "income",
		Amount:      500,
		Currency:    "USD",
		Description: "Updated",
		Category:    "work",
		Date:        time.Now().AddDate(0, 0, 1),
		CreatedAt:   originalCreatedAt, // Keep original
	}

	beforeUpdate := time.Now()
	err := suite.repo.Update(updated)
	afterUpdate := time.Now()

	// Then
	assert.NoError(suite.T(), err)

	// Verify update
	result, err := suite.repo.GetByID(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "income", result.Type)
	assert.Equal(suite.T(), 500.0, result.Amount)
	assert.Equal(suite.T(), "USD", result.Currency)
	assert.Equal(suite.T(), "Updated", result.Description)
	assert.Equal(suite.T(), "work", result.Category)

	// Verify timestamps
	assert.Equal(suite.T(), originalCreatedAt, result.CreatedAt)
	assert.True(suite.T(), result.UpdatedAt.After(originalCreatedAt))
	assert.True(suite.T(), result.UpdatedAt.After(beforeUpdate) || result.UpdatedAt.Equal(beforeUpdate))
	assert.True(suite.T(), result.UpdatedAt.Before(afterUpdate) || result.UpdatedAt.Equal(afterUpdate))
}

func (suite *MemoryTransactionRepositoryTestSuite) TestUpdate_NotFound() {
	// Given
	transaction := &models.Transaction{
		ID:   999,
		Type: "expense",
	}

	// When
	err := suite.repo.Update(transaction)

	// Then
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "transaction not found")
}

func TestMemoryTransactionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(MemoryTransactionRepositoryTestSuite))
}
