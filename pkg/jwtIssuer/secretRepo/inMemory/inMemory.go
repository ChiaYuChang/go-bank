package inMemory

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"gitlab.com/gjerry134679/bank/pkg/jwtIssuer/secretRepo"
)

type InMemoryDB struct {
	Secrets sync.Map
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		Secrets: sync.Map{},
	}
}

func (sr *InMemoryDB) Get(key uuid.UUID) (*secretRepo.Secret, error) {
	srcRaw, ok := sr.Secrets.Load(key)
	if !ok {
		return nil, secretRepo.ErrSecretNotFound
	}

	srct, ok := srcRaw.(*secretRepo.Secret)
	if !ok {
		return nil, secretRepo.ErrBrokenSecret
	}

	// Check time
	t := time.Now()
	if t.Before(srct.ValidAt) {
		return nil, secretRepo.ErrSecretNotYetValid
	}
	if t.After(srct.ExpiredAt) {
		return nil, secretRepo.ErrSecretExpired
	}
	return srct, nil
}

func (sr *InMemoryDB) Create(srct *secretRepo.Secret) error {
	_, err := sr.Get(srct.Key)
	if err == nil {
		return secretRepo.ErrSecretKeyCollision
	}
	sr.Secrets.Store(srct.Key, srct)
	return nil
}

func (sr *InMemoryDB) Delete(key uuid.UUID) error {
	sr.Secrets.Delete(key)
	return nil
}

func (sr *InMemoryDB) Flush() error {
	sr.Secrets = sync.Map{}
	return nil
}
