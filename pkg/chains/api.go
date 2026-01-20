package chains

import (

	// pallets "github.com/wetee-dao/tee-dsecret/pkg/chains/pallets"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/ink.go"
	contracts "github.com/wetee-dao/tee-dsecret/pkg/chains/ink"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

var MainChain MainChainApi

// ChainApi is the interface for the chain
type ChainApi interface {
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

	// secret
	TxCallOfUploadSecret(user types.H160, index uint64, signer types.AccountID) (*types.Call, error)
	DryUploadSecret(user types.H160, index uint64, signer types.AccountID) error

	// disk
	TxCallOfInitDisk(user types.H160, index uint64, hash types.H256, signer types.AccountID) (*types.Call, error)
	DryInitDisk(user types.H160, index uint64, hash types.H256, signer types.AccountID) error

	// TEE call to call
	TEECallToCall(tcall *model.TeeCall, dkgKey types.AccountID) (*types.Call, error)
}

// MainChainApi is the interface for the main chain
type MainChainApi interface {
	ChainApi

	// get chain client
	GetClient() *chain.ChainClient
	GetChainUrls() []string
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
}

// ConnectMainChain 连接主链
func ConnectMainChain(urls []string, pk *model.PrivKey) (MainChainApi, error) {
	var err error

	MainChain, err = contracts.NewContract(urls, pk)
	return MainChain, err
}

func ConnectChain(urls []string, pk *model.PrivKey) (ChainApi, error) {
	var err error

	chain, err := contracts.NewContract(urls, pk)
	return chain, err
}
