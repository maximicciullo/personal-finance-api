package services

import (
	"errors"
	"github.com/maximicciullo/personal-finance-api/internal/models"
	"github.com/maximicciullo/personal-finance-api/internal/repositories"
	"time"
)

type reportService struct {
	repo repositories.TransactionRepository
}

func NewReportService(repo repositories.TransactionRepository) ReportService {
	return &reportService{
		repo: repo,
	}
}

func (s *reportService) GetMonthlyReport(year, month int) (*models.MonthlyReport, error) {
	if year < 1900 || year > time.Now().Year()+10 {
		return nil, errors.New("invalid year")
	}

	if month < 1 || month > 12 {
		return nil, errors.New("month must be between 1 and 12")
	}

	// Calculate date range for the month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	// Get transactions for the month
	transactions, err := s.repo.GetByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	return s.buildMonthlyReport(year, month, transactions), nil
}

func (s *reportService) GetCurrentMonthReport() (*models.MonthlyReport, error) {
	now := time.Now()
	return s.GetMonthlyReport(now.Year(), int(now.Month()))
}

func (s *reportService) buildMonthlyReport(year, month int, transactions []models.Transaction) *models.MonthlyReport {
	totalIncome := make(map[string]float64)
	totalExpense := make(map[string]float64)
	categoryBreakdown := make(map[string]models.CategoryTotal)

	incomeCount := 0
	expenseCount := 0

	for _, transaction := range transactions {
		// Calculate totals by currency
		if transaction.Type == models.TransactionTypeIncome {
			totalIncome[transaction.Currency] += transaction.Amount
			incomeCount++
		} else {
			totalExpense[transaction.Currency] += transaction.Amount
			expenseCount++
		}

		// Category breakdown
		if category, exists := categoryBreakdown[transaction.Category]; exists {
			category.Count++
			if category.Totals[transaction.Currency] == 0 {
				category.Totals[transaction.Currency] = 0
			}
			category.Totals[transaction.Currency] += transaction.Amount
			categoryBreakdown[transaction.Category] = category
		} else {
			categoryBreakdown[transaction.Category] = models.CategoryTotal{
				Count:  1,
				Totals: map[string]float64{transaction.Currency: transaction.Amount},
			}
		}
	}

	// Calculate balance by currency
	balance := make(map[string]float64)
	allCurrencies := s.getAllCurrencies(totalIncome, totalExpense)

	for currency := range allCurrencies {
		balance[currency] = totalIncome[currency] - totalExpense[currency]
	}

	return &models.MonthlyReport{
		Month:        time.Month(month).String(),
		Year:         year,
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		Balance:      balance,
		Transactions: transactions,
		Summary: models.ReportSummary{
			TransactionCount:  len(transactions),
			IncomeCount:       incomeCount,
			ExpenseCount:      expenseCount,
			CategoryBreakdown: categoryBreakdown,
		},
	}
}

func (s *reportService) getAllCurrencies(totalIncome, totalExpense map[string]float64) map[string]bool {
	currencies := make(map[string]bool)

	for currency := range totalIncome {
		currencies[currency] = true
	}

	for currency := range totalExpense {
		currencies[currency] = true
	}

	return currencies
}