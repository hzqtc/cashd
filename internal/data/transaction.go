package data

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
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

func (t *TransactionType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("failed to unmarshal TransactionType: %w", err)
	}

	if s == string(Income) || s == string(Expense) {
		*t = TransactionType(s)
		return nil
	} else {
		return fmt.Errorf("invalid TransactionType: %s", s)
	}
}

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

func (a *AccountType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("failed to unmarshal AccountType: %w", err)
	}

	if s == string(AcctCash) || s == string(AcctBankAccount) || s == string(AcctCreditCard) {
		*a = AccountType(s)
		return nil
	} else {
		return fmt.Errorf("invalid AccountType: %s", s)
	}
}

type Transaction struct {
	Date        time.Time
	Type        TransactionType
	AccountType AccountType
	Account     string
	Category    string
	Amount      float64
	Description string
}

type TransactionField string

var AllTransactionFields = func() []TransactionField {
	fields := []TransactionField{}
	t := reflect.TypeOf(Transaction{})
	for i := range t.NumField() {
		fields = append(fields, TransactionField(t.Field(i).Name))
	}
	return fields
}()

func (t *TransactionField) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("failed to unmarshal TransactionField: %w", err)
	}

	for _, f := range AllTransactionFields {
		if s == string(f) {
			*t = TransactionField(s)
			return nil
		}
	}
	return fmt.Errorf("invalid TransactionField: %s", s)
}

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
		if kw, negative := strings.CutPrefix(kw, negativeKeywordPrefix); negative {
			if t.matchKeyword(kw) {
				return false
			}
		} else if !t.matchKeyword(kw) {
			return false
		}
	}
	return true
}

func (t *Transaction) matchKeyword(kw string) bool {
	if trimmed, hasPrefix := strings.CutPrefix(kw, kwPrefixDate); hasPrefix {
		return t.matchesDate(trimmed)
	}
	if trimmed, hasPrefix := strings.CutPrefix(kw, kwPrefixType); hasPrefix {
		return t.matchesType(trimmed)
	}
	if trimmed, hasPrefix := strings.CutPrefix(kw, kwPrefixAccount); hasPrefix {
		return t.matchesAccount(trimmed)
	}
	if trimmed, hasPrefix := strings.CutPrefix(kw, kwPrefixCategory); hasPrefix {
		return t.matchesCategory(trimmed)
	}
	if trimmed, hasPrefix := strings.CutPrefix(kw, kwPrefixAmount); hasPrefix {
		return t.matchesAmount(trimmed)
	}
	if trimmed, hasPrefix := strings.CutPrefix(kw, kwPrefixDesc); hasPrefix {
		return t.matchesDescription(trimmed)
	}
	// Check all fields for unprefix keywords
	return t.matchesDate(kw) ||
		t.matchesType(kw) ||
		t.matchesAccount(kw) ||
		t.matchesCategory(kw) ||
		t.matchesAmount(kw) ||
		t.matchesDescription(kw)
}

const (
	negativeKeywordPrefix = "-"

	kwPrefixDate     = "d:"
	kwPrefixType     = "t:"
	kwPrefixAccount  = "a:"
	kwPrefixCategory = "c:"
	kwPrefixAmount   = "m:"
	kwPrefixDesc     = "p:"
)

type matchOp string

const (
	noneOp  matchOp = ""
	larger  matchOp = ">"
	smaller matchOp = "<"
)

var (
	dateFormats = []string{"2006-01-02", "2006-01", "2006"}
	// Matches optional < or >, then date in one of the formats anove
	dateRegexp = regexp.MustCompile(`^([<>])?(\d{4}(?:-\d{2}(?:-\d{2})?)?)$`)
	// Matches optional < or >, then a float number
	amountRegexp = regexp.MustCompile(`^([<>])?([0-9]*\.?[0-9]+)?$`)
)

func (t *Transaction) matchesDate(kw string) bool {
	op, dateStr := parseDateKeyword(kw)
	if op == noneOp {
		// No operator, just string matching
		if dateStr != "" {
			return strings.Contains(t.Date.Format(dateFormats[0]), dateStr)
		} else {
			return false
		}
	} else {
		// Date comparison using the operator
		date := parseDate(dateStr)
		if date == nil {
			// Should not happen since date format is ensured by regexp
			return false
		}
		switch op {
		case larger:
			return t.Date.After(*date)
		case smaller:
			return t.Date.Before(*date)
		default:
			panic(fmt.Sprintf("unexpected search keyword operator: %s", op))
		}
	}
}

func parseDateKeyword(kw string) (matchOp, string) {
	matches := dateRegexp.FindStringSubmatch(kw)
	if matches == nil {
		return noneOp, ""
	} else {
		return matchOp(matches[1]), matches[2]
	}
}

func parseDate(dateStr string) *time.Time {
	for _, f := range dateFormats {
		if d, err := time.Parse(f, dateStr); err == nil {
			return &d
		}
	}
	return nil
}

func (t *Transaction) matchesType(kw string) bool {
	return strings.Contains(strings.ToLower(string(t.Type)), kw)
}

func (t *Transaction) matchesAccount(kw string) bool {
	return strings.Contains(strings.ToLower(t.Account), kw)
}

func (t *Transaction) matchesCategory(kw string) bool {
	return strings.Contains(strings.ToLower(t.Category), kw)
}

func (t *Transaction) matchesAmount(kw string) bool {
	op, num := parseAmountKeyword(kw)
	if op == noneOp {
		// No operator, just string matching
		if num != "" {
			return strings.Contains(fmt.Sprintf("%.2f", t.Amount), num)
		} else {
			return false
		}
	} else {
		// Value comparison using the operator
		// Ignoring the error since format is ensured by regexp match
		targetAmount, _ := strconv.ParseFloat(num, 64)
		switch op {
		case larger:
			return t.Amount > targetAmount
		case smaller:
			return t.Amount < targetAmount
		default:
			panic(fmt.Sprintf("unexpected search keyword operator: %s", op))
		}
	}
}

func parseAmountKeyword(kw string) (matchOp, string) {
	matches := amountRegexp.FindStringSubmatch(kw)
	if matches == nil {
		return noneOp, "" // no match at all
	} else {
		return matchOp(matches[1]), matches[2]
	}
}

func (t *Transaction) matchesDescription(kw string) bool {
	return strings.Contains(strings.ToLower(t.Description), kw)
}
