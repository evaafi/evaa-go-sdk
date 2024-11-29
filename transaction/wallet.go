package transaction

import (
	"context"
	"fmt"
	"math/big"

	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/wallet"

	"github.com/evaafi/evaa-go-sdk/config"
)

type Wallet struct {
	config  *config.Config
	wallet  *wallet.Wallet
	builder Builder
}

func NewWallet(config *config.Config, wallet *wallet.Wallet, builder Builder) *Wallet {
	return &Wallet{config: config, wallet: wallet, builder: builder}
}

func (w *Wallet) SendSupply(ctx context.Context, data *SupplyParameters, wait bool) error {
	message, err := w.builder.CreateSupplyBody(data, w.wallet.WalletAddress())
	if err != nil {
		return fmt.Errorf("failed to create message, err: %w", err)
	}

	amount := new(big.Int).Set(data.Amount)
	dst := w.config.MasterAddress
	asset := w.config.Assets[data.Asset.String()]
	if asset.JettonMasterAddress != nil {
		amount = big.NewInt(config.FeeSupplyJetton)
		dst, err = asset.GetJettonWalletAddress(w.wallet.WalletAddress())
		if err != nil {
			return fmt.Errorf("failed to get jetton wallet address, err: %w", err)
		}
	} else {
		amount.Add(amount, big.NewInt(config.FeeSupply))
	}

	return w.wallet.Send(ctx, wallet.SimpleMessage(dst, tlb.FromNanoTON(amount), message), wait)
}

func (w *Wallet) SendWithdraw(ctx context.Context, data *WithdrawParameters, wait bool) error {
	message, err := w.builder.CreateWithdrawBody(data, w.wallet.WalletAddress())
	if err != nil {
		return fmt.Errorf("failed to create message, err: %w", err)
	}

	return w.wallet.Send(ctx, wallet.SimpleMessage(w.config.MasterAddress, tlb.FromNanoTONU(config.FeeWithdraw), message), wait)
}

func (w *Wallet) SendLiquidation(ctx context.Context, data *LiquidationParameters, wait bool) error {
	message, err := w.builder.CreateLiquidationBody(data, w.wallet.WalletAddress())
	if err != nil {
		return fmt.Errorf("failed to create message, err: %w", err)
	}

	amount := new(big.Int).Set(data.LiquidationAmount)
	dst := w.config.MasterAddress
	asset := w.config.Assets[data.LoanAsset.String()]
	if asset.JettonMasterAddress != nil {
		amount = big.NewInt(config.FeeLiquidationJetton)
		dst, err = asset.GetJettonWalletAddress(w.wallet.WalletAddress())
		if err != nil {
			return fmt.Errorf("failed to get jetton wallet address, err: %w", err)
		}
	} else {
		amount.Add(amount, big.NewInt(config.FeeLiquidation))
	}

	return w.wallet.Send(ctx, wallet.SimpleMessage(dst, tlb.FromNanoTON(amount), message), wait)
}
