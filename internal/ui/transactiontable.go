package ui

import (
	"cashd/internal/data"
	"sort"
)

type txnColumn int

const (
	txnColSymbol txnColumn = iota
	txnColDate
	txnColType
	txnColAccount
	txnColCategory
	txnColDesc
	txnColAmount

	totalNumTxnColumns
)

func (c txnColumn) index() int {
	return int(c)
}

func (c txnColumn) rightAligned() bool {
	return c == txnColAmount
}

func (c txnColumn) isSortable() bool {
	return c != txnColSymbol
}

func (c txnColumn) width() int {
	return txnColWidthMap[c]
}

func (c txnColumn) nextColumn() column {
	return column(txnColumn((int(c) + 1) % int(totalNumTxnColumns)))
}

func (c txnColumn) prevColumn() column {
	return column(txnColumn((int(c) - 1 + int(totalNumTxnColumns)) % int(totalNumTxnColumns)))
}

func (c txnColumn) getColumnData(a any) any {
	switch txn := a.(*data.Transaction); c {
	case txnColSymbol:
		return txn.Symbol()
	case txnColDate:
		return txn.Date
	case txnColType:
		return string(txn.Type)
	case txnColAccount:
		return txn.Account
	case txnColCategory:
		return txn.Category
	case txnColDesc:
		return txn.Description
	case txnColAmount:
		return txn.Amount
	default:
		return ""
	}
}

func (c txnColumn) String() string {
	switch c {
	case txnColSymbol:
		return " "
	case txnColDate:
		return "Date"
	case txnColType:
		return "Type"
	case txnColAccount:
		return "Account"
	case txnColCategory:
		return "Category"
	case txnColDesc:
		return "Description"
	case txnColAmount:
		return "Amount"
	default:
		return "Unknown"
	}
}

var txnColWidthMap = map[txnColumn]int{
	txnColSymbol:   symbolColWidth,
	txnColDate:     dateColWidth,
	txnColType:     typeColWidth,
	txnColAccount:  accountColWidth,
	txnColCategory: categoryColWidth,
	txnColDesc:     descColWidth,
	txnColAmount:   amountColWidth,
}

var TxnTableWidth = func() int {
	tableWidth := 0
	for i := range totalNumTxnColumns {
		tableWidth += txnColWidthMap[txnColumn(i)] + 2
	}
	return tableWidth
}()

var TransactionTableConfig = TableConfig{
	columns: func() []column {
		cols := []column{}
		for i := range int(totalNumTxnColumns) {
			cols = append(cols, column(txnColumn(i)))
		}
		return cols
	}(),
	dataProvider:      txnTableDataProvider,
	defaultSortColumn: column(txnColDate),
	defaultSortDir:    sortAsc,
}

func txnTableDataProvider(transactions []*data.Transaction) tableDataSorter {
	result := make([]any, len(transactions))
	for i, txn := range transactions {
		result[i] = txn
	}

	return func(sortCol column, sortDir sortDirection) []any {
		sort.Slice(result, func(i, j int) bool {
			return compareAny(sortCol.getColumnData(result[i]), sortCol.getColumnData(result[j]), sortDir)
		})
		return result
	}
}
