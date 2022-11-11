package jwtIssuer

import (
	"crypto"
	"encoding/json"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"gitlab.com/gjerry134679/bank/pkg/jwtIssuer/secretRepo"
)

type Claims struct {
	KeyId   uuid.UUID
	HashSum string
	jwt.RegisteredClaims
}

type JWTIssuer struct {
	Method        jwt.SigningMethod
	HashSumMethod crypto.Hash
	ServiceName   string
	CheckClaims   jwt.RegisteredClaims
	Repo          secretRepo.JWTSecretRepo
}

func NewJWTIssuer(serviceName, issuer string, hashSumMethod crypto.Hash, repo secretRepo.JWTSecretRepo) *JWTIssuer {
	return &JWTIssuer{
		Method:        jwt.SigningMethodHS256,
		HashSumMethod: hashSumMethod,
		ServiceName:   serviceName,
		Repo:          repo,
		CheckClaims: jwt.RegisteredClaims{
			Issuer: issuer,
		},
	}
}

func (j *JWTIssuer) DeepCopyStdClaims() (jwt.RegisteredClaims, error) {
	var newRClaims jwt.RegisteredClaims
	jsn, err := json.Marshal(j.CheckClaims)
	if err != nil {
		return newRClaims, err
	}
	err = json.Unmarshal(jsn, &newRClaims)
	return newRClaims, err
}

func (j *JWTIssuer) GetJWTSecret(Key uuid.UUID) ([]byte, bool) {
	return j.Repo.GetSecret(Key)
}

func (j *JWTIssuer) RenewCurrKey() error {
	_, _, err := j.Repo.UpdateCurrSecret()
	return err
}

// func (j *JWTIssuer) GenerateSignature(ctx context.Context, appAuth *model.AppAuth) (uuid.UUID, string, error) {
// 	if appAuth.HashedKey == "" || appAuth.HashedSecret == "" {
// 		return uuid.Nil, "", errors.New("hashedkey/hashedsecret field should not be empty")
// 	}
// 	currTime := time.Now()
// 	expiredTime := currTime.Add(j.ExpireAfter)

// 	hasher := j.HashSumMethod.New()
// 	hasher.Write([]byte(appAuth.HashedKey))
// 	hasher.Write(salt)
// 	hasher.Write([]byte(appAuth.HashedSecret))
// 	claims := Claims{
// 		KeyId:   j.CurrJWTSecret.KeyId,
// 		HashSum: hex.EncodeToString(hasher.Sum(nil)),
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
