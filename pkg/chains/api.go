package chains

import (
	// pallets "github.com/wetee-dao/tee-dsecret/pkg/chains/pallets"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/ink.go/util"
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
	GetEpoch() (uint32, uint32, uint32, uint32, [32]byte, error)
	GetNextEpochValidatorList() ([]*model.Validator, error)
	SetNewEpoch(new_key [32]byte, sig [64]byte) error
}

func ConnectMainChain(url string, pk *model.PrivKey) (Chain, error) {
	var err error

	// chain, err = pallets.InitChain(url, pk)
	contractAddress, err := util.HexToH160("0xC2A11E61acC3Bc9598150Fd3086Ea88f8B5c1377")
	if err != nil {
		util.LogWithPurple("HexToH160", err)
		return nil, err
	}

	MainChain, err = contracts.NewContract(url, pk, contractAddress)
	return MainChain, err
}
