package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/cosmos/go-bip39"
	"go.dedis.ch/kyber/v3/suites"
	"wetee.app/dsecret/types"
)

func main() {
	suite := suites.MustFind("Ed25519")
	privateKey, publicKey, err := types.GenerateKeyPair(suite, rand.Reader)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("privateKey:", privateKey.String())
	fmt.Println("publicKey:", publicKey)

	fmt.Println("publicKey.Point():", publicKey.Point())

	sPriv := privateKey.Scalar()
	fmt.Println("sPriv: ", sPriv)
	sPub := suite.Point().Mul(sPriv, nil)
	fmt.Println("sPub: ", sPub)

	bt, err := hex.DecodeString(privateKey.String())
	if err != nil {
		fmt.Println("bt, err := hex.DecodeString(privateKey.String())", err)
		return
	}
	privkey2, err := types.PrivateKeyFromLibp2pBytes(bt)
	if err != nil {
		fmt.Println("privkey2, err := types.PrivateKeyFromBytes(Ed25519, bt)", err)
		return
	}
	fmt.Println("privkey2:", privkey2.String())

}

func seed() {
	// 生成 128 位的随机熵
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		panic(err)
	}

	// 生成助记词
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		panic(err)
	}

	// 打印助记词
	fmt.Println("助记词:", mnemonic)
}
