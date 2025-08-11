package csv

import (
	"cashd/internal/data"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/pflag"
)

type CsvDataSource struct{}

var csvFiles []string

func init() {
	pflag.StringSliceVar(&csvFiles, "csv", []string{}, "CSV file path (can be specified multiple times)")
}

func (s CsvDataSource) LoadTransactions() ([]*data.Transaction, error) {
	allTxns := []*data.Transaction{}
	txnChan := make(chan []*data.Transaction)

	resolvedFilePaths := []string{}
	for _, pattern := range csvFiles {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve glob pattern %s: %w", pattern, err)
		}
		resolvedFilePaths = append(resolvedFilePaths, matches...)
	}
	if len(resolvedFilePaths) == 0 {
		return []*data.Transaction{}, nil
	}

	errChan := make(chan error, len(resolvedFilePaths))

	for _, filePath := range resolvedFilePaths {
		go func(fp string) {
			txns, err := readCsv(fp)
			if err != nil {
				errChan <- err
			} else {
				txnChan <- txns
			}
		}(filePath)
	}

	for i := 0; i < len(resolvedFilePaths); i++ {
		select {
		case txns := <-txnChan:
			allTxns = append(allTxns, txns...)
		case err := <-errChan:
			return nil, err
		}
	}
	sort.Slice(allTxns, func(i, j int) bool {
		return allTxns[i].Date.Before(allTxns[j].Date)
	})
	return allTxns, nil
}

func readCsv(filePath string) ([]*data.Transaction, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	config := getConfig()
	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header from %s: %w", filePath, err)
	}

	if len(config.ColumnIndexes) == 0 {
		// Try to locate column index for each transaction field if not set
		for index, col := range header {
			if field, ok := config.Columns[col]; ok {
				config.ColumnIndexes[field] = index
			}
		}
		// Check each field has an index except for AccountType which can be inferred from Account
		for _, field := range data.AllTransactionFields {
			if _, ok := config.ColumnIndexes[field]; field != "AccountType" && !ok {
				return nil, fmt.Errorf("failed to parse CSV from %s: unable to locate column for transaction field %s", filePath, field)
			}
		}
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read all records from %s: %w", filePath, err)
	}

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
	return len(csvFiles) > 0
}
