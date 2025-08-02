package model

import (
	"cashd/internal/data"
	"cashd/internal/data/ledger"
	"cashd/internal/ui"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type dataLoadingSuccessMsg struct {
	transactions []*data.Transaction
}

type dataLoadingErrorMsg struct {
	err error
}

type Model struct {
	allTransactions  []*data.Transaction
	viewTransactions []*data.Transaction

	errMsg string

	transactionTable ui.TransactionTableModel
	datePicker       ui.DatePickerModel
	navBar           ui.NavBarModel
	summary          ui.SummaryModel
	searchInput      ui.SearchInputModel

	globalQuit     key.Binding
	quit           key.Binding
	activateSearch key.Binding
	clearSearch    key.Binding

	width  int
	height int
}

func NewModel() Model {
	return Model{
		transactionTable: ui.NewTransactionTableModel(),
		datePicker:       ui.NewDatePickerModel(),
		navBar:           ui.NewNavBarModel(),
		summary:          ui.NewSummaryModel(),
		searchInput:      ui.NewSearchInputModel(),
		quit:             key.NewBinding(key.WithKeys("q")),
		globalQuit:       key.NewBinding(key.WithKeys("ctrl+c")),
		activateSearch:   key.NewBinding(key.WithKeys("/")),
		clearSearch:      key.NewBinding(key.WithKeys("esc")),
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
		if key.Matches(msg, m.globalQuit) {
			return m, tea.Quit
		} else if m.searchInput.Focused() {
			m.searchInput, cmd = m.searchInput.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch {
			case key.Matches(msg, m.quit):
				return m, tea.Quit
			case key.Matches(msg, m.activateSearch):
				m.searchInput.Focus()
			case key.Matches(msg, m.clearSearch):
				cmds = append(cmds, m.searchInput.Clear())
			default:
				m.transactionTable, cmd = m.transactionTable.Update(msg)
				cmds = append(cmds, cmd)
				m.datePicker, cmd = m.datePicker.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	case dataLoadingSuccessMsg:
		m.allTransactions = msg.transactions
		m.filterTransactions()
	case dataLoadingErrorMsg:
		m.errMsg = msg.err.Error()
	case ui.DateRangeChangedMsg:
		// This message comes from the date picker when the date range changes
		m.filterTransactions()
	case ui.SearchMsg:
		m.filterTransactions()
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) filterTransactions() {
	startDate, endDate := m.datePicker.SelectedDateRange()
	// m.allTransactions are ordered by date, use binary search to find start, end index
	startIndex := sort.Search(len(m.allTransactions), func(i int) bool {
		d := m.allTransactions[i].Date
		return d.Equal(startDate) || d.After(startDate)
	})
	endIndex := sort.Search(len(m.allTransactions), func(i int) bool {
		d := m.allTransactions[i].Date
		return d.Equal(endDate) || d.After(endDate)
	})

	if startIndex > endIndex {
		panic(fmt.Sprintf("Invalid date range: %s - %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")))
	}

	m.viewTransactions = []*data.Transaction{}
	for i := startIndex; i < endIndex; i++ {
		searchQuery := strings.ToLower(m.searchInput.Value())
		keywords := strings.Fields(searchQuery)
		if t := m.allTransactions[i]; t.Matches(keywords) {
			m.viewTransactions = append(m.viewTransactions, t)
		}
	}

	m.transactionTable.SetTransactions(m.viewTransactions)
	m.summary.SetTransactions(m.viewTransactions)
}

func loadTransactionsCmd() tea.Cmd {
	// TODO: switch data source based on command line flags
	datasource := ledger.LedgerDataSource{}
	return func() tea.Msg {
		transactions, err := datasource.LoadTransactions()
		if err != nil {
			return dataLoadingErrorMsg{err}
		} else {
			return dataLoadingSuccessMsg{transactions}
		}
	}
}
