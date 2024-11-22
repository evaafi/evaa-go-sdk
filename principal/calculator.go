package principal

import (
	"github.com/evaafi/evaa-go-sdk/config"
	"math/big"
)

type Calculator struct {
	assets assetManager
	prices priceProvider
	config *config.Config
}

func NewCalculator(assets assetManager, prices priceProvider, config *config.Config) *Calculator {
	return &Calculator{assets: assets, prices: prices, config: config}
}

func (c *Calculator) BalanceFromPrincipal(principal *big.Int, asset string) *big.Int {
	if principal.Sign() == 0 {
		return big.NewInt(0)
	}
	if principal.Sign() == 1 {
		return mulDiv(principal, c.assets.Data(asset).SRate, c.config.MasterParams.FactorScale)
	}
	return mulDiv(principal, c.assets.Data(asset).BRate, c.config.MasterParams.FactorScale)
}

func (c *Calculator) PrincipalFromBalance(balance *big.Int, asset string) *big.Int {
	if balance.Sign() == 0 {
		return big.NewInt(0)
	}
	if balance.Sign() == 1 {
		return mulDiv(balance, c.config.MasterParams.FactorScale, c.assets.Data(asset).SRate)
	}
	return mulDiv(balance, c.config.MasterParams.FactorScale, c.assets.Data(asset).BRate)
}

func (c *Calculator) ValueFromPrincipal(principal *big.Int, asset string) *big.Int {
	return c.ValueFromBalance(c.BalanceFromPrincipal(principal, asset), asset)
}

func (c *Calculator) PrincipalFromValue(value *big.Int, asset string) *big.Int {
	return c.PrincipalFromBalance(c.BalanceFromValue(value, asset), asset)
}

func (c *Calculator) BalanceFromValue(balance *big.Int, asset string) *big.Int {
	if balance.Sign() == 0 {
		return big.NewInt(0)
	}
	return mulDiv(balance, c.assets.Config(asset).Scale(), c.prices.Get(asset))
}

func (c *Calculator) ValueFromBalance(value *big.Int, asset string) *big.Int {
	if value.Sign() == 0 {
		return big.NewInt(0)
	}
	return mulDiv(value, c.prices.Get(asset), c.assets.Config(asset).Scale())
}
