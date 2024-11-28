package principal

import (
	"fmt"
	"math"
	"math/big"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"

	"github.com/evaafi/evaa-go-sdk/asset"
	"github.com/evaafi/evaa-go-sdk/config"
)

type Service struct {
	config *config.Config
}

func NewService(config *config.Config) *Service {
	return &Service{config: config}
}

type UserBalancer interface {
	Principal(asset string) *big.Int
	Balance(asset string, assetData *asset.Data, applyDust bool, assetConfig *asset.Config) *big.Int
	ChangePrincipal(asset string, volume *big.Int) UserBalancer
	CheckNotInDebtAtAll() bool
}

type assetManager interface {
	Assets() map[string]*big.Int
	Config(asset string) *asset.Config
	Data(asset string) *asset.Data
}

type priceProvider interface {
	Get(asset string) *big.Int
}

type Health struct {
	GreatestCollateralValue *big.Int
	GreatestCollateralAsset *big.Int
	GreatestLoanValue       *big.Int
	GreatestLoanAsset       *big.Int
	TotalDebt               *big.Int
	TotalLimit              *big.Int
	TotalSupply             *big.Int
}

func (s *Service) CalculateHealth(user UserBalancer, assets assetManager, prices priceProvider) *Health {
	type Principal struct {
		assetID *big.Int
		amount  *big.Int
	}

	greatestCollateralValue := new(big.Int)
	var greatestCollateralAsset *big.Int

	greatestLoanValue := new(big.Int)
	var greatestLoanAsset *big.Int

	totalDebt := new(big.Int)
	totalLimit := new(big.Int)

	totalSupply := new(big.Int)

	calculator := NewCalculator(assets, prices, s.config)
	for assetID, assetBigInt := range assets.Assets() {
		balance := user.Balance(assetID, assets.Data(assetID), false, nil)
		if balance.Sign() == 0 {
			continue
		}

		assetConfig := assets.Config(assetID)
		if balance.Sign() == 1 {
			assetWorth := calculator.ValueFromBalance(balance, assetID)

			totalSupply.Add(totalSupply, assetWorth)
			limit := mulDiv(assetWorth, assetConfig.LiquidationThreshold, s.config.MasterParams.AssetLiquidationThresholdScale)
			totalLimit.Add(totalLimit, limit)
			if assetWorth.Cmp(greatestCollateralValue) == 1 {
				greatestCollateralValue.Set(assetWorth)
				greatestCollateralAsset = assetBigInt
			}
		} else if balance.Sign() == -1 {
			assetWorth := calculator.ValueFromBalance(balance.Neg(balance), assetID)

			totalDebt.Add(totalDebt, assetWorth)
			if assetWorth.Cmp(greatestLoanValue) == 1 {
				greatestLoanValue.Set(assetWorth)
				greatestLoanAsset = assetBigInt
			}
		}
	}

	return &Health{
		GreatestCollateralValue: greatestCollateralValue,
		GreatestCollateralAsset: greatestCollateralAsset,
		GreatestLoanValue:       greatestLoanValue,
		GreatestLoanAsset:       greatestLoanAsset,
		TotalDebt:               totalDebt,
		TotalLimit:              totalLimit,
		TotalSupply:             totalSupply,
	}
}

func (h *Health) IsBadDebt(liquidationBonus, liquidationBonusScale *big.Int) bool {
	return new(big.Int).Mul(h.TotalSupply, liquidationBonusScale).Cmp(new(big.Int).Mul(h.TotalDebt, liquidationBonus)) == -1
}

func (h *Health) IsLiquidatable() bool {
	return h.TotalLimit.Cmp(h.TotalDebt) == -1
}

func (s *Service) CalculateMaximumWithdrawAmount(user UserBalancer, assets assetManager, prices priceProvider, asset string) *big.Int {
	cfg := assets.Config(asset)
	data := assets.Data(asset)
	balance := user.Balance(asset, data, true, cfg)
	if balance.Sign() != 1 {
		return mulDiv(s.GetAvailableToBorrow(user, assets, prices), cfg.Scale(), prices.Get(asset))
	}

	if user.CheckNotInDebtAtAll() || cfg.CollateralFactor.Sign() == 0 {
		return balance
	}

	return bigIntMin(balance,
		bigIntMax(big.NewInt(0), new(big.Int).Sub(
			mulDiv(mulDiv(s.GetAvailableToBorrow(user, assets, prices), s.config.MasterParams.AssetCoefficientScale, cfg.CollateralFactor),
				cfg.Scale(), prices.Get(asset)),
			new(big.Int).Div(mulDiv(data.SRate, cfg.Dust, s.config.MasterParams.FactorScale), big.NewInt(2)),
		)))
}

func (s *Service) GetAvailableToBorrow(user UserBalancer, assets assetManager, prices priceProvider) *big.Int {
	calculator := NewCalculator(assets, prices, s.config)

	borrowAmount := new(big.Int)
	borrowLimit := new(big.Int)
	for a := range assets.Assets() {
		b := user.Balance(a, assets.Data(a), false, nil)
		cfg := assets.Config(a)
		if b.Sign() == -1 {
			borrowAmount.Add(borrowAmount, calculator.ValueFromBalance(b.Neg(b), a))
		} else if b.Sign() == 1 {
			borrowLimit.Add(borrowLimit, mulDiv(calculator.ValueFromBalance(b, a), cfg.CollateralFactor, s.config.MasterParams.AssetCoefficientScale))
		}
	}
	return new(big.Int).Sub(borrowLimit, borrowAmount)
}

func (s *Service) AggregatedBalances(user UserBalancer, assets assetManager, prices priceProvider) (*big.Int, *big.Int) {
	health := s.CalculateHealth(user, assets, prices)
	return health.TotalSupply, health.TotalDebt
}

func (s *Service) PredictHealthFactor(user UserBalancer, assets assetManager, prices priceProvider, asset string, amount *big.Int) float64 {
	if asset != "" || (amount != nil && amount.Sign() != 0) {
		user = user.ChangePrincipal(asset, amount)
	}
	health := s.CalculateHealth(user, assets, prices)

	return health.Factor()
}

func (h *Health) Factor() float64 {
	if h.TotalLimit.Sign() <= 0 {
		return 1
	}

	hf, _ := new(big.Rat).Quo(new(big.Rat).SetInt(h.TotalDebt), new(big.Rat).SetInt(h.TotalLimit)).Float64()
	return math.Min(math.Max(0, 1-hf), 1)
}

func (s *Service) CalculateLiquidationData(user UserBalancer, assets assetManager, prices priceProvider) (health *Health, liquidationAmount, collateralAmount *big.Int, ok bool) {
	health = s.CalculateHealth(user, assets, prices)
	if !health.IsLiquidatable() {
		return health, nil, nil, false
	}

	calculator := NewCalculator(assets, prices, s.config)
	loanAssetString := health.GreatestLoanAsset.String()
	loanAssetConfig := assets.Config(loanAssetString)

	collateralAssetString := health.GreatestCollateralAsset.String()
	collateralAssetConfig := assets.Config(collateralAssetString)

	allowedCollateralValue := health.GreatestCollateralValue
	if !health.IsBadDebt(collateralAssetConfig.LiquidationBonus, s.config.MasterParams.AssetLiquidationBonusScale) {
		allowedCollateralValue = bigIntMin(allowedCollateralValue, bigIntMax(
			new(big.Int).Div(allowedCollateralValue, big.NewInt(2)), s.config.MasterParams.CollateralWorthThreshold))
	}

	liquidationValue := bigIntMin(health.GreatestLoanValue, mulDiv(
		allowedCollateralValue,
		s.config.MasterParams.AssetLiquidationBonusScale,
		collateralAssetConfig.LiquidationBonus),
	)
	collateralValue := mulDiv(liquidationValue, collateralAssetConfig.LiquidationBonus, s.config.MasterParams.AssetLiquidationBonusScale)
	collateralAmount = calculator.BalanceFromValue(collateralValue, collateralAssetString)

	liquidationValue = mulDiv(
		liquidationValue,
		s.config.MasterParams.AssetLiquidationReserveFactorScale,
		new(big.Int).Sub(
			s.config.MasterParams.AssetLiquidationReserveFactorScale,
			loanAssetConfig.LiquidationReserveFactor,
		),
	)
	liquidationAmount = calculator.BalanceFromValue(liquidationValue, loanAssetString)

	return health, liquidationAmount, collateralAmount, true
}

func (s *Service) CalculateUserSCAddress(userAddress *address.Address) (*address.Address, error) {
	if userAddress == nil {
		return nil, fmt.Errorf("userAddress cannot be a nil pointer")
	}

	lendingData := cell.BeginCell().
		MustStoreAddr(s.config.MasterAddress).
		MustStoreAddr(userAddress).
		MustStoreUInt(0, 8).
		MustStoreBoolBit(false).
		EndCell()

	stateInit := &tlb.StateInit{
		Data: lendingData,
		Code: s.config.LendingCode,
	}
	stateCell, err := tlb.ToCell(stateInit)
	if err != nil {
		return nil, fmt.Errorf("failed to get state cell: %w", err)
	}

	return address.NewAddress(0, 0, stateCell.Hash()), nil
}

func mulDiv(x, y, z *big.Int) *big.Int {
	return new(big.Int).Div(new(big.Int).Mul(x, y), z)
}

func bigIntMin(first *big.Int, nums ...*big.Int) *big.Int {
	res := first
	for _, num := range nums {
		if num.Cmp(res) < 0 {
			res = num
		}
	}
	return res
}

func bigIntMax(first *big.Int, nums ...*big.Int) *big.Int {
	res := first
	for _, num := range nums {
		if num.Cmp(res) > 0 {
			res = num
		}
	}
	return res
}
