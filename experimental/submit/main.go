package main

import (
	"fmt"
	"os"
	"time"

	stypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/ink.go"
	"golang.org/x/crypto/blake2b"
	types "wetee.app/dsecret/type"
	"wetee.app/dsecret/type/pallet/dsecret"
	gtypes "wetee.app/dsecret/type/pallet/types"
)

func main() {
	client, err := chain.ClientInit("ws://192.168.111.105:30002", true)
	if err != nil {
		panic(err)
	}

	b, err := client.GetBlockNumber()
	if err != nil {
		panic(err)
	}
	fmt.Println(b)

	// 初始化加密套件
	nodeSecret, err := types.PrivateKeyFromLibp2pHex("080112406bce93c01f4b51287b01e55565cf7933cb624b25d478e003ca23446bc3ef83b9d0380163fd5c55a0474b95709da5b31d386da0313bb69bd635618f5cb80f1dde")
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

	err = client.SignAndSubmit(signer, call, true)
	if err != nil {
		panic(err)
	}
}
