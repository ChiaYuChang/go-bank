package models

import "github.com/shopspring/decimal"

func NewDoTransferParams(srcId, dstId int64, amount decimal.Decimal) DoTransferParams {
	return DoTransferParams{
		SrcID:  srcId,
		DstID:  dstId,
		Amount: amount,
	}
}
