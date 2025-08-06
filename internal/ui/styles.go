package ui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

const (
	vPadding = 0
	hPadding = 1
)

const (
	symbolColWidth      = 2
	dateColWidth        = 12
	typeColWidth        = 8
	accountTypeColWidth = 12
	accountColWidth     = 25
	categoryColWidth    = 20
	descColWidth        = 28
	amountColWidth      = 10
	numberColWidth      = 8
)

var (
	highlightColor          = lipgloss.Color("#FFD580")
	highlightForegroudColor = lipgloss.Color("#2E2E2E")
	borderColor             = lipgloss.Color("#909090")

	chartColor1 = lipgloss.Color("#FF9E6D") // coral
	chartColor2 = lipgloss.Color("#A7D5FF") // blue
	chartColor3 = lipgloss.Color("#BFFFD5") // green
	chartColor4 = lipgloss.Color("#FFC4A3") // pink
	chartColor5 = lipgloss.Color("#E9C7FF") // violet

	incomeColor  = chartColor3
	expenseColor = chartColor1
	incomeStyle  = lipgloss.NewStyle().Foreground(incomeColor)
	expenseStyle = lipgloss.NewStyle().Foreground(expenseColor)

	roundedBorder = lipgloss.RoundedBorder()

	baseStyle = lipgloss.NewStyle().
			BorderStyle(roundedBorder).
			BorderForeground(borderColor)

	searchStyle = baseStyle

	keyStyle = lipgloss.NewStyle().
			Foreground(highlightColor)
)

var (
	barChartStyles = [5]lipgloss.Style{
		lipgloss.NewStyle().Foreground(chartColor1).Background(chartColor1),
		lipgloss.NewStyle().Foreground(chartColor2).Background(chartColor2),
		lipgloss.NewStyle().Foreground(chartColor3).Background(chartColor3),
		lipgloss.NewStyle().Foreground(chartColor4).Background(chartColor4),
		lipgloss.NewStyle().Foreground(chartColor5).Background(chartColor5),
	}
)

var (
	tsChartIncomeLineStyle  = incomeStyle
	tsChartExpenseLineStyle = expenseStyle
	tsChartAxisStyle        = lipgloss.NewStyle().Foreground(highlightColor)
	tsChartLabelStyle       = lipgloss.NewStyle().Foreground(borderColor)
)

func getTableStyle() table.Styles {
	tableStyles := table.DefaultStyles()
	tableStyles.Header = tableStyles.Header.
		Foreground(highlightColor).
		BorderStyle(roundedBorder).
		BorderForeground(borderColor).
		BorderBottom(true).
		Bold(true)
	tableStyles.Selected = tableStyles.Selected.
		Foreground(highlightForegroudColor).
		Background(highlightColor).
		Bold(true)
	return tableStyles
}
