package dao

import (
	"context"

	"github.com/shopspring/decimal"
	"gitlab.com/gjerry134679/bank/internal/models"
)

type EntryDao struct {
	engine *models.Queries
}

// switch to entry related doa methods
func (d *Dao) Entry() *EntryDao {
	return &EntryDao{
		engine: d.engine,
	}
}

func (d *EntryDao) Create(ctx context.Context, id int64, amount decimal.Decimal) (models.Entry, error) {
	param := models.CreateEntryParams{
		AccountID: id,
		Amount:    amount,
	}
	return d.engine.CreateEntry(ctx, param)
}

// func (d *EntryDao) BulkCreate(ctx context.Context, accountId []int64, amount []decimal.Decimal) (int64, error) {
// 	if len(accountId) != len(amount) {
// 		return -1, fmt.Errorf(
// 			"array len not match id(%d) amount(%d)",
// 			len(accountId), len(amount),
// 		)
// 	}
// 	ln := len(accountId)
// 	param := make([]models.BulkCreateEntriesParams, ln)
// 	for i := 0; i < ln; i++ {
// 		param[i] = models.BulkCreateEntriesParams{
// 			AccountID: accountId[i],
// 			Amount:    amount[i],
// 		}
// 	}
// 	return d.engine.BulkCreateEntries(ctx, param)
// }

// get target entry
func (d *EntryDao) Get(ctx context.Context, id int64) (models.Entry, error) {
	return d.engine.GetEntry(ctx, id)
}

func (d *EntryDao) GetEntryByAccountId(ctx context.Context, accountId int64) ([]models.Entry, error) {
	return d.engine.GetEntryByAccountId(ctx, accountId)
}

func (d *EntryDao) List(ctx context.Context, limit, offset int32) ([]models.Entry, error) {
	param := models.ListEntriesParams{
		Limit:  limit,
		Offset: offset,
	}
	return d.engine.ListEntries(ctx, param)
}

func (d *EntryDao) HardDelete(ctx context.Context, id int64) error {
	return d.engine.HardDeleteEntry(ctx, id)
}
