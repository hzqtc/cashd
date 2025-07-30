package data

type DataSource interface {
	LoadTransactions() ([]Transaction, error)
}
