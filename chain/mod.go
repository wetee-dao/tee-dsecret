package chain

import (
	chain "github.com/wetee-dao/go-sdk"
)

var ChainClient *chain.ChainClient

func InitChain(url string) error {
	var err error
	ChainClient, err = chain.ClientInit(url, true)
	if err != nil {
		return err
	}
	return nil
}
