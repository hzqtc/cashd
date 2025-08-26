package ui

import (
	"cashd/internal/data"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type column interface {
	index() int
	rightAligned() bool
	isSortable() bool
	width() int
	nextColumn() column
	prevColumn() column
	getColumnData(any) any
	String() string
}

type sortDirection int

const (
	sortAsc sortDirection = iota
	sortDesc
)

// tableDataSorter is a function that returns sorted table data
type tableDataSorter func(sortCol column, sortDir sortDirection) []any

// tableDataProvider is a function that takes transactions as input, and return a TableDataSorter
type tableDataProvider func(transactions []*data.Transaction) tableDataSorter

// Return a unique string as the row's id
type rowIdentifier func(table.Row) string

type tableConfig struct {
	columns           []column
	dataProvider      tableDataProvider
	rowId             rowIdentifier
	defaultSortColumn column
	defaultSortDir    sortDirection
}

type TableSelectionChangedMsg struct {
	TableName string
	Selected  string
}

type SortableTableModel struct {
	name          string
	columns       []column
	dataProvider  tableDataProvider
	dataSorter    tableDataSorter
	rowId         rowIdentifier
	sortColumn    column
	sortDirection sortDirection
	table         table.Model

	sortNext    key.Binding
	sortPrev    key.Binding
	reverseSort key.Binding
}

func newSortableTableModel(name string, config tableConfig) SortableTableModel {
	m := SortableTableModel{
		name:          name,
		columns:       config.columns,
		dataProvider:  config.dataProvider,
		rowId:         config.rowId,
		sortColumn:    config.defaultSortColumn,
		sortDirection: config.defaultSortDir,

		sortNext:    key.NewBinding(key.WithKeys("s")),
		sortPrev:    key.NewBinding(key.WithKeys("S")),
		reverseSort: key.NewBinding(key.WithKeys("r")),
	}
	m.table = table.New(
		table.WithColumns(m.getTableColumns()),
		table.WithFocused(true),
		table.WithStyles(getTableStyle()),
	)
	return m
}

func (m *SortableTableModel) getTableColumns() []table.Column {
	tableCols := []table.Column{}
	for _, col := range m.columns {
		title := col.String()
		width := col.width()
		if col == m.sortColumn {
			if m.sortDirection == sortAsc {
				title = "↑ " + title
			} else {
				title = "↓ " + title
			}
		}
		if col.rightAligned() {
			title = fmt.Sprintf("%*s", width, title)
		}
		tableCols = append(tableCols, table.Column{Title: title, Width: width})
	}
	return tableCols
}

func (m SortableTableModel) Update(msg tea.Msg) (SortableTableModel, tea.Cmd) {
	selected := m.Selected()
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.sortNext):
			m.sortNextColumn()
		case key.Matches(msg, m.sortPrev):
			m.sortPrevColumn()
		case key.Matches(msg, m.reverseSort):
			m.reverseSortDir()
		}
	}
	m.table, _ = m.table.Update(msg)
	if m.Selected() != selected {
		cmd = m.sendSelectionChangedMsg()
	}
	return m, cmd
}

func (m *SortableTableModel) Selected() string {
	if row := m.table.SelectedRow(); row != nil && m.rowId != nil {
		return m.rowId(row)
	} else {
		return ""
	}
}

func (m *SortableTableModel) sendSelectionChangedMsg() tea.Cmd {
	return func() tea.Msg {
		return TableSelectionChangedMsg{
			TableName: m.name,
			Selected:  m.Selected(),
		}
	}
}

func (m *SortableTableModel) sortNextColumn() {
	newCol := m.sortColumn.nextColumn()
	for !newCol.isSortable() {
		newCol = newCol.nextColumn()
	}
	m.sortColumn = newCol
	m.updateSorting()
}

func (m *SortableTableModel) sortPrevColumn() {
	newCol := m.sortColumn.prevColumn()
	for !newCol.isSortable() {
		newCol = newCol.prevColumn()
	}
	m.sortColumn = newCol
	m.updateSorting()
}

func (m *SortableTableModel) reverseSortDir() {
	if m.sortDirection == sortAsc {
		m.sortDirection = sortDesc
	} else {
		m.sortDirection = sortAsc
	}
	m.updateSorting()
}

func (m *SortableTableModel) updateSorting() {
	m.table.SetColumns(m.getTableColumns())
	m.updateRows()
}

func (m SortableTableModel) View() string {
	return baseStyle.Render(m.table.View())
}

func (m *SortableTableModel) SetDimensions(width, height int) {
	m.table.SetWidth(width)
	m.table.SetHeight(height)
}

func (m *SortableTableModel) SetTransactions(transactions []*data.Transaction) tea.Cmd {
	selected := m.Selected()
	m.dataSorter = m.dataProvider(transactions)
	m.updateRows()
	if m.Selected() != selected {
		return m.sendSelectionChangedMsg()
	} else {
		return nil
	}
}

func (m *SortableTableModel) updateRows() {
	if m.dataSorter == nil {
		return
	}
	m.table.SetRows(getTableRows(m.columns, m.dataSorter(m.sortColumn, m.sortDirection)))
}

func getTableRows(cols []column, tableData []any) []table.Row {
	rows := make([]table.Row, len(tableData))
	for i, cat := range tableData {
		row := []string{}
		for _, col := range cols {
			var formattedColData string
			switch colData := col.getColumnData(cat).(type) {
			case string:
				formattedColData = colData
			case int:
				formattedColData = fmt.Sprintf("%d", colData)
			case float64:
				formattedColData = data.FormatMoney(colData)
			case time.Time:
				formattedColData = colData.Format(time.DateOnly)
			default:
				panic(fmt.Sprintf("unexpected table data type: %v", colData))
			}
			if col.rightAligned() {
				formattedColData = fmt.Sprintf("%*s", col.width(), formattedColData)
			}
			row = append(row, formattedColData)
		}
		rows[i] = table.Row(row)
	}
	return rows
}

func compareAny(a, b any, sortDir sortDirection) bool {
	var inOrder bool
	switch a.(type) {
	case string:
		inOrder = a.(string) < b.(string)
	case int:
		inOrder = a.(int) < b.(int)
	case float64:
		inOrder = a.(float64) < b.(float64)
	case time.Time:
		inOrder = a.(time.Time).Before(b.(time.Time))
	default:
		panic(fmt.Sprintf("unexpected table data type: %v", a))
	}

	if sortDir == sortDesc {
		return !inOrder
	} else {
		return inOrder
	}
}
