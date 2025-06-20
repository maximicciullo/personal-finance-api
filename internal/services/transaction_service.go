package services

import (
	"errors"
	"github.com/maximicciullo/personal-finance-api/internal/models"
	"github.com/maximicciullo/personal-finance-api/internal/repositories"
	"time"
)

type transactionService struct {
	repo repositories.TransactionRepository
}

func NewTransactionService(repo repositories.TransactionRepository) TransactionService {
	return &transactionService{
		repo: repo,
	}
}

func (s *transactionService) CreateTransaction(req *models.CreateTransactionRequest) (*models.Transaction, error) {
	// Validate request
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Parse date
	var transactionDate time.Time
	if req.Date != nil {
		var err error
		transactionDate, err = time.Parse("2006-01-02", *req.Date)
		if err != nil {
			return nil, errors.New("invalid date format, use YYYY-MM-DD")
		}
	} else {
		transactionDate = time.Now()
	}

	// Set default currency if not provided
	currency := req.Currency
	if currency == "" {
		currency = models.CurrencyARS
	}

	// Create transaction
	transaction := &models.Transaction{
		Type:        req.Type,
		Amount:      req.Amount,
		Currency:    currency,
		Description: req.Description,
		Category:    req.Category,
		Date:        transactionDate,
	}

	err := s.repo.Create(transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *transactionService) GetTransaction(id int) (*models.Transaction, error) {
	if id <= 0 {
		return nil, errors.New("invalid transaction ID")
	}

	return s.repo.GetByID(id)
}

func (s *transactionService) GetTransactions(filters models.TransactionFilters) ([]models.Transaction, error) {
	return s.repo.GetByFilters(filters)
}

func (s *transactionService) DeleteTransaction(id int) error {
	if id <= 0 {
		return errors.New("invalid transaction ID")
	}

	return s.repo.Delete(id)
}

func (s *transactionService) validateCreateRequest(req *models.CreateTransactionRequest) error {
	if req.Type != models.TransactionTypeExpense && req.Type != models.TransactionTypeIncome {
		return errors.New("type must be 'expense' or 'income'")
	}

	if req.Amount <= 0 {
		return errors.New("amount must be positive")
	}

	if req.Description == "" {
		return errors.New("description is required")
	}

	if req.Category == "" {
		return errors.New("category is required")
	}

	return nil
}