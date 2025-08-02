package ui

import (
	"cashd/internal/data"
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
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

func (c txnColumn) rightAligned() bool {
	return c == txnColAmount
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

type TransactionTableModel struct {
	table table.Model
}

func NewTransactionTableModel() TransactionTableModel {
	columns := []table.Column{}
	for i := range int(totalNumTxnColumns) {
		col := txnColumn(i)
		title := col.String()
		width := txnColWidthMap[col]
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

func (m *TransactionTableModel) SetDimensions(width, height int) {
	m.table.SetWidth(width)
	m.table.SetHeight(height)
}

func (m *TransactionTableModel) SetTransactions(transactions []*data.Transaction) {
	rows := make([]table.Row, len(transactions))
	for i, tx := range transactions {
		rowData := []string{}
		for i := range int(totalNumTxnColumns) {
			col := txnColumn(i)
			colData := getTxnColData(tx, col)
			if col.rightAligned() {
				colData = fmt.Sprintf("%*s", txnColWidthMap[col], colData)
			}
			rowData = append(rowData, colData)
		}
		rows[i] = table.Row(rowData)
	}
	m.table.SetRows(rows)
}

func getTxnColData(tx *data.Transaction, col txnColumn) string {
	switch col {
	case txnColSymbol:
		return tx.Symbol()
	case txnColDate:
		return tx.Date.Format("2006/01/02")
	case txnColType:
		return string(tx.Type)
	case txnColAccount:
		return tx.Account
	case txnColCategory:
		return tx.Category
	case txnColDesc:
		return tx.Description
	case txnColAmount:
		return fmt.Sprintf("$%.2f", tx.Amount)
	default:
		return ""
	}
}
