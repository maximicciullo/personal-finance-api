package services

import (
	"errors"
	"time"

	"github.com/maximicciullo/personal-finance-api/internal/middleware"
	"github.com/maximicciullo/personal-finance-api/internal/models"
	"github.com/maximicciullo/personal-finance-api/internal/repositories"
	"go.uber.org/zap"
)

type transactionService struct {
	repo   repositories.TransactionRepository
	logger *middleware.BusinessLoggerInstance
}

func NewTransactionService(repo repositories.TransactionRepository) TransactionService {
	return &transactionService{
		repo:   repo,
		logger: middleware.BusinessLogger(),
	}
}

func (s *transactionService) CreateTransaction(req *models.CreateTransactionRequest) (*models.Transaction, error) {
	s.logger.Service("CreateTransaction started",
		zap.String("type", req.Type),
		zap.Float64("amount", req.Amount),
		zap.String("currency", req.Currency),
		zap.String("category", req.Category),
	)

	start := time.Now()

	// Validate request
	if err := s.validateCreateRequest(req); err != nil {
		s.logger.Error("service", "CreateTransaction - validation failed", err,
			zap.Any("request", req),
		)
		return nil, err
	}

	s.logger.Service("CreateTransaction - validation passed",
		zap.Duration("validation_duration", time.Since(start)),
	)

	// Parse date
	var transactionDate time.Time
	if req.Date != nil {
		var err error
		transactionDate, err = time.Parse("2006-01-02", *req.Date)
		if err != nil {
			s.logger.Error("service", "CreateTransaction - date parsing failed", err,
				zap.String("date_string", *req.Date),
			)
			return nil, errors.New("invalid date format, use YYYY-MM-DD")
		}
		s.logger.Service("CreateTransaction - custom date parsed",
			zap.Time("transaction_date", transactionDate),
		)
	} else {
		transactionDate = time.Now()
		s.logger.Service("CreateTransaction - using current date",
			zap.Time("transaction_date", transactionDate),
		)
	}

	// Set default currency if not provided
	currency := req.Currency
	if currency == "" {
		currency = models.CurrencyARS
		s.logger.Service("CreateTransaction - using default currency",
			zap.String("default_currency", currency),
		)
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

	s.logger.Service("CreateTransaction - calling repository",
		zap.Any("transaction", transaction),
	)

	repoStart := time.Now()
	err := s.repo.Create(transaction)
	repoDuration := time.Since(repoStart)

	s.logger.Performance("CreateTransaction repository call", repoDuration,
		zap.Bool("success", err == nil),
		zap.Int("transaction_id", transaction.ID),
	)

	if err != nil {
		s.logger.Error("service", "CreateTransaction - repository error", err,
			zap.Any("transaction", transaction),
		)
		return nil, err
	}

	totalDuration := time.Since(start)
	s.logger.Service("CreateTransaction completed successfully",
		zap.Int("transaction_id", transaction.ID),
		zap.Duration("total_duration", totalDuration),
		zap.Duration("repo_duration", repoDuration),
	)

	return transaction, nil
}

func (s *transactionService) GetTransaction(id int) (*models.Transaction, error) {
	s.logger.Service("GetTransaction started",
		zap.Int("transaction_id", id),
	)

	if id <= 0 {
		err := errors.New("invalid transaction ID")
		s.logger.Error("service", "GetTransaction - invalid ID", err,
			zap.Int("transaction_id", id),
		)
		return nil, err
	}

	start := time.Now()
	transaction, err := s.repo.GetByID(id)
	duration := time.Since(start)

	s.logger.Performance("GetTransaction repository call", duration,
		zap.Int("transaction_id", id),
		zap.Bool("success", err == nil),
	)

	if err != nil {
		s.logger.Error("service", "GetTransaction - repository error", err,
			zap.Int("transaction_id", id),
		)
		return nil, err
	}

	s.logger.Service("GetTransaction completed successfully",
		zap.Int("transaction_id", id),
		zap.Duration("duration", duration),
	)

	return transaction, nil
}

func (s *transactionService) GetTransactions(filters models.TransactionFilters) ([]models.Transaction, error) {
	s.logger.Service("GetTransactions started",
		zap.String("type_filter", filters.Type),
		zap.String("category_filter", filters.Category),
		zap.String("currency_filter", filters.Currency),
	)

	start := time.Now()
	transactions, err := s.repo.GetByFilters(filters)
	duration := time.Since(start)

	s.logger.Performance("GetTransactions repository call", duration,
		zap.Int("transaction_count", len(transactions)),
		zap.Bool("success", err == nil),
	)

	if err != nil {
		s.logger.Error("service", "GetTransactions - repository error", err,
			zap.Any("filters", filters),
		)
		return nil, err
	}

	s.logger.Service("GetTransactions completed successfully",
		zap.Int("transaction_count", len(transactions)),
		zap.Duration("duration", duration),
	)

	return transactions, nil
}

func (s *transactionService) DeleteTransaction(id int) error {
	s.logger.Service("DeleteTransaction started",
		zap.Int("transaction_id", id),
	)

	if id <= 0 {
		err := errors.New("invalid transaction ID")
		s.logger.Error("service", "DeleteTransaction - invalid ID", err,
			zap.Int("transaction_id", id),
		)
		return err
	}

	start := time.Now()
	err := s.repo.Delete(id)
	duration := time.Since(start)

	s.logger.Performance("DeleteTransaction repository call", duration,
		zap.Int("transaction_id", id),
		zap.Bool("success", err == nil),
	)

	if err != nil {
		s.logger.Error("service", "DeleteTransaction - repository error", err,
			zap.Int("transaction_id", id),
		)
		return err
	}

	s.logger.Service("DeleteTransaction completed successfully",
		zap.Int("transaction_id", id),
		zap.Duration("duration", duration),
	)

	return nil
}

func (s *transactionService) UpdateTransaction(id int, req *models.UpdateTransactionRequest) (*models.Transaction, error) {
	s.logger.Service("UpdateTransaction started",
		zap.Int("transaction_id", id),
	)

	if id <= 0 {
		err := errors.New("invalid transaction ID")
		s.logger.Error("service", "UpdateTransaction - invalid ID", err,
			zap.Int("transaction_id", id),
		)
		return nil, err
	}

	start := time.Now()

	// Get existing transaction
	existingTransaction, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("service", "UpdateTransaction - transaction not found", err,
			zap.Int("transaction_id", id),
		)
		return nil, err
	}

	s.logger.Service("UpdateTransaction - existing transaction found",
		zap.Int("transaction_id", id),
		zap.String("current_type", existingTransaction.Type),
		zap.Float64("current_amount", existingTransaction.Amount),
	)

	// Validate update request
	if err := s.validateUpdateRequest(req); err != nil {
		s.logger.Error("service", "UpdateTransaction - validation failed", err,
			zap.Any("request", req),
		)
		return nil, err
	}

	// Create updated transaction with merged values
	updatedTransaction := *existingTransaction

	// Update only provided fields
	if req.Type != nil {
		updatedTransaction.Type = *req.Type
		s.logger.Service("UpdateTransaction - updating type",
			zap.String("old_type", existingTransaction.Type),
			zap.String("new_type", *req.Type),
		)
	}

	if req.Amount != nil {
		updatedTransaction.Amount = *req.Amount
		s.logger.Service("UpdateTransaction - updating amount",
			zap.Float64("old_amount", existingTransaction.Amount),
			zap.Float64("new_amount", *req.Amount),
		)
	}

	if req.Currency != nil {
		updatedTransaction.Currency = *req.Currency
		s.logger.Service("UpdateTransaction - updating currency",
			zap.String("old_currency", existingTransaction.Currency),
			zap.String("new_currency", *req.Currency),
		)
	}

	if req.Description != nil {
		updatedTransaction.Description = *req.Description
		s.logger.Service("UpdateTransaction - updating description")
	}

	if req.Category != nil {
		updatedTransaction.Category = *req.Category
		s.logger.Service("UpdateTransaction - updating category",
			zap.String("old_category", existingTransaction.Category),
			zap.String("new_category", *req.Category),
		)
	}

	if req.Date != nil {
		transactionDate, err := time.Parse("2006-01-02", *req.Date)
		if err != nil {
			s.logger.Error("service", "UpdateTransaction - date parsing failed", err,
				zap.String("date_string", *req.Date),
			)
			return nil, errors.New("invalid date format, use YYYY-MM-DD")
		}
		updatedTransaction.Date = transactionDate
		s.logger.Service("UpdateTransaction - updating date",
			zap.Time("old_date", existingTransaction.Date),
			zap.Time("new_date", transactionDate),
		)
	}

	s.logger.Service("UpdateTransaction - calling repository",
		zap.Int("transaction_id", id),
		zap.Any("updated_transaction", updatedTransaction),
	)

	repoStart := time.Now()
	err = s.repo.Update(&updatedTransaction)
	repoDuration := time.Since(repoStart)

	s.logger.Performance("UpdateTransaction repository call", repoDuration,
		zap.Bool("success", err == nil),
		zap.Int("transaction_id", id),
	)

	if err != nil {
		s.logger.Error("service", "UpdateTransaction - repository error", err,
			zap.Int("transaction_id", id),
		)
		return nil, err
	}

	totalDuration := time.Since(start)
	s.logger.Service("UpdateTransaction completed successfully",
		zap.Int("transaction_id", id),
		zap.Duration("total_duration", totalDuration),
		zap.Duration("repo_duration", repoDuration),
	)

	return &updatedTransaction, nil
}

func (s *transactionService) validateCreateRequest(req *models.CreateTransactionRequest) error {
	s.logger.Debug("service", "Validating create request",
		zap.Any("request", req),
	)

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

	s.logger.Debug("service", "Validation completed successfully")
	return nil
}

func (s *transactionService) validateUpdateRequest(req *models.UpdateTransactionRequest) error {
	s.logger.Debug("service", "Validating update request",
		zap.Any("request", req),
	)

	if req.Type != nil {
		if *req.Type != models.TransactionTypeExpense && *req.Type != models.TransactionTypeIncome {
			return errors.New("type must be 'expense' or 'income'")
		}
	}

	if req.Amount != nil && *req.Amount <= 0 {
		return errors.New("amount must be positive")
	}

	if req.Description != nil && *req.Description == "" {
		return errors.New("description cannot be empty")
	}

	if req.Category != nil && *req.Category == "" {
		return errors.New("category cannot be empty")
	}

	s.logger.Debug("service", "Update validation completed successfully")
	return nil
}