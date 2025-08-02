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

	// Top components are fixed
	top := lipgloss.JoinHorizontal(lipgloss.Top,
		m.datePicker.View(),
		m.navBar.View(),
	)

	var body string
	switch m.navBar.ViewMode() {
	case ui.TransactionView:
		table := lipgloss.JoinVertical(lipgloss.Left,
			m.searchInput.View(),
			m.transactionTable.View(),
		)
		body = lipgloss.JoinHorizontal(lipgloss.Top,
			table,
			m.summary.View(),
		)
	case ui.AccountView:
		body = lipgloss.JoinHorizontal(lipgloss.Top,
			m.accountTable.View(),
		)
	case ui.CategoryView:
		body = lipgloss.JoinHorizontal(lipgloss.Top,
			m.categoryTable.View(),
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
	// Transaction view components
	m.searchInput.SetWidth(ui.TxnTableWidth - 4)
	m.transactionTable.SetDimensions(ui.TxnTableWidth, m.height-datePickerHeight-searchInputHeight-vSpacing)
	m.summary.SetDimensions(summaryPanelWidth, m.height-datePickerHeight-vSpacing)
	// Account view components
	m.accountTable.SetDimensions(ui.AccountTableWidth, m.height-datePickerHeight-vSpacing)
	// Category view components
	m.categoryTable.SetDimensions(ui.CategoryTableWidth, m.height-datePickerHeight-vSpacing)
}
