package wesudo

import types "github.com/wetee-dao/go-sdk/pallet/types"

// Execute external transactions as root
// 以 root 账户执行函数
func MakeSudoCall(daoId0 uint64, call1 types.RuntimeCall) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeSudo: true,
		AsWeSudoField0: &types.WeteeSudoPalletCall{
			IsSudo:       true,
			AsSudoDaoId0: daoId0,
			AsSudoCall1:  &call1,
		},
	}
}

// set sudo account
// 设置超级用户账户
func MakeSetSudoAccountCall(daoId0 uint64, sudoAccount1 [32]byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeSudo: true,
		AsWeSudoField0: &types.WeteeSudoPalletCall{
			IsSetSudoAccount:             true,
			AsSetSudoAccountDaoId0:       daoId0,
			AsSetSudoAccountSudoAccount1: sudoAccount1,
		},
	}
}

// close sudo
// 关闭 sudo 功能
func MakeCloseSudoCall(daoId0 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeSudo: true,
		AsWeSudoField0: &types.WeteeSudoPalletCall{
			IsCloseSudo:       true,
			AsCloseSudoDaoId0: daoId0,
		},
	}
}
