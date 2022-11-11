package dbengines

import "database/sql"

type DBSource interface {
	ConnStr() string
	Driver() string
	String() string
	Open(name string) (*sql.Conn, error)
}
