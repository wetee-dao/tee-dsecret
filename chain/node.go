package chain

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/wetee-dao/go-sdk/gen/types"
	"github.com/wetee-dao/go-sdk/gen/weteedsecret"
	"wetee.app/dsecret/util"
)

// RegisterNode register node
// 注册节点
func (c *Chain) RegisterNode(signer *signature.KeyringPair, pubkey []byte) error {
	var bt [32]byte
	copy(bt[:], pubkey)
	call := weteedsecret.MakeRegisterNodeCall(bt)
	return c.client.SignAndSubmit(signer, call, true)
}

// GetNodeList get node list
// 获取节点列表
func (c *Chain) GetNodeList() ([]types.Node, error) {
	ret, err := c.client.QueryMapAll("WeTEEDsecret", "Nodes")
	if err != nil {
		return nil, err
	}

	nodes := make([]types.Node, 0)
	for _, elem := range ret {
		for _, change := range elem.Changes {
			var n types.Node
			if err := codec.Decode(change.StorageData, &n); err != nil {
				util.LogWithRed("codec.Decode", err)
				continue
			}
			nodes = append(nodes, n)
		}
	}

	return nodes, nil
}

// 获取全网当前程序的代码版本
// Get CodeMrenclave
func (c *Chain) GetCodeMrenclave() ([]byte, error) {
	return weteedsecret.GetCodeMrenclaveLatest(c.client.Api.RPC.State)
}

// 获取全网当前程序的签名人
// Get CodeMrsigner
func (c *Chain) GetCodeMrsigner() ([]byte, error) {
	return weteedsecret.GetCodeMrsignerLatest(c.client.Api.RPC.State)
}
