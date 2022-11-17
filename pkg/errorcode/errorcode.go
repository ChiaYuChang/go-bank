package errorcode

import (
	"errors"
	"fmt"
	"net/http"
)

func init() {
	ErrorRepo = NewErrorRepoWithDefaultErrors()
}

type ErrorCode int32

func (ec ErrorCode) String() string {
	return fmt.Sprintf("%04d", ec)
}

type option func(repo errorRepo) error

func WithSuccess() option {
	return func(repo errorRepo) error {
		err := repo.RegisterErr(Success, http.StatusOK, "success")
		return err
	}
}

func WithClientErr() option {
	return func(repo errorRepo) error {
		for _, e := range []struct {
			code   ErrorCode
			status int
			msg    string
		}{
			{EcInvalidParams, http.StatusBadRequest, "invalid parameters"},
			{EcNotFound, http.StatusNotFound, "the server cannot find the requested resource"},
			{EcRequestTimeout, http.StatusRequestTimeout, "request timeout"},
			{EcPayloadTooLarge, http.StatusRequestEntityTooLarge, "payload too large"},
			{EcTooManyRequests, http.StatusTooManyRequests, "too may requests"},
		} {
			err := repo.RegisterErr(e.code, e.status, e.msg)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func WithAuthErr() option {
	return func(repo errorRepo) error {
		for _, e := range []struct {
			code   ErrorCode
			status int
			msg    string
		}{
			{EcUnauthorizedAuthNotExist, http.StatusUnauthorized, "account/key not found"},
			{EcJWTInvalidSignature, http.StatusUnauthorized, "jwt signature is invalid"},
			{EcJWTMalformed, http.StatusUnauthorized, "malformed jwt signature"},
			{EcJWTSignatureExpired, http.StatusUnauthorized, "jwt signature is expired"},
			{EcJWTErrAudience, http.StatusUnauthorized, "error audience (AUD)"},
			{EcJWTErrIssueAt, http.StatusUnauthorized, "error issue at (IAT)"},
			{EcJWTErrIssuer, http.StatusUnauthorized, "error issuer (ISS)"},
			{EcJWTErrNotValidYet, http.StatusUnauthorized, "not yet valid (NBF)"},
			{EcJWTErrId, http.StatusUnauthorized, "error jwt key id"},
			{EcJWTErrAlg, http.StatusUnauthorized, "jwt algorithm not match"},
			{EcJWTInvalidClaims, http.StatusUnauthorized, "invalid jwt claims"},
		} {
			err := repo.RegisterErr(e.code, e.status, e.msg)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func WithServerErr() option {
	return func(repo errorRepo) error {
		for _, e := range []struct {
			code   ErrorCode
			status int
			msg    string
		}{
			{EcInternalSevereError, http.StatusInternalServerError, "internal server error"},
			{EcNotImplemented, http.StatusNotImplemented, "method not implemented"},
		} {
			err := repo.RegisterErr(e.code, e.status, e.msg)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// constants
const StatusUnknown int = 0

// Query success
const Success ErrorCode = 0

// Client error
const (
	EcInvalidParams ErrorCode = 1000 + iota
	EcNotFound
	EcRequestTimeout
	EcPayloadTooLarge
	EcTooManyRequests
)

// Auth error
const (
	// General auth errors
	EcUnauthorizedAuthNotExist ErrorCode = 2000 + iota
	// JWT RelatedError
	EcJWTKeyIdNotFound
	EcJWTInvalidSignature
	EcJWTMalformed
	EcJWTSignatureExpired
	EcJWTErrAudience
	EcJWTErrIssueAt
	EcJWTErrIssuer
	EcJWTErrNotValidYet
	EcJWTErrId
	EcJWTErrAlg
	EcJWTInvalidClaims
)

// Server error
const (
	EcInternalSevereError = 3000 + iota
	EcNotImplemented
)

// API query error

// Global Variables
var ErrErrorCodeHasBeenUsed = errors.New("the given errorcode has been used")
var ErrorRepo errorRepo
