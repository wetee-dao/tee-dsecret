package chain

import (
	"crypto/ed25519"
	"errors"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/wetee-dao/go-sdk/core"
	"github.com/wetee-dao/go-sdk/pallet/dsecret"
	gtypes "github.com/wetee-dao/go-sdk/pallet/types"
	"github.com/wetee-dao/go-sdk/pallet/worker"
	"wetee.app/dsecret/util"

	types "wetee.app/dsecret/type"
)

// RegisterNode register node
// 注册节点
func (c *Chain) RegisterNode(signer *core.Signer, pubkey []byte) error {
	var bt [32]byte
	copy(bt[:], pubkey)

	call := dsecret.MakeRegisterNodeCall(bt)
	return c.SignAndSubmit(signer, call, true)
}

// GetNodes 函数用于获取节点列表，包括 Secret 节点和 Worker 节点，以及转换为自定义的 Node 类型。
func (c *Chain) GetNodes() ([][32]byte, []*gtypes.K8sCluster, []*types.Node, error) {
	// 获取节点列表
	secretNodes, err := c.GetNodeList()
	if err != nil {
		return nil, nil, nil, errors.New("Get node list error:" + err.Error())
	}
	workerNodes, err := c.GetWorkerList()
	if err != nil {
		return nil, nil, nil, errors.New("Get worker list error:" + err.Error())
	}

	nodes := []*types.Node{}
	for _, n := range secretNodes {
		var gopub ed25519.PublicKey = n[:]
		pub, _ := types.PubKeyFromStdPubKey(gopub)
		nodes = append(nodes, &types.Node{
			ID:   pub.String(),
			Type: 1,
		})
	}
	for _, w := range workerNodes {
		var gopub ed25519.PublicKey = w.Account[:]
		pub, _ := types.PubKeyFromStdPubKey(gopub)
		nodes = append(nodes, &types.Node{
			ID: pub.String(),
		})
	}

	return secretNodes, workerNodes, nodes, nil
}

// GetNodeList get node list
// 获取节点列表
func (c *Chain) GetNodeList() ([][32]byte, error) {
	ret, err := c.QueryMapAll("WeTEEDsecret", "Nodes")
	if err != nil {
		return nil, err
	}

	nodes := make([][32]byte, 0)
	for _, elem := range ret {
		for _, change := range elem.Changes {
			n := [32]byte{}
			if err := codec.Decode(change.StorageData, &n); err != nil {
				util.LogError("codec.Decode", err)
				continue
			}
			nodes = append(nodes, n)
		}
	}

	return nodes, nil
}

// GetWorkerList get worker list
// 获取矿工列表
func (c *Chain) GetWorkerList() ([]*gtypes.K8sCluster, error) {
	ret, err := c.QueryMapAll("WeTEEWorker", "K8sClusters")
	if err != nil {
		return nil, err
	}

	// 获取节点列表
	nodes := make([]*gtypes.K8sCluster, 0)
	for _, elem := range ret {
		for _, change := range elem.Changes {
			n := &gtypes.K8sCluster{}
			if err := codec.Decode(change.StorageData, n); err != nil {
				util.LogError("codec.Decode", err)
				continue
			}
			nodes = append(nodes, n)
		}
	}

	return nodes, nil
}

// GetBootPeers get boot peers
func (c *Chain) GetBootPeers() ([]gtypes.P2PAddr, error) {
	return worker.GetBootPeersLatest(c.Api.RPC.State)
}
