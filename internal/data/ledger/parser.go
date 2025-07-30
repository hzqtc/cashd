package ledger

import (
	"bufio"
	"fmt"
	"io"
	"lledger/internal/data"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	assets    = "assets"
	liability = "liability"
	expenses  = "expenses"
	income    = "income"
)

// ParseJournal reads the hledger journal file and parses transactions.
func parseJournal(reader io.ReadCloser) ([]data.Transaction, error) {
	var transactions []data.Transaction
	scanner := bufio.NewScanner(reader)

	// Regex to match transaction header: YYYY-MM-DD Description
	// TODO: add support for '!' or '*' or '(code)' after date
	// TODO: add support for '| note'
	transactionHeaderRegex := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})\s+(.*)$`)
	// Regex to match account line, with the following variations
	// transaction type:category amount, e.g. expenses:Utilities 49.99, income:Cash Back $-47.11
	// account type:account amount, e.g. liability:BoA 123 $-49.99, assets:BoA Checking $47.11
	// TODO: add support for commodity other than '$'
	accountLineRegex := regexp.MustCompile(`^\s+(.+):(.+)\s{2,}\$[-]?([\d,]+\.?\d*)$`)

	var transactionDate time.Time
	var transactionDesc string
	var transactionType data.TransactionType
	var transactionCategory string
	var transactionAccount string
	var transactionAmount float64
	var numPostings int

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			// Skip empty lines and comments
			continue
		}

		if matches := transactionHeaderRegex.FindStringSubmatch(line); len(matches) > 0 {
			// This is a new transaction header
			numPostings = 0
			dateStr := matches[1]
			description := matches[2]

			parsedDate, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
			if err != nil {
				return nil, fmt.Errorf("failed to parse date %q: %w", dateStr, err)
			}
			transactionDate = parsedDate
			transactionDesc = description
		} else if matches := accountLineRegex.FindStringSubmatch(line); len(matches) > 0 {
			// This is a posting for the current transaction
			numPostings++
			typeStr := strings.TrimSpace(matches[1])             // account type (assets, liability) or transactionType (income, expense)
			accountOrCategory := strings.TrimSpace(matches[2])   // account or cateogory
			amountStr := strings.ReplaceAll(matches[3], ",", "") // Remove commas for parsing
			amountStr = strings.ReplaceAll(amountStr, "$", "")   // Remove dollar sign

			amount, err := strconv.ParseFloat(amountStr, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse amount %q: %w", amountStr, err)
			}
			transactionAmount = math.Abs(amount)

			switch typeStr {
			case expenses:
				transactionType = data.Expense
				transactionCategory = accountOrCategory
			case income:
				transactionType = data.Income
				transactionCategory = accountOrCategory
			case assets, liability:
				transactionAccount = accountOrCategory
			}

			// Currently only handling transactions that involves exactly 2 postings
			// TODO: add support for more than 2 postings
			if numPostings == 2 {
				t := data.Transaction{
					Date:        transactionDate,
					Type:        transactionType,
					Account:     transactionAccount,
					Category:    transactionCategory,
					Amount:      transactionAmount,
					Description: transactionDesc,
				}
				transactions = append(transactions, t)
			}
		} else {
			log.Printf("   Skipping line: %s\n", line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error parsing hledger journal: %w", err)
	}

	return transactions, nil
}
