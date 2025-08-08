package csv

import "cashd/internal/data"

type config struct {
	// Mapping for csv column names to Transaction struct fields
	Columns map[string]data.TxnField `json:"columns"`
	// Mapping from csv column index to Transaction struct fields
	ColumnIndexes map[data.TxnField]int `json:"column_indexes"`
	// Formating string in standard go date formats for parsing Transaction.Date
	DateFormats []string `json:"dateFormats"`
	// Value mapping for Tranaction.Type field from csv data to "Income" or "Expense"
	TxnTypeMappings map[string]data.TransactionType `json:"transaction_types"`
	// Value mapping for Tranaction.AccountType field from csv data to "Cash", "Bank Account" or "Credit Card"
	AccountTypeMappings map[string]data.AccountType `json:"account_types"`
	// Mapping from Transaction.Account to Transaction.AccountType if the csv does not have a column for account type
	AccountTypeFromName map[string]data.AccountType `json:"account_type_from_name"`
}

func getConfig() *config {
	// TODO: load from a config file
	return defaultConfig()
}

func defaultConfig() *config {
	c := &config{
		Columns:       map[string]data.TxnField{}, // Would be filled below
		ColumnIndexes: map[data.TxnField]int{},
		DateFormats: []string{
			"2006-01-02",
			"2006-01-02 15:04:05",
			"01/02/2006",
			"01/02/2006 15:04:05",
		},
		TxnTypeMappings: map[string]data.TransactionType{
			"income":  data.Income,
			"inc.":    data.Income,
			"expense": data.Expense,
			"exp.":    data.Expense,
			"exps.":   data.Expense,
		},
		AccountTypeMappings: map[string]data.AccountType{
			"cash":             data.AcctCash,
			"bank account":     data.AcctBankAccount,
			"bankaccount":      data.AcctBankAccount,
			"bank":             data.AcctBankAccount,
			"checking account": data.AcctBankAccount,
			"checking":         data.AcctBankAccount,
			"saving account":   data.AcctBankAccount,
			"saving":           data.AcctBankAccount,
			"credit card":      data.AcctCreditCard,
			"creditcard":       data.AcctCreditCard,
			"credit":           data.AcctCreditCard,
			"cc":               data.AcctCreditCard,
		},
		AccountTypeFromName: map[string]data.AccountType{
			"^cash$":    data.AcctCash,
			"checking$": data.AcctBankAccount,
			"saving$":   data.AcctBankAccount,
		},
	}

	for _, f := range data.TransactionFields {
		c.Columns[string(f)] = f
	}

	return c
}
