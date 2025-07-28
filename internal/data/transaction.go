package data

import (
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

func (t Transaction) Symbol() string {
	switch t.Type {
	case Income:
		return incomeSymbol
	case Expense:
		return expensSymbol
	default:
		return ""
	}
}
