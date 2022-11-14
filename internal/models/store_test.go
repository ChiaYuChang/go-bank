package models_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"gitlab.com/gjerry134679/bank/internal/models"
)

type testerrors struct {
	message string
	details []string
}

func NewTestErrors(msg string) *testerrors {
	return &testerrors{message: msg, details: make([]string, 0)}
}

func (tes *testerrors) Error() string {
	return strings.Join(tes.details, ", ")
}

func (tes *testerrors) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("msg: %s", tes.message))

	if len(tes.details) > 0 {
		sb.WriteString("\n   details:\n")
	}

	for _, d := range tes.details {
		sb.WriteString(fmt.Sprintf("   - %v\n", d))
	}
	return sb.String()
}

func (tes *testerrors) Err(err error) *testerrors {
	tes.message = err.Error()
	return tes
}

func (tes *testerrors) WithDetails(details ...string) *testerrors {
	tes.details = append(tes.details, details...)
	return tes
}

var ErrTestCancel = errors.New("early terminated since cancel")
var ErrRdm = errors.New("random error")

type testTransferTxResuts struct {
	index  int
	err    error
	result *models.TransferTxResult
}

func rdmError() error {
	if rand.Intn(100) < 30 { // 10% of failure
		return ErrRdm
	}
	return nil
}

func mustCreateNewTestTransferRecord(t *testing.T, n int, src, dst, amount int64, status models.Tstatus) []models.Transfer {
	trnfrs := make([]models.Transfer, n)
	for i := 0; i < n; i++ {
		trnfr, err := testQureies.CreateTransferRecord(
			context.Background(),
			models.CreateTransferRecordParams{
				SrcID:  src,
				DstID:  dst,
				Amount: decimal.NewFromInt(amount),
				Status: status,
			})
		require.Nil(t, err)
		trnfrs[i] = trnfr
	}
	return trnfrs
}

func TestCreateTransferRecord(t *testing.T) {
	cr := mustCreateAnTestCurrency(t, 1)[0]
	accs := mustCreateAnTestAccount(t, 2, cr.ID)
	srcAcc := accs[0]
	dstAcc := accs[1]
	amount := decimal.NewFromInt(100)

	trnfr, err := testQureies.CreateTransferRecord(
		context.Background(),
		models.CreateTransferRecordParams{
			SrcID:  srcAcc.ID,
			DstID:  dstAcc.ID,
			Amount: amount,
			Status: models.TstatusFailure,
		})

	require.Nil(t, err)
	require.NotNil(t, trnfr)
	require.Equal(t, srcAcc.ID, trnfr.SrcID)
	require.Equal(t, dstAcc.ID, trnfr.DstID)
	require.Equal(t, models.TstatusFailure, trnfr.Status)
	if !amount.Equal(trnfr.Amount) {
		t.Fatalf("amount not match: want: %d, get: %d", amount, trnfr.Amount)
	}

	if trnfr.ID == 0 {
		t.Fatalf("transfer is should be greater than zero")
	}

	if trnfr.CreatedAt.IsZero() {
		t.Fatalf("created_at should not be zero (%v)", trnfr.CreatedAt)
	}
}

type transferUpdateResult struct {
	err      error
	index    int
	transfer models.Transfer
}

func TestUpdateUpdateTransferStatus(t *testing.T) {
	cr := mustCreateAnTestCurrency(t, 1)[0]
	accs := mustCreateAnTestAccount(t, 2, cr.ID)
	srcAcc := accs[0]
	dstAcc := accs[1]

	n := 10
	trnfrs := mustCreateNewTestTransferRecord(t, n, srcAcc.ID, dstAcc.ID, 100, models.TstatusFailure)

	results := make(chan transferUpdateResult)
	for i, t := range trnfrs {
		go func(index int, tid int64) {
			trnfr, err := testQureies.UpdateTransferStatus(
				context.Background(),
				models.UpdateTransferStatusParams{
					ID:     tid,
					Status: models.TstatusSuccess,
				})
			results <- transferUpdateResult{
				index:    index,
				err:      err,
				transfer: trnfr,
			}
		}(i, t.ID)
	}

	for i := 0; i < n; i++ {
		r := <-results
		require.Nil(t, r.err)
		require.NotNil(t, r.transfer)
		require.Equal(t, r.transfer.ID, trnfrs[r.index].ID)
		require.Equal(t, r.transfer.SrcID, trnfrs[r.index].SrcID)
		require.Equal(t, r.transfer.DstID, trnfrs[r.index].DstID)
		if !r.transfer.Amount.Equal(trnfrs[r.index].Amount) {
			t.Fatalf("amount not match: want: %d, get: %d",
				trnfrs[r.index].Amount, r.transfer.Amount)
		}
		require.Equal(t, r.transfer.CreatedAt, trnfrs[r.index].CreatedAt)
		require.Equal(t, r.transfer.Status, models.TstatusSuccess)
	}
}

func TestTransferTx(t *testing.T) {
	store := models.NewStore(testDB)

	cr := mustCreateAnTestCurrency(t, 1)[0]
	accs := mustCreateAnTestAccount(t, 2, cr.ID)
	srcAcc := accs[0]
	dstAcc := accs[1]

	t.Logf("Id %d, Name %s", cr.ID, cr.Name)
	t.Logf("Id %d, Src Owner %s", srcAcc.ID, srcAcc.Owner)
	t.Logf("Id %d, Dst Owner %s", dstAcc.ID, dstAcc.Owner)

	amount := decimal.NewFromInt(100)

	n := runtime.NumCPU() * 2
	results := make(chan testTransferTxResuts)
	for i := 0; i < n; i++ {
		go func(index int) {
			result, err := store.DoTransferTx(
				context.Background(), models.TransferTxParams{
					SrcId:  srcAcc.ID,
					DstId:  dstAcc.ID,
					Amount: amount,
				})
			results <- testTransferTxResuts{
				index:  index,
				err:    err,
				result: &result}
		}(i)
	}

	errors := make([]*testerrors, 0)
	for i := 0; i < n; i++ {
		r := <-results
		if r.err != nil {
			err := NewTestErrors(r.err.Error())
			errors = append(errors, err)
			continue
		}

		if r.result == nil {
			err := NewTestErrors("the result is empty")
			errors = append(errors, err)
			continue
		}

		dtl := make([]string, 0)
		if srcAcc.ID != r.result.SrcAccount.ID {
			dtl = append(
				dtl, fmt.Sprintf(
					"sourse account not match: want: %d, get: %d",
					srcAcc.ID, r.result.SrcAccount.ID),
			)
		}

		if dstAcc.ID != r.result.DstAccount.ID {
			dtl = append(
				dtl, fmt.Sprintf(
					"destination account not match: want: %d, get: %d",
					dstAcc.ID, r.result.DstAccount.ID),
			)
		}

		if !amount.Equal(r.result.Transfer.Amount) {
			dtl = append(
				dtl, fmt.Sprintf(
					"amount not match: want: %d, get: %d",
					amount, r.result.Transfer.Amount),
			)
		}

		if r.result.Transfer.ID == 0 {
			dtl = append(dtl, "transfer is should be greater than zero")
		}

		if r.result.Transfer.CreatedAt.IsZero() {
			dtl = append(dtl, "created_at should not be zero")
		}

		if r.result.Transfer.Status != models.TstatusSuccess {
			dtl = append(dtl, "status not match")
		}

		if len(dtl) > 0 {
			err := NewTestErrors("field value not match").WithDetails(dtl...)
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		for _, e := range errors {
			t.Log(e.String())
		}
		t.Fatal()
	}
}
