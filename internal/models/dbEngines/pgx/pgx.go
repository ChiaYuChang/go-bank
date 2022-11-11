package pgx

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v4/stdlib"
)

const dbDriver = "postgres"

type DBSource struct {
	Account  string
	Password string
	Address  string
	Port     int
	DBName   string
	Options  map[string]string
	Ctx      context.Context
}

func (dbs DBSource) ConnStr() string {
	s := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?",
		dbs.Account,
		dbs.Password,
		dbs.Address,
		dbs.Port,
		dbs.DBName,
	)

	opts := make([]string, 0)
	for key, val := range dbs.Options {
		opts = append(opts, fmt.Sprintf("%s=%s", key, val))
	}

	s += strings.Join(opts, "&")
	return s
}

func (db DBSource) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Driver: %s:\n", dbDriver))
	sb.WriteString(fmt.Sprintf(" - Account  : %s\n", db.Account))
	sb.WriteString(fmt.Sprintf(" - Password : %s\n", db.Password))
	sb.WriteString(fmt.Sprintf(" - DB Addr  : %s\n", db.Address))
	sb.WriteString(fmt.Sprintf(" - DB Port  : %d\n", db.Port))
	sb.WriteString(fmt.Sprintf(" - DB Name  : %s\n", db.DBName))
	if len(db.Options) > 0 {
		sb.WriteString(" - DB Opts  :")
		for key, val := range db.Options {
			sb.WriteString(fmt.Sprintf("   - %s = %s", key, val))
		}
	}
	return sb.String()
}

func (db DBSource) Driver() string {
	return dbDriver
}

func NewDB(dbSource DBSource) (*sql.DB, error) {
	db, err := sql.Open(dbDriver, dbSource.ConnStr())
	if err != nil {
		return nil, err
	}
	return db, err
}
