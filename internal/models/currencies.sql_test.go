package models_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/gjerry134679/bank/internal/models"
)

func CurrencyIsEq(c1 models.Currency, c2 models.Currency) bool {
	if c1.ID != c2.ID {
		return false
	}

	if c1.Name != c2.Name {
		return false
	}

	if c1.Abbr != c2.Abbr {
		return false
	}

	if c1.CreatedAt.Sub(c2.CreatedAt).Abs() > 1*time.Minute {
		return false
	}
	return true
}

func TestCreatCurrency(t *testing.T) {
	cName := "United States Dollar"
	cAbbr := "USD"
	ct := time.Now()
	c, err := testQureies.CreateCurrency(
		context.Background(),
		models.CreateCurrencyParams{
			Name: cName,
			Abbr: cAbbr,
		})
	if err != nil {
		t.Fatalf("error while calling CreateCurrency method: %v", err)
	}

	assert.Equal(t, cName, c.Name)
	assert.Equal(t, cAbbr, c.Abbr)
	if ct.Sub(c.CreatedAt).Abs() > 1*time.Minute {
		t.Logf("created_time mismatch: want: %v get: %v", ct, c.CreatedAt)
	}
}

func mustCreateAnTestCurrency(t *testing.T, n int) []models.Currency {
	if n > 20 {
		t.Fatalf("n should not be greater than 20")
	}
	cs := make([]models.Currency, n)

	for i := 0; i < n; i++ {
		c, err := testQureies.CreateCurrency(
			context.Background(),
			models.CreateCurrencyParams{
				Name: fmt.Sprintf("Test Currency %d", i),
				Abbr: fmt.Sprintf("T%02d", i),
			})
		if err != nil {
			t.Fatalf("error while calling CreateCurrency method: %v", err)
		}
		cs[i] = c
	}
	return cs
}

func TestGetCurrency(t *testing.T) {
	cs := mustCreateAnTestCurrency(t, 3)

	for _, c1 := range cs {
		t.Logf("retriving currency with id %d", c1.ID)
		c2, err := testQureies.GetCurrency(context.Background(), c1.ID)
		if err != nil {
			t.Fatalf("error while retrieving currency data: %v", err)
		}

		if !CurrencyIsEq(c1, c2) {
			t.Fatal("retrieving data mismatch")
		}
	}
}

func TestUpdateCurrency(t *testing.T) {
	var err error
	cs := mustCreateAnTestCurrency(t, 1)
	c1 := cs[0]

	cName := "New Doggy Coin"
	cAbbr := "NDC"
	_, err = testQureies.UpdateCurrency(
		context.Background(),
		models.UpdateCurrencyParams{
			ID:   c1.ID,
			Name: cName,
			Abbr: cAbbr,
		})
	if err != nil {
		t.Fatalf("error while calling UpdateCurrency method: %v", err)
	}

	c2, err := testQureies.GetCurrency(context.Background(), c1.ID)
	if err != nil {
		t.Fatalf("error while retrieving currency data: %v", err)
	}

	assert.Equal(t, c1.ID, c2.ID)
	assert.NotEqual(t, c1.Name, c2.Name)
	assert.Equal(t, c2.Name, cName)
	assert.NotEqual(t, c1.Abbr, c2.Abbr)
	assert.Equal(t, c2.Abbr, cAbbr)
	assert.Equal(t, c1.CreatedAt, c2.CreatedAt)
	assert.NotEqual(t, c1.UpdatedAt, c2.UpdatedAt)
}

func TestListCurrencies(t *testing.T) {
	mustCreateAnTestCurrency(t, 10)

	cs2, err := testQureies.ListCurrencies(
		context.Background(),
		models.ListCurrenciesParams{Limit: 5, Offset: 0})
	if err != nil {
		t.Fatalf("error while calling ListCurrencies method: %v", err)
	}
	assert.Equal(t, len(cs2), 5)

	cs3, err := testQureies.ListCurrencies(
		context.Background(),
		models.ListCurrenciesParams{Limit: 5, Offset: 2})
	if err != nil {
		t.Fatalf("error while calling ListCurrencies method: %v", err)
	}
	assert.Equal(t, len(cs3), 5)

	for i := 0; i < 3; i++ {
		c2 := cs2[i+2]
		c3 := cs3[i]
		if !CurrencyIsEq(c2, c3) {
			t.Fatal("error while listing currencies")
		}
	}
}

func TestDeleteCurrency(t *testing.T) {
	cs := mustCreateAnTestCurrency(t, 10)

	for _, c := range cs {
		t.Logf("retriving currency with id %d", c.ID)
		err := testQureies.DeleteCurrency(context.Background(), c.ID)
		if err != nil {
			t.Fatalf("error while deleting currency data: %v", err)
		}
	}

	for _, c := range cs {
		t.Logf("retriving currency with id %d", c.ID)
		_, err := testQureies.GetCurrency(context.Background(), c.ID)

		if err == nil {
			t.Fatalf("data with id: %v did not be delete", err)
		}
		if err != sql.ErrNoRows {
			t.Fatalf("error while retrieving currency data: %v", err)
		}
	}
}
