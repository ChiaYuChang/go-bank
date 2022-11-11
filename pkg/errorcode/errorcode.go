package errorcode

import (
	"errors"
	"fmt"
	"net/http"
)

func init() {
	ErrorRepo = make(errorRepo)
	ErrorRepo.WithSuccess()
	ErrorRepo.WithAuthErr()
	ErrorRepo.WithClientErr()
	ErrorRepo.WithServerErr()
}

type ErrorCode int32

func (ec ErrorCode) String() string {
	return fmt.Sprintf("%04d", ec)
}

func (er errorRepo) WithSuccess() {
	er.RegisterErr(Success, http.StatusOK, "success")
}

func (er errorRepo) WithClientErr() {
	er.RegisterErr(EcInvalidParams, http.StatusBadRequest, "invalid parameters")
	er.RegisterErr(EcNotFound, http.StatusNotFound, "the server cannot find the requested resource")
	er.RegisterErr(EcRequestTimeout, http.StatusRequestTimeout, "request timeout")
	er.RegisterErr(EcPayloadTooLarge, http.StatusRequestEntityTooLarge, "payload too large")
	er.RegisterErr(EcTooManyRequests, http.StatusTooManyRequests, "too may requests")
}

func (er errorRepo) WithAuthErr() {
	er.RegisterErr(EcUnauthorizedAuthNotExist, http.StatusUnauthorized, "account/key not found")
	er.RegisterErr(EcJWTInvalidSignature, http.StatusUnauthorized, "jwt signature is invalid")
	er.RegisterErr(EcJWTMalformed, http.StatusUnauthorized, "malformed jwt signature")
	er.RegisterErr(EcJWTSignatureExpired, http.StatusUnauthorized, "jwt signature is expired")
	er.RegisterErr(EcJWTErrAudience, http.StatusUnauthorized, "error audience (AUD)")
	er.RegisterErr(EcJWTErrIssueAt, http.StatusUnauthorized, "error issue at (IAT)")
	er.RegisterErr(EcJWTErrIssuer, http.StatusUnauthorized, "error issuer (ISS)")
	er.RegisterErr(EcJWTErrNotValidYet, http.StatusUnauthorized, "not yet valid (NBF)")
	er.RegisterErr(EcJWTErrId, http.StatusUnauthorized, "error jwt key id")
	er.RegisterErr(EcJWTErrAlg, http.StatusUnauthorized, "jwt algorithm not match")
	er.RegisterErr(EcJWTInvalidClaims, http.StatusUnauthorized, "invalid jwt claims")
}

func (er errorRepo) WithServerErr() {
	er.RegisterErr(EcInternalSevereError, http.StatusInternalServerError, "internal server error")
	er.RegisterErr(EcNotImplemented, http.StatusNotImplemented, "method not implemented")
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
