package adapters

import "database/sql"

// SQLAdapter is an adapter for sql.Open. This is necessary to help unit test some classes.
type SQLAdapter interface {
	Open(driverName, dataSourceName string) (*sql.DB, error)
}
