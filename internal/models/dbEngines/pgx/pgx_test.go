package pgx_test

import (
	"log"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	pgxEngine "gitlab.com/gjerry134679/bank/internal/models/dbEngines/pgx"
	"gitlab.com/gjerry134679/bank/pkg/convert"
)

var dbSource pgxEngine.DBSource

func init() {
	err := godotenv.Load("../../../../.env")
	if err != nil {
		log.Panicf("error while reading .env: %v\n", err)
	}
	dbSource = pgxEngine.DBSource{
		Account:  os.Getenv("DB_TEST_USER_ACC"),
		Password: os.Getenv("DB_TEST_USER_PWD"),
		Address:  os.Getenv("TEST_DB_HOST"),
		Port:     convert.StrTo(os.Getenv("TEST_DB_PORT")).MustInt(),
		DBName:   os.Getenv("TEST_DB_NAME"),
	}
}

func TestConnStr(t *testing.T) {
	src := pgxEngine.DBSource{
		Account:  "admin",
		Password: "password",
		Address:  "127.0.0.1",
		Port:     4321,
		DBName:   "db",
		Options: map[string]string{
			"connect_timeout": "10",
		},
	}
	t.Logf("%v", src.String())
	require.Equal(t, "postgresql://admin:password@127.0.0.1:4321/db?connect_timeout=10", src.ConnStr())
}

func TestNewDB(t *testing.T) {
	errDBSource := pgxEngine.DBSource{
		Account:  "admin",
		Password: "password",
		Address:  "127.0.0.1",
		Port:     43210,
		DBName:   "db",
		Options: map[string]string{
			"connect_timeout": "10",
		},
	}

	require.Equal(t, "pgx", dbSource.Driver())

	db, err := pgxEngine.NewDB(errDBSource)
	require.Error(t, err)
	require.Nil(t, db)

	db, err = dbSource.Open()
	require.NoError(t, err)
	require.NotNil(t, db)
}
