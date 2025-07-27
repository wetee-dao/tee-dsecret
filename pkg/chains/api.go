package chains

import (
	// pallets "github.com/wetee-dao/tee-dsecret/pkg/chains/pallets"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/contracts"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

var MainChain Chain

type Chain interface {
	// get chain client
	GetClient() *chain.ChainClient
	GetSignerAddress() string
	// nodes
	GetBootPeers() ([]model.P2PAddr, error)
	GetNodes() ([]*model.Validator, []*model.PubKey, error)
	GetValidatorList() ([]*model.Validator, error)
	// epoch
	GetEpoch() (uint32, uint32, uint32, uint32, types.H160, error)
	GetNextEpochValidatorList() ([]*model.Validator, error)
	SetNewEpoch(nodeId uint64) error
	TxCallOfSetNextEpoch(nodeId uint64, signer types.AccountID) (*types.Call, error)

	/// query node id
	GetMintWorker(user types.AccountID) (*model.K8sCluster, error)

	// query pods by worker
	GetPodsVersionByWorker(workerId uint64) ([]model.PodVersion, error)
	GetPodsByIds(podIds []uint64) ([]model.Pod, error)
	GetWorker(workerId uint64) (*model.K8sCluster, error)
	ResigerCluster(name []byte, p2p_id [32]byte, ip model.Ip, port uint32, level byte, region_id uint32) error

	// POD
	GetMintInterval() (uint32, error)

	TxCallOfStartPod(nodeId uint64, pod_key types.AccountID, signer types.AccountID) (*types.Call, error)
	DryStartPod(nodeId uint64, pod_key types.AccountID, signer types.AccountID) error

	TxCallOfMintPod(nodeId uint64, hash types.H256, signer types.AccountID) (*types.Call, error)
	DryMintPod(nodeId uint64, hash types.H256, signer types.AccountID) error
}

func ConnectMainChain(url []string, pk *model.PrivKey) (Chain, error) {
	var err error

	MainChain, err = contracts.NewContract(url, pk)
	return MainChain, err
}
