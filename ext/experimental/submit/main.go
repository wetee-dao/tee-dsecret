package main

import (
	"fmt"
	"os"
	"time"

	stypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/dsecret"
	gtypes "github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/types"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"golang.org/x/crypto/blake2b"
)

func main() {
	client, err := chain.ClientInit("ws://127.0.0.1:9944", true)
	if err != nil {
		panic(err)
	}

	b, err := client.GetBlockNumber()
	if err != nil {
		panic(err)
	}
	fmt.Println(b)

	// 初始化加密套件
	var testSecretSeed = "0x7512939e37970b04c2b9a6060b16654473cf0721b71f8e56126ee314cbd0a7e9fe9125a81688d932ea792e9722c777f7696117363f86a107ab9d3681a8c922c8"
	nodeSecret, err := model.PrivateKeyFromHex(testSecretSeed)
	if err != nil {
		fmt.Println("Marshal PKG_PK error:", err)
		os.Exit(1)
	}

	signer, _ := nodeSecret.ToSigner()

	data := []byte("1234567890")
	hash := blake2b.Sum512(data)
	sigbt, err := signer.Sign(hash[:])
	if err != nil {
		panic(err)
	}
	sig2 := stypes.NewSignature(sigbt)
	sig := gtypes.MultiSignature{
		IsEd25519:       true,
		AsEd25519Field0: sig2,
	}

	var account32 [32]byte
	copy(account32[:], signer.Public())
	runtimeCall := dsecret.MakeUploadClusterProofCall(
		1,
		hash[:],
		[][32]byte{account32, account32},
		[]gtypes.MultiSignature{sig, sig},
	)

	time.Sleep(5 * time.Second)

	call, err := (runtimeCall).AsCall()
	if err != nil {
		panic(err)
	}

	err = client.SignAndSubmit(signer, call, true, 0)
	if err != nil {
		panic(err)
	}
}
