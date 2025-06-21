package chains

import (
	pallets "wetee.app/dsecret/chains/pallets"
	"wetee.app/dsecret/internal/model"
)

var MainChain Chain

type Chain interface {
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
	MainChain, err = pallets.InitChain(url, pk)

	return MainChain, err
}
