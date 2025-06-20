package repositories

import (
	"time"

	"github.com/maximicciullo/personal-finance-api/internal/models"
)

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	GetByID(id int) (*models.Transaction, error)
	GetAll() ([]models.Transaction, error)
	GetByFilters(filters models.TransactionFilters) ([]models.Transaction, error)
	GetByDateRange(startDate, endDate time.Time) ([]models.Transaction, error)
	Delete(id int) error
	Update(transaction *models.Transaction) error
}
