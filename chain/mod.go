package chain

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature/ed25519"
	chain "github.com/wetee-dao/go-sdk"
)

var ChainClient *Chain

// Chain
type Chain struct {
	client *chain.ChainClient
	signer *signature.KeyringPair
}

func InitChain(url string, seed string) error {
	client, err := chain.ClientInit(url, true)
	if err != nil {
		return err
	}

	p, err := ed25519.KeyringPairFromSecret(seed, 42)
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
