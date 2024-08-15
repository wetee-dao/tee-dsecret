package chain

import (
	"fmt"

	chain "github.com/wetee-dao/go-sdk"
	"github.com/wetee-dao/go-sdk/core"
	types "wetee.app/dsecret/type"
)

// ChainClient
var ChainClient *Chain

// Chain
type Chain struct {
	client *chain.ChainClient
	signer *core.Signer
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

	ChainClient = &Chain{
		client: client,
		signer: p,
	}
	return nil
}
