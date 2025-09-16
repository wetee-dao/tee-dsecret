package dsecret

import types "github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/types"

// 注册 dkg 节点
// register dkg node
func MakeRegisterNodeCall(sender0 [32]byte, p2pId1 [32]byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsDSecret: true,
		AsDSecretField0: &types.WeteeDsecretPalletCall{
			IsRegisterNode:        true,
			AsRegisterNodeSender0: sender0,
			AsRegisterNodeP2pId1:  p2pId1,
		},
	}
}

// 上传共识节点代码
// update consensus node code
func MakeUploadCodeCall(signature0 []byte, signer1 []byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsDSecret: true,
		AsDSecretField0: &types.WeteeDsecretPalletCall{
			IsUploadCode:           true,
			AsUploadCodeSignature0: signature0,
			AsUploadCodeSigner1:    signer1,
		},
	}
}

// 上传共识节点代码
// update consensus node code
func MakeUploadClusterProofCall(cid0 uint64, report1 []byte, pubs2 [][32]byte, sigs3 []types.MultiSignature) types.RuntimeCall {
	return types.RuntimeCall{
		IsDSecret: true,
		AsDSecretField0: &types.WeteeDsecretPalletCall{
			IsUploadClusterProof:        true,
			AsUploadClusterProofCid0:    cid0,
			AsUploadClusterProofReport1: report1,
			AsUploadClusterProofPubs2:   pubs2,
			AsUploadClusterProofSigs3:   sigs3,
		},
	}
}

// 上传 devloper，report hash 启动应用
func MakeWorkLaunchCall(work0 types.WorkId, report1 types.OptionTByteSlice, deployKey2 [32]byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsDSecret: true,
		AsDSecretField0: &types.WeteeDsecretPalletCall{
			IsWorkLaunch:           true,
			AsWorkLaunchWork0:      work0,
			AsWorkLaunchReport1:    report1,
			AsWorkLaunchDeployKey2: deployKey2,
		},
	}
}

// 设置节点公网服务
// set node pub server
func MakeSetNodePubServerCall(id0 uint64, server1 types.P2PAddr) types.RuntimeCall {
	return types.RuntimeCall{
		IsDSecret: true,
		AsDSecretField0: &types.WeteeDsecretPalletCall{
			IsSetNodePubServer:        true,
			AsSetNodePubServerId0:     id0,
			AsSetNodePubServerServer1: server1,
		},
	}
}

// Set boot peers
// 设置引导节点
func MakeSetBootPeersCall(boots0 []types.P2PAddr) types.RuntimeCall {
	return types.RuntimeCall{
		IsDSecret: true,
		AsDSecretField0: &types.WeteeDsecretPalletCall{
			IsSetBootPeers:       true,
			AsSetBootPeersBoots0: boots0,
		},
	}
}
func MakeSetEpochCall() types.RuntimeCall {
	return types.RuntimeCall{
		IsDSecret: true,
		AsDSecretField0: &types.WeteeDsecretPalletCall{
			IsSetEpoch: true,
		},
	}
}
func MakeSetEpochWithGovCall() types.RuntimeCall {
	return types.RuntimeCall{
		IsDSecret: true,
		AsDSecretField0: &types.WeteeDsecretPalletCall{
			IsSetEpochWithGov: true,
		},
	}
}
