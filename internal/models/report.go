package models

type MonthlyReport struct {
	Month        string             `json:"month"`
	Year         int                `json:"year"`
	TotalIncome  map[string]float64 `json:"total_income"`  // By currency
	TotalExpense map[string]float64 `json:"total_expense"` // By currency
	Balance      map[string]float64 `json:"balance"`       // By currency
	Transactions []Transaction      `json:"transactions"`
	Summary      ReportSummary      `json:"summary"`
}

type ReportSummary struct {
	TransactionCount  int                      `json:"transaction_count"`
	IncomeCount       int                      `json:"income_count"`
	ExpenseCount      int                      `json:"expense_count"`
	CategoryBreakdown map[string]CategoryTotal `json:"category_breakdown"`
}

type CategoryTotal struct {
	Count  int                `json:"count"`
	Totals map[string]float64 `json:"totals"` // By currency
}
