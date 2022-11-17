package errorcode_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gjerry134679/bank/pkg/errorcode"
)

func TestNewError(t *testing.T) {
	EcTest := errorcode.ErrorCode(1)
	status := 200
	message := "test message"

	err := errorcode.NewError(EcTest, status, message)
	require.Equal(t, EcTest, err.Code)
	require.Equal(t, status, err.HttpStatusCode)
	require.Equal(t, message, err.Message)
	require.NotNil(t, err.Details)
}

func TestFromError(t *testing.T) {
	EcTest := errorcode.ErrorCode(1)
	status := 200
	fromerr := errors.New("test message")

	err := errorcode.FromError(fromerr, EcTest, status)

	require.Equal(t, EcTest, err.Code)
	require.Equal(t, status, err.HttpStatusCode)
	require.Equal(t, fromerr.Error(), err.Message)
	require.NotNil(t, err.Details)
}

func TestWithMessage(t *testing.T) {
	EcTest := errorcode.ErrorCode(1)
	status := 200
	message := "test message"

	err := errorcode.NewError(EcTest, status, "")
	require.Zero(t, err.Message)

	err.WithMessage(message)
	require.Equal(t, message, err.Message)
}

func TestWithDeatils(t *testing.T) {
	EcTest := errorcode.ErrorCode(1)
	status := 200
	message := "test message"
	details := []string{"one", "two", "three"}

	err := errorcode.NewError(EcTest, status, message)
	require.NotNil(t, err.Details)

	err.WithDetails(details...)
	for i := 0; i < 3; i++ {
		require.Equal(t, err.Details[i], details[i])
	}
}

func TestMsgAndMsgf(t *testing.T) {
	EcTest := errorcode.ErrorCode(1)
	status := 200
	message := "test message: want [1] get %d"

	err := errorcode.NewError(EcTest, status, message)

	require.Equal(t, message, err.Msg())
	require.Equal(t, fmt.Sprintf(message, []int{2}), err.Msgf(2))
}

func TestClone(t *testing.T) {
	EcTest := errorcode.ErrorCode(1)
	status := 200
	message := "test message"
	details := []string{"one", "two", "three"}

	err := errorcode.NewError(EcTest, status, message)
	err.WithDetails(details...)

	require.Equal(t, EcTest, err.Code)
	require.Equal(t, status, err.HttpStatusCode)
	require.Equal(t, message, err.Message)
	require.NotNil(t, err.Details)

	for i := 0; i < 3; i++ {
		require.Equal(t, err.Details[i], details[i])
	}

	err.Details[1] = "too"
	err.WithMessage("updated message")
	require.NotEqual(t, details[1], err.Details[1])
	require.NotEqual(t, message, err.Message)
}
