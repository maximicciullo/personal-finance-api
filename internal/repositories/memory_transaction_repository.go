package repositories

import (
	"errors"
	"sync"
	"time"

	"github.com/maximicciullo/personal-finance-api/internal/models"
)

type MemoryTransactionRepository struct {
	transactions []models.Transaction
	nextID       int
	mutex        sync.RWMutex
}

func NewMemoryTransactionRepository() *MemoryTransactionRepository {
	return &MemoryTransactionRepository{
		transactions: make([]models.Transaction, 0),
		nextID:       1,
	}
}

func (r *MemoryTransactionRepository) Create(transaction *models.Transaction) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	transaction.ID = r.nextID
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()

	r.transactions = append(r.transactions, *transaction)
	r.nextID++

	return nil
}

func (r *MemoryTransactionRepository) GetByID(id int) (*models.Transaction, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, transaction := range r.transactions {
		if transaction.ID == id {
			return &transaction, nil
		}
	}

	return nil, errors.New("transaction not found")
}

func (r *MemoryTransactionRepository) GetAll() ([]models.Transaction, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Return a copy to avoid concurrent modification
	result := make([]models.Transaction, len(r.transactions))
	copy(result, r.transactions)

	return result, nil
}

func (r *MemoryTransactionRepository) GetByFilters(filters models.TransactionFilters) ([]models.Transaction, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []models.Transaction

	for _, transaction := range r.transactions {
		if !r.matchesFilters(transaction, filters) {
			continue
		}
		result = append(result, transaction)
	}

	return result, nil
}

func (r *MemoryTransactionRepository) GetByDateRange(startDate, endDate time.Time) ([]models.Transaction, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []models.Transaction

	for _, transaction := range r.transactions {
		if transaction.Date.After(startDate.Add(-time.Second)) && transaction.Date.Before(endDate.Add(time.Second)) {
			result = append(result, transaction)
		}
	}

	return result, nil
}

func (r *MemoryTransactionRepository) Delete(id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, transaction := range r.transactions {
		if transaction.ID == id {
			r.transactions = append(r.transactions[:i], r.transactions[i+1:]...)
			return nil
		}
	}

	return errors.New("transaction not found")
}

func (r *MemoryTransactionRepository) Update(transaction *models.Transaction) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, t := range r.transactions {
		if t.ID == transaction.ID {
			transaction.UpdatedAt = time.Now()
			r.transactions[i] = *transaction
			return nil
		}
	}

	return errors.New("transaction not found")
}

func (r *MemoryTransactionRepository) matchesFilters(transaction models.Transaction, filters models.TransactionFilters) bool {
	if filters.Type != "" && transaction.Type != filters.Type {
		return false
	}

	if filters.Category != "" && transaction.Category != filters.Category {
		return false
	}

	if filters.Currency != "" && transaction.Currency != filters.Currency {
		return false
	}

	if filters.FromDate != nil && transaction.Date.Before(*filters.FromDate) {
		return false
	}

	if filters.ToDate != nil && transaction.Date.After(*filters.ToDate) {
		return false
	}

	return true
}
