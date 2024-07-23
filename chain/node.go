package chain

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/wetee-dao/go-sdk/gen/weteedsecret"
)

func RegisterNode(signer *signature.KeyringPair, pubkey []byte) error {
	call := weteedsecret.MakeRegisterNodeCall(pubkey)
	return ChainClient.SignAndSubmit(signer, call, true)
}
