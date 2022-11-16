package models_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"strings"
	"testing"
	"time"

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
	t.Logf("Id %d Src Owner %s Balance %d", srcAcc.ID, srcAcc.Owner, srcAcc.Balance.BigInt().Int64())
	t.Logf("Id %d Dst Owner %s Balance %d", dstAcc.ID, dstAcc.Owner, dstAcc.Balance.BigInt().Int64())

	amountInt := 3
	amount := decimal.NewFromInt(int64(amountInt))

	n := runtime.NumCPU() * 2
	t.Logf("cpu num  %d", n)

	results := make(chan testTransferTxResuts)
	for i := 0; i < n; i++ {
		go func(index int) {
			s := rand.Intn(10)
			time.Sleep(time.Duration(s) * time.Second)

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

	existed := make(map[decimal.Decimal]bool)
	for i := 0; i < n; i++ {
		r := <-results
		// check transfer
		require.NoError(t, r.err)
		require.NotEmpty(t, r.result)
		require.Equal(t, srcAcc.ID, r.result.Transfer.SrcID)
		require.Equal(t, dstAcc.ID, r.result.Transfer.DstID)
		if !amount.Equal(r.result.Transfer.Amount) {
			t.Fatalf("transfer.amount not match: want: %d, get: %d", amount, r.result.Transfer.Amount)
		}
		require.NotEqual(t, r.result.Transfer.ID, 0)
		require.NotZero(t, r.result.Transfer.CreatedAt)
		require.Equal(t, r.result.Transfer.Status, models.TstatusSuccess)

		// check source entry
		require.NotEmpty(t, r.result.ScrEntry)
		require.Equal(t, r.result.ScrEntry.AccountID, srcAcc.ID)
		if !r.result.ScrEntry.Amount.Equal(amount.Neg()) {
			t.Fatalf("src_entry.balance not match: want: %d, get: %d", amount, r.result.ScrEntry.Amount)
		}
		require.NotZero(t, r.result.ScrEntry.ID)
		require.NotZero(t, r.result.ScrEntry.CreatedAt)

		// check destination entry
		require.NotEmpty(t, r.result.DstEntry)
		require.Equal(t, r.result.DstEntry.AccountID, dstAcc.ID)
		if !r.result.DstEntry.Amount.Equal(amount) {
			t.Fatalf("dst_entry.balance not match: want: %d, get: %d", amount, r.result.DstEntry.Amount)
		}
		require.NotZero(t, r.result.DstEntry.ID)
		require.NotZero(t, r.result.DstEntry.CreatedAt)

		// check source account
		require.NotEmpty(t, r.result.SrcAccount)
		require.Equal(t, r.result.SrcAccount.ID, srcAcc.ID)
		require.Equal(t, r.result.SrcAccount.Owner, srcAcc.Owner)
		require.Equal(t, r.result.SrcAccount.Currency, srcAcc.Currency)
		require.Equal(t, r.result.SrcAccount.CreatedAt, srcAcc.CreatedAt)

		// check destination account
		require.NotEmpty(t, r.result.DstAccount)
		require.Equal(t, r.result.DstAccount.ID, dstAcc.ID)
		require.Equal(t, r.result.DstAccount.Owner, dstAcc.Owner)
		require.Equal(t, r.result.DstAccount.Currency, dstAcc.Currency)
		require.Equal(t, r.result.DstAccount.CreatedAt, dstAcc.CreatedAt)

		diff1 := srcAcc.Balance.Sub(r.result.SrcAccount.Balance)
		diff2 := dstAcc.Balance.Sub(r.result.DstAccount.Balance)
		if !diff1.Equal(diff2.Neg()) {
			t.Fatalf(
				"difference between accounts not match: src: %d, dst: %d",
				diff1.BigInt().Int64(),
				diff2.BigInt().Int64(),
			)
		}

		k := diff1.Div(amount)
		if k.LessThan(decimal.NewFromInt(1)) {
			kf, _ := k.Float64()
			t.Fatalf("k should be within [1, n], get: %f", kf)
		}
		if k.GreaterThan(decimal.NewFromInt(int64(n))) {
			kf, _ := k.Float64()
			t.Fatalf("k should be within [1, n], get: %f", kf)
		}
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	uSrcAcc, err := testQureies.GetAccount(context.Background(), srcAcc.ID)
	t.Logf("updated src_acc.balance: %d", uSrcAcc.Balance.BigInt().Int64())
	require.NoError(t, err)

	uDstAcc, err := testQureies.GetAccount(context.Background(), dstAcc.ID)
	t.Logf("updated dst_acc.balance: %d", uDstAcc.Balance.BigInt().Int64())
	require.NoError(t, err)

	t.Logf(
		"Src: Balance %d -> %d (diff want: %3d get %3d)",
		srcAcc.Balance.BigInt().Int64(),
		uSrcAcc.Balance.BigInt().Int64(),
		-1*n*amountInt,
		uSrcAcc.Balance.Sub(srcAcc.Balance).BigInt().Int64(),
	)

	t.Logf(
		"Dst: Balance %d -> %d (diff want: %3d get %3d)",
		dstAcc.Balance.BigInt().Int64(),
		uDstAcc.Balance.BigInt().Int64(),
		n*amountInt,
		uDstAcc.Balance.Sub(dstAcc.Balance).BigInt().Int64(),
	)

	if !srcAcc.Balance.Sub(amount.Mul(decimal.NewFromInt(int64(n)))).Equal(uSrcAcc.Balance) {
		log.Fatalf("src_acc balance error")
	}

	if !dstAcc.Balance.Add(amount.Mul(decimal.NewFromInt(int64(n)))).Equal(uDstAcc.Balance) {
		log.Fatalf("dst_acc balance error")
	}
}

func TestTransferTxDeadlock(t *testing.T) {
	store := models.NewStore(testDB)

	cr := mustCreateAnTestCurrency(t, 1)[0]
	accs := mustCreateAnTestAccount(t, 2, cr.ID)
	acc1 := accs[0]
	acc2 := accs[1]

	t.Logf("Id %d, Name %s", cr.ID, cr.Name)
	t.Logf("Id %d Src Owner %s Balance %d", acc1.ID, acc1.Owner, acc1.Balance.BigInt().Int64())
	t.Logf("Id %d Dst Owner %s Balance %d", acc2.ID, acc2.Owner, acc2.Balance.BigInt().Int64())

	amountInt := 3
	amount := decimal.NewFromInt(int64(amountInt))

	n := runtime.NumCPU() * 2
	t.Logf("cpu num  %d", n)

	results := make(chan testTransferTxResuts)
	for i := 0; i < n; i++ {
		var src, dst models.Account
		if i%2 == 0 {
			src = acc1
			dst = acc2
		} else {
			src = acc2
			dst = acc1
		}

		go func(index int, src, dst models.Account) {
			_, err := store.DoTransferTx(
				context.Background(), models.TransferTxParams{
					SrcId:  src.ID,
					DstId:  dst.ID,
					Amount: amount,
				})
			results <- testTransferTxResuts{
				index:  index,
				err:    err,
				result: nil}
		}(i, src, dst)
	}

	for i := 0; i < n; i++ {
		r := <-results
		if r.err != nil {
			if r.err == models.ErrSelfLoop {
				continue
			}
			t.Fatalf("error: %v", r.err)
		}
	}

	uAcc1, err := testQureies.GetAccount(context.Background(), acc1.ID)
	t.Logf(
		"updated src_acc.balance: %d (expected: %d)",
		uAcc1.Balance.BigInt().Int64(),
		acc1.Balance.BigInt().Int64(),
	)
	require.NoError(t, err)

	uAcc2, err := testQureies.GetAccount(context.Background(), acc2.ID)
	t.Logf("updated dst_acc.balance: %d", uAcc2.Balance.BigInt().Int64())
	require.NoError(t, err)

	if !acc1.Balance.Equal(uAcc1.Balance) {
		log.Fatalf("src_acc balance error")
	}

	if !acc2.Balance.Equal(uAcc2.Balance) {
		log.Fatalf("dst_acc balance error")
	}
}

func TestTransferTxSelfLoop(t *testing.T) {
	store := models.NewStore(testDB)

	cr := mustCreateAnTestCurrency(t, 1)[0]
	acc := mustCreateAnTestAccount(t, 1, cr.ID)[0]

	amount := decimal.NewFromInt(int64(3))

	n := 5
	t.Logf("cpu num  %d", n)
	results := make(chan testTransferTxResuts)
	for i := 0; i < n; i++ {
		go func(index int, src, dst models.Account) {
			_, err := store.DoTransferTx(
				context.Background(), models.TransferTxParams{
					SrcId:  src.ID,
					DstId:  dst.ID,
					Amount: amount,
				})
			results <- testTransferTxResuts{
				index:  index,
				err:    err,
				result: nil}
		}(i, acc, acc)
	}

	for i := 0; i < n; i++ {
		r := <-results
		if r.err != nil {
			if r.err == models.ErrSelfLoop {
				continue
			}
			t.Fatalf("error: %v", r.err)
		}
	}

	uAcc, err := testQureies.GetAccount(context.Background(), acc.ID)
	require.NoError(t, err)

	if !acc.Balance.Equal(uAcc.Balance) {
		log.Fatalf("account balance error")
	}

}
