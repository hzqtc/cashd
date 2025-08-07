package data

type DataSource interface {
	// Returned transactions must be ordered by date, from earliest to oldest
	LoadTransactions() ([]*Transaction, error)

	// Whether the data source is preferred
	Preferred() bool

	// Whether the data source is enabled
	Enabled() bool
}
