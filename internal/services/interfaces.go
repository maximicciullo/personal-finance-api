package services

import (
	"github.com/maximicciullo/personal-finance-api/internal/models"
)

type TransactionService interface {
	CreateTransaction(req *models.CreateTransactionRequest) (*models.Transaction, error)
	GetTransaction(id int) (*models.Transaction, error)
	GetTransactions(filters models.TransactionFilters) ([]models.Transaction, error)
	DeleteTransaction(id int) error
}

type ReportService interface {
	GetMonthlyReport(year, month int) (*models.MonthlyReport, error)
	GetCurrentMonthReport() (*models.MonthlyReport, error)
}