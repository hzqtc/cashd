package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	numLines := flag.Int("lines", 1000, "Number of CSV lines to generate")
	flag.Parse()

	file, err := os.Create("sample.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Period", "Accounts", "Category", "Note", "USD", "Income/Expense"}
	if err := writer.Write(header); err != nil {
		panic(err)
	}

	// Sample data pools
	accounts := []string{"Cash", "Checking", "Savings", "Credit Card"}
	incomeCategories := []string{"Salary", "Bonus", "Sales", "Stock"}
	expenseCategories := []string{"Food", "Transport", "Utilities", "Entertainment", "Misc"}
	types := []string{"Income", "Expense"}

	for i := range *numLines {
		date := randomDate()
		account := accounts[rand.Intn(len(accounts))]
		typ := types[rand.Intn(len(types))]
		var category string
		if typ == "Income" {
			category = incomeCategories[rand.Intn(len(incomeCategories))]
		} else {
			category = expenseCategories[rand.Intn(len(expenseCategories))]
		}
		description := fmt.Sprintf("Sample transaction %d", i+1)
		amount := fmt.Sprintf("%.2f", rand.Float64()*1000)

		record := []string{date, account, category, description, amount, typ}
		if err := writer.Write(record); err != nil {
			panic(err)
		}
	}

	fmt.Printf("CSV file 'sample.csv' created with %d lines\n", *numLines)
}

func randomDate() string {
	// Generate random date between Jan 1, 2020 and today
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Now()

	diff := end.Sub(start)
	randSec := time.Duration(rand.Int63n(int64(diff)))

	randomTime := start.Add(randSec)
	return randomTime.Format("2006-01-02")
}
