package dao

import (
	"context"

	"gitlab.com/gjerry134679/bank/internal/models"
)

type CurrencyDao struct {
	engine *models.Queries
}

// switch to currency related doa methods
func (d *Dao) Currency() *CurrencyDao {
	return &CurrencyDao{
		engine: d.engine,
	}
}

// create a currecy
func (d *CurrencyDao) Create(ctx context.Context, name, abbr string) {
	param := models.CreateCurrencyParams{
		Name: name,
		Abbr: abbr,
	}
	d.engine.CreateCurrency(ctx, param)
}

// create currencies
// func (d *CurrencyDao) BulkCreate(ctx context.Context, name, abbr []string) (int64, error) {
// 	if len(name) != len(abbr) {
// 		return -1, fmt.Errorf(
// 			"array len not match name(%d) abbr(%d)",
// 			len(name), len(abbr),
// 		)
// 	}
// 	ln := len(name)
// 	param := make([]models.BulkCreateCurrenciesParams, ln)
// 	for i := 0; i < ln; i++ {
// 		param[i] = models.BulkCreateCurrenciesParams{
// 			Name: name[i],
// 			Abbr: abbr[i],
// 		}
// 	}
// 	return d.engine.BulkCreateCurrencies(ctx, param)
// }

// update target currency data
func (d *CurrencyDao) Update(ctx context.Context, id int32, name, abbr string) (models.Currency, error) {
	param := models.UpdateCurrencyParams{
		ID:   id,
		Name: name,
		Abbr: abbr,
	}
	return d.engine.UpdateCurrency(ctx, param)
}

// get target currency
func (d *CurrencyDao) Get(ctx context.Context, id int32) (models.Currency, error) {
	return d.engine.GetCurrency(ctx, id)
}

// list currencies in db
func (d *CurrencyDao) List(ctx context.Context, limit, offset int32) ([]models.Currency, error) {
	param := models.ListCurrenciesParams{
		Limit:  limit,
		Offset: offset,
	}
	return d.engine.ListCurrencies(ctx, param)
}

// soft delete (update deleted_at column) target currency data
func (d *CurrencyDao) Delete(ctx context.Context, id int32) error {
	return d.engine.DeleteCurrency(ctx, id)
}

// hard delete target currency data
func (d *CurrencyDao) HardDelete(ctx context.Context, id int32) error {
	return d.engine.HardDeleteCurrency(ctx, id)
}
