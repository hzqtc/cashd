package ui

import (
	"cashd/internal/data"
	"sort"

	"github.com/charmbracelet/bubbles/table"
)

type categoryColumn int

const (
	catColSymbol categoryColumn = iota
	catColType
	catColName
	catColNumTxns
	catColAmount

	totalNumCatColumns
)

func (c categoryColumn) index() int {
	return int(c)
}

func (c categoryColumn) rightAligned() bool {
	return c == catColNumTxns || c == catColAmount
}

func (c categoryColumn) isSortable() bool {
	return c != catColSymbol
}

func (c categoryColumn) width() int {
	return categoryColWidthMap[c]
}

func (c categoryColumn) nextColumn() column {
	return column(categoryColumn((int(c) + 1) % int(totalNumCatColumns)))
}

func (c categoryColumn) prevColumn() column {
	return column(categoryColumn((int(c) - 1 + int(totalNumCatColumns)) % int(totalNumCatColumns)))
}

func (c categoryColumn) getColumnData(a any) any {
	switch category := a.(*categoryInfo); c {
	case catColSymbol:
		return category.symbol
	case catColType:
		return string(category.catType)
	case catColName:
		return category.name
	case catColNumTxns:
		return category.numTxns
	case catColAmount:
		return category.amount
	default:
		return ""
	}
}

func (c categoryColumn) String() string {
	switch c {
	case catColSymbol:
		return " "
	case catColType:
		return "Type"
	case catColName:
		return "Category"
	case catColNumTxns:
		return "Txn #"
	case catColAmount:
		return "Amount"
	default:
		return "Unknown"
	}
}

var categoryColWidthMap = map[categoryColumn]int{
	catColSymbol:  symbolColWidth,
	catColType:    typeColWidth,
	catColName:    categoryColWidth,
	catColNumTxns: numberColWidth,
	catColAmount:  amountColWidth,
}

var CategoryTableWidth = func() int {
	tableWidth := 0
	for i := range totalNumCatColumns {
		tableWidth += categoryColWidthMap[categoryColumn(i)] + 2
	}
	return tableWidth
}()

var CategoryTableConfig = TableConfig{
	columns: func() []column {
		cols := []column{}
		for i := range int(totalNumCatColumns) {
			cols = append(cols, column(categoryColumn(i)))
		}
		return cols
	}(),
	dataProvider:      categoryTableDataProvider,
	rowId:             func(row table.Row) string { return row[catColName] },
	defaultSortColumn: column(catColName),
	defaultSortDir:    sortAsc,
}

type categoryInfo struct {
	catType data.TransactionType
	symbol  string
	name    string
	numTxns int
	amount  float64
}

func categoryTableDataProvider(transactions []*data.Transaction) tableDataSorter {
	categories := getCategoryInfo(transactions)
	result := make([]any, len(categories))
	for i, cat := range categories {
		result[i] = cat
	}

	return func(sortCol column, sortDir sortDirection) []any {
		sort.Slice(result, func(i, j int) bool {
			return compareAny(sortCol.getColumnData(result[i]), sortCol.getColumnData(result[j]), sortDir)
		})
		return result
	}
}

// Get category-level stats by aggregating transactions
func getCategoryInfo(transactions []*data.Transaction) []*categoryInfo {
	categoryMap := make(map[string]*categoryInfo)
	for _, tx := range transactions {
		cat, exist := categoryMap[tx.Category]
		if !exist {
			cat = &categoryInfo{
				symbol:  tx.Symbol(),
				catType: tx.Type,
				name:    tx.Category,
			}
			categoryMap[tx.Category] = cat
		}
		cat.numTxns++
		cat.amount += tx.Amount
	}

	categories := []*categoryInfo{}
	for _, c := range categoryMap {
		categories = append(categories, c)
	}
	return categories
}
