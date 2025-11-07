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
	case USDT, DOGS, NOT, USDE:
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
	case PT_tsUSDe_01Sep2025, PT_tsUSDe_18Dec2025:
		data = cell.BeginCell().
			MustStoreCoins(0).
			MustStoreAddr(walletAddress).
			MustStoreAddr(c.JettonMasterAddress).
			EndCell()
	case TSUSDE:
		data = cell.BeginCell().
			MustStoreUInt(0, 4).
			MustStoreCoins(0).
			MustStoreAddr(walletAddress).
			MustStoreAddr(c.JettonMasterAddress).
			MustStoreCoins(0).
			MustStoreUInt(0, 64).
			EndCell()
	default:
		data = cell.BeginCell().
			MustStoreCoins(0).
			MustStoreAddr(walletAddress).
			MustStoreAddr(c.JettonMasterAddress).
			MustStoreRef(c.JettonWalletCode).
			EndCell()
	}
	// Special handling for PT_tsUSDe_18Dec2025 - uses close address calculation
	if c.Name.ID() == PT_tsUSDe_18Dec2025.ID() {
		stateInit := cell.BeginCell().
			MustStoreUInt(1, 1).
			MustStoreUInt(8, 5).
			MustStoreUInt(0, 1).
			MustStoreMaybeRef(c.JettonWalletCode).
			MustStoreMaybeRef(data).
			MustStoreUInt(0, 1).
			EndCell()

		walletAddr := address.NewAddress(0, 0, stateInit.Hash())
		closeAddr, err := calculateCloseAddress(walletAddr, c.JettonMasterAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate close address: %w", err)
		}
		return closeAddr, nil
	}

	walletStateInitCell, err := tlb.ToCell(&tlb.StateInit{Code: c.JettonWalletCode, Data: data})
	if err != nil {
		return nil, fmt.Errorf("failed to convert StateInit to Cell, err: %w", err)
	}

	walletAddr := address.NewAddress(0, 0, walletStateInitCell.Hash())
	return walletAddr, nil
}

const (
	jettonAddressBits = 256
	closeRewriteBits  = 8
)

// calculateCloseAddress calculates a "close address" by combining bits from source and target addresses.
// It takes the lower (jettonAddressBits - closeRewriteBits) bits from source and the upper closeRewriteBits bits from target.
// This is used for PT_tsUSDe_18Dec2025 to create a special wallet address format.
func calculateCloseAddress(source, target *address.Address) (*address.Address, error) {
	if source == nil || target == nil {
		return nil, errors.New("address is nil-pointer")
	}

	sourceData := source.Data()
	targetData := target.Data()

	const hashBytes = jettonAddressBits / 8
	if len(sourceData) != hashBytes || len(targetData) != hashBytes {
		return nil, fmt.Errorf("unexpected address data length")
	}

	curHash := new(big.Int).SetBytes(sourceData)
	desiredHash := new(big.Int).SetBytes(targetData)

	mask := maskForBits(jettonAddressBits - closeRewriteBits)
	maxValue := maskForBits(jettonAddressBits)
	antimask := new(big.Int).Xor(mask, maxValue)

	newHash := new(big.Int).Or(
		new(big.Int).And(curHash, mask),
		new(big.Int).And(desiredHash, antimask),
	)

	hashBytesData, err := bigIntToFixedBytes(newHash, hashBytes)
	if err != nil {
		return nil, err
	}

	return address.NewAddress(0, 0, hashBytesData), nil
}

// maskForBits creates a bitmask with the specified number of bits set to 1.
// For example, maskForBits(8) returns 0xFF (255).
func maskForBits(bits int) *big.Int {
	if bits <= 0 {
		return big.NewInt(0)
	}

	return new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), uint(bits)), big.NewInt(1))
}

// bigIntToFixedBytes converts a big.Int to a byte slice of fixed size.
// The value is right-aligned (padded with zeros on the left if necessary).
func bigIntToFixedBytes(value *big.Int, size int) ([]byte, error) {
	if value == nil {
		return nil, errors.New("big.Int value is nil")
	}
	if value.BitLen() > size*8 {
		return nil, fmt.Errorf("value exceeds %d bits", size*8)
	}

	buf := make([]byte, size)
	value.FillBytes(buf)
	return buf, nil
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
			USDE.ID(): {
				Name:                USDE,
				ID:                  USDE.Sha256Hash(),
				Decimals:            6,
				JettonMasterAddress: address.MustParseAddr(USDEJettonAddress),
				JettonWalletCode:    getCellFromHex(CodeWalletUSDe),
			},
			TSUSDE.ID(): {
				Name:                TSUSDE,
				ID:                  TSUSDE.Sha256Hash(),
				Decimals:            6,
				JettonMasterAddress: address.MustParseAddr(TSUSDEJettonAddress),
				JettonWalletCode:    getCellFromHex(CodeWalletTsUsde),
			},
		},
		LendingCode: getCellFromHex(CodeLending),
	}
}

func GetAltsMainnetConfig() *Config {
	return &Config{
		MasterAddress: address.MustParseAddr(AltsMainnet),
		MasterVersion: AltsVersion,
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
			CATI.ID(): {
				Name:                CATI,
				ID:                  CATI.Sha256Hash(),
				Decimals:            9,
				JettonMasterAddress: address.MustParseAddr(CATIJettonAddress),
				JettonWalletCode:    getCellFromHex(CodeWalletCati),
			},
			DOGS.ID(): {
				Name:                DOGS,
				ID:                  DOGS.Sha256Hash(),
				Decimals:            9,
				JettonMasterAddress: address.MustParseAddr(DOGSJettonAddress),
				JettonWalletCode:    getCellFromHex(CodeWalletDogs),
			},
			NOT.ID(): {
				Name:                NOT,
				ID:                  NOT.Sha256Hash(),
				Decimals:            9,
				JettonMasterAddress: address.MustParseAddr(NOTJettonAddress),
				JettonWalletCode:    getCellFromHex(CodeWalletNot),
			},
			STON.ID(): {
				Name:                STON,
				ID:                  STON.Sha256Hash(),
				Decimals:            9,
				JettonMasterAddress: address.MustParseAddr(STONJettonAddress),
				JettonWalletCode:    getCellFromHex(CodeWalletSTON),
			},
		},
		LendingCode: getCellFromHex(CodeLending),
	}
}

func GetStableMainnetConfig() *Config {
	return &Config{
		MasterAddress: address.MustParseAddr(StableMainnet),
		MasterVersion: StableVersion,
		MasterParams:  GetMasterParams(),
		Oracles: []*OracleNFT{
			{ID: 0, Address: "0xd3a8c0b9fd44fd25a49289c631e3ac45689281f2f8cf0744400b4c65bed38e5d"},
			{ID: 1, Address: "0x2c21cabdaa89739de16bde7bc44e86401fac334a3c7e55305fe5e7563043e191"},
			{ID: 2, Address: "0x2eb258ce7b5d02466ab8a178ad8b0ba6ffa7b58ef21de3dc3b6dd359a1e16af0"},
			{ID: 3, Address: "0xf9a0769954b4430bca95149fb3d876deb7799d8f74852e0ad4ccc5778ce68b52"},
		},
		MinimalOracles: 3,
		Assets: map[string]*AssetConfig{
			USDT.ID(): {
				Name:                USDT,
				ID:                  USDT.Sha256Hash(),
				Decimals:            6,
				JettonMasterAddress: address.MustParseAddr(USDTJettonAddress),
				JettonWalletCode:    getCellFromHex(CodeWalletUSDT),
			},
			USDE.ID(): {
				Name:                USDE,
				ID:                  USDE.Sha256Hash(),
				Decimals:            6,
				JettonMasterAddress: address.MustParseAddr(USDEJettonAddress),
				JettonWalletCode:    getCellFromHex(CodeWalletUSDe),
			},
			TSUSDE.ID(): {
				Name:                TSUSDE,
				ID:                  TSUSDE.Sha256Hash(),
				Decimals:            6,
				JettonMasterAddress: address.MustParseAddr(TSUSDEJettonAddress),
				JettonWalletCode:    getCellFromHex(CodeWalletTsUsde),
			},
			PT_tsUSDe_01Sep2025.ID(): {
				Name:                PT_tsUSDe_01Sep2025,
				ID:                  PT_tsUSDe_01Sep2025.Sha256Hash(),
				Decimals:            9,
				JettonMasterAddress: address.MustParseAddr(PT_tsUSDe_01Sep2025JettonAddress),
				JettonWalletCode:    getCellFromHex(CodeWalletPT_tsUSDe_01Sep2025),
			},
			PT_tsUSDe_18Dec2025.ID(): {
				Name:                PT_tsUSDe_18Dec2025,
				ID:                  PT_tsUSDe_18Dec2025.Sha256Hash(),
				Decimals:            9,
				JettonMasterAddress: address.MustParseAddr(PT_tsUSDe_18Dec2025JettonAddress),
				JettonWalletCode:    getCellFromHex(CodeWalletPT_tsUSDe_18Dec2025),
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
			USDE.ID(): {
				Name:                USDE,
				ID:                  USDE.Sha256Hash(),
				Decimals:            6,
				JettonMasterAddress: address.MustParseAddr(USDEJettonAddress),
				JettonWalletCode:    getCellFromHex(CodeWalletUSDe),
			},
			TSUSDE.ID(): {
				Name:                TSUSDE,
				ID:                  TSUSDE.Sha256Hash(),
				Decimals:            6,
				JettonMasterAddress: address.MustParseAddr(TSUSDEJettonAddress),
				JettonWalletCode:    getCellFromHex(CodeWalletTsUsde),
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
			TONUSDT_STONFI.ID(): {
				Name:                TONUSDT_STONFI,
				ID:                  TONUSDT_STONFI.Sha256Hash(),
				Decimals:            9,
				JettonMasterAddress: address.MustParseAddr(TONUSDTStonfiJettonAddress),
				JettonWalletCode:    getCellFromHex(CodeWalletTonUsdtSton),
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
