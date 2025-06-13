package pallets

import (
	"crypto/ed25519"
	"errors"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	chain "github.com/wetee-dao/ink.go"
	"wetee.app/dsecret/chains/pallets/generated/dsecret"
	gtypes "wetee.app/dsecret/chains/pallets/generated/types"
	"wetee.app/dsecret/internal/model"
	"wetee.app/dsecret/internal/util"
)

// RegisterNode register node
// 注册节点
func (c *Chain) RegisterNode(signer *chain.Signer, pubkey []byte) error {
	var bt [32]byte
	copy(bt[:], pubkey)

	runtimeCall := dsecret.MakeRegisterNodeCall(bt)

	call, err := (runtimeCall).AsCall()
	if err != nil {
		return errors.New("(runtimeCall).AsCall() error: " + err.Error())
	}

	return c.SignAndSubmit(signer, call, true)
}

// GetNodes 函数用于获取节点列表，包括 Secret 节点和 Worker 节点，以及转换为自定义的 Node 类型
func (c *Chain) GetNodes() ([][32]byte, []*model.Node, error) {
	// 获取节点列表
	secretNodes, err := c.GetNodeList()
	if err != nil {
		return nil, nil, errors.New("Get node list error:" + err.Error())
	}
	workerNodes, err := c.GetWorkerList()
	if err != nil {
		return nil, nil, errors.New("Get worker list error:" + err.Error())
	}

	nodes := []*model.Node{}
	for _, n := range secretNodes {
		var gopub ed25519.PublicKey = n[:]
		pub, _ := model.PubKeyFromStdPubKey(gopub)
		nodes = append(nodes, &model.Node{
			ID:   *pub,
			Type: 1,
		})
	}
	for _, w := range workerNodes {
		var gopub ed25519.PublicKey = w.Account[:]
		pub, _ := model.PubKeyFromStdPubKey(gopub)
		nodes = append(nodes, &model.Node{
			ID: *pub,
		})
	}

	return secretNodes, nodes, nil
}

// GetNodeList get node list
// 获取节点列表
func (c *Chain) GetNodeList() ([][32]byte, error) {
	ret, err := c.QueryMapAll("DSecret", "Nodes")
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
	ret, err := c.QueryMapAll("Worker", "K8sClusters")
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
func (c *Chain) GetBootPeers() ([]model.P2PAddr, error) {
	peers := []model.P2PAddr{}
	key, err := dsecret.MakeBootPeersStorageKey()
	if err != nil {
		return nil, err
	}

	var isSome bool
	isSome, err = c.Api.RPC.State.GetStorageLatest(key, &peers)
	if err != nil {
		return nil, err
	}

	if !isSome {
		err = codec.Decode(dsecret.BootPeersResultDefaultBytes, &peers)
		if err != nil {
			return nil, err
		}
	}

	return peers, nil
}
