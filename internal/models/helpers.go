package models

import "github.com/shopspring/decimal"

func NewDoTransferParams(srcId, dstId int64, amount decimal.Decimal) DoTransferParams {
	return DoTransferParams{
		Column1: srcId,
		Column2: amount,
		Column3: dstId,
		Column4: amount,
	}
}
