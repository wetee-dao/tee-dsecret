package asset

import (
	types1 "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	types "github.com/wetee-dao/go-sdk/pallet/types"
)

// create we asset.
// 创建 WETEE 资产
func MakeCreateAssetCall(metadata0 types.AssetMeta, initAmount1 types1.U128) types.RuntimeCall {
	return types.RuntimeCall{
		IsAsset: true,
		AsAssetField0: &types.WeteeAssetsPalletCall{
			IsCreateAsset:            true,
			AsCreateAssetMetadata0:   metadata0,
			AsCreateAssetInitAmount1: initAmount1,
		},
	}
}
func MakeDeleteAssetCall(assetId0 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsAsset: true,
		AsAssetField0: &types.WeteeAssetsPalletCall{
			IsDeleteAsset:         true,
			AsDeleteAssetAssetId0: assetId0,
		},
	}
}

// Users destroy their own assets.
// 销毁资产
func MakeBurnCall(assetId0 uint64, amount1 types1.U128) types.RuntimeCall {
	return types.RuntimeCall{
		IsAsset: true,
		AsAssetField0: &types.WeteeAssetsPalletCall{
			IsBurn:         true,
			AsBurnAssetId0: assetId0,
			AsBurnAmount1:  amount1,
		},
	}
}

// This function transfers the given amount from the source to the destination.
//
// # Arguments
//
// * `amount` - The amount to transfer
// * `source` - The source account
// * `destination` - The destination account
// 转移资产
func MakeTransferCall(dest0 types.MultiAddress, assetId1 uint64, amount2 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsAsset: true,
		AsAssetField0: &types.WeteeAssetsPalletCall{
			IsTransfer:         true,
			AsTransferDest0:    dest0,
			AsTransferAssetId1: assetId1,
			AsTransferAmount2:  amount2,
		},
	}
}
func MakeParachainAssetRegisterCall(paraId0 uint32, generalKey1 []byte, metadata2 types.AssetMeta) types.RuntimeCall {
	return types.RuntimeCall{
		IsAsset: true,
		AsAssetField0: &types.WeteeAssetsPalletCall{
			IsParachainAssetRegister:            true,
			AsParachainAssetRegisterParaId0:     paraId0,
			AsParachainAssetRegisterGeneralKey1: generalKey1,
			AsParachainAssetRegisterMetadata2:   metadata2,
		},
	}
}
func MakeSetChainIdCall(paraId0 uint32) types.RuntimeCall {
	return types.RuntimeCall{
		IsAsset: true,
		AsAssetField0: &types.WeteeAssetsPalletCall{
			IsSetChainId:        true,
			AsSetChainIdParaId0: paraId0,
		},
	}
}
func MakeDeleteParachainForAssetCall(assetId0 uint64, paraId1 uint32) types.RuntimeCall {
	return types.RuntimeCall{
		IsAsset: true,
		AsAssetField0: &types.WeteeAssetsPalletCall{
			IsDeleteParachainForAsset:         true,
			AsDeleteParachainForAssetAssetId0: assetId0,
			AsDeleteParachainForAssetParaId1:  paraId1,
		},
	}
}
func MakeSetParachainForAssetCall(assetId0 uint64, paraId1 uint32, generalKey2 []byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsAsset: true,
		AsAssetField0: &types.WeteeAssetsPalletCall{
			IsSetParachainForAsset:            true,
			AsSetParachainForAssetAssetId0:    assetId0,
			AsSetParachainForAssetParaId1:     paraId1,
			AsSetParachainForAssetGeneralKey2: generalKey2,
		},
	}
}
