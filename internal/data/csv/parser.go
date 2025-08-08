package csv

import (
	"cashd/internal/data"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func parseCsvRecord(segments []string, config *config) *data.Transaction {
	txn := data.Transaction{}
	// Get reflect.Value of the struct pointer
	v := reflect.ValueOf(&txn).Elem()

	for _, f := range data.AllTransactionFields {
		index, ok := config.ColumnIndexes[f]
		if !ok {
			if f == "AccountType" {
				// Allow AccountType to be missing from CSV since we'd try to infer it from AccountName
				continue
			} else {
				panic(fmt.Sprintf("parse CSV failed: transaction field %s not found", f))
			}
		}
		if index >= len(segments) {
			panic(fmt.Sprintf("parse CSV failed: transaction field %s's not found", f))
		}
		value := segments[index]

		field := v.FieldByName(string(f))
		switch field.Kind() {
		case reflect.String:
			if field.Type().Name() == "TransactionType" {
				txnType, ok := config.TxnTypeMappings[strings.ToLower(value)]
				if !ok {
					panic(fmt.Sprintf("parse CSV failed: transaction type %s not recognized", value))
				}
				field.Set(reflect.ValueOf(data.TransactionType(txnType)))
			} else if field.Type().Name() == "AccountType" {
				txnType, ok := config.AccountTypeMappings[strings.ToLower(value)]
				if !ok {
					panic(fmt.Sprintf("parse CSV failed: account type %s not recognized", value))
				}
				field.Set(reflect.ValueOf(data.TransactionType(txnType)))
			} else {
				field.SetString(value)
			}
		case reflect.Float64:
			num, err := strconv.ParseFloat(value, 64)
			if err != nil {
				panic(fmt.Sprintf("parse CSV failed: %s is not float", value))
			}
			field.SetFloat(num)
		case reflect.Struct:
			// Special case: check if it's time.Time
			if field.Type() == reflect.TypeOf(time.Time{}) {
				var parsed time.Time
				var err error
				for _, format := range config.DateFormats {
					parsed, err = time.Parse(format, value)
					if err == nil {
						break
					}
				}
				if err != nil {
					panic(fmt.Sprintf("parse CSV failed: unsupported date format %s", value))
				}
				field.Set(reflect.ValueOf(parsed))
			} else {
				panic(fmt.Sprintf("parse CSV failed: unsupported struct type %s", field.Type()))
			}
		default:
			panic(fmt.Sprintf("parse CSV failed: unsupported field type %s", field.Kind()))
		}
	}

	// Account type missing, try parsing from account name
	if txn.AccountType == "" {
		if len(config.AccountTypeFromName) == 0 {
			panic("parse CSV failed: account name to type mapping is empty")
		}

		for accountNamePattern, accountType := range config.AccountTypeFromName {
			re := regexp.MustCompile(accountNamePattern)
			if re.MatchString(strings.ToLower(txn.Account)) {
				txn.AccountType = accountType
				break
			}
		}
		if txn.AccountType == "" {
			// Default to credit card
			log.Printf("Did not find a match in account name to type mapping for %s, default to credit card", txn.Account)
			txn.AccountType = data.AcctCreditCard
		}
	}

	if txn.IsValid() {
		return &txn
	} else {
		panic(fmt.Sprintf("parse CSV failed: transaction is incomplete: %v", txn))
	}
}
