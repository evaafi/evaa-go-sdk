package asset

import (
	"fmt"
	"github.com/evaafi/evaa-go-sdk/config"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"maps"
	"math/big"
	"sync"
	"time"
)

type Parser struct {
	keys map[string]*big.Int

	mtx    sync.RWMutex
	config map[string]*Config
	data   map[string]*Data
}

func NewParser(config *config.Config) *Parser {
	keys := make(map[string]*big.Int, len(config.Assets))
	for _, asset := range config.Assets {
		keys[asset.ID.String()] = asset.ID
	}
	return &Parser{keys: keys}
}

func (p *Parser) Assets() map[string]*big.Int {
	return maps.Clone(p.keys)
}

func (p *Parser) Config(asset string) *Config {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	return p.config[asset]
}

func (p *Parser) Data(asset string) *Data {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	return p.data[asset]
}

func (p *Parser) SetInfo(data *cell.Dictionary, config *cell.Dictionary) error {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if data != nil {
		err := p.setData(data)
		if err != nil {
			return fmt.Errorf("setData err: %w", err)
		}
	}

	if config != nil {
		err := p.setConfig(config)
		if err != nil {
			return fmt.Errorf("setConfig err: %w", err)
		}
	}

	return nil
}

func (p *Parser) setConfig(config *cell.Dictionary) error {
	p.config = make(map[string]*Config, len(p.keys))
	for asset, id := range p.keys {
		assetConfig, err := config.LoadValue(uIntSliceKey(id))
		if err != nil {
			return err
		}

		assetConfigValueRef, err := assetConfig.LoadRef()
		if err != nil {
			return err
		}

		p.config[asset] = newConfig(assetConfig, assetConfigValueRef)
	}
	return nil
}

func (p *Parser) setData(data *cell.Dictionary) error {
	p.data = make(map[string]*Data, len(p.keys))
	for asset, id := range p.keys {
		assetData, err := data.LoadValue(uIntSliceKey(id))
		if err != nil {
			return err
		}

		p.data[asset] = newData(assetData)
	}
	return nil
}

func uIntSliceKey(id *big.Int) *cell.Cell {
	return cell.BeginCell().MustStoreBigUInt(id, 256).EndCell()
}

func (p *Parser) UpdateCurrentRates(forward int64) *Parser {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	ts := time.Now().Unix() + forward
	data := make(map[string]*Data, len(p.keys))
	for asset := range p.keys {
		data[asset], _, _ = p.calculateCurrentRates(asset, ts)
	}

	return &Parser{
		keys:   p.keys,
		config: p.config,
		data:   data,
	}
}

func (p *Parser) CalculateCurrentRates(asset string) (*Data, *big.Int, *big.Int) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	return p.calculateCurrentRates(asset, time.Now().Unix())
}

func (p *Parser) calculateCurrentRates(asset string, ts int64) (_ *Data, supplyInterest, borrowInterest *big.Int) {
	assetData := p.data[asset]
	timeElapsed := ts - assetData.LastAccrual.Int64()
	if timeElapsed <= 0 {
		return assetData, big.NewInt(0), big.NewInt(0)
	}
	totalSupply := mulDiv(assetData.SRate, assetData.TotalSupply, big.NewInt(1e12))
	totalBorrow := mulDiv(assetData.BRate, assetData.TotalBorrow, big.NewInt(1e12))

	utilization := new(big.Int)
	if totalSupply.Sign() != 0 {
		utilization = mulDiv(totalBorrow, big.NewInt(1e12), totalSupply)
	}

	assetConfig := p.config[asset]

	if utilization.Cmp(assetConfig.TargetUtilization) != 1 {
		borrowInterest = new(big.Int).Add(assetConfig.BaseBorrowRate,
			mulDiv(assetConfig.BorrowRateSlopeLow, utilization, big.NewInt(1e12)))
	} else {
		borrowInterest = new(big.Int).Add(assetConfig.BaseBorrowRate, new(big.Int).Add(
			mulDiv(assetConfig.BorrowRateSlopeLow, assetConfig.TargetUtilization, big.NewInt(1e12)),
			mulDiv(assetConfig.BorrowRateSlopeHigh, new(big.Int).Sub(utilization, assetConfig.TargetUtilization), big.NewInt(1e12)),
		))
	}
	supplyInterest = mulDiv(mulDiv(borrowInterest, utilization, big.NewInt(1e12)),
		new(big.Int).Sub(big.NewInt(10_000), assetConfig.ReserveFactor), big.NewInt(10_000))

	timeElapsedBigInt := big.NewInt(timeElapsed)

	return &Data{
		SRate:       new(big.Int).Add(assetData.SRate, mulDiv(assetData.SRate, new(big.Int).Mul(supplyInterest, timeElapsedBigInt), big.NewInt(1e12))),
		BRate:       new(big.Int).Add(assetData.BRate, mulDiv(assetData.BRate, new(big.Int).Mul(borrowInterest, timeElapsedBigInt), big.NewInt(1e12))),
		TotalSupply: assetData.TotalSupply,
		TotalBorrow: assetData.TotalBorrow,
		LastAccrual: big.NewInt(ts),
	}, supplyInterest, borrowInterest
}

func mulDiv(x, y, z *big.Int) *big.Int {
	return new(big.Int).Div(new(big.Int).Mul(x, y), z)
}
