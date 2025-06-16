package chains

import (
	pallets "wetee.app/dsecret/chains/pallets"
	"wetee.app/dsecret/internal/model"
)

var ChainIns MainChain

type MainChain interface {
	// nodes
	GetBootPeers() ([]model.P2PAddr, error)
	GetNodes() ([]*model.Validator, []*model.PubKey, error)
	GetValidatorList() ([]*model.Validator, error)
	// epoch
	GetEpoch() (uint32, uint32, uint32, error)
}

func ConnectMainChain(url string, pk *model.PrivKey) (MainChain, error) {
	var err error
	ChainIns, err = pallets.InitChain(url, pk)

	return ChainIns, err
}
