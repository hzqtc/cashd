package ui

import (
	"fmt"
	"lledger/internal/data"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type SummaryModel struct {
	width        int
	height       int
	totalIncome  float64
	totalExpense float64
	transactions []*data.Transaction
}

func NewSummaryModel() SummaryModel {
	return SummaryModel{}
}

func (m *SummaryModel) SetDimensions(width, height int) {
	m.width = width
	m.height = height
}

func (m *SummaryModel) SetTransactions(transactions []*data.Transaction) {
	m.transactions = transactions
	m.totalIncome = 0
	m.totalExpense = 0
	for _, tx := range transactions {
		if tx.Type == data.Income {
			m.totalIncome += tx.Amount
		} else {
			m.totalExpense += tx.Amount
		}
	}
}

func (m SummaryModel) View() string {
	if len(m.transactions) == 0 {
		return ""
	}

	var s strings.Builder
	s.WriteString(fmt.Sprintf("Total income: $%.2f\n", m.totalIncome))
	s.WriteString(fmt.Sprintf("Total expenses: $%.2f\n", m.totalExpense))
	s.WriteString("\nTop income categories:\n")
	s.WriteString(m.topCategoriesView(data.Income))
	s.WriteString("\nTop income accounts:\n")
	s.WriteString(m.topAccountsView(data.Income))
	s.WriteString("\nTop expense catgories:\n")
	s.WriteString(m.topCategoriesView(data.Expense))
	s.WriteString("\nTop expense accounts:\n")
	s.WriteString(m.topAccountsView(data.Expense))

	return lipgloss.NewStyle().Padding(2).Render(s.String())
}

func (m SummaryModel) topCategoriesView(txnType data.TransactionType) string {
	expenseByCategory := make(map[string]float64)
	for _, tx := range m.transactions {
		if tx.Type == txnType {
			expenseByCategory[tx.Category] += tx.Amount
		}
	}

	return m.renderSummarySection(expenseByCategory)
}

func (m SummaryModel) topAccountsView(txnType data.TransactionType) string {
	expenseByAccount := make(map[string]float64)
	for _, tx := range m.transactions {
		if tx.Type == txnType {
			expenseByAccount[tx.Account] += tx.Amount
		}
	}

	return m.renderSummarySection(expenseByAccount)
}

func (m SummaryModel) renderSummarySection(data map[string]float64) string {
	type kv struct {
		Key   string
		Value float64
	}

	var sorted []kv
	for k, v := range data {
		sorted = append(sorted, kv{k, v})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	var s strings.Builder
	for i, item := range sorted {
		if i >= 5 {
			break
		}
		s.WriteString(fmt.Sprintf("%s: $%.2f\n", item.Key, item.Value))
	}

	return s.String()
}

