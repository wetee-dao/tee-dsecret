package pallets

import (
	"crypto/ed25519"
	"errors"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/dsecret"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/types"
	gtypes "github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/types"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/util"
)

// RegisterNode register node
// 注册节点
func (c *Chain) RegisterNode(signer *chain.Signer, vid []byte, pid []byte) error {
	var bt [32]byte
	copy(bt[:], vid)

	var pidBt [32]byte
	copy(pidBt[:], pid)

	runtimeCall := dsecret.MakeRegisterNodeCall(bt, pidBt)

	call, err := (runtimeCall).AsCall()
	if err != nil {
		return errors.New("(runtimeCall).AsCall() error: " + err.Error())
	}

	return c.SignAndSubmit(signer, call, true, 0)
}

// GetNodes 函数用于获取节点列表，包括 Secret 节点和 Worker 节点，以及转换为自定义的 Node 类型
func (c *Chain) GetNodes() ([]*model.Validator, []*model.PubKey, error) {
	// 获取节点列表
	secretNodes, err := c.GetValidatorList()
	if err != nil {
		return nil, nil, errors.New("Get node list error:" + err.Error())
	}
	workerNodes, err := c.GetWorkerList()
	if err != nil {
		return nil, nil, errors.New("Get worker list error:" + err.Error())
	}

	nodes := []*model.PubKey{}
	for _, n := range secretNodes {
		nodes = append(nodes, &n.P2pId)
	}

	for _, w := range workerNodes {
		var gopub ed25519.PublicKey = w.Account[:]
		pub, _ := model.PubKeyFromStdPubKey(gopub)
		nodes = append(nodes, pub)
	}

	return secretNodes, nodes, nil
}

// GetNodeList get node list
// 获取节点列表
func (c *Chain) GetValidatorList() ([]*model.Validator, error) {
	ret, err := c.QueryMapAll("DSecret", "Validators")
	if err != nil {
		return nil, err
	}

	nodes := make([]types.Validator, 0)
	for _, elem := range ret {
		for _, change := range elem.Changes {
			n := types.Validator{}
			if err := codec.Decode(change.StorageData, &n); err != nil {
				util.LogError("codec.Decode", err)
				continue
			}
			nodes = append(nodes, n)
		}
	}

	validators := make([]*model.Validator, 0, len(nodes))
	for _, n := range nodes {
		validators = append(validators, &model.Validator{
			ValidatorId: *model.PubKeyFromByte(n.ValidatorId[:]),
			P2pId:       *model.PubKeyFromByte(n.P2pId[:]),
		})
	}

	return validators, nil
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
	isSome, err = c.Api().RPC.State.GetStorageLatest(key, &peers)
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

func (c *Chain) GetEpoch() (uint32, uint32, uint32, error) {
	epoch, err := dsecret.GetEpochLatest(c.Api().RPC.State)
	if err != nil {
		return 0, 0, 0, err
	}

	lastEpochBlock, err := dsecret.GetLastEpochBlockLatest(c.Api().RPC.State)
	if err != nil {
		return 0, 0, 0, err
	}

	now, err := c.GetBlockNumber()
	if err != nil {
		return 0, 0, 0, err
	}

	return epoch, lastEpochBlock, uint32(now), nil
}
