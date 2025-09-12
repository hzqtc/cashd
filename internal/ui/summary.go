package ui

import (
	"cashd/internal/data"
	"fmt"
	"sort"
	"strings"

	"github.com/NimbleMarkets/ntcharts/barchart"
	"github.com/NimbleMarkets/ntcharts/canvas/runes"
	"github.com/charmbracelet/lipgloss"
)

const (
	barChartHeight = 1

	maxSummaryEntries = 5
)

type summaryEntry struct {
	key   string
	value float64
}

type SummaryModel struct {
	width  int
	height int

	transactions  []*data.Transaction
	incomeTxnNum  int
	expenseTxnNum int
	totalIncome   float64
	totalExpense  float64

	topIncomeCategories  []summaryEntry
	topIncomeAccounts    []summaryEntry
	topExpenseCategories []summaryEntry
	topExpenseAccounts   []summaryEntry

	incomeCategoryChart  barchart.Model
	incomeAccountChart   barchart.Model
	expenseCategoryChart barchart.Model
	expenseAccountChart  barchart.Model
}

func NewSummaryModel() SummaryModel {
	return SummaryModel{}
}

func (m *SummaryModel) SetDimensions(width, height int) {
	m.width = width
	m.height = height
	m.resizeCharts()
	m.updateCharts()
}

func (m *SummaryModel) SetTransactions(transactions []*data.Transaction) {
	m.transactions = transactions

	m.incomeTxnNum = 0
	m.expenseTxnNum = 0
	m.totalIncome = 0
	m.totalExpense = 0
	for _, tx := range m.transactions {
		if tx.Type == data.Income {
			m.incomeTxnNum++
			m.totalIncome += tx.Amount
		} else {
			m.expenseTxnNum++
			m.totalExpense += tx.Amount
		}
	}

	m.topIncomeCategories, m.incomeCategoryChart = m.getTopCategories(data.Income)
	m.topIncomeAccounts, m.incomeAccountChart = m.getTopAccounts(data.Income)
	m.topExpenseCategories, m.expenseCategoryChart = m.getTopCategories(data.Expense)
	m.topExpenseAccounts, m.expenseAccountChart = m.getTopAccounts(data.Expense)

	m.updateCharts()
}

func (m *SummaryModel) updateCharts() {
	m.incomeCategoryChart.Draw()
	m.incomeAccountChart.Draw()
	m.expenseCategoryChart.Draw()
	m.expenseAccountChart.Draw()
}

func (m *SummaryModel) resizeCharts() {
	m.incomeCategoryChart.Resize(m.width-2*hPadding, barChartHeight)
	m.incomeAccountChart.Resize(m.width-2*hPadding, barChartHeight)
	m.expenseCategoryChart.Resize(m.width-2*hPadding, barChartHeight)
	m.expenseAccountChart.Resize(m.width-2*hPadding, barChartHeight)
}

func (m SummaryModel) View() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("Income transactions: %d\n", m.incomeTxnNum))
	s.WriteString(fmt.Sprintf("Expense transactions: %d\n", m.expenseTxnNum))
	s.WriteString(fmt.Sprintf("Total income: $%s\n", data.FormatMoney(m.totalIncome)))
	s.WriteString(fmt.Sprintf("Total expenses: $%s\n", data.FormatMoney(m.totalExpense)))

	if len(m.transactions) > 0 {
		s.WriteString("\nTop income categories:\n")
		s.WriteString(m.renderSummarySection(m.topIncomeCategories, m.incomeCategoryChart))

		s.WriteString("\n\nTop income accounts:\n")
		s.WriteString(m.renderSummarySection(m.topIncomeAccounts, m.incomeAccountChart))

		s.WriteString("\n\nTop expense catgories:\n")
		s.WriteString(m.renderSummarySection(m.topExpenseCategories, m.expenseCategoryChart))

		s.WriteString("\n\nTop expense accounts:\n")
		s.WriteString(m.renderSummarySection(m.topExpenseAccounts, m.expenseAccountChart))
	}

	return lipgloss.NewStyle().
		Border(getRoundedBorderWithTitle("Summary", m.width)).
		BorderForeground(borderColor).
		Width(m.width).
		Height(m.height).
		Padding(vPadding, hPadding).
		Render(s.String())
}

func (m SummaryModel) getTopCategories(txnType data.TransactionType) ([]summaryEntry, barchart.Model) {
	expenseByCategory := make(map[string]float64)
	for _, tx := range m.transactions {
		if tx.Type == txnType {
			expenseByCategory[tx.Category] += tx.Amount
		}
	}

	entries := sortAndTruncate(expenseByCategory)
	return entries, getBarChartModel(m.width-2*hPadding, entries)
}

func (m SummaryModel) getTopAccounts(txnType data.TransactionType) ([]summaryEntry, barchart.Model) {
	expenseByAccount := make(map[string]float64)
	for _, tx := range m.transactions {
		if tx.Type == txnType {
			expenseByAccount[tx.Account] += tx.Amount
		}
	}

	entries := sortAndTruncate(expenseByAccount)
	return entries, getBarChartModel(m.width-2*hPadding, entries)
}

// Sort by value in reverse order and keep the top 5 entries
func sortAndTruncate(input map[string]float64) []summaryEntry {
	var sorted []summaryEntry
	for k, v := range input {
		sorted = append(sorted, summaryEntry{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].value > sorted[j].value
	})
	if len(sorted) > maxSummaryEntries {
		// Sum [maxSummaryEntries-1:] to be "Everything else"
		sumOfRemaining := 0.0
		for _, s := range sorted[maxSummaryEntries-1:] {
			sumOfRemaining += s.value
		}
		// Keep the first items as-is
		sorted = sorted[:maxSummaryEntries-1]
		// Add Everything else as the last item
		sorted = append(sorted, summaryEntry{"Everything else", sumOfRemaining})
	}

	if len(sorted) > maxSummaryEntries {
		panic("assertion failed: each summary section should not have more than 5 items")
	}

	return sorted
}

func getBarChartModel(width int, data []summaryEntry) barchart.Model {
	barValues := []barchart.BarValue{}
	for i, item := range data {
		barValues = append(barValues, barchart.BarValue{Name: item.key, Value: item.value, Style: barChartStyles[i]})
	}
	return barchart.New(
		width,
		barChartHeight,
		barchart.WithDataSet([]barchart.BarData{{Values: barValues}}),
		barchart.WithNoAxis(),
		barchart.WithHorizontalBars(),
	)
}

func (m SummaryModel) renderSummarySection(entries []summaryEntry, chart barchart.Model) string {
	var s strings.Builder
	for i, item := range entries {
		s.WriteString(fmt.Sprintf(
			"%s %s: $%s\n",
			barChartStyles[i].Render(string(runes.FullBlock)), // Bar legend
			item.key,
			data.FormatMoney(item.value),
		))
	}
	s.WriteString("\n")

	return s.String() + chart.View()
}
