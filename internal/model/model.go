package model

import (
	"lledger/internal/data"
	"lledger/internal/ui"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	ENV_LEDGER_FILE = "LEDGER_FILE"
)

type dataLoadingSuccessMsg struct {
	transactions []data.Transaction
}

type dataLoadingErrorMsg struct {
	err error
}

type Model struct {
	allTransactions      []data.Transaction
	filteredTransactions []*data.Transaction

	errMsg string

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
	case dataLoadingSuccessMsg:
		m.allTransactions = msg.transactions
		startDate, endDate := m.datePicker.GetSelectedDateRange()
		m.filterTransactions(startDate, endDate)
	case dataLoadingErrorMsg:
		m.errMsg = msg.err.Error()
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
	m.filteredTransactions = []*data.Transaction{}
	for _, tx := range m.allTransactions {
		if (tx.Date.After(startDate) || tx.Date.Equal(startDate)) && tx.Date.Before(endDate) {
			m.filteredTransactions = append(m.filteredTransactions, &tx)
		}
	}
	m.transactionTable.SetTransactions(m.filteredTransactions)
}

func loadTransactionsCmd() tea.Cmd {
	return func() tea.Msg {
		transactions, err := data.LoadTransactions()
		if err != nil {
			return dataLoadingErrorMsg{err}
		} else {
			return dataLoadingSuccessMsg{transactions}
		}
	}
}
