package models_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"gitlab.com/gjerry134679/bank/internal/models"
)

const dbDriver string = "pgx"

var dbSource string
var testDB *sql.DB
var testQureies *models.Queries

func init() {
	var err error
	err = godotenv.Load("../../.env")
	if err != nil {
		log.Panicf("error while reading .env: %v\n", err)
	}

	dbSource = fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		os.Getenv("DB_TEST_USER_ACC"),
		os.Getenv("DB_TEST_USER_PWD"),
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_PORT"),
		os.Getenv("TEST_DB_NAME"),
	)
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Panicf("error while calling sql.Open: %v\n", err)
	}
	testQureies = models.New(testDB)
}
