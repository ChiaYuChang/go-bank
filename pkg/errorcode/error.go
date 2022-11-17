package errorcode

import (
	"fmt"
	"strings"
)

type Error struct {
	Code           ErrorCode `json:"code"`
	HttpStatusCode int       `json:"http_status_code"`
	Message        string    `json:"message"`
	Details        []string  `json:"detals"`
}

// New error struct
func NewError(code ErrorCode, httpStatusCode int, message string) *Error {
	return &Error{
		Code:           code,
		HttpStatusCode: httpStatusCode,
		Message:        message,
		Details:        make([]string, 0),
	}
}

func FromError(err error, code ErrorCode, httpStatusCode int) *Error {
	return NewError(code, httpStatusCode, err.Error())
}

// Append details to the Error struct
func (e *Error) WithDetails(details ...string) *Error {
	e.Details = append(e.Details, details...)
	return e
}

// Append messages to the Error struct
func (e *Error) WithMessage(msg string) *Error {
	e.Message = msg
	return e
}

// Return formated message
func (e *Error) Msgf(a ...any) string {
	return fmt.Sprintf(e.Message, a)
}

// Return message
func (e *Error) Msg() string {
	return e.Message
}

// Error interface
func (e *Error) Error() string {
	return fmt.Sprintf(
		"code: %04d (%03d) msg: %s details: [%s]",
		e.Code, e.HttpStatusCode, e.Message, strings.Join(e.Details, ", "),
	)
}

// String interface
func (e *Error) String() string {
	s := fmt.Sprintf(
		"code: %04d (%03d) msg: %s\ndetails:\n",
		e.Code, e.HttpStatusCode, e.Message,
	)

	for i, d := range e.Details {
		s += fmt.Sprintf(" - %d. %s\n", i, d)
	}
	return s
}

func (e *Error) Clone() *Error {
	newE := NewError(e.Code, e.HttpStatusCode, e.Message)
	copy(newE.Details, e.Details)
	return newE
}
