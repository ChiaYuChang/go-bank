package jwtIssuer

import (
	"crypto"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gitlab.com/gjerry134679/bank/pkg/jwtIssuer/secretRepo"
)

type option func(ji *JWTIssuer) error

var ErrFieldHasAlreadyAsigned = errors.New("the field has already been asigned")

func IsHasAssignedErr(err error) bool { return err == ErrFieldHasAlreadyAsigned }

func WithServiceName(name string) option {
	return func(ji *JWTIssuer) error {
		ji.ServiceName = name
		return nil
	}
}

func WithJWTSignMethod(method jwt.SigningMethod) option {
	return func(ji *JWTIssuer) error {
		ji.method = method
		return nil
	}
}

func WithSecretKeyRenewTime(t time.Duration) option {
	return func(ji *JWTIssuer) error {
		ji.renewfreq = t
		ji.ticker = time.NewTicker(t)
		return nil
	}
}

func WithSecretRepository(repo *secretRepo.JWTSecretRepo) option {
	return func(ji *JWTIssuer) error {
		ji.repo = repo
		return nil
	}
}

func WithClaimsIssuer(issuer string) option {
	return func(ji *JWTIssuer) error {
		ji.Add(fiss, issuer)
		return nil
	}
}

func WithClaimSubject(subject string) option {
	return func(ji *JWTIssuer) error {
		ji.Add(fsub, subject)
		return nil
	}
}

func WithClaimsAudience(audience []string) option {
	return func(ji *JWTIssuer) error {
		for _, aud := range audience {
			ji.Add(faud, aud)
		}
		return nil
	}
}

func WithClaimsId(id string) option {
	return func(ji *JWTIssuer) error {
		ji.Add(fjti, id)
		return nil
	}
}

func WithInfoHashMethod(method crypto.Hash) option {
	return func(ji *JWTIssuer) error {
		ji.hashMethod = method
		return nil
	}
}

func WithJWTValidAfter(validAfter string) option {
	return func(ji *JWTIssuer) error {
		d, err := time.ParseDuration(validAfter)
		if err != nil {
			return err
		}

		if d < 0 {
			return fmt.Errorf("valid after should be greater or equal zero (get %v)", d)
		}
		ji.validAfter = &d
		return nil
	}
}

func WithJWTExpiredAfter(expiredAfter string) option {
	return func(ji *JWTIssuer) error {
		d, err := time.ParseDuration(expiredAfter)
		if err != nil {
			return err
		}

		if d < 0 {
			return fmt.Errorf("expired after should be greater or equal zero (get %v)", d)
		}
		ji.validAfter = &d
		return nil
	}
}

func WithAuthChecker(authChecker func(account, password string) bool) option {
	return func(ji *JWTIssuer) error {
		ji.authChecker = authChecker
		return nil
	}
}
