package model

import (
	"cashd/internal/ui"
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	datePickerHeight  = 5
	searchInputHeight = 3
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
	top := m.datePicker.View()
	return lipgloss.JoinVertical(lipgloss.Left,
		top,
		body,
	)
}

func (m *Model) updateLayout() {
	m.datePicker.SetWidth(ui.PreferredTableWidth)
	m.searchInput.SetWidth(ui.PreferredTableWidth - 4)
	m.transactionTable.SetDimensions(ui.PreferredTableWidth, m.height-datePickerHeight-searchInputHeight)
	m.summary.SetDimensions(m.width-ui.PreferredTableWidth-4, m.height-datePickerHeight)
}
