package store

import (
	types1 "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	types "github.com/wetee-dao/go-sdk/pallet/types"
)

func MakeRegisterAppCall(name0 []byte, meta1 []byte, ty2 types.AppType, images3 []types.Image, run4 types.TEEVersion) types.RuntimeCall {
	return types.RuntimeCall{
		IsStore: true,
		AsStoreField0: &types.WeteeStorePalletCall{
			IsRegisterApp:        true,
			AsRegisterAppName0:   name0,
			AsRegisterAppMeta1:   meta1,
			AsRegisterAppTy2:     ty2,
			AsRegisterAppImages3: images3,
			AsRegisterAppRun4:    run4,
		},
	}
}
func MakeUnregisterAppCall(appId0 types1.U128) types.RuntimeCall {
	return types.RuntimeCall{
		IsStore: true,
		AsStoreField0: &types.WeteeStorePalletCall{
			IsUnregisterApp:       true,
			AsUnregisterAppAppId0: appId0,
		},
	}
}
func MakeAddAppVersionCall(appId0 types1.U128, version1 uint16, images2 []types.Image) types.RuntimeCall {
	return types.RuntimeCall{
		IsStore: true,
		AsStoreField0: &types.WeteeStorePalletCall{
			IsAddAppVersion:         true,
			AsAddAppVersionAppId0:   appId0,
			AsAddAppVersionVersion1: version1,
			AsAddAppVersionImages2:  images2,
		},
	}
}
func MakeReviewAppCall(appId0 types1.U128, tokenId1 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsStore: true,
		AsStoreField0: &types.WeteeStorePalletCall{
			IsReviewApp:         true,
			AsReviewAppAppId0:   appId0,
			AsReviewAppTokenId1: tokenId1,
		},
	}
}

// Set boot peers
// 设置引导节点
func MakeInitMintCall() types.RuntimeCall {
	return types.RuntimeCall{
		IsStore: true,
		AsStoreField0: &types.WeteeStorePalletCall{
			IsInitMint: true,
		},
	}
}
