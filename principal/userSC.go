package principal

import (
	"github.com/evaafi/evaa-go-sdk/asset"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"maps"
	"math/big"
)

type UserSC struct {
	address *address.Address

	codeVersion   uint64
	userAddress   *address.Address
	masterAddress *address.Address
	principals    map[string]*big.Int
	userState     int64
	rewards       *cell.Dictionary
	backupCell1   *cell.Cell
	backupCell2   *cell.Cell
}

func NewUserSC(addr *address.Address) *UserSC {
	return &UserSC{address: addr, principals: map[string]*big.Int{}}
}

func (u *UserSC) Address() *address.Address {
	return u.address
}

func (u *UserSC) UserAddress() *address.Address {
	return u.userAddress
}

func (u *UserSC) MasterAddress() *address.Address {
	return u.masterAddress
}

func (u *UserSC) CodeVersion() uint64 {
	return u.codeVersion
}

func (u *UserSC) Principals() map[string]*big.Int {
	return maps.Clone(u.principals)
}

func (u *UserSC) Rewards() *cell.Dictionary {
	return u.rewards
}

func (u *UserSC) BackupCell1() *cell.Cell {
	return u.backupCell1
}

func (u *UserSC) BackupCell2() *cell.Cell {
	return u.backupCell2
}

func (u *UserSC) UserState() int64 {
	return u.userState
}

func (u *UserSC) CheckNotInDebtAtAll() bool {
	for _, principal := range u.principals {
		if principal.Sign() == -1 {
			return false
		}
	}
	return true
}

func (u *UserSC) Principal(asset string) *big.Int {
	principal, ok := u.principals[asset]
	if !ok {
		return big.NewInt(0)
	}

	return new(big.Int).Set(principal)
}

func (u *UserSC) Balance(asset string, assetData *asset.Data, applyDust bool, assetConfig *asset.Config) *big.Int {
	principal := u.Principal(asset)
	if principal.Sign() == 0 {
		return big.NewInt(0)
	}

	balance := new(big.Int)
	if principal.Sign() == 1 {
		if applyDust && assetConfig != nil && principal.Cmp(assetConfig.Dust) == -1 {
			return big.NewInt(0)
		}
		balance.Mul(principal, assetData.SRate)
	} else {
		balance.Mul(principal, assetData.BRate)
	}

	return new(big.Int).Div(balance, big.NewInt(1e12)) // FactorScale
}

func (u *UserSC) ChangePrincipal(asset string, amount *big.Int) UserBalancer {
	principals := maps.Clone(u.principals)
	principals[asset] = new(big.Int).Add(u.Principal(asset), amount)
	return &UserSC{
		address:       u.address,
		codeVersion:   u.codeVersion,
		userAddress:   u.userAddress,
		masterAddress: u.masterAddress,
		principals:    principals,
		userState:     u.userState,
	}
}

func (u *UserSC) SetAccData(userData *cell.Cell) (UserBalancer, error) {
	userSlice := userData.BeginParse()
	u.codeVersion = userSlice.MustLoadCoins()
	u.masterAddress = userSlice.MustLoadAddr()
	u.userAddress = userSlice.MustLoadAddr()
	principalsDict := userSlice.MustLoadDict(256)
	if !principalsDict.IsEmpty() {
		kvs, err := principalsDict.LoadAll()
		if err != nil {
			return nil, err
		}

		for _, kv := range kvs {
			value, err := kv.Value.LoadBigInt(64)
			if err != nil {
				return nil, err
			}

			u.principals[kv.Key.MustLoadBigUInt(256).String()] = value
		}
	}

	u.userState = userSlice.MustLoadInt(64)

	// Deprecated?
	//if (userData.bitsLeft > 32) {
	//	trackingSupplyIndex = userData.loadUintBig(64);
	//	trackingBorrowIndex = userData.loadUintBig(64);
	//	dutchAuctionStart = userData.loadUint(32);
	//	backupCell = loadMyRef(userSlice);
	//} else {
	u.rewards = userSlice.MustLoadDict(256)
	if slice := userSlice.MustLoadMaybeRef(); slice != nil {
		u.backupCell1 = slice.MustToCell()
	}
	if slice := userSlice.MustLoadMaybeRef(); slice != nil {
		u.backupCell2 = slice.MustToCell()
	}
	//}
	return u, nil
}

func (u *UserSC) SetData(userData *ton.ExecutionResult) (UserBalancer, error) {
	u.codeVersion = userData.MustInt(0).Uint64()
	u.masterAddress = userData.MustSlice(1).MustLoadAddr()
	u.userAddress = userData.MustSlice(2).MustLoadAddr()

	if !userData.MustIsNil(3) {
		if principalsDict := userData.MustCell(3).AsDict(256); !principalsDict.IsEmpty() {
			kvs, err := principalsDict.LoadAll()
			if err != nil {
				return nil, err
			}

			for _, kv := range kvs {
				value, err := kv.Value.LoadBigInt(64)
				if err != nil {
					return nil, err
				}

				u.principals[kv.Key.MustLoadBigUInt(256).String()] = value
			}
		}
	}

	u.userState = userData.MustInt(4).Int64()
	u.rewards = userData.MustCell(5).AsDict(256)
	if !userData.MustIsNil(6) {
		u.backupCell1 = userData.MustCell(6)
	}
	if !userData.MustIsNil(7) {
		u.backupCell2 = userData.MustCell(7)
	}
	//}
	return u, nil
}
