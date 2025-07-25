package model

import "github.com/charmbracelet/lipgloss"

func (m Model) View() string {
	// Given the available vertical space, calculate the new height for the transaction table.
	// The date picker has a fixed height of 3, and we want to leave 2 rows for the header and footer.
	h := m.height - 5
	w := m.width - 5
	m.transactionTable.SetDimensions(h, w)

	return lipgloss.JoinVertical(lipgloss.Left,
		m.datePicker.View(),
		m.transactionTable.View(),
	)
}
