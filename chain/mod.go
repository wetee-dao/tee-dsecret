package chain

import (
	"crypto/ed25519"
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

	bt, err := pk.Raw()
	if err != nil {
		return err
	}

	var ed25519Key ed25519.PrivateKey = bt
	p, err := core.Ed25519PairFromPk(ed25519Key, 42)
	if err != nil {
		return err
	}
	fmt.Println("Node chain pubkey:", p.Address)

	ChainClient = &Chain{
		client: client,
		signer: &p,
	}
	return nil
}
