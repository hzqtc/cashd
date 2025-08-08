package model

import (
	"cashd/internal/ui"
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	datePickerHeight  = 5
	searchInputHeight = 3
	insightsHeight    = 14
	vSpacing          = 2
)

func (m Model) View() string {
	if m.errMsg != "" {
		return fmt.Sprintf("An error occurred: %s\nPress 'ctrl + c' to quit.", m.errMsg)
	}

	var top, body string
	top = lipgloss.JoinHorizontal(lipgloss.Top,
		m.navBar.View(),
		m.datePicker.View(),
	)

	switch m.navBar.ViewMode() {
	case ui.TransactionView:
		body = lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.JoinVertical(lipgloss.Left,
				m.searchInput.View(),
				m.transactionTable.View(),
			),
			m.summary.View(),
		)
	case ui.AccountView:
		body = lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Top,
				m.accountTable.View(),
				m.accountInsights.View(),
			),
			m.accountChart.View(),
		)
	case ui.CategoryView:
		body = lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Top,
				m.categoryTable.View(),
				m.categoryInsights.View(),
			),
			m.categoryChart.View(),
		)
	}

	views := []string{top, body}
	if m.help.Visible() {
		views = append(views, m.help.View())
	}

	return lipgloss.JoinVertical(lipgloss.Left, views...)
}

func (m *Model) updateLayout() {
	m.navBar.SetWidth(ui.NavBarWidth)
	m.datePicker.SetWidth(m.width - ui.NavBarWidth - 4)

	m.help.SetWidth(m.width - 2)
	var bottomHeight int
	if m.help.Visible() {
		bottomHeight = lipgloss.Height(m.help.View())
	} else {
		bottomHeight = 0
	}
	bodyHeight := m.height - datePickerHeight - bottomHeight - vSpacing
	// Transaction view components
	m.searchInput.SetWidth(ui.TxnTableWidth - 4)
	m.transactionTable.SetDimensions(ui.TxnTableWidth, bodyHeight-searchInputHeight)
	m.summary.SetDimensions(max(30, m.width-ui.TxnTableWidth-4), bodyHeight)
	// Account view components
	m.accountTable.SetDimensions(ui.AccountTableWidth, insightsHeight)
	m.accountInsights.SetDimension(max(30, m.width-ui.AccountTableWidth-4), insightsHeight)
	m.accountChart.SetDimension(m.width-4, bodyHeight-m.accountInsights.Height()-2)
	// Category view components
	m.categoryTable.SetDimensions(ui.CategoryTableWidth, insightsHeight)
	m.categoryInsights.SetDimension(max(30, m.width-ui.CategoryTableWidth-4), insightsHeight)
	m.categoryChart.SetDimension(m.width-4, bodyHeight-m.categoryInsights.Height()-2)
}
