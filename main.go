package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gitlab.com/gjerry134679/bank/third_party/config_generator/sqlc"
	"gitlab.com/gjerry134679/bank/third_party/config_generator/sqlc/gen"
	"gopkg.in/yaml.v3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error while reading .env file: %v", err)
	}
	config := sqlc.NewSQLCConfig()
	postgresSQL := sqlc.NewSQL(
		os.Getenv("SQLC_SCHEMA"),
		os.Getenv("SQLC_QUERIES"),
		sqlc.SQLCEngn(os.Getenv("SQLC_ENGINE")),
	)
	goGen := gen.NewGoGen(
		os.Getenv("SQLC_GEN_PKG"),
		os.Getenv("SQLC_GEN_OUT"),
	)
	postgresSQL.Gen[sqlc.Golang] = goGen
	config.SQL = append(config.SQL, postgresSQL)
	yml, err := yaml.Marshal(config)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Printf("%s", string(yml))
}
