package model

import (
	"fmt"
	"lledger/internal/ui"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.errMsg != "" {
		return fmt.Sprintf("An error occurred: %s\nPress 'q' to quit.", m.errMsg)
	}

	body := lipgloss.JoinHorizontal(lipgloss.Top,
		m.transactionTable.View(),
		m.summary.View(),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		m.datePicker.View(),
		body,
	)
}

func (m *Model) updateLayout() {
	m.transactionTable.SetDimensions(ui.PreferredTableWidth, m.height-5)
	m.summary.SetDimensions(m.width-ui.PreferredTableWidth-4, m.height-5)
}
