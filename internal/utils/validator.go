package utils

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// ISO 4217 currency codes pattern
	currencyPattern = regexp.MustCompile(`^[A-Z]{3}$`)
)

func ValidateCurrency(currency string) error {
	if currency == "" {
		return nil // Optional field
	}

	currency = strings.ToUpper(currency)
	if !currencyPattern.MatchString(currency) {
		return errors.New("currency must be a valid 3-letter ISO code (e.g., USD, ARS, EUR)")
	}

	return nil
}

func ValidateTransactionType(transactionType string) error {
	if transactionType != "expense" && transactionType != "income" {
		return errors.New("transaction type must be 'expense' or 'income'")
	}

	return nil
}

func ValidateAmount(amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	return nil
}

func ValidateRequiredString(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New(fieldName + " is required")
	}

	return nil
}