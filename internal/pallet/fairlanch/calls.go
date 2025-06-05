package fairlanch

import (
	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	types1 "github.com/wetee-dao/go-sdk/pallet/types"
)

// 质押 vtoken
// vtoken stake
func MakeVStakingCall(vassetId0 uint64, vamount1 types.U128) types1.RuntimeCall {
	return types1.RuntimeCall{
		IsFairlanch: true,
		AsFairlanchField0: &types1.WeteeFairlanchPalletCall{
			IsVStaking:          true,
			AsVStakingVassetId0: vassetId0,
			AsVStakingVamount1:  vamount1,
		},
	}
}

// 取消 vtoken 质押
func MakeVUnstakingCall(vassetId0 uint64, amount1 types.U128) types1.RuntimeCall {
	return types1.RuntimeCall{
		IsFairlanch: true,
		AsFairlanchField0: &types1.WeteeFairlanchPalletCall{
			IsVUnstaking:          true,
			AsVUnstakingVassetId0: vassetId0,
			AsVUnstakingAmount1:   amount1,
		},
	}
}

// 设置 economic 质押比例
func MakeSetEconomicsCall(assetId0 uint64, rewardRate1 byte, quota2 types.U128) types1.RuntimeCall {
	return types1.RuntimeCall{
		IsFairlanch: true,
		AsFairlanchField0: &types1.WeteeFairlanchPalletCall{
			IsSetEconomics:            true,
			AsSetEconomicsAssetId0:    assetId0,
			AsSetEconomicsRewardRate1: rewardRate1,
			AsSetEconomicsQuota2:      quota2,
		},
	}
}
func MakeRegisterVtokenCall(vassetId0 uint64, assetId1 uint64, vassetPool2 types.U128, assetPool3 types.U128) types1.RuntimeCall {
	return types1.RuntimeCall{
		IsFairlanch: true,
		AsFairlanchField0: &types1.WeteeFairlanchPalletCall{
			IsRegisterVtoken:            true,
			AsRegisterVtokenVassetId0:   vassetId0,
			AsRegisterVtokenAssetId1:    assetId1,
			AsRegisterVtokenVassetPool2: vassetPool2,
			AsRegisterVtokenAssetPool3:  assetPool3,
		},
	}
}
func MakeSetVtokenRateCall(vassetId0 uint64, assetId1 uint64, vassetPool2 types.U128, assetPool3 types.U128) types1.RuntimeCall {
	return types1.RuntimeCall{
		IsFairlanch: true,
		AsFairlanchField0: &types1.WeteeFairlanchPalletCall{
			IsSetVtokenRate:            true,
			AsSetVtokenRateVassetId0:   vassetId0,
			AsSetVtokenRateAssetId1:    assetId1,
			AsSetVtokenRateVassetPool2: vassetPool2,
			AsSetVtokenRateAssetPool3:  assetPool3,
		},
	}
}
func MakeDeleteEconomicsCall(assetId0 uint64) types1.RuntimeCall {
	return types1.RuntimeCall{
		IsFairlanch: true,
		AsFairlanchField0: &types1.WeteeFairlanchPalletCall{
			IsDeleteEconomics:         true,
			AsDeleteEconomicsAssetId0: assetId0,
		},
	}
}
func MakeSetEpochCall() types1.RuntimeCall {
	return types1.RuntimeCall{
		IsFairlanch: true,
		AsFairlanchField0: &types1.WeteeFairlanchPalletCall{
			IsSetEpoch: true,
		},
	}
}
