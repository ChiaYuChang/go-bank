package models_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"gitlab.com/gjerry134679/bank/internal/models"
	dbengines "gitlab.com/gjerry134679/bank/internal/models/dbEngines"
	pgxEngine "gitlab.com/gjerry134679/bank/internal/models/dbEngines/pgx"
	"gitlab.com/gjerry134679/bank/pkg/convert"
)

var dbSource dbengines.DBSource
var testDB *sql.DB
var testQureies *models.Queries

func init() {
	var err error
	err = godotenv.Load("../../.env")
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
	testDB, err = dbSource.Open()
	if err != nil {
		log.Panicf("error while calling sql.Open: %v\n", err)
	}
	testQureies = models.New(testDB)
}

func TestPrep(t *testing.T) {
	pq, err := models.Prepare(context.Background(), testDB)
	require.NoError(t, err)
	require.NotNil(t, pq)

	ct := time.Now()
	cName := "Ethereum"
	cAbbr := "ETH"
	c, err := pq.CreateCurrency(
		context.Background(),
		models.CreateCurrencyParams{
			Name: cName,
			Abbr: cAbbr,
		})
	if err != nil {
		t.Fatalf("error while creating currency: %v", err)
	}

	require.Equal(t, cName, c.Name)
	require.Equal(t, cAbbr, c.Abbr)
	if ct.Sub(c.CreatedAt).Abs() > 1*time.Minute {
		t.Logf("created_time mismatch: want: %v get: %v", ct, c.CreatedAt)
	}
}
