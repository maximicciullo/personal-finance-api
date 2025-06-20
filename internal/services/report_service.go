package services

import (
	"errors"
	"time"

	"github.com/maximicciullo/personal-finance-api/internal/middleware"
	"github.com/maximicciullo/personal-finance-api/internal/models"
	"github.com/maximicciullo/personal-finance-api/internal/repositories"
	"go.uber.org/zap"
)

type reportService struct {
	repo   repositories.TransactionRepository
	logger *middleware.BusinessLoggerInstance
}

func NewReportService(repo repositories.TransactionRepository) ReportService {
	return &reportService{
		repo:   repo,
		logger: middleware.BusinessLogger(),
	}
}

func (s *reportService) GetMonthlyReport(year, month int) (*models.MonthlyReport, error) {
	s.logger.Service("GetMonthlyReport started",
		zap.Int("year", year),
		zap.Int("month", month),
	)

	if year < 1900 || year > time.Now().Year()+10 {
		err := errors.New("invalid year")
		s.logger.Error("service", "GetMonthlyReport - invalid year", err,
			zap.Int("year", year),
		)
		return nil, err
	}

	if month < 1 || month > 12 {
		err := errors.New("month must be between 1 and 12")
		s.logger.Error("service", "GetMonthlyReport - invalid month", err,
			zap.Int("month", month),
		)
		return nil, err
	}

	// Calculate date range for the month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	s.logger.Service("GetMonthlyReport - date range calculated",
		zap.Time("start_date", startDate),
		zap.Time("end_date", endDate),
	)

	// Get transactions for the month
	repoStart := time.Now()
	transactions, err := s.repo.GetByDateRange(startDate, endDate)
	repoDuration := time.Since(repoStart)

	s.logger.Performance("GetMonthlyReport repository call", repoDuration,
		zap.Int("transaction_count", len(transactions)),
		zap.Bool("success", err == nil),
	)

	if err != nil {
		s.logger.Error("service", "GetMonthlyReport - repository error", err,
			zap.Int("year", year),
			zap.Int("month", month),
		)
		return nil, err
	}

	s.logger.Service("GetMonthlyReport - building report",
		zap.Int("transaction_count", len(transactions)),
	)

	buildStart := time.Now()
	report := s.buildMonthlyReport(year, month, transactions)
	buildDuration := time.Since(buildStart)

	s.logger.Performance("GetMonthlyReport report building", buildDuration,
		zap.Int("currencies_count", len(report.TotalIncome)),
		zap.Int("categories_count", len(report.Summary.CategoryBreakdown)),
	)

	totalDuration := time.Since(repoStart)
	s.logger.Service("GetMonthlyReport completed successfully",
		zap.Int("year", year),
		zap.Int("month", month),
		zap.Duration("total_duration", totalDuration),
		zap.Duration("repo_duration", repoDuration),
		zap.Duration("build_duration", buildDuration),
	)

	return report, nil
}

func (s *reportService) GetCurrentMonthReport() (*models.MonthlyReport, error) {
	now := time.Now()
	s.logger.Service("GetCurrentMonthReport started",
		zap.Int("current_year", now.Year()),
		zap.Int("current_month", int(now.Month())),
	)
	
	return s.GetMonthlyReport(now.Year(), int(now.Month()))
}

func (s *reportService) buildMonthlyReport(year, month int, transactions []models.Transaction) *models.MonthlyReport {
	s.logger.Debug("service", "Building monthly report",
		zap.Int("year", year),
		zap.Int("month", month),
		zap.Int("transaction_count", len(transactions)),
	)

	totalIncome := make(map[string]float64)
	totalExpense := make(map[string]float64)
	categoryBreakdown := make(map[string]models.CategoryTotal)

	incomeCount := 0
	expenseCount := 0

	for _, transaction := range transactions {
		s.logger.Debug("service", "Processing transaction",
			zap.Int("transaction_id", transaction.ID),
			zap.String("type", transaction.Type),
			zap.Float64("amount", transaction.Amount),
			zap.String("currency", transaction.Currency),
			zap.String("category", transaction.Category),
		)

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

	s.logger.Debug("service", "Calculating balances",
		zap.Int("currencies_count", len(allCurrencies)),
	)

	for currency := range allCurrencies {
		balance[currency] = totalIncome[currency] - totalExpense[currency]
		s.logger.Debug("service", "Currency balance calculated",
			zap.String("currency", currency),
			zap.Float64("income", totalIncome[currency]),
			zap.Float64("expense", totalExpense[currency]),
			zap.Float64("balance", balance[currency]),
		)
	}

	report := &models.MonthlyReport{
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

	s.logger.Debug("service", "Monthly report built successfully",
		zap.String("month", report.Month),
		zap.Int("year", report.Year),
		zap.Int("total_transactions", report.Summary.TransactionCount),
		zap.Int("income_transactions", report.Summary.IncomeCount),
		zap.Int("expense_transactions", report.Summary.ExpenseCount),
	)

	return report
}

func (s *reportService) getAllCurrencies(totalIncome, totalExpense map[string]float64) map[string]bool {
	currencies := make(map[string]bool)

	for currency := range totalIncome {
		currencies[currency] = true
	}

	for currency := range totalExpense {
		currencies[currency] = true
	}

	s.logger.Debug("service", "All currencies extracted",
		zap.Strings("currencies", func() []string {
			keys := make([]string, 0, len(currencies))
			for k := range currencies {
				keys = append(keys, k)
			}
			return keys
		}()),
	)

	return currencies
}