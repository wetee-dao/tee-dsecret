package main

import (
	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cometbft/cometbft/p2p"

	"github.com/cometbft/cometbft/privval"
)

func main() {
	if _, err := p2p.LoadOrGenNodeKey("node_key.json"); err != nil {
		panic(err)
	}

	pv, err := privval.GenFilePV("priv_validator_key.json", "priv_validator_state.json", func() (crypto.PrivKey, error) { //nolint: unparam
		return ed25519.GenPrivKey(), nil
	})
	if err != nil {
		panic(err)
	}

	pv.Save()
}
