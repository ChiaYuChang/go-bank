package models

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/shopspring/decimal"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

type QueryCallBackFun func(*Queries) error

func (s *Store) exectTx(ctx context.Context, fn QueryCallBackFun) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		// rollback if queries failed
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	// commit transaction
	return tx.Commit()
}

type TransferTxParams struct {
	Id     int64           `json:"id"`
	SrcId  int64           `json:"src_id"`
	DstId  int64           `json:"dst_id"`
	Amount decimal.Decimal `json:"amount"`
}

type TransferTxResult struct {
	Transfer   Transfer `json:"transfer"`
	SrcAccount Account  `json:"src_account"`
	DstAccount Account  `json:"dst_account"`
	ScrEntry   Entry    `json:"src_entry"`
	DstEntry   Entry    `json:"dst_entry"`
}

func (s *Store) DoTransferTx(ctx context.Context, params TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	var err error
	result.Transfer, err = s.Queries.CreateTransferRecord(ctx, CreateTransferRecordParams{
		SrcID:  params.SrcId,
		DstID:  params.DstId,
		Amount: params.Amount,
		Status: TstatusFailure,
	})
	if err != nil {
		return result, fmt.Errorf("error while creating transfer record: %v", err)
	}
	params.Id = result.Transfer.ID
	return s.TransferTx(ctx, params)
}

func (s *Store) TransferTx(ctx context.Context, params TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.exectTx(ctx, func(q *Queries) error {
		var err error

		// Create Source Account Entry
		result.ScrEntry, err = q.CreateEntry(
			ctx, CreateEntryParams{
				AccountID: params.SrcId,
				Amount:    params.Amount.Neg(),
			})
		if err != nil {
			return fmt.Errorf("error while creating entry: %v", err)
		}

		// Create Destination Account Entry
		result.DstEntry, err = q.CreateEntry(
			ctx, CreateEntryParams{
				AccountID: params.DstId,
				Amount:    params.Amount,
			})
		if err != nil {
			return fmt.Errorf("error while creating entry: %v", err)
		}

		result.SrcAccount, err = q.AccountBalanceWithdraw(
			ctx, AccountBalanceWithdrawParams{
				ID:      params.SrcId,
				Balance: params.Amount,
			})
		if err != nil {
			return fmt.Errorf("error while withdrawing from source account: %v", err)
		}

		result.DstAccount, err = q.AccountBalanceDeposit(
			ctx, AccountBalanceDepositParams{
				ID:      params.DstId,
				Balance: params.Amount,
			})
		if err != nil {
			return fmt.Errorf("error while depositing to destination account: %v", err)
		}
		// err = q.DoTransfer(ctx, NewDoTransferParams(
		// 	params.SrcId,
		// 	params.DstId,
		// 	params.Amount))
		// if err != nil {
		// 	return err
		// }

		// Create transfer record
		result.Transfer, err = q.UpdateTransferStatus(
			ctx, UpdateTransferStatusParams{
				params.Id, TstatusSuccess,
			})
		if err != nil {
			return fmt.Errorf("error while updating status: %v", err)
		}

		return nil
	})

	return result, err
}
