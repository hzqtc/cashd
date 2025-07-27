package ui

import (
	"fmt"
	"lledger/internal/data"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type columnName int

const (
	colUnknown columnName = -1
)

const (
	colSymbol columnName = iota
	colDate
	colType
	colAccount
	colCategory
	colDesc
	colAmount

	totalNumColumns
)

func (c columnName) rightAligned() bool {
	return c == colAmount
}

func (c columnName) String() string {
	switch c {
	case colSymbol:
		return " "
	case colDate:
		return "Date"
	case colType:
		return "Type"
	case colAccount:
		return "Account"
	case colCategory:
		return "Category"
	case colDesc:
		return "Desc"
	case colAmount:
		return "Amount"
	default:
		return "Unknown"
	}
}

var colWidthMap = map[columnName]int{
	colSymbol:   2,
	colDate:     12,
	colType:     10,
	colAccount:  25,
	colCategory: 15,
	colDesc:     20,
	colAmount:   12,
}

type TransactionTableModel struct {
	table table.Model
}

func NewTransactionTableModel() TransactionTableModel {
	columns := []table.Column{}
	for i := range int(totalNumColumns) {
		col := columnName(i)
		title := col.String()
		width := colWidthMap[col]
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
	return TransactionTableModel{table: t}
}

func (m TransactionTableModel) Update(msg tea.Msg) (TransactionTableModel, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m TransactionTableModel) View() string {
	return baseStyle.Render(m.table.View())
}

func (m *TransactionTableModel) SetDimensions(height, width int) {
	m.table.SetHeight(height)
	m.table.SetWidth(width)
}

func (m *TransactionTableModel) SetTransactions(transactions []*data.Transaction) {
	rows := make([]table.Row, len(transactions))
	for i, tx := range transactions {
		rowData := []string{}
		for i := range int(totalNumColumns) {
			col := columnName(i)
			colData := getColData(tx, col)
			if col.rightAligned() {
				colData = fmt.Sprintf("%*s", colWidthMap[col], colData)
			}
			rowData = append(rowData, colData)
		}
		rows[i] = table.Row(rowData)
	}
	m.table.SetRows(rows)
}

func getColData(tx *data.Transaction, col columnName) string {
	switch col {
	case colSymbol:
		return tx.Symbol()
	case colDate:
		return tx.Date.Format("2006/01/02")
	case colType:
		return string(tx.Type)
	case colAccount:
		return tx.Account
	case colCategory:
		return tx.Category
	case colDesc:
		return tx.Description
	case colAmount:
		return fmt.Sprintf("$%.2f", tx.Amount)
	default:
		return ""
	}
}
