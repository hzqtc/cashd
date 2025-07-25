package ui

import (
	"fmt"
	"lledger-tui/pkg/journal"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type TransactionTableModel struct {
	table table.Model
}

func NewTransactionTableModel() TransactionTableModel {
	columns := []table.Column{
		{Title: "Date", Width: 10},
		{Title: "Type", Width: 8},
		{Title: "Account", Width: 30},
		{Title: "Category", Width: 15},
		{Title: "Amount", Width: 12},
	}

	t := table.New(table.WithColumns(columns))

	t.SetStyles(table.Styles{
		Header: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(true),
		Selected: lipgloss.NewStyle().
			Background(lipgloss.Color("235")).
			Foreground(lipgloss.Color("255")).
			Bold(false),
	})

	return TransactionTableModel{table: t}
}

func (m TransactionTableModel) Init() tea.Cmd {
	return nil
}

func (m TransactionTableModel) Update(msg tea.Msg) (TransactionTableModel, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m TransactionTableModel) View() string {
	return baseStyle.Render(m.table.View())
}

func (m *TransactionTableModel) SetHeight(height int) {
	m.table.SetHeight(height)
}

func (m *TransactionTableModel) SetTransactions(transactions []journal.Transaction) {
	rows := make([]table.Row, len(transactions))
	for i, tx := range transactions {
		rows[i] = table.Row{
			tx.Date.Format("2006/01/02"),
			string(tx.Type),
			tx.Account,
			tx.Category,
			fmt.Sprintf("%.2f", tx.Amount),
		}
	}
	m.table.SetRows(rows)
}
