package dbengines

import (
	"database/sql"
)

func NewEngine(dbSource DBSource) (*sql.DB, error) {
	engine, err := sql.Open(dbSource.Driver(), dbSource.ConnStr())
	if err != nil {
		return nil, err
	}
	return engine, nil
}
