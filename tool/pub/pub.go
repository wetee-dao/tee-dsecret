package main

import (
	"crypto/ed25519"

	"github.com/vedhavyas/go-subkey/v2"
	types "wetee.app/dsecret/type"
)

func main() {
	pubkey := "5FSwTPcYPutz4WKSqFJAS8jCt9eJGKdB5zbUvpUqwyVByguF"
	_, bt, _ := subkey.SS58Decode(pubkey)

	var goppub ed25519.PublicKey = bt
	println(goppub)

	pub, err := types.PubKeyFromStdPubKey(goppub)
	if err != nil {
		panic(err)
	}
	println(pub)
}
