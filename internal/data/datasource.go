package data

type DataSource interface {
	// Returned transactions must be ordered by date, from earliest to oldest
	LoadTransactions() ([]*Transaction, error)
}
