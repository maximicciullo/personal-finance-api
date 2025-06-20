package repositories

import (
	"errors"
	"sync"
	"time"

	"github.com/maximicciullo/personal-finance-api/internal/middleware"
	"github.com/maximicciullo/personal-finance-api/internal/models"
	"go.uber.org/zap"
)

type MemoryTransactionRepository struct {
	transactions []models.Transaction
	nextID       int
	mutex        sync.RWMutex
	logger       *middleware.BusinessLoggerInstance
}

func NewMemoryTransactionRepository() *MemoryTransactionRepository {
	return &MemoryTransactionRepository{
		transactions: make([]models.Transaction, 0),
		nextID:       1,
		logger:       middleware.BusinessLogger(),
	}
}

func (r *MemoryTransactionRepository) Create(transaction *models.Transaction) error {
	r.logger.Repository("Create transaction started",
		zap.String("type", transaction.Type),
		zap.Float64("amount", transaction.Amount),
		zap.String("currency", transaction.Currency),
		zap.String("category", transaction.Category),
	)

	r.mutex.Lock()
	defer r.mutex.Unlock()

	start := time.Now()

	transaction.ID = r.nextID
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()

	r.transactions = append(r.transactions, *transaction)
	r.nextID++

	duration := time.Since(start)
	r.logger.Performance("Create transaction", duration,
		zap.Int("transaction_id", transaction.ID),
		zap.Int("total_transactions", len(r.transactions)),
	)

	r.logger.Repository("Create transaction completed successfully",
		zap.Int("transaction_id", transaction.ID),
		zap.Int("next_id", r.nextID),
		zap.Duration("duration", duration),
	)

	return nil
}

func (r *MemoryTransactionRepository) GetByID(id int) (*models.Transaction, error) {
	r.logger.Repository("GetByID started",
		zap.Int("transaction_id", id),
	)

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	start := time.Now()
	searched := 0

	for _, transaction := range r.transactions {
		searched++
		if transaction.ID == id {
			duration := time.Since(start)
			r.logger.Performance("GetByID transaction found", duration,
				zap.Int("transaction_id", id),
				zap.Int("search_iterations", searched),
			)

			r.logger.Repository("GetByID completed successfully",
				zap.Int("transaction_id", id),
				zap.Duration("duration", duration),
			)
			return &transaction, nil
		}
	}

	duration := time.Since(start)
	err := errors.New("transaction not found")
	
	r.logger.Performance("GetByID transaction not found", duration,
		zap.Int("transaction_id", id),
		zap.Int("searched_count", searched),
	)

	r.logger.Error("repository", "GetByID - transaction not found", err,
		zap.Int("transaction_id", id),
		zap.Int("total_transactions", len(r.transactions)),
	)

	return nil, err
}

func (r *MemoryTransactionRepository) GetAll() ([]models.Transaction, error) {
	r.logger.Repository("GetAll started")

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	start := time.Now()

	// Return a copy to avoid concurrent modification
	result := make([]models.Transaction, len(r.transactions))
	copy(result, r.transactions)

	duration := time.Since(start)
	r.logger.Performance("GetAll transactions", duration,
		zap.Int("transaction_count", len(result)),
	)

	r.logger.Repository("GetAll completed successfully",
		zap.Int("transaction_count", len(result)),
		zap.Duration("duration", duration),
	)

	return result, nil
}

func (r *MemoryTransactionRepository) GetByFilters(filters models.TransactionFilters) ([]models.Transaction, error) {
	r.logger.Repository("GetByFilters started",
		zap.String("type_filter", filters.Type),
		zap.String("category_filter", filters.Category),
		zap.String("currency_filter", filters.Currency),
	)

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	start := time.Now()
	var result []models.Transaction
	processed := 0

	for _, transaction := range r.transactions {
		processed++
		if r.matchesFilters(transaction, filters) {
			result = append(result, transaction)
			r.logger.Debug("repository", "Transaction matches filters",
				zap.Int("transaction_id", transaction.ID),
				zap.String("type", transaction.Type),
				zap.String("category", transaction.Category),
			)
		}
	}

	duration := time.Since(start)
	r.logger.Performance("GetByFilters search", duration,
		zap.Int("total_transactions", len(r.transactions)),
		zap.Int("processed_transactions", processed),
		zap.Int("filtered_count", len(result)),
	)

	r.logger.Repository("GetByFilters completed successfully",
		zap.Int("filtered_count", len(result)),
		zap.Duration("duration", duration),
	)

	return result, nil
}

func (r *MemoryTransactionRepository) GetByDateRange(startDate, endDate time.Time) ([]models.Transaction, error) {
	r.logger.Repository("GetByDateRange started",
		zap.Time("start_date", startDate),
		zap.Time("end_date", endDate),
	)

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	start := time.Now()
	var result []models.Transaction
	processed := 0

	for _, transaction := range r.transactions {
		processed++
		if transaction.Date.After(startDate.Add(-time.Second)) && transaction.Date.Before(endDate.Add(time.Second)) {
			result = append(result, transaction)
			r.logger.Debug("repository", "Transaction matches date range",
				zap.Int("transaction_id", transaction.ID),
				zap.Time("transaction_date", transaction.Date),
			)
		}
	}

	duration := time.Since(start)
	r.logger.Performance("GetByDateRange search", duration,
		zap.Int("total_transactions", len(r.transactions)),
		zap.Int("processed_transactions", processed),
		zap.Int("filtered_count", len(result)),
		zap.Duration("date_range", endDate.Sub(startDate)),
	)

	r.logger.Repository("GetByDateRange completed successfully",
		zap.Int("filtered_count", len(result)),
		zap.Duration("duration", duration),
	)

	return result, nil
}

func (r *MemoryTransactionRepository) Delete(id int) error {
	r.logger.Repository("Delete started",
		zap.Int("transaction_id", id),
	)

	r.mutex.Lock()
	defer r.mutex.Unlock()

	start := time.Now()
	searched := 0

	for i, transaction := range r.transactions {
		searched++
		if transaction.ID == id {
			// Store transaction info before deletion for logging
			deletedTransaction := transaction
			
			r.transactions = append(r.transactions[:i], r.transactions[i+1:]...)
			
			duration := time.Since(start)
			r.logger.Performance("Delete transaction", duration,
				zap.Int("transaction_id", id),
				zap.Int("search_iterations", searched),
				zap.Int("remaining_transactions", len(r.transactions)),
			)

			r.logger.Repository("Delete completed successfully",
				zap.Int("transaction_id", id),
				zap.Int("position", i),
				zap.String("deleted_type", deletedTransaction.Type),
				zap.Float64("deleted_amount", deletedTransaction.Amount),
				zap.Duration("duration", duration),
			)
			return nil
		}
	}

	duration := time.Since(start)
	err := errors.New("transaction not found")
	
	r.logger.Performance("Delete transaction not found", duration,
		zap.Int("transaction_id", id),
		zap.Int("searched_count", searched),
	)

	r.logger.Error("repository", "Delete - transaction not found", err,
		zap.Int("transaction_id", id),
		zap.Int("total_transactions", len(r.transactions)),
	)

	return err
}

func (r *MemoryTransactionRepository) Update(transaction *models.Transaction) error {
	r.logger.Repository("Update started",
		zap.Int("transaction_id", transaction.ID),
		zap.String("type", transaction.Type),
		zap.Float64("amount", transaction.Amount),
	)

	r.mutex.Lock()
	defer r.mutex.Unlock()

	start := time.Now()
	searched := 0

	for i, t := range r.transactions {
		searched++
		if t.ID == transaction.ID {
			// Store old values for logging
			oldTransaction := t
			
			transaction.UpdatedAt = time.Now()
			r.transactions[i] = *transaction
			
			duration := time.Since(start)
			r.logger.Performance("Update transaction", duration,
				zap.Int("transaction_id", transaction.ID),
				zap.Int("search_iterations", searched),
			)

			r.logger.Repository("Update completed successfully",
				zap.Int("transaction_id", transaction.ID),
				zap.Int("position", i),
				zap.Float64("old_amount", oldTransaction.Amount),
				zap.Float64("new_amount", transaction.Amount),
				zap.Duration("duration", duration),
			)
			return nil
		}
	}

	duration := time.Since(start)
	err := errors.New("transaction not found")
	
	r.logger.Performance("Update transaction not found", duration,
		zap.Int("transaction_id", transaction.ID),
		zap.Int("searched_count", searched),
	)

	r.logger.Error("repository", "Update - transaction not found", err,
		zap.Int("transaction_id", transaction.ID),
		zap.Int("total_transactions", len(r.transactions)),
	)

	return err
}

func (r *MemoryTransactionRepository) matchesFilters(transaction models.Transaction, filters models.TransactionFilters) bool {
	r.logger.Debug("repository", "Checking transaction against filters",
		zap.Int("transaction_id", transaction.ID),
		zap.String("transaction_type", transaction.Type),
		zap.String("transaction_category", transaction.Category),
		zap.String("transaction_currency", transaction.Currency),
		zap.String("filter_type", filters.Type),
		zap.String("filter_category", filters.Category),
		zap.String("filter_currency", filters.Currency),
	)

	if filters.Type != "" && transaction.Type != filters.Type {
		r.logger.Debug("repository", "Transaction filtered out by type",
			zap.Int("transaction_id", transaction.ID),
			zap.String("transaction_type", transaction.Type),
			zap.String("filter_type", filters.Type),
		)
		return false
	}

	if filters.Category != "" && transaction.Category != filters.Category {
		r.logger.Debug("repository", "Transaction filtered out by category",
			zap.Int("transaction_id", transaction.ID),
			zap.String("transaction_category", transaction.Category),
			zap.String("filter_category", filters.Category),
		)
		return false
	}

	if filters.Currency != "" && transaction.Currency != filters.Currency {
		r.logger.Debug("repository", "Transaction filtered out by currency",
			zap.Int("transaction_id", transaction.ID),
			zap.String("transaction_currency", transaction.Currency),
			zap.String("filter_currency", filters.Currency),
		)
		return false
	}

	if filters.FromDate != nil && transaction.Date.Before(*filters.FromDate) {
		r.logger.Debug("repository", "Transaction filtered out by from_date",
			zap.Int("transaction_id", transaction.ID),
			zap.Time("transaction_date", transaction.Date),
			zap.Time("filter_from_date", *filters.FromDate),
		)
		return false
	}

	if filters.ToDate != nil && transaction.Date.After(*filters.ToDate) {
		r.logger.Debug("repository", "Transaction filtered out by to_date",
			zap.Int("transaction_id", transaction.ID),
			zap.Time("transaction_date", transaction.Date),
			zap.Time("filter_to_date", *filters.ToDate),
		)
		return false
	}

	r.logger.Debug("repository", "Transaction matches all filters",
		zap.Int("transaction_id", transaction.ID),
	)

	return true
}