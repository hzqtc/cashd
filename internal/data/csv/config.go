package csv

import (
	"cashd/internal/data"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/spf13/pflag"
)

var csvConfigFlag string

func init() {
	pflag.StringVar(&csvConfigFlag, "csv-config", "", "CSV configuration json file path")
}

type config struct {
	Columns             map[string]data.TransactionField `json:"columns"`
	ColumnIndexes       map[data.TransactionField]int    `json:"column_indexes"`
	DateFormats         []string                         `json:"date_formats"`
	TxnTypeMappings     map[string]data.TransactionType  `json:"transaction_types"`
	AccountTypeMappings map[string]data.AccountType      `json:"account_types"`
	AccountTypeFromName map[string]data.AccountType      `json:"account_type_from_name"`
}

func getConfig() *config {
	if csvConfigFlag == "" {
		return defaultConfig
	}

	fileContent, err := os.ReadFile(csvConfigFlag)
	if err != nil {
		panic(fmt.Sprintf("Error reading CSV config file: %v", err))
	}

	var c config
	err = json.Unmarshal(fileContent, &c)
	if err != nil {
		panic(fmt.Sprintf("Error unmarshaling CSV config: %v", err))
	}

	cv := reflect.ValueOf(&c).Elem()
	dv := reflect.ValueOf(defaultConfig).Elem()
	for i := 0; i < cv.NumField(); i++ {
		f := cv.Field(i)
		if f.IsZero() {
			// Backfill missing config fields using the default
			f.Set(dv.Field(i))
		}
	}

	return &c
}

var defaultConfig = func() *config {
	c := &config{
		ColumnIndexes: map[data.TransactionField]int{},
		DateFormats: []string{
			time.DateOnly,
			time.DateTime,
		},
		TxnTypeMappings: map[string]data.TransactionType{
			"income":  data.Income,
			"inc.":    data.Income,
			"expense": data.Expense,
			"exp.":    data.Expense,
			"exps.":   data.Expense,
		},
		AccountTypeMappings: map[string]data.AccountType{},
		AccountTypeFromName: map[string]data.AccountType{
			"^cash$":             data.AcctCash,
			"checking$":          data.AcctBankAccount,
			"saving(s)?$":        data.AcctBankAccount,
			"credit(\\s?card)?$": data.AcctCreditCard,
		},
	}

	c.Columns = map[string]data.TransactionField{}
	for _, f := range data.AllTransactionFields {
		c.Columns[string(f)] = f
	}

	return c
}()
