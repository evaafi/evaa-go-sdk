package principal

import (
	"encoding/hex"
	"math/big"
	"reflect"
	"testing"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/tvm/cell"

	"github.com/evaafi/evaa-go-sdk/config"
)

func TestCalculateUserSCAddress(t *testing.T) {
	service := NewService(config.GetMainMainnetConfig())
	userSCAddress, err := service.CalculateUserSCAddress(address.MustParseAddr("UQBlB6eFlc-to_YqCabuBtSWFyY8uYm7Y6G39ADdiKvzi389"))
	if err != nil {
		t.Fatalf("failed to CalculateUserSCAddress, err: %s", err)
	}
	if userSCAddress.String() != "EQBHgCET1SV9Y2_dbnJBDczB4eIUvelocCsuQf3tAelZCniF" {
		t.Errorf("userSCAddress want %s, got %s", "EQBHgCET1SV9Y2_dbnJBDczB4eIUvelocCsuQf3tAelZCniF", userSCAddress.String())
	}
}

func TestUserCS_SetAccData(t *testing.T) {
	masterAddress := address.MustParseAddr("EQC8rUZqR_pWV1BylWUlPNBzyiTYVoBEmQkMIQDZXICfnuRr")
	userAddress := address.MustParseAddr("UQBlB6eFlc-to_YqCabuBtSWFyY8uYm7Y6G39ADdiKvzi389")
	//dataBoc, err := hex.DecodeString("b5ee9c7241020701000111000299106801795a8cd48ff4acaea0e52aca4a79a0e79449b0ad00893212184201b2b9013f3d0001aafb2a1e6a6ad3c6137e7af73ce8a8a63086c27f1edb9cbcebbeaaf70ccb6aa00000000000000012010202012003040063a0092acd1d210c89e60645732fbd1f555f843e0b54a5f1305e085b583fc169a13a000000000000000000000000000000001002012005060053bfe548035e9fd81e9aaed777fc9d925f4857d5371e50039ffab907c5429649993c7ffffffffffff629400052bfb313e2f57ba870af34480350c789b0987d15b43a53172bfce294de21e7d724e70000000000135cbd0052bf895668e908644f30322b997de8faaafc21f05aa52f8982f042dac1fe0b4d09d00000005d54c31d38bf59b66f")
	dataBoc, err := hex.DecodeString("b5ee9c7201020f010001e9000299106801795a8cd48ff4acaea0e52aca4a79a0e79449b0ad00893212184201b2b9013f3d001941e9e16573eb68fd8a8269bb81b52585c98f2e626ed8e86dfd0037622afce2e0000000000000001201020201200304020120090a02012005060053bfe548035e9fd81e9aaed777fc9d925f4857d5371e50039ffab907c5429649993c7ffffb739c3d3cdac002012007080052bf895668e908644f30322b997de8faaafc21f05aa52f8982f042dac1fe0b4d09d00001c3b91faab2470051bf748433fcbcc1ac75e54798fb9cdfd8d368b8d6ae3092f4c291cf8465590f7b14000cdb460a300f750051bf6627c5eaf750e15e689006a18f136130fa2b6874a62e57f9c529bc43cfae49ce000af9207f047f710201200b0c0063bfe548035e9fd81e9aaed777fc9d925f4857d5371e50039ffab907c5429649993c00000000000000000000000000000000400201200d0e0062bf895668e908644f30322b997de8faaafc21f05aa52f8982f042dac1fe0b4d09d0000000000000000000000000000000000061bf748433fcbcc1ac75e54798fb9cdfd8d368b8d6ae3092f4c291cf8465590f7b14000000000000000000000000000000010061bf6627c5eaf750e15e689006a18f136130fa2b6874a62e57f9c529bc43cfae49ce00000000000000000000000000000001")
	if err != nil {
		t.Fatalf("%s", err)
	}
	data, err := cell.FromBOC(dataBoc)
	if err != nil {
		t.Fatalf("%s", err)
	}
	userSCAddress := address.MustParseAddr("EQBHgCET1SV9Y2_dbnJBDczB4eIUvelocCsuQf3tAelZCniF")
	user := NewUserSC(userSCAddress)
	_, err = user.SetAccData(data)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if user.CodeVersion() != 6 {
		t.Errorf("CodeVersion want %d, got %d", 6, user.CodeVersion())
	}
	if !user.Address().Equals(userSCAddress) {
		t.Errorf("Address want %s, got %s", userSCAddress, user.Address())
	}
	if !user.UserAddress().Equals(userAddress) {
		t.Errorf("UserAddress want %s, got %s", userSCAddress, user.Address())
	}
	if !user.MasterAddress().Equals(masterAddress) {
		t.Errorf("UserAddress want %s, got %s", masterAddress, user.MasterAddress())
	}
	if !user.UserAddress().Equals(userAddress) {
		t.Errorf("UserAddress want %s, got %s", userAddress, user.UserAddress())
	}
	principalMap := map[string]*big.Int{
		"11876925370864614464799087627157805050745321306404563164673853337929163193738": big.NewInt(1809396792821690),
		"23103091784861387372100043848078515239542568751939923972799733728526040769767": big.NewInt(1544333866188728),
		"33171510858320790266247832496974106978700190498800858393089426423762035476944": big.NewInt(496674844357191),
		"91621667903763073563570557639433445791506232618002614896981036659302854767224": big.NewInt(-10002031281739),
	}
	if !reflect.DeepEqual(user.Principals(), principalMap) {
		t.Errorf("Principals want %#v, got %#v", user.Principals(), principalMap)
	}
	if user.UserState() != 0 {
		t.Errorf("CodeVersion want %d, got %d", 6, user.UserState())
	}
	rewardsKV, err := user.Rewards().LoadAll()
	if err != nil {
		t.Errorf("Rewards().LoadAll(), err %s", err)
	}
	if len(rewardsKV) != 4 {
		t.Errorf("len(rewardsKV) want %d, got %d", 4, len(rewardsKV))
	}
	//for _, kv := range rewardsKV {
	//	t.Logf("%v - %v", kv.Value, kv.Key)
	//}
}

func TestUserCS_SetData(t *testing.T) {
	masterAddress := address.MustParseAddr("EQC8rUZqR_pWV1BylWUlPNBzyiTYVoBEmQkMIQDZXICfnuRr")
	userAddress := address.MustParseAddr("UQBlB6eFlc-to_YqCabuBtSWFyY8uYm7Y6G39ADdiKvzi389")
	principalsBoc, err := hex.DecodeString("b5ee9c720101070100bc00020120010202012003040053bfe548035e9fd81e9aaed777fc9d925f4857d5371e50039ffab907c5429649993c7ffffb739c3d3cdac002012005060052bf895668e908644f30322b997de8faaafc21f05aa52f8982f042dac1fe0b4d09d00001c3b91faab2470051bf748433fcbcc1ac75e54798fb9cdfd8d368b8d6ae3092f4c291cf8465590f7b14000ce8ac3aab7cb10051bf6627c5eaf750e15e689006a18f136130fa2b6874a62e57f9c529bc43cfae49ce000af9207f047f71")
	if err != nil {
		t.Fatalf("%s", err)
	}
	principals, err := cell.FromBOC(principalsBoc)
	if err != nil {
		t.Fatalf("%s", err)
	}
	rewardsBoc, err := hex.DecodeString("b5ee9c720101070100dc00020120010202012003040063bfe548035e9fd81e9aaed777fc9d925f4857d5371e50039ffab907c5429649993c000000000000000000000000000000004002012005060062bf895668e908644f30322b997de8faaafc21f05aa52f8982f042dac1fe0b4d09d0000000000000000000000000000000000061bf748433fcbcc1ac75e54798fb9cdfd8d368b8d6ae3092f4c291cf8465590f7b14000000000000000000000000000000010061bf6627c5eaf750e15e689006a18f136130fa2b6874a62e57f9c529bc43cfae49ce00000000000000000000000000000001")
	if err != nil {
		t.Fatalf("%s", err)
	}
	rewards, err := cell.FromBOC(rewardsBoc)
	if err != nil {
		t.Fatalf("%s", err)
	}
	data := []any{
		big.NewInt(6),
		cell.BeginCell().MustStoreAddr(masterAddress).EndCell().BeginParse(),
		cell.BeginCell().MustStoreAddr(userAddress).EndCell().BeginParse(),
		principals,
		big.NewInt(0),
		rewards,
		interface{}(nil),
		interface{}(nil),
	}
	userSCAddress := address.MustParseAddr("EQBHgCET1SV9Y2_dbnJBDczB4eIUvelocCsuQf3tAelZCniF")
	user := NewUserSC(userSCAddress)
	_, err = user.SetData(ton.NewExecutionResult(data))
	if err != nil {
		t.Fatalf("%s", err)
	}
	if user.CodeVersion() != 6 {
		t.Errorf("CodeVersion want %d, got %d", 6, user.CodeVersion())
	}
	if !user.Address().Equals(userSCAddress) {
		t.Errorf("Address want %s, got %s", userSCAddress, user.Address())
	}
	if !user.UserAddress().Equals(userAddress) {
		t.Errorf("UserAddress want %s, got %s", userSCAddress, user.Address())
	}
	if !user.MasterAddress().Equals(masterAddress) {
		t.Errorf("UserAddress want %s, got %s", masterAddress, user.MasterAddress())
	}
	if !user.UserAddress().Equals(userAddress) {
		t.Errorf("UserAddress want %s, got %s", userAddress, user.UserAddress())
	}
	principalMap := map[string]*big.Int{
		"11876925370864614464799087627157805050745321306404563164673853337929163193738": big.NewInt(1816763068431960),
		"23103091784861387372100043848078515239542568751939923972799733728526040769767": big.NewInt(1544333866188728),
		"33171510858320790266247832496974106978700190498800858393089426423762035476944": big.NewInt(496674844357191),
		"91621667903763073563570557639433445791506232618002614896981036659302854767224": big.NewInt(-10002031281739),
	}
	if !reflect.DeepEqual(user.Principals(), principalMap) {
		t.Errorf("Principals want %#v, got %#v", user.Principals(), principalMap)
	}
	if user.UserState() != 0 {
		t.Errorf("CodeVersion want %d, got %d", 6, user.UserState())
	}
	rewardsKV, err := user.Rewards().LoadAll()
	if err != nil {
		t.Errorf("Rewards().LoadAll(), err %s", err)
	}
	if len(rewardsKV) != 4 {
		t.Errorf("len(rewardsKV) want %d, got %d", 4, len(rewardsKV))
	}
	//for _, kv := range rewardsKV {
	//	t.Logf("%v - %v", kv.Value, kv.Key)
	//}
}
