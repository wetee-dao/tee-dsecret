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
	TxCallOfSetNextEpoch(nodeId uint64, signer chain.SignerType) (*types.Call, error)
}

func ConnectMainChain(url string, pk *model.PrivKey) (Chain, error) {
	var err error

	MainChain, err = contracts.NewContract(url, pk)
	return MainChain, err
}
