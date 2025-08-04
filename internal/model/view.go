package model

import (
	"cashd/internal/ui"
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	datePickerHeight  = 5
	searchInputHeight = 3
	vSpacing          = 2
)

func (m Model) View() string {
	if m.errMsg != "" {
		return fmt.Sprintf("An error occurred: %s\nPress 'q' to quit.", m.errMsg)
	}

	var top, body string
	top = lipgloss.JoinHorizontal(lipgloss.Top,
		m.datePicker.View(),
		m.navBar.View(),
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
		body = lipgloss.JoinHorizontal(lipgloss.Top,
			m.accountTable.View(),
			// TODO: add an account highlight view panel
			m.accountChart.View(),
		)
	case ui.CategoryView:
		body = lipgloss.JoinHorizontal(lipgloss.Top,
			m.categoryTable.View(),
			// TODO: add a category highlight view panel
			m.categoryChart.View(),
		)
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		top,
		body,
	)
}

func (m *Model) updateLayout() {
	summaryPanelWidth := max(30, m.width-ui.TxnTableWidth-4)
	if summaryPanelWidth >= ui.NavBarWidth {
		// Try to align top components to transaction view components since that's the default view
		m.datePicker.SetWidth(ui.TxnTableWidth)
		m.navBar.SetWidth(summaryPanelWidth)
	} else {
		// However, prefer to keep navbar a singleline
		m.navBar.SetWidth(ui.NavBarWidth)
		m.datePicker.SetWidth(m.width - ui.NavBarWidth - 4)
	}

	bodyHeight := m.height - datePickerHeight - vSpacing
	// Transaction view components
	m.searchInput.SetWidth(ui.TxnTableWidth - 4)
	m.transactionTable.SetDimensions(ui.TxnTableWidth, bodyHeight-searchInputHeight)
	m.summary.SetDimensions(summaryPanelWidth, bodyHeight)
	// Account view components
	m.accountTable.SetDimensions(ui.AccountTableWidth, bodyHeight)
	m.accountChart.SetDimension(m.width-ui.AccountTableWidth-4, bodyHeight)
	// Category view components
	m.categoryTable.SetDimensions(ui.CategoryTableWidth, bodyHeight)
	m.categoryChart.SetDimension(m.width-ui.CategoryTableWidth-4, bodyHeight)
}
