package dbengines_test

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	dbengines "gitlab.com/gjerry134679/bank/internal/models/dbEngines"
	pgxEngine "gitlab.com/gjerry134679/bank/internal/models/dbEngines/pgx"
	"gitlab.com/gjerry134679/bank/pkg/convert"
)

func TestPgxEngine(t *testing.T) {
	var src dbengines.DBSource = pgxEngine.DBSource{
		Account:  "admin",
		Password: "password",
		Address:  "127.0.0.1",
		Port:     4321,
		DBName:   "db",
		Options: map[string]string{
			"connect_timeout": "10",
		},
	}

	require.Equal(t, "pgx", src.Driver())

	db, err := dbengines.NewEngine(src)
	require.Error(t, err)
	require.Nil(t, db)

	db, err = src.Open()
	require.Error(t, err)
	require.Nil(t, db)

	err = godotenv.Load("../../../.env")
	if err != nil {
		log.Panicf("error while reading .env: %v\n", err)
	}
	src = pgxEngine.DBSource{
		Account:  os.Getenv("DB_TEST_USER_ACC"),
		Password: os.Getenv("DB_TEST_USER_PWD"),
		Address:  os.Getenv("TEST_DB_HOST"),
		Port:     convert.StrTo(os.Getenv("TEST_DB_PORT")).MustInt(),
		DBName:   os.Getenv("TEST_DB_NAME"),
	}
	db, err = src.Open()
	require.NoError(t, err)
	require.NotNil(t, db)
}
