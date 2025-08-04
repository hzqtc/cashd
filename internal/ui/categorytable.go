package ui

import (
	"cashd/internal/data"
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
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

func (c categoryColumn) rightAligned() bool {
	return c == catColNumTxns || c == catColAmount
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
		return "Txn Num"
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

type categoryInfo struct {
	catType data.TransactionType
	symbol  string
	name    string
	numTxns int
	amount  float64
}

type CategoryTableModel struct {
	categories []*categoryInfo
	table      table.Model
}

func NewCategoryTableModel() CategoryTableModel {
	columns := []table.Column{}
	for i := range int(totalNumCatColumns) {
		col := categoryColumn(i)
		title := col.String()
		width := categoryColWidthMap[col]
		if col.rightAligned() {
			title = fmt.Sprintf("%*s", width, title)
		}
		columns = append(columns, table.Column{Title: title, Width: width})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithStyles(getTableStyle()),
	)
	return CategoryTableModel{table: t}
}

func (m CategoryTableModel) Update(msg tea.Msg) (CategoryTableModel, tea.Cmd) {
	// TODO: send msg when selected row changes
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m CategoryTableModel) View() string {
	return baseStyle.Render(m.table.View())
}

func (m *CategoryTableModel) SetDimensions(width, height int) {
	m.table.SetWidth(width)
	m.table.SetHeight(height)
}

func (m *CategoryTableModel) SetTransactions(transactions []*data.Transaction) {
	m.categories = getCategoryInfo(transactions)
	rows := make([]table.Row, len(m.categories))
	for i, c := range m.categories {
		rowData := []string{}
		for i := range int(totalNumCatColumns) {
			col := categoryColumn(i)
			colData := getCategoryColData(c, col)
			if col.rightAligned() {
				colData = fmt.Sprintf("%*s", categoryColWidthMap[col], colData)
			}
			rowData = append(rowData, colData)
		}
		rows[i] = table.Row(rowData)
	}
	m.table.SetRows(rows)
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

func getCategoryColData(c *categoryInfo, col categoryColumn) string {
	switch col {
	case catColSymbol:
		return c.symbol
	case catColType:
		return string(c.catType)
	case catColName:
		return c.name
	case catColNumTxns:
		return fmt.Sprintf("%d", c.numTxns)
	case catColAmount:
		return fmt.Sprintf("%.2f", c.amount)
	default:
		return ""
	}
}
