package matrix

import (
	types1 "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	types "github.com/wetee-dao/go-sdk/pallet/types"
)

// Create a N o de
// 从一个通证池,创建一个节点
func MakeCreateNodeCall(name0 []byte, desc1 []byte, purpose2 []byte, metaData3 []byte, imApi4 []byte, bg5 []byte, logo6 []byte, img7 []byte, homeUrl8 []byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsMatrix: true,
		AsMatrixField0: &types.WeteeMatrixPalletCall{
			IsCreateNode:          true,
			AsCreateNodeName0:     name0,
			AsCreateNodeDesc1:     desc1,
			AsCreateNodePurpose2:  purpose2,
			AsCreateNodeMetaData3: metaData3,
			AsCreateNodeImApi4:    imApi4,
			AsCreateNodeBg5:       bg5,
			AsCreateNodeLogo6:     logo6,
			AsCreateNodeImg7:      img7,
			AsCreateNodeHomeUrl8:  homeUrl8,
		},
	}
}

// update node info
// 更新节点信息
func MakeUpdateNodeCall(nodeId0 types1.U128, name1 types.OptionTByteSlice, desc2 types.OptionTByteSlice, purpose3 types.OptionTByteSlice, metaData4 types.OptionTByteSlice, imApi5 types.OptionTByteSlice, bg6 types.OptionTByteSlice, logo7 types.OptionTByteSlice, img8 types.OptionTByteSlice, homeUrl9 types.OptionTByteSlice, status10 types.OptionTStatus1) types.RuntimeCall {
	return types.RuntimeCall{
		IsMatrix: true,
		AsMatrixField0: &types.WeteeMatrixPalletCall{
			IsUpdateNode:          true,
			AsUpdateNodeNodeId0:   nodeId0,
			AsUpdateNodeName1:     name1,
			AsUpdateNodeDesc2:     desc2,
			AsUpdateNodePurpose3:  purpose3,
			AsUpdateNodeMetaData4: metaData4,
			AsUpdateNodeImApi5:    imApi5,
			AsUpdateNodeBg6:       bg6,
			AsUpdateNodeLogo7:     logo7,
			AsUpdateNodeImg8:      img8,
			AsUpdateNodeHomeUrl9:  homeUrl9,
			AsUpdateNodeStatus10:  status10,
		},
	}
}
