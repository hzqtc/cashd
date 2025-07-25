package ui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	highlightColor          = lipgloss.Color("#FFD580")
	highlightForegroudColor = lipgloss.Color("#2E2E2E")
	borderColor             = lipgloss.Color("240")
	focusedBorderColor      = highlightColor
	incomeColor             = lipgloss.Color("#22C55E")
	expenseColor            = lipgloss.Color("#FBBF24")

	roundedBorder = lipgloss.RoundedBorder()

	baseStyle = lipgloss.NewStyle().
			BorderStyle(roundedBorder).
			BorderForeground(borderColor)

	headerStyle = lipgloss.NewStyle().
			Foreground(highlightColor).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Padding(1 /* top */, 2 /* horizontal */, 0 /* bottom */)

	keyStyle = lipgloss.NewStyle().
			Foreground(highlightColor)

	searchStyle = baseStyle.
			Margin(1 /* top */, 0 /* horizontal */, 0 /* bottom */)

	tableStyle = baseStyle

	datepickerStyle = lipgloss.NewStyle().
			BorderStyle(roundedBorder).
			BorderForeground(borderColor).
			Padding(0, 1)
)

func getTableStyle() table.Styles {
	return table.Styles{
		Header: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(borderColor).
			BorderBottom(true).
			Bold(true),
		Selected: lipgloss.NewStyle().
			Background(highlightColor).
			Foreground(highlightForegroudColor).
			Bold(false),
	}
}
