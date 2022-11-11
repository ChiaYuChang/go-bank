package models_test

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.com/gjerry134679/bank/internal/models"
)

func AccountIsEq(a1, a2 models.Account) (bool, []string) {
	fields := make([]string, 0)
	if a1.ID != a2.ID {
		fields = append(fields, fmt.Sprintf("Id want: %v, get: %v", a1.ID, a2.ID))
	}

	if a1.Owner != a2.Owner {
		fields = append(fields, fmt.Sprintf("owner want: %v, get: %v", a1.Owner, a2.Owner))
	}

	if !a1.Balance.Equal(a2.Balance) {
		fields = append(fields, fmt.Sprintf("balance want: %v, get: %v", a1.Balance, a2.Balance))
	}

	if a1.CreatedAt.Sub(a2.CreatedAt) > 1*time.Minute {
		fields = append(fields, fmt.Sprintf("created_at want: %v, get: %v", a1.CreatedAt, a2.CreatedAt))
	}

	if len(fields) > 0 {
		return false, fields
	}
	return true, fields
}

func TestCreateAccount(t *testing.T) {
	c, err := testQureies.CreateCurrency(
		context.Background(),
		models.CreateCurrencyParams{
			Name: "Ethereum",
			Abbr: "ETH",
		})
	if err != nil {
		t.Fatalf("error while creating currency: %v", err)
	}

	owner := "owner_001"
	balance := decimal.NewFromInt(1000)
	acc, err := testQureies.CreateAccount(
		context.Background(),
		models.CreateAccountParams{
			Owner:    owner,
			Balance:  balance,
			Currency: c.ID,
		})
	if err != nil {
		t.Fatalf("error while creating account: %v", err)
	}

	assert.Equal(t, owner, acc.Owner)
	assert.Equal(t, balance, acc.Balance)
	assert.Equal(t, c.ID, acc.Currency)
	assert.IsType(t, c.DeletedAt, sql.NullTime{})
}

func mustCreateAnTestAccount(t *testing.T, n int, cid int32) []models.Account {
	if n > 10 {
		t.Fatalf("n should not be greater than 10")
	}
	accs := make([]models.Account, n)

	for i := 0; i < n; i++ {
		acc, err := testQureies.CreateAccount(
			context.Background(),
			models.CreateAccountParams{
				Owner:    fmt.Sprintf("Owner%2d", i),
				Balance:  decimal.NewFromInt(int64(i * 243)),
				Currency: cid,
			})
		if err != nil {
			t.Fatalf("error while calling CreateCurrency method: %v", err)
		}
		accs[i] = acc
	}
	return accs
}

func TestGetAccount(t *testing.T) {
	cs := mustCreateAnTestCurrency(t, 1)
	c := cs[0]

	accs := mustCreateAnTestAccount(t, 1, c.ID)
	acc1 := accs[0]

	acc2, err := testQureies.GetAccount(
		context.Background(), acc1.ID)
	if err != nil {
		t.Fatalf("error while creating account: %v", err)
	}

	ok, ef := AccountIsEq(acc1, acc2)
	if !ok {
		t.Fatalf("retrieving data mismatch in fields: %s", strings.Join(ef, ", "))
	}
}

func TestUpdateAccount(t *testing.T) {
	cs := mustCreateAnTestCurrency(t, 1)
	c := cs[0]

	accs := mustCreateAnTestAccount(t, 1, c.ID)
	acc1 := accs[0]

	acc2, err := testQureies.UpdateAccount(
		context.Background(),
		models.UpdateAccountParams{
			ID:        acc1.ID,
			Owner:     "owner_001_updated",
			Balance:   decimal.NewFromInt(200),
			Currency:  c.ID,
			UpdatedAt: time.Now(),
		})
	if err != nil {
		t.Fatalf("error while updating account: %v", err)
	}

	assert.Equal(t, acc1.ID, acc2.ID)
	assert.NotEqual(t, acc1.Owner, acc2.Owner)
	assert.Equal(t, acc2.Owner, "owner_001_updated")
	assert.NotEqual(t, acc1.Balance, acc2.Balance)
	if !acc2.Balance.Equal(decimal.NewFromInt(200)) {
		t.Fatalf("balance field did not update")
	}
	assert.Equal(t, acc1.CreatedAt, acc2.CreatedAt)
	assert.Equal(t, acc1.Currency, acc2.Currency)
	assert.NotEqual(t, acc1.UpdatedAt, acc2.UpdatedAt)
}

func TestListAccount(t *testing.T) {
	cs := mustCreateAnTestCurrency(t, 2)
	c := cs[0]

	mustCreateAnTestAccount(t, 10, c.ID)

	acclst1, err := testQureies.ListAccounts(
		context.Background(),
		models.ListAccountsParams{
			Limit:  10,
			Offset: 0,
		})
	if err != nil {
		t.Fatalf("error while listing account: %v", err)
	}

	acclst2, err := testQureies.ListAccounts(
		context.Background(),
		models.ListAccountsParams{
			Limit:  10,
			Offset: 5,
		})
	if err != nil {
		t.Fatalf("error while listing account: %v", err)
	}

	for i := 0; i < 5; i++ {
		a1 := acclst1[i+5]
		a2 := acclst2[i]
		ok, fs := AccountIsEq(a1, a2)
		if !ok {
			t.Fatalf("error while listing account: %s", fs)
		}
	}
}

func TestDeleteAccound(t *testing.T) {
	cs := mustCreateAnTestCurrency(t, 1)
	c := cs[0]

	accs := mustCreateAnTestAccount(t, 1, c.ID)
	acc := accs[0]

	err := testQureies.DeleteAccount(context.Background(), acc.ID)
	if err != nil {
		t.Fatalf("error while deleting account: %v", err)
	}

	_, err = testQureies.GetAccount(context.Background(), acc.ID)
	if err == nil {
		t.Fatalf("data with id: %v did not be delete", err)
	}
	if err != sql.ErrNoRows {
		t.Fatalf("error while retrieving account data: %v", err)
	}
}
