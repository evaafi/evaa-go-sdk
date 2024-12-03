package config

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"math/big"
)

const (
	OpcodeSupply         = 0x1
	OpcodeWithdraw       = 0x2
	OpcodeLiquidate      = 0x3
	OpcodeJettonTransfer = 0xf8a7ea5
	OpcodeOnchainGetter  = 0x9998
)

const (
	FactorScale                        = 1e12
	AssetCoefficientScale              = 10_000
	AssetPriceScale                    = 1e9
	AssetReserveFactorScale            = 10_000
	AssetLiquidationReserveFactorScale = 10_000
	AssetOriginationFeeScale           = 1e9
	AssetLiquidationThresholdScale     = 10_000
	AssetLiquidationBonusScale         = 10_000
	AssetSRateScale                    = 1e12
	AssetBRateScale                    = 1e12
	CollateralWorthThreshold           = 100 * AssetPriceScale // literally 100$
)

const (
	FeeSupply               = 3e8
	FeeWithdraw             = 35e7
	FeeSupplyJetton         = 35e7
	FeeSupplyJettonFWD      = 3e8
	FeeLiquidation          = 8e8
	FeeLiquidationJetton    = 1
	FeeLiquidationJettonFWD = 8e8
)

const (
	CodeLending               = "b5ee9c72c1010e0100fd000d12182a555a6065717691969efd0114ff00f4a413f4bcf2c80b010202c8050202039f740403001ff2f8276a2687d2018fd201800f883b840051d38642c678b64e4400780e58fc10802faf07f80e59fa801e78b096664c02078067c07c100627a7978402014807060007a0ddb0c60201c709080013a0fd007a026900aa90400201200b0a0031b8e1002191960aa00b9e2ca007f4042796d225e8019203f6010201200d0c000bf7c147d2218400b9d10e86981fd201840b07f8138d809797976a2687d2029116382f970fd9178089910374daf81b619fd20182c7883b8701981684100627910eba56001797a6a6ba610fd8200e8768f76a9f6aa00cc2a32a8292878809bef2f1889f883bbcdeb86f01"
	CodeWalletUSDT            = "b5ee9c72010101010023000842028f452d7a4dfd74066b682365177259ed05734435be76b5fd4bd5d8af2b7c3d68"
	CodeWalletTSTON           = "b5ee9c724101010100230008420212bebb0dc8e202b7e26f721e2547e16bb9ebaec934f657d19f22e76d62bec8786eb568d4"
	CodeWalletSTTON           = "b5ee9c7201021201000362000114ff00f4a413f4bcf2c80b0102016202030202cc0405001ba0f605da89a1f401f481f481a8610201d40607020120090a01cf0831c02497c138007434c0c05c6c2544d7c0fc03783e903e900c7e800c5c75c87e800c7e800c1cea6d0000b4c7c8608403e29fa96ea54c4d167c02b808608405e351466ea58c511100fc02f80860841657c1ef2ea54c4d167c03380517c1300138c08c2103fcbc200800113e910c30003cb85360007ced44d0fa00fa40fa40d43010235f03018208989680a16d801072226eb32091719170e203c8cb055006cf165004fa02cb6a039358cc019130e201c901fb000201580b0c020148101101f100f4cffe803e90087c007b51343e803e903e90350c144da8548ab1c17cb8b04a30bffcb8b0950d109c150804d50500f214013e809633c58073c5b33248b232c044bd003d0032c032483e401c1d3232c0b281f2fff274013e903d010c7e800835d270803cb8b11de0063232c1540233c59c3e8085f2dac4f3200d02f73b51343e803e903e90350c0234cffe80145468017e903e9014d6f1c1551cdb5c150804d50500f214013e809633c58073c5b33248b232c044bd003d0032c0327e401c1d3232c0b281f2fff274140371c1472c7cb8b0c2be80146a2860822625a019ad8228608239387028062849e5c412440e0dd7c138c34975c2c0600e0f009e8210178d4519c8cb1f19cb3f5007fa0222cf165006cf1625fa025003cf16c95005cc2391729171e25008a813a0820a625a00a014bcf2e2c504c98040fb001023c85004fa0258cf1601cf16ccc9ed5400705279a018a182107362d09cc8cb1f5230cb3f58fa025007cf165007cf16c9718010c8cb0524cf165006fa0215cb6a14ccc971fb0010241023007cc30023c200b08e218210d53276db708010c8cb055008cf165004fa0216cb6a12cb1f12cb3fc972fb0093356c21e203c85004fa0258cf1601cf16ccc9ed5400c90c3b51343e803e903e90350c01b4cffe800c145128548df1c17cb8b04970bffcb8b0812082e4e1c02fbcb8b160841ef765f7b232c7c532cfd63e808873c5b25c60063232c14933c59c3e80b2dab33260103ec01004f214013e809633c58073c5b3327b55200081200835c87b51343e803e903e90350c0134c7c8608405e351466e80a0841ef765f7ae84ac7cb83234cfcc7e800c04e81408f214013e809633c58073c5b3327b5520"
	CodeWalletSTTONTestnet    = "b5ee9c7201021201000362000114ff00f4a413f4bcf2c80b0102016202030202cc0405001ba0f605da89a1f401f481f481a8610201d40607020120090a01cf0831c02497c138007434c0c05c6c2544d7c0fc03383e903e900c7e800c5c75c87e800c7e800c1cea6d0000b4c7c8608403e29fa96ea54c4d167c027808608405e351466ea58c511100fc02b80860841657c1ef2ea54c4d167c02f80517c1300138c08c2103fcbc200800113e910c30003cb85360007ced44d0fa00fa40fa40d43010235f03018208989680a16d801072226eb32091719170e203c8cb055006cf165004fa02cb6a039358cc019130e201c901fb000201200b0c0081d40106b90f6a2687d007d207d206a1802698f90c1080bc6a28cdd0141083deecbef5d0958f97064699f98fd001809d02811e428027d012c678b00e78b6664f6aa401f1503d33ffa00fa4021f001ed44d0fa00fa40fa40d4305136a1522ac705f2e2c128c2fff2e2c254344270542013541403c85004fa0258cf1601cf16ccc922c8cb0112f400f400cb00c920f9007074c8cb02ca07cbffc9d004fa40f40431fa0020d749c200f2e2c4778018c8cb055008cf1670fa0217cb6b13cc80d0201200e0f009e8210178d4519c8cb1f19cb3f5007fa0222cf165006cf1625fa025003cf16c95005cc2391729171e25008a813a0820a625a00a014bcf2e2c504c98040fb001023c85004fa0258cf1601cf16ccc9ed5402f73b51343e803e903e90350c0234cffe80145468017e903e9014d6f1c1551cdb5c150804d50500f214013e809633c58073c5b33248b232c044bd003d0032c0327e401c1d3232c0b281f2fff274140371c1472c7cb8b0c2be80146a2860822625a019ad8228608239387028062849e5c412440e0dd7c138c34975c2c060101100c90c3b51343e803e903e90350c01b4cffe800c145128548df1c17cb8b04970bffcb8b0812082e4e1c02fbcb8b160841ef765f7b232c7c532cfd63e808873c5b25c60063232c14933c59c3e80b2dab33260103ec01004f214013e809633c58073c5b3327b552000705279a018a182107362d09cc8cb1f5230cb3f58fa025007cf165007cf16c9718010c8cb0524cf165006fa0215cb6a14ccc971fb0010241023007cc30023c200b08e218210d53276db708010c8cb055008cf165004fa0216cb6a12cb1f12cb3fc972fb0093356c21e203c85004fa0258cf1601cf16ccc9ed54"
	CodeWalletStandard        = "b5ee9c7201021301000385000114ff00f4a413f4bcf2c80b0102016202030202cb0405001ba0f605da89a1f401f481f481a9a30201ce06070201580a0b02f70831c02497c138007434c0c05c6c2544d7c0fc07783e903e900c7e800c5c75c87e800c7e800c1cea6d0000b4c7c076cf16cc8d0d0d09208403e29fa96ea68c1b088d978c4408fc06b809208405e351466ea6cc1b08978c840910c03c06f80dd6cda0841657c1ef2ea7c09c6c3cb4b01408eebcb8b1807c073817c160080900113e910c30003cb85360005c804ff833206e953080b1f833de206ef2d29ad0d30731d3ffd3fff404d307d430d0fa00fa00fa00fa00fa00fa00300008840ff2f00201580c0d020148111201f70174cfc0407e803e90087c007b51343e803e903e903534544da8548b31c17cb8b04ab0bffcb8b0950d109c150804d50500f214013e809633c58073c5b33248b232c044bd003d0032c032481c007e401d3232c084b281f2fff274013e903d010c7e800835d270803cb8b13220060072c15401f3c59c3e809dc072dae00e02f33b51343e803e903e90353442b4cfc0407e80145468017e903e9014d771c1551cdbdc150804d50500f214013e809633c58073c5b33248b232c044bd003d0032c0325c007e401d3232c084b281f2fff2741403f1c147ac7cb8b0c33e801472a84a6d8206685401e8062849a49b1578c34975c2c070c00870802c200f1000aa13ccc88210178d4519580a02cb1fcb3f5007fa0222cf165006cf1625fa025003cf16c95005cc2391729171e25007a813a008aa005004a017a014bcf2e2c501c98040fb004300c85004fa0258cf1601cf16ccc9ed5400725269a018a1c882107362d09c2902cb1fcb3f5007fa025004cf165007cf16c9c8801001cb0527cf165004fa027101cb6a13ccc971fb0050421300748e23c8801001cb055006cf165005fa027001cb6a8210d53276db580502cb1fcb3fc972fb00925b33e24003c85004fa0258cf1601cf16ccc9ed5400eb3b51343e803e903e9035344174cfc0407e800870803cb8b0be903d01007434e7f440745458a8549631c17cb8b049b0bffcb8b0b220841ef765f7960100b2c7f2cfc07e8088f3c58073c584f2e7f27220060072c148f3c59c3e809c4072dab33260103ec01004f214013e809633c58073c5b3327b55200087200835c87b51343e803e903e9035344134c7c06103c8608405e351466e80a0841ef765f7ae84ac7cbd34cfc04c3e800c04e81408f214013e809633c58073c5b3327b5520"
	CodeWalletStandardTestnet = "b5ee9c7241021101000323000114ff00f4a413f4bcf2c80b0102016202030202cc0405001ba0f605da89a1f401f481f481a8610201d40607020120080900c30831c02497c138007434c0c05c6c2544d7c0fc03383e903e900c7e800c5c75c87e800c7e800c1cea6d0000b4c7e08403e29fa954882ea54c4d167c0278208405e3514654882ea58c511100fc02b80d60841657c1ef2ea4d67c02f817c12103fcbc2000113e910c1c2ebcb853600201200a0b0083d40106b90f6a2687d007d207d206a1802698fc1080bc6a28ca9105d41083deecbef09dd0958f97162e99f98fd001809d02811e428027d012c678b00e78b6664f6aa401f1503d33ffa00fa4021f001ed44d0fa00fa40fa40d4305136a1522ac705f2e2c128c2fff2e2c254344270542013541403c85004fa0258cf1601cf16ccc922c8cb0112f400f400cb00c920f9007074c8cb02ca07cbffc9d004fa40f40431fa0020d749c200f2e2c4778018c8cb055008cf1670fa0217cb6b13cc80c0201200d0e009e8210178d4519c8cb1f19cb3f5007fa0222cf165006cf1625fa025003cf16c95005cc2391729171e25008a813a08209c9c380a014bcf2e2c504c98040fb001023c85004fa0258cf1601cf16ccc9ed5402f73b51343e803e903e90350c0234cffe80145468017e903e9014d6f1c1551cdb5c150804d50500f214013e809633c58073c5b33248b232c044bd003d0032c0327e401c1d3232c0b281f2fff274140371c1472c7cb8b0c2be80146a2860822625a019ad822860822625a028062849e5c412440e0dd7c138c34975c2c0600f1000d73b51343e803e903e90350c01f4cffe803e900c145468549271c17cb8b049f0bffcb8b08160824c4b402805af3cb8b0e0841ef765f7b232c7c572cfd400fe8088b3c58073c5b25c60063232c14933c59c3e80b2dab33260103ec01004f214013e809633c58073c5b3327b552000705279a018a182107362d09cc8cb1f5230cb3f58fa025007cf165007cf16c9718010c8cb0524cf165006fa0215cb6a14ccc971fb0010241023007cc30023c200b08e218210d53276db708010c8cb055008cf165004fa0216cb6a12cb1f12cb3fc972fb0093356c21e203c85004fa0258cf1601cf16ccc9ed5495eaedd7"
	CodeWalletTonUsdtDeDust   = "b5ee9c7241021201000334000114ff00f4a413f4bcf2c80b0102016202110202cc03060201d4040500c30831c02497c138007434c0c05c6c2544d7c0fc02f83e903e900c7e800c5c75c87e800c7e800c1cea6d0000b4c7e08403e29fa954882ea54c4d167c0238208405e3514654882ea58c511100fc02780d60841657c1ef2ea4d67c02b817c12103fcbc2000113e910c1c2ebcb85360020148070e020120080a01f100f4cffe803e90087c007b51343e803e903e90350c144da8548ab1c17cb8b04a30bffcb8b0950d109c150804d50500f214013e809633c58073c5b33248b232c044bd003d0032c032483e401c1d3232c0b281f2fff274013e903d010c7e800835d270803cb8b11de0063232c1540233c59c3e8085f2dac4f3200900ae8210178d4519c8cb1f19cb3f5007fa0222cf165006cf1625fa025003cf16c95005cc2391729171e25008a813a08208e4e1c0aa008208989680a0a014bcf2e2c504c98040fb001023c85004fa0258cf1601cf16ccc9ed5403f73b51343e803e903e90350c0234cffe80145468017e903e9014d6f1c1551cdb5c150804d50500f214013e809633c58073c5b33248b232c044bd003d0032c0327e401c1d3232c0b281f2fff274140371c1472c7cb8b0c2be80146a2860822625a020822625a004ad8228608239387028062849f8c3c975c2c070c008e00b0c0d00705279a018a182107362d09cc8cb1f5230cb3f58fa025007cf165007cf16c9718010c8cb0524cf165006fa0215cb6a14ccc971fb0010241023000e10491038375f040076c200b08e218210d53276db708010c8cb055008cf165004fa0216cb6a12cb1f12cb3fc972fb0093356c21e203c85004fa0258cf1601cf16ccc9ed540201200f1000db3b51343e803e903e90350c01f4cffe803e900c145468549271c17cb8b049f0bffcb8b0a0823938702a8005a805af3cb8b0e0841ef765f7b232c7c572cfd400fe8088b3c58073c5b25c60063232c14933c59c3e80b2dab33260103ec01004f214013e809633c58073c5b3327b55200083200835c87b51343e803e903e90350c0134c7e08405e3514654882ea0841ef765f784ee84ac7cb8b174cfcc7e800c04e81408f214013e809633c58073c5b3327b5520001ba0f605da89a1f401f481f481a861f0a7d84b"
	CodeWalletTonStorm        = "b5ee9c7241021301000380000114ff00f4a413f4bcf2c80b0102016202100202cd030f04add106380792000e8698180b8d8474289af81ed9e707d207d2018fd0018b8eb90fd0018fd001839d4da0001698f90c10807c53f52dd4742989a2ced9e7010c1080bc6a28cdd474318a22201ed9e701a9041082caf83de5d40405080c00888020d721ed44d0fa00fa40fa40d74c04d31f8200fff0228210178d4519ba0382107bdd97deba13b112f2f4d33f31fa003013a05023c85004fa0258cf1601cf16ccc9ed5401f603d33ffa00fa4021f007ed44d0fa00fa40fa40d74c5136a1522ac705f2e2c128c2fff2e2c2545255705202541304c85004fa0258cf1601cf16ccc921c8cb0113f40012f400cb00c97021f90074c8cb0212cb07cbffc9d004fa40f40431fa0020d749c200f2e2c48210178d4519c8cb1f1acb3f5008fa0223cf1601060176cf1626fa025007cf162591729171e2500aa815a0820a43d580a016bcf2e2c57007c9103741408040db3c4300c85004fa0258cf1601cf16ccc9ed5407002e778018c8cb055005cf165005fa0213cb6bccccc901fb0002f6ed44d0fa00fa40fa40d74c08d33ffa005151a005fa40fa40547c25705202541304c85004fa0258cf1601cf16ccc921c8cb0113f40012f400cb00c97021f90074c8cb0212cb07cbffc9d031536cc7050dc7051cb1f2e2c30afa0051a8a12195104a395f04e30d048208989680b60972fb0225d70b01c30003c20013090a014c521aa018a182107362d09cc8cb1f5240cb3f5003fa0201cf165008cf16c954253071db3c10350e0158b08e948210d53276dbc8cb1f14cb3f704055810082db3c923333e25003c85004fa0258cf1601cf16ccc9ed540b0030708018c8cb055004cf165004fa0212cb6a01cf17c901fb0002a08e843059db3ce06c22ed44d0fa00fa40fa40d74c10235f030182106d8e5e3cba8ea75202c705f2e2c1820898968070fb0201d70b3f8210d53276dbc8cb1fcb3f7001c912810082db3ce05f03840ff2f00d0e01bc30ed44d0fa00fa40fa40d74c06d33ffa00fa40305151a15248c705f2e2c126c2fff2e2c20582100ee6b280b9f2d2c782107bdd97dec8cb1fcb3f5004fa0221cf1658cf167001c952308040db3c4013c85004fa0258cf1601cf16ccc9ed540e002c718018c8cb055004cf165004fa0212cb6accc901fb000011f7d22186000797163402016611120021b5fe7da89a1f401f481f481ae9826be070001bb7605da89a1f401f481f481ae990c6f31b9d"
	CodeWalletUsdtStorm       = CodeWalletTonStorm
)

func getCellFromHex(lendingCode string) *cell.Cell {
	codeCell, err := GetCellFromHex(lendingCode)
	if err != nil {
		panic(err)
	}
	return codeCell
}

func GetCellFromHex(data string) (*cell.Cell, error) {
	boc, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}
	code, err := cell.FromBOC(boc)
	if err != nil {
		return nil, err
	}
	return code, nil
}

const (
	MasterMainnet = "EQC8rUZqR_pWV1BylWUlPNBzyiTYVoBEmQkMIQDZXICfnuRr"
	MasterTestnet = "EQDLsg3w-iBj26Gww7neYoJAxiT2t77Zo8ro56b0yuHsPp3C"
	LpMainnet     = "EQBIlZX2URWkXCSg3QF2MJZU-wC5XkBoLww-hdWk2G37Jc6N"
)

const (
	MainnetVersion = 6
	TestnetVersion = 1
	LpVersion      = 3
)

type Asset string

const (
	TON            Asset = "TON"
	USDT           Asset = "USDT"
	JUSDT          Asset = "jUSDT"
	JUSDC          Asset = "jUSDC"
	STTON          Asset = "stTON"
	TSTON          Asset = "tsTON"
	TONUSDT_DEDUST Asset = "TONUSDT_DEDUST"
	TONUSDT_STONFI Asset = "TONUSDT_STONFI"
	TON_STORM      Asset = "TON_STORM"
	USDT_STORM     Asset = "USDT_STORM"
)

func (a Asset) Sha256Hash() *big.Int {
	h := sha256.New()
	_, err := h.Write([]byte(a))
	if err != nil {
		panic(err)
	}
	return big.NewInt(0).SetBytes(h.Sum(nil))
}

func (a Asset) ID() string {
	return a.Sha256Hash().String()
}

const (
	USDTJettonAddress  = "EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"
	JUSDTJettonAddress = "EQBynBO23ywHy_CgarY9NK9FTz0yDsG82PtcbSTQgGoXwiuA"
	JUSDCJettonAddress = "EQB-MPwrd1G6WKNkLz_VnV6WqBDd142KMQv-g1O-8QUA3728"
	STTONJettonAddress = "EQDNhy-nxYFgUqzfUzImBEP67JqsyMIcyk2S5_RwNNEYku0k"
	TSTONJettonAddress = "EQC98_qAmNEptUtPc7W6xdHh_ZHrBUFpw5Ft_IzNU20QAJav"

	TONUSDTDeDustJettonAddress = "EQA-X_yo3fzzbDbJ_0bzFWKqtRuZFIRa1sJsveZJ1YpViO3r"
	TONUSDTStonfiJettonAddress = "EQD8TJ8xEWB1SpnRE4d89YO3jl0W0EiBnNS4IBaHaUmdfizE"
	TONStormJettonAddress      = "EQCNY2AQ3ZDYwJAqx_nzl9i9Xhd_Ex7izKJM6JTxXRnO6n1F"
	USDTStormJettonAddress     = "EQCup4xxCulCcNwmOocM9HtDYPU8xe0449tQLp6a-5BLEegW"

	JUSDTJettonAddressTestnet = "kQBe4gtSQMxM5RpMYLr4ydNY72F8JkY-icZXG1NJcsju8XM7"
	JUSDCJettonAddressTestnet = "kQDaY5yUatYnHei73HBqRX_Ox9LK2XnR7XuCY9MFC2INbfYI"
	STTONJettonAddressTestnet = "kQC3Duw3dg8k98xf5S7Bm7YOWVJ5QW8hm3iLqFfJfa_g9h07"
)
