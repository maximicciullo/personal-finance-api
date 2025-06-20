package models

import "time"

type Transaction struct {
	ID          int       `json:"id"`
	Type        string    `json:"type"` // "expense" or "income"
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"` // "ARS", "USD", etc.
	Description string    `json:"description"`
	Category    string    `json:"category"` // "food", "salary", "rent", etc.
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTransactionRequest struct {
	Type        string  `json:"type" binding:"required,oneof=expense income"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Currency    string  `json:"currency"`
	Description string  `json:"description" binding:"required"`
	Category    string  `json:"category" binding:"required"`
	Date        *string `json:"date,omitempty"` // Optional, format: YYYY-MM-DD
}

type TransactionFilters struct {
	Type     string
	Category string
	Currency string
	FromDate *time.Time
	ToDate   *time.Time
}

const (
	TransactionTypeExpense = "expense"
	TransactionTypeIncome  = "income"
)

const (
	CurrencyARS = "ARS"
	CurrencyUSD = "USD"
	CurrencyEUR = "EUR"
)
