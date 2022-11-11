package secretRepo

import (
	crand "crypto/rand"
	"errors"
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
	randGenMethod string
}

func SetRandomSeed(seed int64) {
	mRSeed = seed
	mRSource = mrand.NewSource(mRSeed)
}

func InitRandomGenerator() {
	mRand = mrand.New(mRSource)
}

func NewJWTSecretRepo(srctLen int, lifeSpan time.Duration, validAfter time.Duration, db SecretDB) (*JWTSecretRepo, error) {
	srctRepo := &JWTSecretRepo{
		SecretLen:     srctLen,
		LifeSpan:      lifeSpan,
		ValidAfter:    validAfter,
		DB:            db,
		randGenMethod: "math",
	}
	_, _, err := srctRepo.UpdateCurrSecret()
	return srctRepo, err
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
	if sr.randGenMethod == "math" {
		randindex, err = genRdmIntWMRand(sr.SecretLen)
	}

	if sr.randGenMethod == "crypto" {
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
