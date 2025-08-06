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
		return fmt.Sprintf("An error occurred: %s\nPress 'q' to quit.", m.errMsg)
	}

	var top, body string
	top = lipgloss.JoinHorizontal(lipgloss.Top,
		m.navBar.View(),
		m.datePicker.View(),
	)
	// TODO: add a help panel at the bottom
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

	return lipgloss.JoinVertical(lipgloss.Left,
		top,
		body,
	)
}

func (m *Model) updateLayout() {
	m.navBar.SetWidth(ui.NavBarWidth)
	m.datePicker.SetWidth(m.width - ui.NavBarWidth - 4)

	bodyHeight := m.height - datePickerHeight - vSpacing
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
