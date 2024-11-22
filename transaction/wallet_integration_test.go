//go:build integration
// +build integration

package transaction_test

import (
	"context"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/jetton"
	"github.com/xssnick/tonutils-go/ton/wallet"

	"github.com/evaafi/evaa-go-sdk/asset"
	"github.com/evaafi/evaa-go-sdk/config"
	"github.com/evaafi/evaa-go-sdk/price"
	"github.com/evaafi/evaa-go-sdk/principal"
	"github.com/evaafi/evaa-go-sdk/transaction"
)

const emptyWalletSeedEnvFatalMsg = "WALLET_SEED not found in environment"
const emptyWalletSeed2EnvFatalMsg = "WALLET_SEED2 not found in environment"

var _seed = os.Getenv("WALLET_SEED")
var _seed2 = os.Getenv("WALLET_SEED2")

func TestWallet_SendWithdraw(t *testing.T) {
	if _seed == "" {
		t.Fatal(emptyWalletSeedEnvFatalMsg)
	}
	seed := strings.Split(_seed, " ")

	cfg := config.GetMasterTestnetConfig()

	client := liteclient.NewConnectionPool()
	err := client.AddConnectionsFromConfigUrl(context.Background(), "https://ton.org/testnet-global.config.json")
	if err != nil {
		t.Fatal(err)
	}
	api := ton.NewAPIClient(client)
	masterchainInfo, err := api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	assetService := asset.NewParser(cfg)
	getAssetsDataResult, err := api.WaitForBlock(masterchainInfo.SeqNo).RunGetMethod(context.Background(), masterchainInfo, cfg.MasterAddress, "getAssetsData")
	if err != nil {
		t.Fatal(err)
	}
	getAssetsConfigResult, err := api.WaitForBlock(masterchainInfo.SeqNo).RunGetMethod(context.Background(), masterchainInfo, cfg.MasterAddress, "getAssetsConfig")
	if err != nil {
		t.Fatal(err)
	}
	err = assetService.SetInfo(getAssetsDataResult.MustCell(0).AsDict(256), getAssetsConfigResult.MustCell(0).AsDict(256))
	if err != nil {
		t.Fatal(err)
	}

	priceService := price.NewService(cfg, nil)
	prices, err := priceService.GetPrices(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	w, err := wallet.FromSeed(api, seed, wallet.ConfigV5R1Final{
		NetworkGlobalID: wallet.TestnetGlobalID,
		Workchain:       0,
	})
	t.Log(w.WalletAddress())
	transactionService := transaction.NewWallet(cfg, w, transaction.NewBuilder(cfg))
	principalService := principal.NewService(cfg)

	userSCAddress, err := principalService.CalculateUserSCAddress(w.WalletAddress())
	if err != nil {
		t.Fatal(err)
	}
	userSCAccount, err := api.WaitForBlock(masterchainInfo.SeqNo).GetAccount(context.Background(), masterchainInfo, userSCAddress)
	if err != nil {
		t.Fatal(err)
	}

	if !userSCAccount.IsActive || true {
		jettonWallet, err := jetton.NewJettonMasterClient(api, cfg.Assets[config.STTON.ID()].JettonMasterAddress).GetJettonWallet(context.Background(), w.WalletAddress())
		if err != nil {
			t.Fatal(err)
		}
		balance, err := jettonWallet.GetBalance(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if balance.Sign() > 0 {
			t.Log("supply", balance)
			err = transactionService.SendSupply(context.Background(), &transaction.SupplyParameters{
				Asset:            config.STTON.Sha256Hash(),
				QueryID:          0,
				IncludeUserCode:  true,
				Amount:           balance,
				AmountToTransfer: big.NewInt(0),
			}, true)
			if err != nil {
				t.Fatal(err)
			}
			time.Sleep(10 * time.Second)
		}

		masterchainInfo, err = api.CurrentMasterchainInfo(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		userSCAccount, err = api.WaitForBlock(masterchainInfo.SeqNo).GetAccount(context.Background(), masterchainInfo, userSCAddress)
		if err != nil {
			t.Fatal(err)
		}
		prices, err = priceService.GetPrices(context.Background())
		if err != nil {
			t.Fatal(err)
		}
	}

	userSC, err := principal.NewUserSC(w.WalletAddress()).SetAccData(userSCAccount.Data)
	if err != nil {
		t.Fatal(err)
	}
	if health, liquidationAmount, minCollateralAmount, ok := principalService.CalculateLiquidationData(userSC, assetService, prices); ok {
		t.Log("liq", liquidationAmount, minCollateralAmount)
		liqBody := &transaction.LiquidationBaseData{
			BorrowerAddress:     w.WalletAddress(),
			LoanAsset:           health.GreatestLoanAsset,
			CollateralAsset:     health.GreatestCollateralAsset,
			MinCollateralAmount: minCollateralAmount.Sub(minCollateralAmount, big.NewInt(10)),
			LiquidationAmount:   liquidationAmount,
		}
		if _seed2 == "" {
			t.Fatal(emptyWalletSeed2EnvFatalMsg)
		}
		seed2 := strings.Split(_seed2, " ")
		liqW, err := wallet.FromSeed(api, seed2, wallet.V4R2)
		if err != nil {
			t.Fatal(err)
		}
		err = transaction.NewWallet(cfg, liqW, transaction.NewBuilder(cfg)).SendLiquidation(context.Background(), &transaction.LiquidationParameters{
			LiquidationBaseData: liqBody,
			IncludeUserCode:     true,
			PriceData:           prices.Data(),
		}, true)
		if err != nil {
			t.Fatal(err)
		}
		return
	}

	amount := principalService.CalculateMaximumWithdrawAmount(userSC, assetService, prices, config.TON.ID())
	t.Log("withdraw", amount)

	err = transactionService.SendWithdraw(context.Background(), &transaction.WithdrawParameters{
		Asset:           config.TON.Sha256Hash(),
		Amount:          amount,
		PriceData:       prices.Data(),
		IncludeUserCode: true,
	}, true)
	if err != nil {
		t.Fatal(err)
	}
}
