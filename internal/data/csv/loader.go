package csv

import (
	"cashd/internal/data"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

type CsvDataSource struct{}

var csvFileFlag string

func init() {
	pflag.StringVar(&csvFileFlag, "csv", "", "CSV file path")
}

func (s CsvDataSource) LoadTransactions() ([]*data.Transaction, error) {
	file, err := os.Open(csvFileFlag)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	config := getConfig()
	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		return []*data.Transaction{}, fmt.Errorf("failed to read CSV: %w", err)
	}
	if len(config.ColumnIndexes) == 0 {
		// Try to locate column index for each transaction field if not set
		for index, col := range header {
			if field, ok := config.Columns[col]; ok {
				config.ColumnIndexes[field] = index
			}
		}
		// Check all fields have index
		for _, field := range data.TransactionFields {
			if _, ok := config.ColumnIndexes[field]; field != "AccountType" && !ok {
				return []*data.Transaction{}, fmt.Errorf("failed to parse CSV: unable to locate column for transaction field %s", field)
			}
		}
	}

	records, _ := reader.ReadAll()
	txns := make([]*data.Transaction, len(records))
	for i, rec := range records {
		txns[i] = parseCsvRecord(rec, config)
	}

	return txns, nil
}

func (s CsvDataSource) Preferred() bool {
	return s.Enabled()
}

func (s CsvDataSource) Enabled() bool {
	return csvFileFlag != ""
}
