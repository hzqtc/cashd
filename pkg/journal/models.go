package journal

import (
	"time"
)

type TransactionType string

const (
	Income  TransactionType = "Income"
	Expense TransactionType = "Expense"
)

type Transaction struct {
	Date        time.Time
	Type        TransactionType
	Account     string
	Category    string
	Amount      float64
	Description string
}
