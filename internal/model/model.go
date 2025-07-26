package model

import (
	"fmt"
	"lledger/internal/data"
	"lledger/internal/ui"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	ENV_LEDGER_FILE = "LEDGER_FILE"
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
	return tea.Batch(loadTransactionsCmd())
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
		m.filterTransactions(msg.Start, msg.End)
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
		if (tx.Date.After(startDate) || tx.Date.Equal(startDate)) && tx.Date.Before(endDate) {
			m.filteredTransactions = append(m.filteredTransactions, tx)
		}
	}
	m.transactionTable.SetTransactions(m.filteredTransactions)
}

func loadTransactionsCmd() tea.Cmd {
	return func() tea.Msg {
		filePath := os.Getenv(ENV_LEDGER_FILE)
		if filePath == "" {
			return fmt.Errorf("%s environment variable not set", ENV_LEDGER_FILE)
		}

		transactions, err := data.ParseJournal(filePath)
		if err != nil {
			return err
		}
		return transactions
	}
}
