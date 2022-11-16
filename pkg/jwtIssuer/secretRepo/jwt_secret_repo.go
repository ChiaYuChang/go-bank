package secretRepo

import (
	crand "crypto/rand"
	"errors"
	"fmt"
	mrand "math/rand"
	"time"

	"github.com/google/uuid"
)

var DefaultJWTSecretRepo JWTSecretRepo
var ErrSecretKeyCollision = errors.New("secret key has been used")
var ErrBrokenSecret = errors.New("secret could not be parsed")
var ErrSecretNotFound = errors.New("secret not found")
var ErrCannotGenSecret = errors.New("error while generating secret")
var ErrSecretExpired = errors.New("secret has expired")
var ErrSecretNotYetValid = errors.New("secret not yet valid")

// var salt = []byte("B04ZS75TS7F1MQ4LOQVX")
var alnums = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var mRSeed int64 = 20221111
var mRSource mrand.Source
var mRand *mrand.Rand

type KeyRandomGenerateMethod string

var GMMath KeyRandomGenerateMethod = "math/rand"
var GMCrypto KeyRandomGenerateMethod = "crypto/rand"

type Secret struct {
	Key       uuid.UUID `json:"key"`
	Secret    []byte    `json:"secret"`
	CreateAt  time.Time `json:"created_at"`
	ValidAt   time.Time `json:"valid_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

type SecretDB interface {
	Get(key uuid.UUID) (*Secret, error)
	Create(srct *Secret) error
	Delete(key uuid.UUID) error
	Flush() error
}

type JWTSecretRepo struct {
	SecretLen     int
	LifeSpan      time.Duration
	ValidAfter    time.Duration
	CurrKey       uuid.UUID
	CurrSecret    []byte
	DB            SecretDB
	randGenMethod KeyRandomGenerateMethod
}

func SetRandomSeed(seed int64) {
	mRSeed = seed
	mRSource = mrand.NewSource(mRSeed)
}

func InitRandomGenerator() {
	mRand = mrand.New(mRSource)
}

func NewJWTSecretRepo(db SecretDB, opts ...option) (*JWTSecretRepo, error) {
	srctRepo := &JWTSecretRepo{DB: db}
	for _, opt := range opts {
		err := opt(srctRepo)
		if err != nil {
			return srctRepo, err
		}
	}

	if srctRepo.SecretLen <= 0 {
		srctRepo.SecretLen = 256
	}

	if srctRepo.LifeSpan == 0 {
		srctRepo.LifeSpan = 1 * time.Hour
	}

	if srctRepo.ValidAfter < 0 {
		srctRepo.ValidAfter = 0
	}

	_, _, err := srctRepo.UpdateCurrSecret()
	return srctRepo, err
}

type option func(sr *JWTSecretRepo) error

func WithSecretLength(sln int) option {
	return func(sr *JWTSecretRepo) error {
		if sln <= 0 {
			return fmt.Errorf("secret length should be greather than zero (get %d)", sln)
		}
		sr.SecretLen = sln
		return nil
	}
}

func WithSecretKeyLifeSpan(lifespan string) option {
	return func(sr *JWTSecretRepo) error {
		d, err := time.ParseDuration(lifespan)
		if err != nil {
			return err
		}

		if d <= 1*time.Second {
			return fmt.Errorf("life span of a secret key should be greather than 1 sec (get %v)", d)
		}
		sr.LifeSpan = d
		return nil
	}
}

func WithSecretKeyValidAfter(validAfter string) option {
	return func(sr *JWTSecretRepo) error {
		d, err := time.ParseDuration(validAfter)
		if err != nil {
			return err
		}
		sr.ValidAfter = d
		return nil
	}
}

func WithSecretRandomGenerateMethod(method KeyRandomGenerateMethod) option {
	return func(sr *JWTSecretRepo) error {
		sr.randGenMethod = method
		return nil
	}
}

func (sr *JWTSecretRepo) GetSecret(key uuid.UUID) ([]byte, bool) {
	if key == sr.CurrKey {
		return sr.CurrSecret, true
	}

	srct, err := sr.DB.Get(key)
	if err != nil {
		if err == ErrSecretExpired {
			sr.DB.Delete(key)
		}
		return nil, false
	}
	return srct.Secret, true
}

func genRdmIntWCRand(sLen int) ([]uint8, error) {
	randindex := make([]uint8, sLen)
	_, err := crand.Read(randindex)
	if err != nil {
		return nil, ErrCannotGenSecret
	}
	return randindex, nil
}

func genRdmIntWMRand(sLen int) ([]uint8, error) {
	randindex := make([]uint8, sLen)
	for i := 0; i < sLen; i++ {
		randindex[i] = uint8(mRand.Uint32())
	}
	return randindex, nil
}

func (sr *JWTSecretRepo) UpdateCurrSecret() (uuid.UUID, []byte, error) {
	var err error
	var randindex []uint8

	switch sr.randGenMethod {
	case GMMath:
		randindex, err = genRdmIntWMRand(sr.SecretLen)
	case GMCrypto:
		randindex, err = genRdmIntWCRand(sr.SecretLen)
	}

	// randindex, err := genRdmIntWCRand(sr.SecretLen)
	if err != nil {
		return uuid.Nil, nil, err
	}

	randbytes := make([]byte, sr.SecretLen)
	for i, r := range randindex {
		randbytes[i] = alnums[r%uint8(len(alnums))]
	}

	ct := time.Now()
	vt := ct.Add(sr.ValidAfter)
	et := vt.Add(sr.LifeSpan)
	srct := &Secret{
		Key:       uuid.New(),
		Secret:    randbytes,
		CreateAt:  ct,
		ValidAt:   vt,
		ExpiredAt: et,
	}
	err = sr.DB.Create(srct)
	if err != nil {
		return uuid.Nil, nil, err
	}
	sr.CurrKey = srct.Key
	sr.CurrSecret = srct.Secret

	return sr.CurrKey, sr.CurrSecret, nil
}
