package asset

import (
	"math/big"

	"github.com/xssnick/tonutils-go/tvm/cell"
)

// Data represents the structure of asset data
type Data struct {
	SRate               *big.Int
	BRate               *big.Int
	TotalSupply         *big.Int
	TotalBorrow         *big.Int
	LastAccrual         *big.Int
	Balance             *big.Int
	TrackingSupplyIndex *big.Int
	TrackingBorrowIndex *big.Int
	AwaitedSupply       *big.Int
}

func newData(assetData *cell.Slice) *Data {
	return &Data{
		SRate:               assetData.MustLoadBigUInt(64),
		BRate:               assetData.MustLoadBigUInt(64),
		TotalSupply:         assetData.MustLoadBigUInt(64),
		TotalBorrow:         assetData.MustLoadBigUInt(64),
		LastAccrual:         assetData.MustLoadBigUInt(32),
		Balance:             assetData.MustLoadBigUInt(64),
		TrackingSupplyIndex: assetData.MustLoadBigUInt(64),
		TrackingBorrowIndex: assetData.MustLoadBigUInt(64),
		AwaitedSupply:       assetData.MustLoadBigUInt(64),
	}
}
