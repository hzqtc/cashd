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

const (
	tableTransaction = "Transaction"
	tableAccount     = "Account"
	tableCategory    = "Category"
)

type Model struct {
	allTransactions  []*data.Transaction
	viewTransactions []*data.Transaction

	errMsg string

	datePicker       ui.DatePickerModel
	navBar           ui.NavBarModel
	searchInput      ui.SearchInputModel
	transactionTable ui.SortableTableModel
	summary          ui.SummaryModel
	accountTable     ui.SortableTableModel
	accountChart     ui.TimeSeriesChartModel
	categoryTable    ui.SortableTableModel
	categoryChart    ui.TimeSeriesChartModel

	globalQuit     key.Binding
	quit           key.Binding
	activateSearch key.Binding
	clearSearch    key.Binding

	width  int
	height int
}

func NewModel() Model {
	return Model{
		transactionTable: ui.NewSortableTableModel(tableTransaction, ui.TransactionTableConfig),
		datePicker:       ui.NewDatePickerModel(),
		navBar:           ui.NewNavBarModel(),
		summary:          ui.NewSummaryModel(),
		searchInput:      ui.NewSearchInputModel(),
		accountTable:     ui.NewSortableTableModel(tableAccount, ui.AccountTableConfig),
		accountChart:     ui.NewTimeSeriesChartModel(),
		categoryTable:    ui.NewSortableTableModel(tableCategory, ui.CategoryTableConfig),
		categoryChart:    ui.NewTimeSeriesChartModel(),

		globalQuit:     key.NewBinding(key.WithKeys("ctrl+c")),
		activateSearch: key.NewBinding(key.WithKeys("/")),
		clearSearch:    key.NewBinding(key.WithKeys("esc")),
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
		}
		// Send key to the active view
		switch m.navBar.ViewMode() {
		case ui.TransactionView:
			cmds = append(cmds, m.processTransactionViewKeys(msg))
		case ui.AccountView:
			cmds = append(cmds, m.processAccountViewKeys(msg))
		case ui.CategoryView:
			cmds = append(cmds, m.processCategoryViewKeys(msg))
		}
		// Global components always process key events
		m.datePicker, cmd = m.datePicker.Update(msg)
		cmds = append(cmds, cmd)
		m.navBar, cmd = m.navBar.Update(msg)
		cmds = append(cmds, cmd)
	case dataLoadingSuccessMsg:
		m.allTransactions = msg.transactions
		m.updateDatePickerLimits()
		cmds = append(cmds, m.filterTransactions())
		m.updateAccountTimeSeriesCharts()
		m.updateCategoryTimeSeriesCharts()
	case dataLoadingErrorMsg:
		m.errMsg = msg.err.Error()
	case ui.DateRangeChangedMsg:
		cmds = append(cmds, m.filterTransactions())
	case ui.DateIncrementChangedMsg:
		m.updateAccountTimeSeriesCharts()
		m.updateCategoryTimeSeriesCharts()
	case ui.TableSelectionChangedMsg:
		switch msg.TableName {
		case tableAccount:
			m.updateAccountTimeSeriesCharts()
		case tableCategory:
			m.updateCategoryTimeSeriesCharts()
		}
	case ui.SearchMsg:
		m.searchTransactions()
	case ui.NavigationMsg:
		// No logic change
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) processTransactionViewKeys(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd

	if key.Matches(msg, m.clearSearch) {
		m.searchInput.Blur()
		return m.searchInput.Clear()
	} else if m.searchInput.Focused() {
		m.searchInput, cmd = m.searchInput.Update(msg)
		return cmd
	} else {
		switch {
		case key.Matches(msg, m.activateSearch):
			m.searchInput.Focus()
		default:
			m.transactionTable, cmd = m.transactionTable.Update(msg)
			return cmd
		}
	}
	return nil
}

func (m *Model) processAccountViewKeys(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	m.accountTable, cmd = m.accountTable.Update(msg)
	return cmd
}

func (m *Model) processCategoryViewKeys(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	m.categoryTable, cmd = m.categoryTable.Update(msg)
	return cmd
}

func (m *Model) updateDatePickerLimits() {
	if txnCount := len(m.allTransactions); txnCount == 0 {
		return
	} else {
		// m.allTransactions are sorted by date in asending order
		// So simply use the first transaction's date as the min and the last transaction's date as the max
		m.datePicker.SetLimits(
			m.allTransactions[0].Date,
			m.allTransactions[txnCount-1].Date,
		)
	}
}

func (m *Model) filterTransactions() tea.Cmd {
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

	// Get transactions within date range
	m.viewTransactions = []*data.Transaction{}
	for i := startIndex; i < endIndex; i++ {
		t := m.allTransactions[i]
		m.viewTransactions = append(m.viewTransactions, t)
	}

	var cmds []tea.Cmd
	cmds = append(cmds, m.accountTable.SetTransactions(m.viewTransactions))
	cmds = append(cmds, m.categoryTable.SetTransactions(m.viewTransactions))
	// Trigger a search msg to handle existing search queries
	cmds = append(cmds, func() tea.Msg { return ui.SearchMsg{} })
	return tea.Batch(cmds...)
}

func (m *Model) searchTransactions() {
	var matchingTransactions []*data.Transaction
	searchQuery := strings.ToLower(m.searchInput.Value())
	keywords := strings.Fields(searchQuery)
	if len(keywords) == 0 {
		matchingTransactions = m.viewTransactions
	} else {
		matchingTransactions = []*data.Transaction{}
		for _, t := range m.viewTransactions {
			if t.Matches(keywords) {
				matchingTransactions = append(matchingTransactions, t)
			}
		}
	}
	// Search query only applies to transaction view (transaction table & summary)
	m.transactionTable.SetTransactions(matchingTransactions)
	m.summary.SetTransactions(matchingTransactions)
}

func (m *Model) updateAccountTimeSeriesCharts() {
	if m.accountTable.Selected() == "" {
		return
	}
	entries := aggregateByAccount(m.allTransactions, m.datePicker.Inc(), m.accountTable.Selected())
	m.accountChart.SetEntries(
		fmt.Sprintf("%s: %s", m.accountTable.Selected(), m.getTimeSeriesRange(entries)),
		entries,
		m.datePicker.Inc(),
	)
	// TODO: scroll chart to the selected date in datepicker
}

func (m *Model) updateCategoryTimeSeriesCharts() {
	if m.categoryTable.Selected() == "" {
		return
	}
	entries := aggregateByCategory(m.allTransactions, m.datePicker.Inc(), m.categoryTable.Selected())
	m.categoryChart.SetEntries(
		fmt.Sprintf("%s: %s", m.categoryTable.Selected(), m.getTimeSeriesRange(entries)),
		entries,
		m.datePicker.Inc(),
	)
	// TODO: scroll chart to the selected date in datepicker
}

func (m *Model) getTimeSeriesRange(entries []*ui.TsChartEntry) string {
	if len(entries) == 0 {
		return ""
	}
	firstDate := entries[0].Date
	lastDate := entries[len(entries)-1].Date
	return fmt.Sprintf(
		"%s - %s",
		m.datePicker.Inc().FormatDateShorter(firstDate),
		m.datePicker.Inc().FormatDateShorter(lastDate),
	)
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
