package model

import (
	"fmt"
	"lledger/internal/data"
	"lledger/internal/ui"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	allTransactions      []data.Transaction
	filteredTransactions []data.Transaction

	transactionTable ui.TransactionTableModel
	datePicker       ui.DatePickerModel

	width  int
	height int
}

func NewModel() Model {
	return Model{
		transactionTable: ui.NewTransactionTableModel(),
		datePicker:       ui.NewDatePickerModel(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(loadTransactionsCmd(), m.datePicker.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	case []data.Transaction:
		m.allTransactions = msg
		startDate, endDate := m.datePicker.GetSelectedDateRange()
		m.filterTransactions(startDate, endDate)
	case ui.DateRangeChangedMsg:
		// This message comes from the date picker when the date range changes
		m.filterTransactions(msg.StartDate, msg.EndDate)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// Update sub-models
	m.transactionTable, cmd = m.transactionTable.Update(msg)
	cmds = append(cmds, cmd)

	m.datePicker, cmd = m.datePicker.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) filterTransactions(startDate, endDate time.Time) {
	m.filteredTransactions = []data.Transaction{}
	for _, tx := range m.allTransactions {
		if (tx.Date.After(startDate) || tx.Date.Equal(startDate)) && (tx.Date.Before(endDate) || tx.Date.Equal(endDate)) {
			m.filteredTransactions = append(m.filteredTransactions, tx)
		}
	}
	m.transactionTable.SetTransactions(m.filteredTransactions)
}

func (m Model) View() string {
	// Given the available vertical space, calculate the new height for the transaction table.
	// The date picker has a fixed height of 3, and we want to leave 2 rows for the header and footer.
	height := m.height - 5
	m.transactionTable.SetHeight(height)

	return lipgloss.JoinVertical(lipgloss.Left,
		m.datePicker.View(),
		m.transactionTable.View(),
	)
}

func loadTransactionsCmd() tea.Cmd {
	return func() tea.Msg {
		filePath := os.Getenv("LEDGER_FILE")
		if filePath == "" {
			return fmt.Errorf("LEDGER_FILE environment variable not set")
		}

		transactions, err := data.ParseJournal(filePath)
		if err != nil {
			return err
		}
		return transactions
	}
}
