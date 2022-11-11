package dao

import (
	"gitlab.com/gjerry134679/bank/internal/models"
)

type TransferDao struct {
	engine *models.Queries
}

// switch to transfer related doa methods
func (d *Dao) Transfer() *TransferDao {
	return &TransferDao{
		engine: d.engine,
	}
}

// func (d *TransferDao) Transfer(ctx context.Context, srcId, dstId int64, amount decimal.Decimal) error {
// 	ctrparam := models.CreateTransferRecordParams{
// 		SrcID:  srcId,
// 		DstID:  dstId,
// 		Amount: amount,
// 		Status: models.TstatusCreated,
// 	}
// 	tr, err := d.engine.CreateTransferRecord(ctx, ctrparam)
// 	if err != nil {
// 		utrparam := models.UpdateTransferStatusParams{
// 			ID:     tr.ID,
// 			Status: models.TstatusFailure,
// 		}
// 		d.engine.UpdateTransferStatus(ctx, utrparam)
// 		return err
// 	}

// 	dtrparam := models.NewDoTransferParamsParams(srcId, dstId, amount)
// 	err = d.engine.DoTransfer(ctx, dtrparam)
// 	if err != nil {
// 		utrparam := models.UpdateTransferStatusParams{
// 			ID:     tr.ID,
// 			Status: models.TstatusFailure,
// 		}
// 		d.engine.UpdateTransferStatus(ctx, utrparam)
// 		return err
// 	}

// 	utrparam := models.UpdateTransferStatusParams{
// 		ID:     tr.ID,
// 		Status: models.TstatusSuccess,
// 	}
// 	_, err = d.engine.UpdateTransferStatus(ctx, utrparam)
// 	return err
// }
