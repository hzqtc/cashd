package data

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
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
func ParseJournal(filePath string) ([]Transaction, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open journal file: %w", err)
	}
	defer file.Close()

	var transactions []Transaction
	scanner := bufio.NewScanner(file)

	// Regex to match transaction header: YYYY-MM-DD Description
	transactionHeaderRegex := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})\s+(.*)$`)
	// Regex to match account line, with the following variations
	// transaction type:category amount, e.g. expenses:Utilities 49.99, income:Cash Back $-47.11
	// account type:account amount, e.g. liability:BoA 123 $-49.99, assets:BoA Checking $47.11
	accountLineRegex := regexp.MustCompile(`^\s+(.+):(.+)\s{8}(\$[-]?[\d,]+\.?\d*)$`)

	var transactionDate time.Time
	var transactionDesc string
	var transactionType TransactionType
	var transactionCategory string
	var transactionAccount string
	var transactionAmount float64
	var transactionLines int

	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("Processing line: %s\n", line) // Debug print

		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			// Skip empty lines and comments
			continue
		}

		if matches := transactionHeaderRegex.FindStringSubmatch(line); len(matches) > 0 {
			// This is a new transaction header
			dateStr := matches[1]
			description := matches[2]

			parsedDate, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse date %q: %w", dateStr, err)
			}
			transactionDate = parsedDate
			transactionDesc = description
			log.Printf("  Matched transaction header. Date: %s, Desc: %s\n", dateStr, description) // Debug print
		} else if matches := accountLineRegex.FindStringSubmatch(line); len(matches) > 0 {
			// This is an account or category line for the current transaction
			transactionLines++
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
				transactionType = Expense
				transactionCategory = accountOrCategory
			case income:
				transactionType = Income
				transactionCategory = accountOrCategory
			case assets, liability:
				transactionAccount = accountOrCategory
			}

			// Each transaction have 2 lines
			if transactionLines == 2 {
				t := Transaction{
					Date:        transactionDate,
					Type:        transactionType,
					Account:     transactionAccount,
					Category:    transactionCategory,
					Amount:      transactionAmount,
					Description: transactionDesc,
				}
				transactions = append(transactions, t)
				log.Printf("  Matched account line: %v\n", t) // Debug print
				transactionLines = 0
			}
		} else {
			log.Printf("  No regex match for line.\n") // Debug print
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading journal file: %w", err)
	}

	return transactions, nil
}
