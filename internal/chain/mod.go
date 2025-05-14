package chain

import (
	"fmt"

	chain "github.com/wetee-dao/go-sdk"
	types "wetee.app/dsecret/type"
)

// ChainClient
var ChainIns *Chain

// Chain
type Chain struct {
	*chain.ChainClient
	signer *chain.Signer
}

func InitChain(url string, pk *types.PrivKey) error {
	client, err := chain.ClientInit(url, true)
	if err != nil {
		return err
	}

	p, err := pk.ToSigner()
	if err != nil {
		return err
	}
	fmt.Println("Node chain pubkey:", p.Address)

	ChainIns = &Chain{
		ChainClient: client,
		signer:      p,
	}
	return nil
}
