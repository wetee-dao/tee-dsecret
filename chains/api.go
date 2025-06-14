package chains

import (
	pallets "wetee.app/dsecret/chains/pallets"
	"wetee.app/dsecret/internal/model"
)

var ChainIns MainChain

type MainChain interface {
	GetBootPeers() ([]model.P2PAddr, error)
	GetNodes() ([][32]byte, []*model.Node, error)
}

func ConnectMainChain(url string, pk *model.PrivKey) (MainChain, error) {
	var err error
	ChainIns, err = pallets.InitChain(url, pk)

	return ChainIns, err
}
