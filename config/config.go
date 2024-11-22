package config

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

type MasterParams struct {
	FactorScale                        *big.Int
	AssetCoefficientScale              *big.Int
	AssetPriceScale                    *big.Int
	AssetReserveFactorScale            *big.Int
	AssetLiquidationReserveFactorScale *big.Int
	AssetOriginationFeeScale           *big.Int
	AssetLiquidationThresholdScale     *big.Int
	AssetLiquidationBonusScale         *big.Int
	AssetSRateScale                    *big.Int
	AssetBRateScale                    *big.Int
	CollateralWorthThreshold           *big.Int
}

type AssetConfig struct {
	Name                Asset
	ID                  *big.Int
	Decimals            int
	JettonMasterAddress *address.Address
	JettonWalletCode    *cell.Cell
}

func (c *AssetConfig) GetJettonWalletAddress(walletAddress *address.Address) (*address.Address, error) {
	if walletAddress == nil {
		return nil, errors.New("user wallet address is nil-pointer")
	}
	if c.JettonMasterAddress == nil {
		return nil, errors.New("asset master address is nil-pointer")
	}
	if c.JettonWalletCode == nil {
		return nil, errors.New("asset wallet code is nil-pointer")
	}

	var data *cell.Cell
	switch c.Name {
	case TON:
		return nil, errors.New("asset is TON")
	case USDT:
		data = cell.BeginCell().
			MustStoreUInt(0, 4).
			MustStoreCoins(0).
			MustStoreAddr(walletAddress).
			MustStoreAddr(c.JettonMasterAddress).
			EndCell()
	case TSTON:
		data = cell.BeginCell().
			MustStoreCoins(0).
			MustStoreAddr(walletAddress).
			MustStoreAddr(c.JettonMasterAddress).
			MustStoreRef(c.JettonWalletCode).
			MustStoreCoins(0).
			MustStoreUInt(0, 48).
			EndCell()
	default:
		data = cell.BeginCell().
			MustStoreCoins(0).
			MustStoreAddr(walletAddress).
			MustStoreAddr(c.JettonMasterAddress).
			MustStoreRef(c.JettonWalletCode).
			EndCell()
	}
	walletStateInitCell, err := tlb.ToCell(&tlb.StateInit{Code: c.JettonWalletCode, Data: data})
	if err != nil {
		return nil, fmt.Errorf("failed to convert StateInit to Cell, err: %w", err)
	}
	return address.NewAddress(0, 0, walletStateInitCell.Hash()), nil
}

type OracleNFT struct {
	ID      uint64
	Address string
}

type Config struct {
	MasterAddress  *address.Address
	MasterVersion  int64
	MasterParams   *MasterParams
	Oracles        []*OracleNFT
	MinimalOracles int
	Assets         map[string]*AssetConfig
	LendingCode    *cell.Cell
}

func GetMainMainnetConfig() *Config {
	return &Config{
		MasterAddress: address.MustParseAddr(MasterMainnet),
		MasterVersion: MainnetVersion,
		MasterParams:  GetMasterParams(),
		Oracles: []*OracleNFT{
			{ID: 0, Address: "0xd3a8c0b9fd44fd25a49289c631e3ac45689281f2f8cf0744400b4c65bed38e5d"},
			{ID: 1, Address: "0x2c21cabdaa89739de16bde7bc44e86401fac334a3c7e55305fe5e7563043e191"},
			{ID: 2, Address: "0x2eb258ce7b5d02466ab8a178ad8b0ba6ffa7b58ef21de3dc3b6dd359a1e16af0"},
			{ID: 3, Address: "0xf9a0769954b4430bca95149fb3d876deb7799d8f74852e0ad4ccc5778ce68b52"},
		},
		MinimalOracles: 3,
		Assets: map[string]*AssetConfig{
			TON.Sha256Hash().String(): {
				Name:                TON,
				ID:                  TON.Sha256Hash(),
				Decimals:            9,
				JettonMasterAddress: nil,
			},
			USDT.ID(): {
				Name:                USDT,
				ID:                  USDT.Sha256Hash(),
				Decimals:            6,
				JettonMasterAddress: address.MustParseAddr(USDTJettonAddress),
				JettonWalletCode:    getCellFromHex(CodeWalletUSDT),
			},
			JUSDT.ID(): {
				Name:                JUSDT,
				ID:                  JUSDT.Sha256Hash(),
				Decimals:            6,
				JettonMasterAddress: address.MustParseAddr(JUSDTJettonAddress),
				JettonWalletCode:    getCellFromHex(CodeWalletStandard),
			},
			JUSDC.ID(): {
				Name:                JUSDC,
				ID:                  JUSDC.Sha256Hash(),
				Decimals:            6,
				JettonMasterAddress: address.MustParseAddr(JUSDCJettonAddress),
				JettonWalletCode:    getCellFromHex(CodeWalletStandard),
			},
			STTON.ID(): {
				Name:                STTON,
				ID:                  STTON.Sha256Hash(),
				Decimals:            9,
				JettonMasterAddress: address.MustParseAddr(STTONJettonAddress),
				JettonWalletCode:    getCellFromHex(CodeWalletSTTON),
			},
			TSTON.ID(): {
				Name:                TSTON,
				ID:                  TSTON.Sha256Hash(),
				Decimals:            9,
				JettonMasterAddress: address.MustParseAddr(TSTONJettonAddress),
				JettonWalletCode:    getCellFromHex(CodeWalletTSTON),
			},
		},
		LendingCode: getCellFromHex(CodeLending),
	}
}

func GetMasterTestnetConfig() *Config {
	return &Config{
		MasterAddress: address.MustParseAddr(MasterTestnet),
		MasterVersion: TestnetVersion,
		MasterParams:  GetMasterParams(),
		Oracles: []*OracleNFT{
			{ID: 0, Address: "0xd3a8c0b9fd44fd25a49289c631e3ac45689281f2f8cf0744400b4c65bed38e5d"},
			{ID: 1, Address: "0x2c21cabdaa89739de16bde7bc44e86401fac334a3c7e55305fe5e7563043e191"},
			{ID: 2, Address: "0x2eb258ce7b5d02466ab8a178ad8b0ba6ffa7b58ef21de3dc3b6dd359a1e16af0"},
			{ID: 3, Address: "0xf9a0769954b4430bca95149fb3d876deb7799d8f74852e0ad4ccc5778ce68b52"},
		},
		MinimalOracles: 3,
		Assets: map[string]*AssetConfig{
			TON.ID(): {
				Name:                TON,
				ID:                  TON.Sha256Hash(),
				Decimals:            9,
				JettonMasterAddress: nil,
			},
			JUSDT.ID(): {
				Name:                JUSDT,
				ID:                  JUSDT.Sha256Hash(),
				Decimals:            6,
				JettonMasterAddress: address.MustParseAddr(JUSDTJettonAddressTestnet),
				JettonWalletCode:    getCellFromHex(CodeWalletStandardTestnet),
			},
			JUSDC.ID(): {
				Name:                JUSDC,
				ID:                  JUSDC.Sha256Hash(),
				Decimals:            6,
				JettonMasterAddress: address.MustParseAddr(JUSDCJettonAddressTestnet),
				JettonWalletCode:    getCellFromHex(CodeWalletStandardTestnet),
			},
			STTON.ID(): {
				Name:                STTON,
				ID:                  STTON.Sha256Hash(),
				Decimals:            9,
				JettonMasterAddress: address.MustParseAddr(STTONJettonAddressTestnet),
				JettonWalletCode:    getCellFromHex(CodeWalletSTTONTestnet),
			},
		},
		LendingCode: getCellFromHex(CodeLending),
	}
}

func GetLpMainnetConfig() *Config {
	return &Config{
		MasterAddress: address.MustParseAddr(LpMainnet),
		MasterVersion: LpVersion,
		MasterParams:  GetMasterParams(),
		Oracles: []*OracleNFT{
			{ID: 0, Address: "0xd3a8c0b9fd44fd25a49289c631e3ac45689281f2f8cf0744400b4c65bed38e5d"},
			{ID: 1, Address: "0x2c21cabdaa89739de16bde7bc44e86401fac334a3c7e55305fe5e7563043e191"},
			{ID: 2, Address: "0x2eb258ce7b5d02466ab8a178ad8b0ba6ffa7b58ef21de3dc3b6dd359a1e16af0"},
			{ID: 3, Address: "0xf9a0769954b4430bca95149fb3d876deb7799d8f74852e0ad4ccc5778ce68b52"},
		},
		MinimalOracles: 3,
		Assets: map[string]*AssetConfig{
			TON.ID(): {
				Name:                TON,
				ID:                  TON.Sha256Hash(),
				Decimals:            9,
				JettonMasterAddress: nil,
			},
			USDT.ID(): {
				Name:                USDT,
				ID:                  USDT.Sha256Hash(),
				Decimals:            6,
				JettonMasterAddress: address.MustParseAddr(USDTJettonAddress),
				JettonWalletCode:    getCellFromHex(CodeWalletUSDT),
			},
			TON_STORM.ID(): {
				Name:                TON_STORM,
				ID:                  TON_STORM.Sha256Hash(),
				Decimals:            9,
				JettonMasterAddress: address.MustParseAddr(TONStormJettonAddress),
				JettonWalletCode:    getCellFromHex(CodeWalletTonStorm),
			},
			USDT_STORM.ID(): {
				Name:                USDT_STORM,
				ID:                  USDT_STORM.Sha256Hash(),
				Decimals:            9,
				JettonMasterAddress: address.MustParseAddr(USDTStormJettonAddress),
				JettonWalletCode:    getCellFromHex(CodeWalletUsdtStorm),
			},
			TONUSDT_DEDUST.ID(): {
				Name:                TONUSDT_DEDUST,
				ID:                  TONUSDT_DEDUST.Sha256Hash(),
				Decimals:            9,
				JettonMasterAddress: address.MustParseAddr(TONUSDTDeDustJettonAddress),
				JettonWalletCode:    getCellFromHex(CodeWalletTonUsdtDeDust),
			},
		},
		LendingCode: getCellFromHex(CodeLending),
	}
}

func GetMasterParams() *MasterParams {
	return &MasterParams{
		FactorScale:                        big.NewInt(FactorScale),
		AssetCoefficientScale:              big.NewInt(AssetCoefficientScale),
		AssetPriceScale:                    big.NewInt(AssetPriceScale),
		AssetReserveFactorScale:            big.NewInt(AssetReserveFactorScale),
		AssetLiquidationReserveFactorScale: big.NewInt(AssetLiquidationReserveFactorScale),
		AssetOriginationFeeScale:           big.NewInt(AssetOriginationFeeScale),
		AssetLiquidationThresholdScale:     big.NewInt(AssetLiquidationThresholdScale),
		AssetLiquidationBonusScale:         big.NewInt(AssetLiquidationBonusScale),
		AssetSRateScale:                    big.NewInt(AssetSRateScale),
		AssetBRateScale:                    big.NewInt(AssetBRateScale),
		CollateralWorthThreshold:           big.NewInt(CollateralWorthThreshold),
	}
}
