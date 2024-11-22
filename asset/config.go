package asset

import (
	"github.com/xssnick/tonutils-go/tvm/cell"
	"math/big"
)

// Config represents the configuration for an asset
type Config struct {
	Oracle                   *big.Int
	Decimals                 *big.Int
	CollateralFactor         *big.Int
	LiquidationThreshold     *big.Int
	LiquidationBonus         *big.Int
	BaseBorrowRate           *big.Int
	BorrowRateSlopeLow       *big.Int
	BorrowRateSlopeHigh      *big.Int
	SupplyRateSlopeLow       *big.Int
	SupplyRateSlopeHigh      *big.Int
	TargetUtilization        *big.Int
	OriginationFee           *big.Int
	Dust                     *big.Int
	MaxTotalSupply           *big.Int
	ReserveFactor            *big.Int
	LiquidationReserveFactor *big.Int
	MinPrincipalForRewards   *big.Int
	BaseTrackingSupplySpeed  *big.Int
	BaseTrackingBorrowSpeed  *big.Int
}

func newConfig(assetConfig *cell.Slice, assetConfigValueRef *cell.Slice) *Config {
	return &Config{
		Oracle:                   assetConfig.MustLoadBigUInt(256),
		Decimals:                 assetConfig.MustLoadBigUInt(8),
		CollateralFactor:         assetConfigValueRef.MustLoadBigUInt(16),
		LiquidationThreshold:     assetConfigValueRef.MustLoadBigUInt(16),
		LiquidationBonus:         assetConfigValueRef.MustLoadBigUInt(16),
		BaseBorrowRate:           assetConfigValueRef.MustLoadBigUInt(64),
		BorrowRateSlopeLow:       assetConfigValueRef.MustLoadBigUInt(64),
		BorrowRateSlopeHigh:      assetConfigValueRef.MustLoadBigUInt(64),
		SupplyRateSlopeLow:       assetConfigValueRef.MustLoadBigUInt(64),
		SupplyRateSlopeHigh:      assetConfigValueRef.MustLoadBigUInt(64),
		TargetUtilization:        assetConfigValueRef.MustLoadBigUInt(64),
		OriginationFee:           assetConfigValueRef.MustLoadBigUInt(64),
		Dust:                     assetConfigValueRef.MustLoadBigUInt(64),
		MaxTotalSupply:           assetConfigValueRef.MustLoadBigUInt(64),
		ReserveFactor:            assetConfigValueRef.MustLoadBigUInt(16),
		LiquidationReserveFactor: assetConfigValueRef.MustLoadBigUInt(16),
		MinPrincipalForRewards:   assetConfigValueRef.MustLoadBigUInt(64),
		BaseTrackingSupplySpeed:  assetConfigValueRef.MustLoadBigUInt(64),
		BaseTrackingBorrowSpeed:  assetConfigValueRef.MustLoadBigUInt(64),
	}
}

func (c *Config) Scale() *big.Int {
	return new(big.Int).Exp(big.NewInt(10), c.Decimals, nil)
}
