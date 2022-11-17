package pgx

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v4/stdlib"
)

const dbDriver = "pgx"

type DBSource struct {
	Account  string            `json:"account"`
	Password string            `json:"password"`
	Address  string            `json:"address"`
	Port     int               `json:"port"`
	DBName   string            `json:"db_name"`
	Options  map[string]string `json:"options"`
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

func (dbs DBSource) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Driver: %s:\n", dbDriver))
	sb.WriteString(fmt.Sprintf(" - Account  : %s\n", dbs.Account))
	sb.WriteString(fmt.Sprintf(" - Password : %s\n", dbs.Password))
	sb.WriteString(fmt.Sprintf(" - DB Addr  : %s\n", dbs.Address))
	sb.WriteString(fmt.Sprintf(" - DB Port  : %d\n", dbs.Port))
	sb.WriteString(fmt.Sprintf(" - DB Name  : %s\n", dbs.DBName))
	if len(dbs.Options) > 0 {
		sb.WriteString(" - DB Opts  :\n")
		for key, val := range dbs.Options {
			sb.WriteString(fmt.Sprintf("   - %s = %s", key, val))
		}
	}
	return sb.String()
}

func (dbs DBSource) Marshal() ([]byte, error) {
	return json.Marshal(dbs)
}

func (dbs DBSource) Driver() string {
	return dbDriver
}

func (dbs DBSource) Open() (*sql.DB, error) {
	db, err := sql.Open(dbs.Driver(), dbs.ConnStr())
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, err
}

func NewDB(dbSource DBSource) (*sql.DB, error) {
	return dbSource.Open()
}
