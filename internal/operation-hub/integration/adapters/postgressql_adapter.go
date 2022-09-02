package adapters

import (
	"database/sql"
)

type PostgresSQLAdapter interface {
	OpenConnection() (*sql.DB, error)
}
