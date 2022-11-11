package models

import (
	"os"
	"testing"
)

func TestTransferTx(t *testing.T) {
	t.Logf("db pwd: %s", os.Getenv("DB_TEST_USER_PWD"))
}
