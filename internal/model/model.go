package model

import (
	"cashd/internal/data"
	"cashd/internal/data/csv"
	"cashd/internal/data/ledger"
	"cashd/internal/date"
	"cashd/internal/ui"
	"fmt"
	"sort"
	"strings"
	"time"

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

	loadingScreen    ui.LoadingScreenModel
	datePicker       ui.DatePickerModel
	navBar           ui.NavBarModel
	searchInput      ui.SearchInputModel
	transactionTable ui.SortableTableModel
	summary          ui.SummaryModel
	accountTable     ui.SortableTableModel
	accountInsights  ui.InsightsModel
	accountChart     ui.TimeSeriesChartModel
	categoryTable    ui.SortableTableModel
	categoryInsights ui.InsightsModel
	categoryChart    ui.TimeSeriesChartModel
	help             ui.HelpModel

	globalQuit     key.Binding
	activateSearch key.Binding
	clearSearch    key.Binding
	toggleHelp     key.Binding

	width  int
	height int
}

func NewModel() Model {
	return Model{
		loadingScreen:    ui.NewLoadingScreenModel(),
		transactionTable: ui.NewTransactionTableModel(),
		datePicker:       ui.NewDatePickerModel(),
		navBar:           ui.NewNavBarModel(),
		summary:          ui.NewSummaryModel(),
		searchInput:      ui.NewSearchInputModel(),
		accountTable:     ui.NewAccountTableModel(),
		accountInsights:  ui.NewInsightsModel(),
		accountChart:     ui.NewTimeSeriesChartModel(),
		categoryTable:    ui.NewCategoryTableModel(),
		categoryInsights: ui.NewInsightsModel(),
		categoryChart:    ui.NewTimeSeriesChartModel(),
		help:             ui.NewHelpModel(),

		globalQuit:     key.NewBinding(key.WithKeys("ctrl+c")),
		activateSearch: key.NewBinding(key.WithKeys("/")),
		clearSearch:    key.NewBinding(key.WithKeys("esc")),
		toggleHelp:     key.NewBinding(key.WithKeys("?")),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.loadingScreen.Init(), loadTransactions())
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
			return m, m.processSearchInputKeys(msg)
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
		cmds = append(cmds, m.loadingScreen.Stop())
		m.allTransactions = msg.transactions
		m.updateDatePickerLimits()
		cmds = append(cmds, m.filterTransactions())
		m.onSelectedAccountChanged()
		m.onSelectedCategoryChanged()

	case dataLoadingErrorMsg:
		cmds = append(cmds, m.loadingScreen.Stop())
		m.errMsg = msg.err.Error()

	case ui.DateRangeChangedMsg:
		cmds = append(cmds, m.filterTransactions())
		m.updateAccountInsights()
		m.updateCategoryInsights()

	case ui.DateIncrementChangedMsg:
		m.onSelectedAccountChanged()
		m.onSelectedCategoryChanged()

	case ui.TableSelectionChangedMsg:
		switch msg.TableName {
		case ui.AccountTableName:
			m.onSelectedAccountChanged()
		case ui.CategoryTableName:
			m.onSelectedCategoryChanged()
		}

	case ui.SearchMsg:
		m.searchTransactions()

	case ui.NavigationMsg:
		// Handled in view.go; nothing to do here

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()

	default:
		if m.loadingScreen.Handles(msg) {
			m.loadingScreen, cmd = m.loadingScreen.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) processSearchInputKeys(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	m.updateLayout()
	return cmd
}

func (m *Model) processTransactionViewKeys(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd

	if key.Matches(msg, m.clearSearch) {
		m.searchInput.Blur()
		return m.searchInput.Clear()
	} else {
		switch {
		case key.Matches(msg, m.activateSearch):
			m.searchInput.Focus()
		case key.Matches(msg, m.toggleHelp):
			m.help.ToggleVisibility()
			m.updateLayout()
		default:
			m.transactionTable, cmd = m.transactionTable.Update(msg)
			return cmd
		}
	}
	return nil
}

func (m *Model) processAccountViewKeys(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.toggleHelp):
		m.help.ToggleVisibility()
		m.updateLayout()
	default:
		var cmd tea.Cmd
		m.accountTable, cmd = m.accountTable.Update(msg)
		return cmd
	}
	return nil
}

func (m *Model) processCategoryViewKeys(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.toggleHelp):
		m.help.ToggleVisibility()
		m.updateLayout()
	default:
		var cmd tea.Cmd
		m.categoryTable, cmd = m.categoryTable.Update(msg)
		return cmd
	}
	return nil
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
		panic(fmt.Sprintf("Invalid date range: %s - %s", startDate.Format(time.DateOnly), endDate.Format(time.DateOnly)))
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
	searchQuery := strings.TrimSpace(m.searchInput.Value())

	if searchQuery == "" {
		matchingTransactions = m.viewTransactions
	} else {
		matchingTransactions = []*data.Transaction{}
		subQueries := strings.Split(searchQuery, " OR ")
		for _, t := range m.viewTransactions {
			for _, query := range subQueries {
				keywords := strings.Fields(strings.ToLower(query))
				if len(keywords) > 0 && t.Matches(keywords) {
					// If any of the sub-queries match, the transaction  a match
					matchingTransactions = append(matchingTransactions, t)
					break
				}
			}
		}
	}

	// Search query only applies to transaction view (transaction table & summary)
	m.transactionTable.SetTransactions(matchingTransactions)
	m.summary.SetTransactions(matchingTransactions)
}

func (m *Model) onSelectedAccountChanged() {
	if m.accountTable.Selected() == "" {
		return
	}

	entries := aggregateByAccount(m.allTransactions, m.datePicker.Inc(), m.accountTable.Selected())
	m.accountChart.SetEntries(
		getTimeSeriesChartName(m.datePicker.Inc(), m.accountTable.Selected()),
		entries,
		m.datePicker.Inc(),
	)

	m.updateAccountInsights()
}

func (m *Model) onSelectedCategoryChanged() {
	if m.categoryTable.Selected() == "" {
		return
	}

	entries := aggregateByCategory(m.allTransactions, m.datePicker.Inc(), m.categoryTable.Selected())
	m.categoryChart.SetEntries(
		getTimeSeriesChartName(m.datePicker.Inc(), m.categoryTable.Selected()),
		entries,
		m.datePicker.Inc(),
	)

	m.updateCategoryInsights()
}

func getTimeSeriesChartName(inc date.Increment, name string) string {
	incStr := string(inc)
	if inc == date.AllTime {
		incStr = string(date.Annually)
	}
	return fmt.Sprintf("%s time series: %s", incStr, name)
}

func (m *Model) updateAccountInsights() {
	m.accountInsights.SetTransactionsWithAccount(m.viewTransactions, m.accountTable.Selected())
	m.accountInsights.SetName(fmt.Sprintf("%s insights: %s", m.accountTable.Selected(), m.datePicker.ViewDateRange()))

	m.updateLayout()
}

func (m *Model) updateCategoryInsights() {
	m.categoryInsights.SetTransactionsWithCategory(m.viewTransactions, m.categoryTable.Selected())
	m.categoryInsights.SetName(fmt.Sprintf("%s insights: %s", m.categoryTable.Selected(), m.datePicker.ViewDateRange()))

	m.updateLayout()
}

func loadTransactions() tea.Cmd {
	datasources := []data.DataSource{ledger.LedgerDataSource{}, csv.CsvDataSource{}}
	for _, ds := range datasources {
		if ds.Preferred() {
			return loadTransactionsFromDataSource(ds)
		}
	}
	for _, ds := range datasources {
		if ds.Enabled() {
			return loadTransactionsFromDataSource(ds)
		}
	}
	return func() tea.Msg {
		return dataLoadingErrorMsg{fmt.Errorf("No available data source")}
	}
}

func loadTransactionsFromDataSource(ds data.DataSource) tea.Cmd {
	return func() tea.Msg {
		transactions, err := ds.LoadTransactions()
		if err != nil {
			return dataLoadingErrorMsg{err}
		} else {
			return dataLoadingSuccessMsg{transactions}
		}
	}
}
