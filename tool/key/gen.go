package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/cosmos/go-bip39"
	gtypes "github.com/wetee-dao/go-sdk/pallet/types"
	"go.dedis.ch/kyber/v3/suites"
	types "wetee.app/dsecret/type"
	"wetee.app/dsecret/util"
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

	fmt.Println("libp2p pub:", privkey2.GetPublic().String())

	pub, _ := publicKey.Std()
	fmt.Println("pub", pub)

	var accpunt [32]byte
	copy(accpunt[:], pub.(ed25519.PublicKey))
	b := gtypes.P2PAddr{
		Ip: gtypes.Ip{
			Ipv4: gtypes.OptionTUint32{
				IsNone: true,
			},
			Ipv6: gtypes.OptionTU128{
				IsNone: true,
			},
			Domain: gtypes.OptionTByteSlice{
				IsSome:       true,
				AsSomeField0: []byte("xiaobai.asyou.me"),
			},
		},
		Port: 1234,
		Id:   accpunt,
	}

	pubHex := hex.EncodeToString(b.Id[:])
	pub2, err := types.PublicKeyFromHex(pubHex)

	fmt.Println("pub2", pub2)
	n := &types.Node{
		ID: pub2.String(),
	}
	d := util.GetUrlFromIp(b.Ip)
	fmt.Println(n.PeerID().String())
	url := d + "/tcp/" + fmt.Sprint(b.Port) + "/p2p/" + n.PeerID().String()
	fmt.Println(url)
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
