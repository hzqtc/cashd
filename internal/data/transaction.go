package data

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type TransactionType string

const (
	// TODO: support more transaction types
	Income  TransactionType = "Income"
	Expense TransactionType = "Expense"
)

const (
	incomeSymbol = "󱙹"
	expensSymbol = ""
)

type AccountType string

const (
	// TODO: support more account types
	AcctCash        AccountType = "Cash"
	AcctBankAccount AccountType = "Bank Account"
	AcctCreditCard  AccountType = "Credit Card"
	AcctOverall     AccountType = ""
)

const (
	cashSymbol = "󰄔"
	bankSymbol = "󰁰"
	cardSymbol = "󰆛"
)

type Transaction struct {
	Date        time.Time
	Type        TransactionType
	AccountType AccountType
	Account     string
	Category    string
	Amount      float64
	Description string
}

type TxnField string

var TransactionFields = func() []TxnField {
	fields := []TxnField{}
	t := reflect.TypeOf(Transaction{})
	for i := range t.NumField() {
		fields = append(fields, TxnField(t.Field(i).Name))
	}
	return fields
}()

func (t *Transaction) IsValid() bool {
	return !t.Date.IsZero() &&
		t.Type != "" &&
		t.AccountType != "" &&
		t.Account != "" &&
		t.Category != "" &&
		t.Amount > 0 &&
		t.Description != ""
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

func (t *Transaction) AccountSymbol() string {
	switch t.AccountType {
	case AcctCash:
		return cashSymbol
	case AcctBankAccount:
		return bankSymbol
	case AcctCreditCard:
		return cardSymbol
	default:
		return ""
	}
}

func (t *Transaction) Matches(kws []string) bool {
	for _, kw := range kws {
		// Requires the transaction to contain ALL keywords
		// So we can return false on any unmatched keyword
		if !strings.Contains(t.Date.Format("2006-01-02"), kw) &&
			!strings.Contains(strings.ToLower(string(t.Type)), kw) &&
			!strings.Contains(strings.ToLower(t.Account), kw) &&
			!strings.Contains(strings.ToLower(t.Category), kw) &&
			!strings.Contains(fmt.Sprintf("%.2f", t.Amount), kw) &&
			!strings.Contains(strings.ToLower(t.Description), kw) {
			return false
		}
	}
	return true
}
