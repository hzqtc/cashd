package ui

import (
	"cashd/internal/data"
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type accountColumn int

const (
	acctColSymbol accountColumn = iota
	acctColType
	acctColName
	acctColIncome
	acctColExpense

	totalNumAcctColumns
)

func (c accountColumn) rightAligned() bool {
	return c == acctColIncome || c == acctColExpense
}

func (c accountColumn) String() string {
	switch c {
	case acctColSymbol:
		return " "
	case acctColType:
		return "Type"
	case acctColName:
		return "Name"
	case acctColIncome:
		return "Income"
	case acctColExpense:
		return "Expense"
	default:
		return "Unknown"
	}
}

var accountColWidthMap = map[accountColumn]int{
	acctColSymbol:  2,
	acctColType:    12,
	acctColName:    25,
	acctColIncome:  12,
	acctColExpense: 12,
}

var AccountTableWidth = func() int {
	tableWidth := 0
	for i := range totalNumAcctColumns {
		tableWidth += accountColWidthMap[accountColumn(i)] + 2
	}
	return tableWidth
}()

type accountInfo struct {
	accountType data.AccountType
	symbol      string
	name        string
	income      float64
	expense     float64
}

type AccountTableModel struct {
	accounts []*accountInfo
	table    table.Model
}

func NewAccountTableModel() AccountTableModel {
	columns := []table.Column{}
	for i := range int(totalNumAcctColumns) {
		col := accountColumn(i)
		title := col.String()
		width := accountColWidthMap[col]
		if col.rightAligned() {
			title = fmt.Sprintf("%*s", width, title)
		}
		columns = append(columns, table.Column{Title: title, Width: width})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithStyles(getTableStyle()),
	)
	return AccountTableModel{table: t}
}

func (m AccountTableModel) Update(msg tea.Msg) (AccountTableModel, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m AccountTableModel) View() string {
	return baseStyle.Render(m.table.View())
}

func (m *AccountTableModel) SetDimensions(width, height int) {
	m.table.SetWidth(width)
	m.table.SetHeight(height)
}

func (m *AccountTableModel) SetTransactions(transactions []*data.Transaction) {
	m.accounts = getAccountInfo(transactions)
	rows := make([]table.Row, len(m.accounts))
	for i, acct := range m.accounts {
		rowData := []string{}
		for i := range int(totalNumAcctColumns) {
			col := accountColumn(i)
			colData := getColData(acct, col)
			if col.rightAligned() {
				colData = fmt.Sprintf("%*s", accountColWidthMap[col], colData)
			}
			rowData = append(rowData, colData)
		}
		rows[i] = table.Row(rowData)
	}
	m.table.SetRows(rows)
}

// Get account-level stats by aggregating transactions
func getAccountInfo(transactions []*data.Transaction) []*accountInfo {
	var totalIncome, totalExpense float64
	accountMap := make(map[string]*accountInfo)
	for _, tx := range transactions {
		account, exist := accountMap[tx.Account]
		if !exist {
			account = &accountInfo{
				symbol:      tx.AccountSymbol(),
				accountType: tx.AccountType,
				name:        tx.Account,
			}
			accountMap[tx.Account] = account
		}
		if tx.Type == data.Income {
			account.income += tx.Amount
			totalIncome += tx.Amount
		} else {
			account.expense += tx.Amount
			totalExpense += tx.Amount
		}
	}

	accounts := []*accountInfo{
		// Add a pseudo account for "Overall" income and expense
		{
			symbol:      "",
			accountType: data.AcctOverall,
			name:        "Overall",
			income:      totalIncome,
			expense:     totalExpense,
		},
	}
	for _, a := range accountMap {
		accounts = append(accounts, a)
	}
	return accounts
}

func getColData(acct *accountInfo, col accountColumn) string {
	switch col {
	case acctColSymbol:
		return acct.symbol
	case acctColType:
		return string(acct.accountType)
	case acctColName:
		return acct.name
	case acctColIncome:
		return fmt.Sprintf("%.2f", acct.income)
	case acctColExpense:
		return fmt.Sprintf("%.2f", acct.expense)
	default:
		return ""
	}
}
