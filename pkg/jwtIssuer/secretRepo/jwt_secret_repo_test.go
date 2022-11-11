package secretRepo_test

import (
	"errors"
	"log"
	"runtime"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"gitlab.com/gjerry134679/bank/pkg/jwtIssuer/secretRepo"
	"gitlab.com/gjerry134679/bank/pkg/jwtIssuer/secretRepo/inMemory"
)

var jwtSrctDB secretRepo.SecretDB
var jwtSrctRepo *secretRepo.JWTSecretRepo

func init() {
	var err error
	jwtSrctDB = inMemory.NewInMemoryDB()
	secretRepo.SetRandomSeed(1111)
	secretRepo.InitRandomGenerator()

	jwtSrctRepo, err = secretRepo.NewJWTSecretRepo(
		16, 2*time.Second, 2*time.Second, jwtSrctDB,
	)
	if err != nil {
		log.Fatalf("error while creating jwt secret repo")
	}
}

func TestGetSecret(t *testing.T) {
	secretRepo.SetRandomSeed(1111)
	secretRepo.InitRandomGenerator()

	t.Log("query by random uuid...")
	_, ok := jwtSrctRepo.GetSecret(uuid.New())
	if ok {
		t.Fatal("should not retrieve secret invalid secret key")
	}

	n := runtime.NumCPU() * 2
	t.Logf("cpu num: %d", n/2)
	echn := make(chan error, 1)

	key1, srct1, err := jwtSrctRepo.UpdateCurrSecret()
	if err != nil {
		t.Fatalf("error while updating current secret: %v", err)
	}
	for i := 0; i < n; i++ {
		go func(ch chan<- error) {
			srct2, ok := jwtSrctRepo.GetSecret(key1)
			if !ok {
				echn <- errors.New("error while retrieve secret")
			}
			assert.Equal(t, srct1, srct2)
			echn <- nil
		}(echn)
	}
	for i := 0; i < n; i++ {
		e := <-echn
		if e != nil {
			t.Fatal(e)
		}
	}

	key2, srct3, err := jwtSrctRepo.UpdateCurrSecret()
	if err != nil {
		t.Fatalf("error while updating current secret: %v", err)
	}
	for i := 0; i < n; i++ {
		go func(ch chan<- error) {
			srct4, ok := jwtSrctRepo.GetSecret(key2)
			if !ok {
				echn <- errors.New("error while retrieve secret")
			}
			assert.Equal(t, srct3, srct4)
			echn <- nil
		}(echn)
	}
	for i := 0; i < n; i++ {
		e := <-echn
		if e != nil {
			t.Fatal(e)
		}
	}

	// test not valid yet
	for i := 0; i < n; i++ {
		go func(ch chan<- error) {
			_, ok := jwtSrctRepo.GetSecret(key1)
			if ok {
				ch <- errors.New("should not retrieve secret not yet valid")
			}
			echn <- nil
		}(echn)
	}
	for i := 0; i < n; i++ {
		e := <-echn
		if e != nil {
			t.Fatal(e)
		}
	}

	t.Log("wait until secret valid...")
	time.Sleep(2 * time.Second)
	srct5, ok := jwtSrctRepo.GetSecret(key1)
	if !ok {
		t.Fatal("could not retrieve secret")
	}
	assert.Equal(t, srct1, srct5)

	// test expired time
	t.Log("wait until secret expired...")
	time.Sleep(2 * time.Second)
	for i := 0; i < n; i++ {
		go func(ch chan<- error) {
			_, ok := jwtSrctRepo.GetSecret(key1)
			if ok {
				ch <- errors.New("should not retrieve secret expired")
			}
			echn <- nil
		}(echn)
	}

	for i := 0; i < n; i++ {
		e := <-echn
		if e != nil {
			t.Fatal(e)
		}
	}
}

func TestUpdateCurrSecret(t *testing.T) {
	secretRepo.SetRandomSeed(1111)
	secretRepo.InitRandomGenerator()

	jwtSrctRepo.UpdateCurrSecret()
	assert.Equal(
		t, string(jwtSrctRepo.CurrSecret), "TQXuSADf0OXp46HC",
	)

	jwtSrctRepo.UpdateCurrSecret()
	assert.Equal(
		t, string(jwtSrctRepo.CurrSecret), "hFEyf00EkVij5MPe",
	)

	n := runtime.NumCPU() * 5
	t.Logf("cpu num: %d", n/2)
	echn := make(chan error, 1)

	for i := 0; i < n; i++ {
		go func(ch chan<- error) {
			_, _, err := jwtSrctRepo.UpdateCurrSecret()
			echn <- err
		}(echn)
	}

	for i := 0; i < n; i++ {
		e := <-echn
		if e != nil {
			t.Fatal(e)
		}
	}
}
