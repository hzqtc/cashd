package ui

import (
	"cashd/internal/data"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/NimbleMarkets/ntcharts/canvas/runes"
	"github.com/charmbracelet/lipgloss"
)

type insight struct {
	income         float64
	expense        float64
	incomeTxnNum   int
	expenseTxnNum  int
	topIncomeTxns  []*data.Transaction
	topExpenseTxns []*data.Transaction
}

const (
	topTxnNum = 3
)

type InsightsModel struct {
	name   string
	ins    insight
	width  int
	height int
}

func NewInsightsModel() InsightsModel {
	return InsightsModel{}
}

func (m *InsightsModel) SetTransactionsWithAccount(transactions []*data.Transaction, account string) {
	m.updateInsights(transactions, func(t *data.Transaction) bool {
		return account == AccountNameTotal || t.Account == account
	})
}

func (m *InsightsModel) SetTransactionsWithCategory(transactions []*data.Transaction, category string) {
	m.updateInsights(transactions, func(t *data.Transaction) bool {
		return t.Category == category
	})
}

func (m *InsightsModel) updateInsights(transactions []*data.Transaction, match func(*data.Transaction) bool) {
	m.ins = insight{}
	incomeTxns := []*data.Transaction{}
	expenseTxns := []*data.Transaction{}

	for _, t := range transactions {
		if match(t) {
			switch t.Type {
			case data.Income:
				m.ins.income += t.Amount
				incomeTxns = append(incomeTxns, t)
			case data.Expense:
				m.ins.expense += t.Amount
				expenseTxns = append(expenseTxns, t)
			}
		}
	}

	m.ins.incomeTxnNum = len(incomeTxns)
	m.ins.expenseTxnNum = len(expenseTxns)

	sort.Slice(incomeTxns, func(i, j int) bool {
		return incomeTxns[i].Amount > incomeTxns[j].Amount
	})
	sort.Slice(expenseTxns, func(i, j int) bool {
		return expenseTxns[i].Amount > expenseTxns[j].Amount
	})
	if m.ins.incomeTxnNum <= topTxnNum {
		m.ins.topIncomeTxns = incomeTxns
	} else {
		m.ins.topIncomeTxns = incomeTxns[:topTxnNum]
	}
	if m.ins.expenseTxnNum <= topTxnNum {
		m.ins.topExpenseTxns = expenseTxns
	} else {
		m.ins.topExpenseTxns = expenseTxns[:topTxnNum]
	}
}

func (m *InsightsModel) SetDimension(width, height int) {
	m.width = width
	m.height = height
}

func (m *InsightsModel) SetName(name string) {
	m.name = name
}

func (m *InsightsModel) Height() int {
	return lipgloss.Height(m.View())
}

func (m InsightsModel) View() string {
	var s strings.Builder

	if m.ins.incomeTxnNum > 0 {
		s.WriteString("\n")
		s.WriteString(fmt.Sprintf(
			"%s Income: $%s in %d transactions\n",
			incomeStyle.Render(string(runes.FullBlock)),
			data.FormatMoney(m.ins.income),
			m.ins.incomeTxnNum,
		))
		s.WriteString("Top transactions:\n")
		for _, t := range m.ins.topIncomeTxns {
			s.WriteString(formatTransaction(t))
		}
	}

	if m.ins.expenseTxnNum > 0 {
		s.WriteString("\n")
		s.WriteString(fmt.Sprintf(
			"%s Expense: $%s in %d transactions\n",
			expenseStyle.Render(string(runes.FullBlock)),
			data.FormatMoney(m.ins.expense),
			m.ins.expenseTxnNum,
		))
		s.WriteString("Top transactions:\n")
		for _, t := range m.ins.topExpenseTxns {
			s.WriteString(formatTransaction(t))
		}
	}

	return lipgloss.NewStyle().
		Border(getRoundedBorderWithTitle(m.name, m.width)).
		BorderForeground(borderColor).
		Width(m.width).
		Height(m.height).
		Padding(vPadding, hPadding).
		Render(s.String())
}

func formatTransaction(t *data.Transaction) string {
	return fmt.Sprintf(
		"%s %*s %s\n",
		t.Date.Format(time.DateOnly),
		amountColWidth,
		fmt.Sprintf("$%s", t.FormattedAmount()),
		t.Description,
	)
}
