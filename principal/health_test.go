package principal

import (
	"math/big"
	"testing"

	"github.com/evaafi/evaa-go-sdk/asset"
	"github.com/evaafi/evaa-go-sdk/config"
)

type Prices map[string]*big.Int

func (p Prices) Get(asset string) *big.Int {
	return p[asset]
}

func TestService_GetAvailableToBorrow(t *testing.T) {
	cfg := config.GetMainMainnetConfig()
	parser := getAssetsService(t, cfg)
	prices := Prices(map[string]*big.Int{
		"11876925370864614464799087627157805050745321306404563164673853337929163193738": big.NewInt(4575000000),
		"81203563022592193867903899252711112850180680126331353892172221352147647262515": big.NewInt(998760000),
		"59636546167967198470134647008558085436004969028957957410318094280110082891718": big.NewInt(999908900),
		"33171510858320790266247832496974106978700190498800858393089426423762035476944": big.NewInt(5030456829),
		"23103091784861387372100043848078515239542568751939923972799733728526040769767": big.NewInt(5003470805),
		"91621667903763073563570557639433445791506232618002614896981036659302854767224": big.NewInt(998760000),
	})
	service := NewService(cfg)
	user := NewUserSC(nil)
	tonAsset := config.TON.ID()
	usdtAsset := config.USDT.ID()
	targetAsset := usdtAsset
	user.principals = map[string]*big.Int{
		tonAsset:  big.NewInt(1350457583812),
		usdtAsset: big.NewInt(0),
	}
	healthFactor := service.PredictHealthFactor(user, parser, prices, "", nil)
	if healthFactor != 1 {
		t.Errorf("healthFactor want %d, got %f", 1, healthFactor)
	}
	availableToBorrow := service.GetAvailableToBorrow(user, parser, prices)
	if availableToBorrow.Cmp(big.NewInt(3623455705047)) != 0 {
		t.Errorf("availableToBorrow want %d, got %s", 3623455705047, availableToBorrow)
	}

	change := NewCalculator(parser, prices, cfg).PrincipalFromValue(new(big.Int).Neg(availableToBorrow), targetAsset)
	if change.Cmp(big.NewInt(-4171290650)) != 0 {
		t.Errorf("change want %d, got %s", -4171290649, change)
	}
	healthFactorPredict := service.PredictHealthFactor(user, parser, prices, targetAsset, change)
	if healthFactorPredict > 0.077 {
		t.Errorf("healthFactorPredict want %d, got %f", 0, healthFactorPredict)
	}
	user2 := user.ChangePrincipal(targetAsset, change)
	healthFactor2 := service.PredictHealthFactor(user2, parser, prices, targetAsset, big.NewInt(0))
	if healthFactor2 != healthFactorPredict {
		t.Errorf("healthFactor2 want %f, got %f", healthFactorPredict, healthFactor2)
	}
	availableToBorrow2 := service.GetAvailableToBorrow(user2, parser, prices)
	if new(big.Int).Abs(availableToBorrow2).Cmp(new(big.Int).Abs(big.NewInt(1555))) == 1 {
		t.Errorf("availableToBorrow2 want %d, got %s", 1555, availableToBorrow2)
	}
}

func TestService_CalculateMaximumWithdrawAmount(t *testing.T) {
	cfg := config.GetMainMainnetConfig()
	parser := getAssetsService(t, cfg)
	prices := Prices(map[string]*big.Int{
		"11876925370864614464799087627157805050745321306404563164673853337929163193738": big.NewInt(4975000000),
		"81203563022592193867903899252711112850180680126331353892172221352147647262515": big.NewInt(998760000),
		"59636546167967198470134647008558085436004969028957957410318094280110082891718": big.NewInt(999908900),
		"33171510858320790266247832496974106978700190498800858393089426423762035476944": big.NewInt(5230456829),
		"23103091784861387372100043848078515239542568751939923972799733728526040769767": big.NewInt(5193470805),
		"91621667903763073563570557639433445791506232618002614896981036659302854767224": big.NewInt(998760000),
	})
	service := NewService(cfg)
	user := NewUserSC(nil)
	tonAsset := (config.TON).ID()
	usdtAsset := (config.USDT).ID()
	targetAsset := usdtAsset
	user.principals = map[string]*big.Int{
		tonAsset:  big.NewInt(1350457583812),
		usdtAsset: big.NewInt(-4519473935),
	}
	healthFactor := service.PredictHealthFactor(user, parser, prices, "", nil)
	if healthFactor <= 0.08 {
		t.Errorf("healthFactor want %d, got %f", 0, healthFactor)
	}
	maximumWithdrawAmount := service.CalculateMaximumWithdrawAmount(user, parser, prices, targetAsset)
	if maximumWithdrawAmount.Cmp(big.NewInt(14367925)) != 0 {
		t.Errorf("maximumWithdrawAmount TON want %d, got %s", 14367925, maximumWithdrawAmount)
	}

	change := mulDiv(new(big.Int).Neg(maximumWithdrawAmount), service.config.MasterParams.FactorScale, parser.Data(targetAsset).BRate)

	healthFactorPredict := service.PredictHealthFactor(user, parser, prices, targetAsset, change)
	if healthFactorPredict >= 0.77 {
		t.Errorf("healthFactorPredict want %d, got %f", 0, healthFactorPredict)
	}

	user2 := user.ChangePrincipal(targetAsset, change)
	healthFactor2 := service.PredictHealthFactor(user2, parser, prices, targetAsset, big.NewInt(0))
	if healthFactor2 != healthFactorPredict {
		t.Errorf("healthFactor2 want %f, got %f", healthFactorPredict, healthFactor2)
	}
	maximumWithdrawAmount2 := service.CalculateMaximumWithdrawAmount(user2, parser, prices, targetAsset)
	if maximumWithdrawAmount2.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("maximumWithdrawAmount2 TON want %d, got %s", 0, maximumWithdrawAmount2)
	}
}

func TestService_CalculateLiquidationData(t *testing.T) {
	cfg := config.GetMainMainnetConfig()
	parser := getAssetsService(t, cfg)
	prices := Prices(map[string]*big.Int{
		"11876925370864614464799087627157805050745321306404563164673853337929163193738": big.NewInt(4575000000),
		"81203563022592193867903899252711112850180680126331353892172221352147647262515": big.NewInt(998760000),
		"59636546167967198470134647008558085436004969028957957410318094280110082891718": big.NewInt(999908900),
		"33171510858320790266247832496974106978700190498800858393089426423762035476944": big.NewInt(5030456829),
		"23103091784861387372100043848078515239542568751939923972799733728526040769767": big.NewInt(5003470805),
		"91621667903763073563570557639433445791506232618002614896981036659302854767224": big.NewInt(998760000),
	})
	service := NewService(cfg)
	user := NewUserSC(nil)
	user.principals = map[string]*big.Int{
		(config.TON).ID():  big.NewInt(1350457583812),
		(config.USDT).ID(): big.NewInt(-4519473935),
	}
	if hf := service.PredictHealthFactor(user, parser, prices, "", big.NewInt(0)); hf != 0 {
		t.Errorf("PredictHealthFactor want %d, got %f", 0, hf)
	}
	health, liquidationAmount, minCollateralAmount, ok := service.CalculateLiquidationData(user, parser, prices)
	if health.TotalSupply.Cmp(big.NewInt(5032577368122)) != 0 {
		t.Errorf("TotalSupply want %d, got %s", 5032577368122, health.TotalSupply)
	}
	if health.TotalDebt.Cmp(big.NewInt(3925910466047)) != 0 {
		t.Errorf("TotalDebt want %d, got %s", 3925910466047, health.TotalDebt)
	}
	if health.GreatestCollateralValue.Cmp(big.NewInt(5032577368122)) != 0 {
		t.Errorf("GreatestCollateralValue want %d, got %s", 5032577368122, health.GreatestCollateralValue)
	}
	if health.GreatestLoanValue.Cmp(big.NewInt(3925910466047)) != 0 {
		t.Errorf("GreatestLoanValue want %d, got %s", 3925910466047, health.GreatestLoanValue)
	}
	if !health.IsLiquidatable() || !ok {
		t.Errorf("IsLiquidatable want %v, got %v", true, false)
	}
	if health.IsBadDebt(parser.Config(health.GreatestCollateralAsset.String()).LiquidationBonus, cfg.MasterParams.AssetLiquidationBonusScale) {
		t.Errorf("IsBadDebt want %v, got %v", false, true)
	}
	if liquidationAmount == nil || liquidationAmount.Cmp(big.NewInt(2375549479)) != 0 {
		t.Errorf("liquidationAmount want %d, got %s", 2375549479, liquidationAmount)
	}
	if minCollateralAmount == nil || minCollateralAmount.Cmp(big.NewInt(550008455532)) != 0 {
		t.Errorf("minCollateralAmount want %d, got %s", 550008455532, minCollateralAmount)
	}
}

func getAssetsService(t *testing.T, cfg *config.Config) *asset.Parser {
	parser := asset.NewParser(cfg)
	assetsDataCell, _ := config.GetCellFromHex("b5ee9c7241020b0100028300020120010202012003040201200506020120070800cabf895668e908644f30322b997de8faaafc21f05aa52f8982f042dac1fe0b4d09d0000000ba85859fd4000000bc4aa2bec00007de36d82a0c2b00008ca20ad5529a671f51b90005dd9121feaa18000000000000000000000000000000000000000000000000020120090a00cabf8a9006bd3fb03d355daeeff93b24be90afaa6e3ca0073ff5720f8a852c933278000000c319365650000000ca80c5df17000018070112552e00000c648b3cc706671f52a20000099c5958d4bd00000000000000000000000000000000000000000006b6c000c9bf748433fcbcc1ac75e54798fb9cdfd8d368b8d6ae3092f4c291cf8465590f7b140000017b4e135da2000001825d6af8ea00191e1744cd410e000950b177ace73ece3ea1b6000cdbbf61a9a52c00000000000000000000000000000000000000000000000100c9bf6627c5eaf750e15e689006a18f136130fa2b6874a62e57f9c529bc43cfae49ce000001753bedb04e00000178963de082000f1b122b25f50c0005a70bd5bccffece3ea5440007923142840bf200000000000000000000000000000000000000000000000100c9bf47b22d8d0a21004209a3eeb54d9c61d63c8ef5dbc1a701ddc4311c1cacb03f8c000001580ac5064400000361ca46bd9e00000000b462d90c000000008911a282ce3ea438000000029a968de000000000000000000000000000000000000000000000000100c9bf670f2d046c32f2b194958abd36b7c71cd118ec635f0990ceac863e9350f1de6600000159b144454c00000190922227de00000008f6f1847400000005a845a7b0ce3ea43800000006bb1b06540000000000000000000000000000000000000000000000017909da03")
	assetConfigCell, _ := config.GetCellFromHex("b5ee9c72410210010003d30002012001020201200304020120050602012007080184bf895668e908644f30322b997de8faaafc21f05aa52f8982f042dac1fe0b4d09d04e9fed5bfb7d79a2078297995f3d85b4badeac8c0d9eab82d3751bf9bc92754a09090201200a0b0184bf8a9006bd3fb03d355daeeff93b24be90afaa6e3ca0073ff5720f8a852c933278ff90c4242be03df8242edc274e9d0503676ffc6d8f19ae1d4fbed137859a71bd060c0183bf748433fcbcc1ac75e54798fb9cdfd8d368b8d6ae3092f4c291cf8465590f7b14348433fcbcc1ac75e54798fb9cdfd8d368b8d6ae3092f4c291cf8465590f7b14130d0183bf6627c5eaf750e15e689006a18f136130fa2b6874a62e57f9c529bc43cfae49cf221696bd2c37ea80895e9e3b7139b528b03ded985958ba1a038fbdcaae773886130e00d419641bbc2bc000000000000001d80000000000000320000000000005f37000000000000000000000000000000000000000d18c2e280000000000000f4240000000000016e3600016345785d8a00007d000c80000000000000000000000000000000000000000000000000183bf47b22d8d0a21004209a3eeb54d9c61d63c8ef5dbc1a701ddc4311c1cacb03f8d083b317c116bed042ed86b6c6fcdd323a9c29c18dce490f5da574311698ead7e0d0f0183bf670f2d046c32f2b194958abd36b7c71cd118ec635f0990ceac863e9350f1de66e0d8ba516e25e406ead68c7185ad5dc5b8dd976d0252ee40c6c8971b6d7085c20d0f00d41ce81f40296800000000000000000000000000001c6b000000000001731800000000000000000000000000000000000000ba43b7400000000000002dc6c000000000000027100000221b262dd80007d000b400000000000000000000000000000000000000000000000000d41c201e782a3000000000000001d80000000000000320000000000005f37000000000000000000000000000000000000000d18c2e280000000000000f4240000000000016e360000000000000000007d000c800000000000000000000000000000000000000000000000000d419c81c202b5c00000000000001d80000000000000320000000000005f37000000000000000000000000000000000000000d18c2e280000000000000f4240000000000016e3600016345785d8a00007d000c800000000000000000000000000000000000000000000000000d400001d4c2bc000000000000000000000000000001c6b000000000001731800000000000000000000000000000000000000ba43b7400000000000002dc6c0000000000000271000000000000f4240271000c8000000000000000000000000000000000000000000000000b0d7b158")
	err := parser.SetInfo(assetsDataCell.AsDict(256), assetConfigCell.AsDict(256))
	if err != nil {
		t.Fatalf("parser.SetInfo err: %s", err)
	}
	return parser
}
