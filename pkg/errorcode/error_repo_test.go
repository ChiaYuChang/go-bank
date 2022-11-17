package errorcode_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gjerry134679/bank/pkg/errorcode"
)

func TestNewRepoWith(t *testing.T) {
	repo, err := errorcode.NewErrorRepo(
		errorcode.WithSuccess(),
		errorcode.WithAuthErr(),
		errorcode.WithClientErr(),
		errorcode.WithServerErr(),
	)

	require.NoError(t, err)
	require.NotNil(t, repo)

	repo = errorcode.NewErrorRepoWithDefaultErrors()
	require.NotNil(t, repo)
}

func TestRegisterErr(t *testing.T) {
	repo, err := errorcode.NewErrorRepo()
	require.NoError(t, err)
	require.NotNil(t, repo)

	var EcTest errorcode.ErrorCode = 1
	var msg string = "error for testing"
	err = repo.RegisterErr(EcTest, 200, msg)
	require.NoError(t, err)

	err = repo.RegisterErr(EcTest, 400, msg)
	require.ErrorIs(t, err, errorcode.ErrErrorCodeHasBeenUsed)
}

func TestRegisterErrFromErr(t *testing.T) {
	repo, err := errorcode.NewErrorRepo()
	require.NoError(t, err)
	require.NotNil(t, repo)

	var EcTest errorcode.ErrorCode = 1
	testerr := errors.New("error for testing")
	err = repo.RegisterErrFromErr(testerr, EcTest, 200)
	require.NoError(t, err)

	err = repo.RegisterErrFromErr(testerr, EcTest, 400)
	require.ErrorIs(t, err, errorcode.ErrErrorCodeHasBeenUsed)
}

func TestGetErr(t *testing.T) {
	repo, err := errorcode.NewErrorRepo()
	require.NoError(t, err)
	require.NotNil(t, repo)

	var EcEmpty errorcode.ErrorCode = 0
	var EcTest errorcode.ErrorCode = 1
	var msg string = "error for testing"
	var status int = 200

	getErr, ok := repo.GetErr(EcEmpty)
	require.False(t, ok)
	require.Nil(t, getErr)

	testErr := &errorcode.Error{
		Code:           EcTest,
		HttpStatusCode: status,
		Message:        msg,
		Details:        make([]string, 0),
	}
	repo[EcTest] = testErr

	getErr, ok = repo.GetErr(EcTest)
	require.True(t, ok)
	require.Equal(t, testErr.Code, getErr.Code)
	require.Equal(t, testErr.HttpStatusCode, getErr.HttpStatusCode)
	require.Equal(t, testErr.Message, getErr.Message)
	require.Equal(t, testErr.String(), getErr.String())

	// should be clone error
	getErr.Message = "updated message"
	require.Equal(t, testErr.Code, getErr.Code)
	require.Equal(t, testErr.HttpStatusCode, getErr.HttpStatusCode)
	require.NotEqual(t, testErr.Message, getErr.Message)

	getErr = repo.MustGetErr(EcTest)
	require.Equal(t, testErr.Code, getErr.Code)
	require.Equal(t, testErr.HttpStatusCode, getErr.HttpStatusCode)
	require.Equal(t, testErr.Message, getErr.Message)
	require.Equal(t, testErr.String(), getErr.String())
}
