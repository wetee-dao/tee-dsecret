package chains

import (
	// pallets "github.com/wetee-dao/tee-dsecret/chains/pallets"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/ink.go/util"
	"github.com/wetee-dao/tee-dsecret/chains/contracts"
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
	GetEpoch() (uint32, uint32, uint32, error)
}

func ConnectMainChain(url string, pk *model.PrivKey) (Chain, error) {
	var err error

	// chain, err = pallets.InitChain(url, pk)

	contractAddress, err := util.HexToH160("0x2F6991f9eF07B521d3b45e831aE816E13cc0e4c5")
	if err != nil {
		util.LogWithPurple("HexToH160", err)
		return nil, err
	}

	MainChain, err = contracts.NewContract(url, pk, contractAddress)
	return MainChain, err
}
