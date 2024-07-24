package main

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/go-bip39"
	"github.com/libp2p/go-libp2p/core/crypto"
	"wetee.app/dsecret/types"
)

func main() {
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

	privateKey, err := types.PrivateKeyFromPhrase(mnemonic, "")
	if err != nil {
		fmt.Println(err)
		return
	}
	publicKey := privateKey.GetPublic()

	fmt.Println("privateKey:", privateKey.String())
	fmt.Println("publicKey:", publicKey)
	fmt.Println("publicKey.Point():", publicKey.Point())
	// sPriv := privateKey.Scalar()
	// fmt.Println("sPriv: ", sPriv)
	// sPub := suite.Point().Mul(sPriv, nil)
	// fmt.Println("sPub: ", sPub)
	// bt, err := hex.DecodeString(privateKey.String())
	// if err != nil {
	// 	fmt.Println("bt, err := hex.DecodeString(privateKey.String())", err)
	// 	return
	// }
	// privkey2, err := types.PrivateKeyFromBytes(bt)
	// if err != nil {
	// 	fmt.Println("privkey2, err := types.PrivateKeyFromBytes(Ed25519, bt)", err)
	// 	return
	// }
	// fmt.Println("privkey2:", privkey2.String())

	// pair, err := ed25519sig.KeyringPairFromSecret(mnemonic, 42)
	// fmt.Println("pair:", hex.EncodeToString(pair.PublicKey))

	// pair2, err := signature.KeyringPairFromSecret(mnemonic, 42)
	// fmt.Println("pair2:", hex.EncodeToString(pair2.PublicKey))

	priv, pub, err := crypto.GenerateKeyPair(
		crypto.Ed25519,
		-1,
	)

	b, _ := crypto.MarshalPrivateKey(priv)
	b2, _ := pub.Raw()
	fmt.Println("priv:", hex.EncodeToString(b))
	fmt.Println("pub:", hex.EncodeToString(b2))
}
