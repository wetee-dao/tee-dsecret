package main

import (
	"encoding/hex"
	"fmt"

	"github.com/libp2p/go-libp2p/core/crypto"
	"go.dedis.ch/kyber/v3/suites"
)

func main() {
	priv2, pub2, err := crypto.GenerateKeyPair(
		crypto.Ed25519, // Select your key type. Ed25519 are nice short
		-1,             // Select key length when possible (i.e. RSA).
	)
	if err != nil {
		panic(err)
	}

	bt, _ := pub2.Raw()
	pubStr := hex.EncodeToString(bt)
	fmt.Println("pubStr2:", pubStr)
	fmt.Println("pub2:", pub2)
	fmt.Println("priv2:", priv2)

	bt2, err := crypto.MarshalPrivateKey(priv2)
	if err != nil {
		panic(err)
	}
	fmt.Println("priv2:", hex.EncodeToString(bt2))
	fmt.Println("priv2 bt2:", bt2)

	fmt.Println("---------------------------------------------")

	// 初始化加密套件。
	suite := suites.MustFind("Ed25519")
	nodeSecret := suite.Scalar().Pick(suite.RandomStream())
	nodePublic := suite.Point().Mul(nodeSecret, nil)

	fmt.Println("priv:", nodeSecret.String())
	fmt.Println("pub:", nodePublic.String())

}
