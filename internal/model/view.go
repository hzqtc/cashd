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

	table := lipgloss.JoinVertical(lipgloss.Left,
		m.searchInput.View(),
		m.transactionTable.View(),
	)
	body := lipgloss.JoinHorizontal(lipgloss.Top,
		table,
		m.summary.View(),
	)
	top := lipgloss.JoinHorizontal(lipgloss.Top,
		m.datePicker.View(),
		m.navBar.View(),
	)
	return lipgloss.JoinVertical(lipgloss.Left,
		top,
		body,
	)
}

func (m *Model) updateLayout() {
	// TODO: Hide right side panel if not enough width
	m.datePicker.SetWidth(ui.TxnTableWidth)
	m.searchInput.SetWidth(ui.TxnTableWidth - 4)
	m.transactionTable.SetDimensions(ui.TxnTableWidth, m.height-datePickerHeight-searchInputHeight-vSpacing)
	m.summary.SetDimensions(m.width-ui.TxnTableWidth-4, m.height-datePickerHeight-vSpacing)
	m.navBar.SetWidth(m.width - ui.TxnTableWidth - 4)
}
