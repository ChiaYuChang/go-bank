package dao

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/gjerry134679/bank/internal/models"
)

type AccountDao struct {
	engine *models.Queries
}

// switch to account related doa methods
func (d *Dao) Account() *AccountDao {
	return &AccountDao{
		engine: d.engine,
	}
}

// create a account
func (d *AccountDao) Create(ctx context.Context, owner string, balance decimal.Decimal, currency int32) (models.Account, error) {
	param := models.CreateAccountParams{
		Owner:    owner,
		Balance:  balance,
		Currency: currency,
	}
	return d.engine.CreateAccount(ctx, param)
}

// create accounts
// func (d *AccountDao) BulkCreate(ctx context.Context, owners []string, balances []decimal.Decimal, curremcies []int32) (int64, error) {
// 	if len(owners) != len(balances) || len(owners) != len(curremcies) {
// 		return -1, fmt.Errorf(
// 			"array len not match owners(%d) balances(%d) currencies(%d)",
// 			len(owners), len(balances), len(curremcies),
// 		)
// 	}
// 	ln := len(owners)
// 	param := make([]models.BulkCreateAccountsParams, ln)
// 	for i := 0; i < ln; i++ {
// 		param[i] = models.BulkCreateAccountsParams{
// 			Owner:    owners[i],
// 			Balance:  balances[i],
// 			Currency: curremcies[i],
// 		}
// 	}
// 	return d.engine.BulkCreateAccounts(ctx, param)
// }

// update account
func (d *AccountDao) Update(ctx context.Context, id int64, owner string, balance decimal.Decimal, currency int32) (models.Account, error) {
	param := models.UpdateAccountParams{
		ID:        id,
		Owner:     owner,
		Balance:   balance,
		Currency:  currency,
		UpdatedAt: time.Now(),
	}
	return d.engine.UpdateAccount(ctx, param)
}

// get target account
func (d *AccountDao) Get(ctx context.Context, id int64) (models.Account, error) {
	return d.engine.GetAccount(ctx, id)
}

// assign account balance value
func (d *AccountDao) UpdateBalance(ctx context.Context, id int64, balance decimal.Decimal) (models.Account, error) {
	param := models.UpdateAccountBalanceParams{
		ID:        id,
		Balance:   balance,
		UpdatedAt: time.Now(),
	}
	return d.engine.UpdateAccountBalance(ctx, param)
}

// withdraw from balance
func (d *AccountDao) BalanceWithdraw(ctx context.Context, id int64, amount decimal.Decimal) (models.Account, error) {
	param := models.AccountBalanceWithdrawParams{
		ID:      id,
		Balance: amount,
	}
	return d.engine.AccountBalanceWithdraw(ctx, param)
}

// deposit to balance
func (d *AccountDao) BalanceDeposit(ctx context.Context, id int64, amount decimal.Decimal) (models.Account, error) {
	param := models.AccountBalanceDepositParams{
		ID:      id,
		Balance: amount,
	}
	return d.engine.AccountBalanceDeposit(ctx, param)
}

// list accounts in db
func (d *AccountDao) List(ctx context.Context, limit, offset int32) ([]models.Account, error) {
	param := models.ListAccountsParams{
		Limit:  limit,
		Offset: offset,
	}
	return d.engine.ListAccounts(ctx, param)
}

// soft delete (update deleted_at column) target account in db
func (d *AccountDao) Delete(ctx context.Context, id int64) error {
	return d.engine.DeleteAccount(ctx, id)
}

// hard delete target account in db
func (d *AccountDao) HardDelete(ctx context.Context, id int64) error {
	return d.engine.DeleteAccount(ctx, id)
}
