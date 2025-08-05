package ui

import (
	"cashd/internal/data"
	"fmt"

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
	getColumnData(any) string
	String() string
}

type sortDirection int

const (
	sortAsc sortDirection = iota
	sortDesc
)

// tableDataSorter is a function that returns sorted table data
type tableDataSorter func(sortCol column, sortDir sortDirection) []table.Row

// tableDataProvider is a function that takes transactions as input, and return a TableDataSorter
type tableDataProvider func(transactions []*data.Transaction) tableDataSorter

// Return a unique string as the row's id
type rowIdentifier func(table.Row) string

type TableConfig struct {
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
}

func NewSortableTableModel(name string, config TableConfig) SortableTableModel {
	m := SortableTableModel{
		name:          name,
		columns:       config.columns,
		dataProvider:  config.dataProvider,
		rowId:         config.rowId,
		sortColumn:    config.defaultSortColumn,
		sortDirection: config.defaultSortDir,
	}
	t := table.New(
		table.WithColumns(m.getTableColumns()),
		table.WithFocused(true),
		table.WithStyles(getTableStyle()),
	)
	m.table = t
	return m
}

func (m SortableTableModel) getTableColumns() []table.Column {
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
	// TODO: define key bindings using bubbletea.key
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			m.sortNextColumn()
		case "S":
			m.sortPrevColumn()
		case "r":
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
	if row := m.table.SelectedRow(); row != nil {
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

func (m *SortableTableModel) SetTransactions(transactions []*data.Transaction) {
	m.dataSorter = m.dataProvider(transactions)
	m.updateRows()
}

func (m *SortableTableModel) updateRows() {
	// TODO: trigger TableSelectionChangedMsg
	if m.dataSorter == nil {
		return
	}
	m.table.SetRows(m.dataSorter(m.sortColumn, m.sortDirection))
}
