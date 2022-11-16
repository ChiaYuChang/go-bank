package jwtIssuer

import (
	"context"
	"crypto"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"gitlab.com/gjerry134679/bank/pkg/jwtIssuer/secretRepo"
	"gitlab.com/gjerry134679/bank/pkg/jwtIssuer/secretRepo/inMemory"
)

type CtxKey interface {
	Key() ctxkey
	String() string
}

type ctxkey struct{ key string }

type Auth struct {
	Account  string
	Password string
}

func (a Auth) Key() ctxkey {
	return ctxkey{key: a.String()}
}

func (a Auth) String() string {
	return "CtxKeyAuth"
}

type Claims struct {
	KeyId   uuid.UUID
	HashSum string
	jwt.RegisteredClaims
}

type JWTIssuer struct {
	method       jwt.SigningMethod
	hashMethod   crypto.Hash
	ServiceName  string
	validAfter   *time.Duration
	expiredAfter *time.Duration
	renewfreq    time.Duration
	done         chan bool
	ticker       *time.Ticker
	repo         *secretRepo.JWTSecretRepo
	authChecker  func(account, password string) bool
	checker
}

func NewJWTIssuer(opts ...option) (*JWTIssuer, error) {
	ji := &JWTIssuer{done: make(chan bool, 1)}
	for _, opt := range opts {
		err := opt(ji)
		if err != nil {
			return ji, err
		}
	}

	if ji.method == nil {
		ji.method = jwt.SigningMethodHS256
	}

	if ji.ServiceName == "" {
		ji.ServiceName = "my-jwt"
	}

	if ji.ticker == nil {
		ji.renewfreq = 1 * time.Hour
		ji.ticker = time.NewTicker(ji.renewfreq)
	}

	if ji.repo == nil {
		ji.repo, _ = secretRepo.NewJWTSecretRepo(inMemory.NewInMemoryDB())
	}

	if ji.authChecker == nil {
		ji.authChecker = func(account, password string) bool { return true }
	}

	if len(ji.checker.issuers) == 0 {
		ji.Add(fiss, "issuer")
	}

	if len(ji.subjects) == 0 {
		ji.Add(fsub, "gerneral")
	}

	if ji.expiredAfter == nil {
		t := 3 * time.Hour
		ji.expiredAfter = &t
	}

	if ji.validAfter == nil {
		t := 0 * time.Second
		ji.validAfter = &t
	}

	go func() {
		for {
			select {
			case <-ji.done:
				return
			case <-ji.ticker.C:
				ji.RenewCurrKey()
			}
		}
	}()

	return ji, nil
}

func (ji *JWTIssuer) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Service name     : %v\n", ji.ServiceName))
	sb.WriteString(fmt.Sprintf("JWT Algorithm    : %v\n", ji.method.Alg()))
	sb.WriteString(fmt.Sprintf("Info Hash Method : %v\n", ji.hashMethod.String()))
	sb.WriteString(fmt.Sprintf("Renew  Frequency : %v\n", ji.renewfreq))
	sb.WriteString("JWT:\n")
	sb.WriteString("  - Issuers (iss):\n")
	for k := range ji.issuers {
		sb.WriteString(fmt.Sprintf("    - %s:\n", k))
	}
	sb.WriteString("  - Audiences (aud):\n")
	for k := range ji.audiences {
		sb.WriteString(fmt.Sprintf("    - %s:\n", k))
	}
	sb.WriteString("  - Subjects (sub):\n")
	for k := range ji.subjects {
		sb.WriteString(fmt.Sprintf("    - %s:\n", k))
	}
	sb.WriteString("  - JWT Id (jti):\n")
	for k := range ji.jwtids {
		sb.WriteString(fmt.Sprintf("    - %s:\n", k))
	}
	sb.WriteString(fmt.Sprintf("  - Valid After   : %v\n", ji.validAfter.String()))
	sb.WriteString(fmt.Sprintf("  - Expired After : %v\n", ji.validAfter.String()))
	return sb.String()
}

func (ji *JWTIssuer) GetJWTSecret(Key uuid.UUID) ([]byte, bool) {
	return ji.repo.GetSecret(Key)
}

func (ji *JWTIssuer) RenewCurrKey() error {
	_, _, err := ji.repo.UpdateCurrSecret()
	return err
}

func (ji *JWTIssuer) ResetTicker(d time.Duration) {
	ji.renewfreq = d
	ji.ticker.Reset(d)
}

func (ji *JWTIssuer) GracefulStop() {
	ji.ticker.Stop()
	ji.done <- true
}

func (ji *JWTIssuer) GenerateSignature(ctx context.Context) (uuid.UUID, string, error) {
	var err error
	var secretKey uuid.UUID = uuid.Nil
	var signature string = ""
	var auth Auth

	auth = ctx.Value(auth.Key()).(Auth)
	if !ji.authChecker(auth.Account, auth.Password) {
		return secretKey, signature, err
	}

	// rclms := jwt.RegisteredClaims{
	// 	Issuer: ji,
	// }
	// currTime := time.Now()
	// expiredTime := currTime.Add(ji.repo.ValidAfter).Add(ji.repo.ValidAfter)
	return secretKey, signature, err
}

// 	claims := Claims{
// 		KeyId: ji.repo.CurrKey,
// 	}
// 	rclms, err := j.DeepCopyStdClaims()
// 	if err != nil {
// 		err = fmt.Errorf("error while copying registered claims: %v", err)
// 		span.SetStatus(codes.Error, "JWT-DeepCopyStdClaims error")
// 		span.RecordError(err)
// 		return uuid.Nil, "", err
// 	}
// 	claims.RegisteredClaims = rclms
// 	claims.IssuedAt = jwt.NewNumericDate(currTime)
// 	claims.ExpiresAt = jwt.NewNumericDate(expiredTime)

// 	token := jwt.NewWithClaims(j.Method, claims)
// 	signature, err := token.SignedString(
// 		[]byte(Must(j.GetJWTSecret(j.CurrJWTSecret.KeyId)).Secret),
// 	)
// 	if err != nil {
// 		err = fmt.Errorf("error while creating signed token: %v", err)
// 		span.SetStatus(codes.Error, "JWT-NewWithClaims error")
// 		span.RecordError(err)
// 		return uuid.Nil, "", err
// 	}
// 	return j.CurrJWTSecret.KeyId, signature, nil
// }

// func (j *JWTIssuer) ParseSignature(ctx context.Context, keyId uuid.UUID, signature string) (*Claims, error) {
// 	_, span := otel.Tracer(j.ServiceName).Start(ctx, "JWT-ParseSignature")
// 	defer span.End()

// 	vertfytoken, err := jwt.ParseWithClaims(signature, &Claims{},
// 		func(t *jwt.Token) (interface{}, error) {
// 			var err error
// 			var jwtSecret *model.JWTSecret
// 			// vertfy the algorithm in header
// 			if t.Method.Alg() != j.Method.Alg() {
// 				return nil, ErrAlg
// 			}

// 			if keyId == j.CurrJWTSecret.KeyId {
// 				return []byte(j.CurrJWTSecret.Secret), nil
// 			}

// 			jwtSecret, err = j.GetJWTSecret(keyId)
// 			if err != nil {
// 				return nil, ErrKeyId
// 			}

// 			return []byte(jwtSecret.Secret), nil
// 		})
// 	if err != nil {
// 		span.SetStatus(codes.Error, "jwt.ParseWithClaims error")
// 		span.RecordError(err)
// 		return nil, err
// 	}

// 	if !vertfytoken.Valid {
// 		span.SetStatus(codes.Error, "invalid JWT-Token")
// 		span.RecordError(SignatureInvalid)
// 		return nil, SignatureInvalid
// 	}

// 	claims, ok := vertfytoken.Claims.(*Claims)
// 	if !ok {
// 		span.SetStatus(codes.Error, "claims could not be asserted")
// 		span.RecordError(ClaimsInvalid)
// 		return nil, ClaimsInvalid
// 	}

// 	if claims.Issuer != j.CheckClaims.Issuer {
// 		span.SetStatus(codes.Error, "invalid issuer")
// 		span.RecordError(ErrIssuer)
// 		return nil, ErrIssuer
// 	}

// 	return claims, nil
// }

// func (j *JWTIssuer) UpdateCurrJWTSecret(s *model.JWTSecret, err error) (*model.JWTSecret, error) {
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Println(*s)
// 	j.CurrJWTSecret = s
// 	_, err = j.DBEngine.JWTSecret().Create(*s)
// 	return s, err
// }

// func (j *JWTIssuer) ParseGinContext(ctx *gin.Context) (uuid.UUID, string, error) {
// 	var keyId uuid.UUID
// 	var valid bool = true
// 	var err error
// 	errs := app.ValidErrors{}

// 	keyIdStr, ok := ctx.GetQuery("key_id")
// 	if !ok {
// 		keyIdStr = ctx.GetHeader("X-Key-Id")
// 	}
// 	if keyIdStr != "" {
// 		keyId, err = uuid.Parse(keyIdStr)
// 		if err != nil {
// 			errs = append(errs, &app.ValidError{Key: "key_id", Message: "failed on the 'uuidv4' tag"})
// 			keyId = uuid.Nil
// 			valid = false
// 		}
// 	} else {
// 		errs = append(errs, &app.ValidError{Key: "key_id", Message: "failed on the 'required' tag"})
// 		keyId = uuid.Nil
// 		valid = false
// 	}

// 	signature, ok := ctx.GetQuery("signature")
// 	if !ok {
// 		signature = ctx.GetHeader("X-Signature")
// 	}
// 	if signature == "" {
// 		errs = append(errs, &app.ValidError{Key: "signature", Message: "failed on the 'required' tag"})
// 		valid = false
// 	}

// 	if valid {
// 		return keyId, signature, nil
// 	}

// 	return keyId, signature, errs
// }
