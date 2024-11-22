package transaction

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tvm/cell"

	"github.com/evaafi/evaa-go-sdk/config"
)

type Builder interface {
	CreateSupplyBody(data *SupplyParameters, myAddress *address.Address) (*cell.Cell, error)
	CreateWithdrawBody(data *WithdrawParameters, myAddress *address.Address) (*cell.Cell, error)
	CreateLiquidationBody(data *LiquidationParameters, myAddress *address.Address) (*cell.Cell, error)
}

type builder struct {
	config *config.Config
}

func NewBuilder(config *config.Config) Builder {
	return &builder{config: config}
}

type SupplyParameters struct {
	Asset            *big.Int
	QueryID          uint64
	IncludeUserCode  bool
	Amount           *big.Int
	AmountToTransfer *big.Int
	Payload          *cell.Cell
	ResponseAddress  *address.Address
	ForwardAmount    *big.Int
}

func (s *builder) CreateSupplyBody(data *SupplyParameters, myAddress *address.Address) (*cell.Cell, error) {
	assetData, ok := s.config.Assets[data.Asset.String()]
	if !ok {
		return nil, errors.New("unknown asset")
	}
	includeUserCode := int64(-1)
	if !data.IncludeUserCode {
		includeUserCode = 0
	}
	payload := data.Payload
	if payload == nil {
		payload = cell.BeginCell().EndCell()
	}

	if data.Amount == nil {
		return nil, fmt.Errorf("wrong amount value")
	}

	if data.AmountToTransfer == nil {
		data.AmountToTransfer = big.NewInt(0)
	}
	if !data.AmountToTransfer.IsUint64() {
		return nil, fmt.Errorf("wrong amountToTransfer value")
	}

	if assetData.JettonMasterAddress == nil {
		if !data.Amount.IsUint64() {
			return nil, fmt.Errorf("wrong amount value")
		}
		return cell.BeginCell().
			MustStoreUInt(config.OpcodeSupply, 32).
			MustStoreUInt(data.QueryID, 64).
			MustStoreInt(includeUserCode, 2).
			MustStoreBigUInt(data.Amount, 64).
			MustStoreAddr(myAddress).
			MustStoreBigUInt(data.AmountToTransfer, 64).
			MustStoreRef(payload).
			EndCell(), nil
	}
	responseAddress := data.ResponseAddress
	if responseAddress == nil {
		responseAddress = myAddress
	}
	forwardAmount := data.ForwardAmount
	if forwardAmount == nil {
		forwardAmount = big.NewInt(config.FeeSupplyJettonFWD)
	}
	return cell.BeginCell().
		MustStoreUInt(config.OpcodeJettonTransfer, 32).
		MustStoreUInt(data.QueryID, 64).
		MustStoreBigCoins(data.Amount).
		MustStoreAddr(s.config.MasterAddress).
		MustStoreAddr(responseAddress).
		MustStoreBoolBit(false).
		MustStoreBigCoins(forwardAmount).
		MustStoreBoolBit(true).
		MustStoreRef(cell.BeginCell().
			MustStoreUInt(config.OpcodeSupply, 32).
			MustStoreInt(includeUserCode, 2).
			MustStoreAddr(myAddress).
			MustStoreBigUInt(data.AmountToTransfer, 64).
			MustStoreRef(payload).
			EndCell()).
		EndCell(), nil
}

type WithdrawParameters struct {
	QueryID          uint64
	Asset            *big.Int
	Amount           *big.Int
	IncludeUserCode  bool
	AmountToTransfer *big.Int
	Payload          *cell.Cell
	PriceData        *cell.Cell
}

func (s *builder) CreateWithdrawBody(data *WithdrawParameters, myAddress *address.Address) (*cell.Cell, error) {
	includeUserCode := int64(-1)
	if !data.IncludeUserCode {
		includeUserCode = 0
	}
	payload := data.Payload
	if payload == nil {
		payload = cell.BeginCell().EndCell()
	}
	if data.Amount == nil || !data.Amount.IsUint64() {
		return nil, fmt.Errorf("wrong amount value")
	}
	if data.AmountToTransfer == nil {
		data.AmountToTransfer = big.NewInt(0)
	}
	if !data.AmountToTransfer.IsUint64() {
		return nil, fmt.Errorf("wrong amountToTransfer value")
	}
	if data.PriceData == nil {
		return nil, fmt.Errorf("wrong priceData value")
	}
	return cell.BeginCell().
		MustStoreUInt(config.OpcodeWithdraw, 32).
		MustStoreUInt(data.QueryID, 64).
		MustStoreBigUInt(data.Asset, 256).
		MustStoreBigUInt(data.Amount, 64).
		MustStoreAddr(myAddress).
		MustStoreInt(includeUserCode, 2).
		MustStoreBigUInt(data.AmountToTransfer, 64).
		MustStoreRef(payload).
		MustStoreRef(data.PriceData).
		EndCell(), nil
}

type LiquidationBaseData struct {
	BorrowerAddress     *address.Address
	LoanAsset           *big.Int
	CollateralAsset     *big.Int
	MinCollateralAmount *big.Int
	LiquidationAmount   *big.Int
}

type LiquidationParameters struct {
	*LiquidationBaseData
	PriceData            *cell.Cell
	QueryID              uint64
	ResponseAddress      *address.Address
	IncludeUserCode      bool
	Payload              *cell.Cell
	PayloadForwardAmount *big.Int
	ForwardAmount        *big.Int
}

func NewLiquidationParameters(liquidationBaseData *LiquidationBaseData, priceData *cell.Cell) *LiquidationParameters {
	return &LiquidationParameters{LiquidationBaseData: liquidationBaseData, PriceData: priceData}
}
func (l *LiquidationParameters) SetQueryID(q uint64) *LiquidationParameters {
	l.QueryID = q
	return l
}
func (l *LiquidationParameters) SetIncludeUserCode(b bool) *LiquidationParameters {
	l.IncludeUserCode = b
	return l
}
func (l *LiquidationParameters) SetResponseAddress(responseAddress *address.Address) *LiquidationParameters {
	l.ResponseAddress = responseAddress
	return l
}
func (l *LiquidationParameters) SetForwardAmount(forwardAmount *big.Int) *LiquidationParameters {
	l.ForwardAmount = forwardAmount
	return l
}
func (l *LiquidationParameters) SetPayload(payload *cell.Cell, payloadForwardAmount *big.Int) *LiquidationParameters {
	l.Payload = payload
	l.PayloadForwardAmount = payloadForwardAmount
	return l
}

func (s *builder) CreateLiquidationBody(data *LiquidationParameters, myAddress *address.Address) (*cell.Cell, error) {
	assetData, ok := s.config.Assets[data.LoanAsset.String()]
	if !ok {
		return nil, errors.New("unknown asset")
	}
	payload := cell.BeginCell().MustStoreUInt(0, 64).MustStoreRef(cell.BeginCell().EndCell()).EndCell()
	if data.Payload != nil && data.PayloadForwardAmount.Sign() == 1 {
		if !data.PayloadForwardAmount.IsUint64() {
			return nil, fmt.Errorf("wrong payloadForwardAmount value")
		}
		payload = cell.BeginCell().MustStoreBigUInt(data.PayloadForwardAmount, 64).MustStoreRef(data.Payload).EndCell()
	}
	includeUserCode := int64(-1)
	if !data.IncludeUserCode {
		includeUserCode = 0
	}
	if data.MinCollateralAmount == nil {
		return nil, fmt.Errorf("wrong minCollateralAmount value")
	}
	if data.LiquidationAmount == nil || !data.LiquidationAmount.IsUint64() {
		return nil, fmt.Errorf("wrong liquidationAmount value")
	}
	if data.PriceData == nil {
		return nil, fmt.Errorf("wrong priceData value")
	}
	if assetData.JettonMasterAddress == nil {
		if !data.MinCollateralAmount.IsUint64() {
			return nil, fmt.Errorf("wrong amount value")
		}
		return cell.BeginCell().
			MustStoreUInt(config.OpcodeLiquidate, 32).
			MustStoreUInt(data.QueryID, 64).
			MustStoreAddr(data.BorrowerAddress).
			MustStoreAddr(myAddress).
			MustStoreBigUInt(data.CollateralAsset, 256).
			MustStoreBigUInt(data.MinCollateralAmount, 64).
			MustStoreInt(includeUserCode, 2).
			MustStoreBigUInt(data.LiquidationAmount, 64).
			MustStoreRef(payload).
			MustStoreRef(data.PriceData).
			EndCell(), nil
	}
	responseAddress := data.ResponseAddress
	if responseAddress == nil {
		responseAddress = myAddress
	}
	forwardAmount := data.ForwardAmount
	if forwardAmount == nil {
		forwardAmount = big.NewInt(config.FeeLiquidationJettonFWD)
	}
	return cell.BeginCell().
		MustStoreUInt(config.OpcodeJettonTransfer, 32).
		MustStoreUInt(data.QueryID, 64).
		MustStoreBigCoins(data.LiquidationAmount).
		MustStoreAddr(s.config.MasterAddress).
		MustStoreAddr(responseAddress).
		MustStoreBoolBit(false).
		MustStoreBigCoins(forwardAmount).
		MustStoreBoolBit(true).
		MustStoreRef(cell.BeginCell().
			MustStoreUInt(config.OpcodeLiquidate, 32).
			MustStoreAddr(data.BorrowerAddress).
			MustStoreAddr(myAddress).
			MustStoreBigUInt(data.CollateralAsset, 256).
			MustStoreBigUInt(data.MinCollateralAmount, 64).
			MustStoreInt(includeUserCode, 2).
			MustStoreUInt(0, 64).
			MustStoreRef(payload).
			MustStoreRef(data.PriceData).
			EndCell()).
		EndCell(), nil
}
