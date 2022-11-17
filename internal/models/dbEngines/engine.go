package dbengines

import (
	"database/sql"
)

type DBSource interface {
	ConnStr() string
	Driver() string
	String() string
	Open() (*sql.DB, error)
	Marshal() ([]byte, error)
}

func NewEngine(dbSource DBSource) (*sql.DB, error) {
	return dbSource.Open()
}
