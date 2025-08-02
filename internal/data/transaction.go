package data

import (
	"fmt"
	"strings"
	"time"
)

type TransactionType string

const (
	Income  TransactionType = "Income"
	Expense TransactionType = "Expense"
)

const (
	incomeSymbol = "󱙹"
	expensSymbol = ""
)

type Transaction struct {
	Date        time.Time
	Type        TransactionType
	Account     string
	Category    string
	Amount      float64
	Description string
}

func (t *Transaction) Symbol() string {
	switch t.Type {
	case Income:
		return incomeSymbol
	case Expense:
		return expensSymbol
	default:
		return ""
	}
}

func (t *Transaction) Matches(kws []string) bool {
	for _, kw := range kws {
		// Requires the transaction to contain ALL keywords
		// So we can return false on any unmatched keyword
		if !strings.Contains(t.Date.Format("2006-01-02"), kw) &&
			!strings.Contains(strings.ToLower(t.Account), kw) &&
			!strings.Contains(strings.ToLower(t.Category), kw) &&
			!strings.Contains(fmt.Sprintf("%.2f", t.Amount), kw) &&
			!strings.Contains(strings.ToLower(t.Description), kw) {
			return false
		}
	}
	return true
}
