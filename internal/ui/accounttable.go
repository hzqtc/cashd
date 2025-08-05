package ui

import (
	"cashd/internal/data"
	"fmt"
	"sort"

	"github.com/charmbracelet/bubbles/table"
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

func (c accountColumn) index() int {
	return int(c)
}

func (c accountColumn) rightAligned() bool {
	return c == acctColIncome || c == acctColExpense
}

func (c accountColumn) isSortable() bool {
	return c != acctColSymbol
}

func (c accountColumn) width() int {
	return accountColWidthMap[c]
}

func (c accountColumn) nextColumn() column {
	return column(accountColumn((int(c) + 1) % int(totalNumAcctColumns)))
}

func (c accountColumn) prevColumn() column {
	return column(accountColumn((int(c) - 1 + int(totalNumAcctColumns)) % int(totalNumAcctColumns)))
}

func (c accountColumn) getColumnData(a any) string {
	account, ok := a.(*accountInfo)
	if !ok {
		panic(fmt.Sprintf("can't convert to accountInfo: %v", a))
	}
	switch c {
	case acctColSymbol:
		return account.symbol
	case acctColType:
		return string(account.accountType)
	case acctColName:
		return account.name
	case acctColIncome:
		return fmt.Sprintf("%.2f", account.income)
	case acctColExpense:
		return fmt.Sprintf("%.2f", account.expense)
	default:
		return ""
	}
}

func (c accountColumn) String() string {
	switch c {
	case acctColSymbol:
		return " "
	case acctColType:
		return "Type"
	case acctColName:
		return "Account"
	case acctColIncome:
		return "Income"
	case acctColExpense:
		return "Expense"
	default:
		return "Unknown"
	}
}

var accountColWidthMap = map[accountColumn]int{
	acctColSymbol:  symbolColWidth,
	acctColType:    accountTypeColWidth,
	acctColName:    accountColWidth,
	acctColIncome:  amountColWidth,
	acctColExpense: amountColWidth,
}

const AccountNameTotal = "Total Accounts"

var AccountTableWidth = func() int {
	tableWidth := 0
	for i := range totalNumAcctColumns {
		tableWidth += accountColWidthMap[accountColumn(i)] + 2
	}
	return tableWidth
}()

var AccountTableConfig = TableConfig{
	columns: func() []column {
		cols := []column{}
		for i := range int(totalNumAcctColumns) {
			cols = append(cols, column(accountColumn(i)))
		}
		return cols
	}(),
	dataProvider:      accountTableDataProvider,
	rowId:             func(row table.Row) string { return row[acctColName] },
	defaultSortColumn: column(acctColName),
	defaultSortDir:    sortAsc,
}

type accountInfo struct {
	accountType data.AccountType
	symbol      string
	name        string
	income      float64
	expense     float64
}

func accountTableDataProvider(transactions []*data.Transaction) tableDataSorter {
	accounts := getAccountInfo(transactions)

	return func(sortCol column, sortDir sortDirection) []table.Row {
		sort.Slice(accounts, func(i, j int) bool {
			a, b := accounts[i], accounts[j]
			// Make sure AccountTotal stay on top of the table
			if a.name == AccountNameTotal {
				return true
			} else if b.name == AccountNameTotal {
				return false
			}
			// Compare the rest accounts by specified column
			inOrder := sortCol.getColumnData(a) < sortCol.getColumnData(b)
			if sortDir == sortDesc {
				return !inOrder
			} else {
				return inOrder
			}
		})

		rows := make([]table.Row, len(accounts))
		for i, acct := range accounts {
			row := []string{}
			for i := range int(totalNumAcctColumns) {
				col := accountColumn(i)
				colData := col.getColumnData(acct)
				if col.rightAligned() {
					colData = fmt.Sprintf("%*s", accountColWidthMap[col], colData)
				}
				row = append(row, colData)
			}
			rows[i] = table.Row(row)
		}
		return rows
	}
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
			name:        AccountNameTotal,
			income:      totalIncome,
			expense:     totalExpense,
		},
	}
	for _, a := range accountMap {
		accounts = append(accounts, a)
	}
	return accounts
}
