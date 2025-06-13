package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	stypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/cosmos/go-bip39"
	oed "github.com/oasisprotocol/curve25519-voi/primitives/ed25519"
	inkutil "github.com/wetee-dao/ink.go/util"
	"go.dedis.ch/kyber/v4/suites"
	"wetee.app/dsecret/internal/model"
	"wetee.app/dsecret/internal/util"
)

func main() {
	_, privGo, _ := ed25519.GenerateKey(rand.Reader)
	privOasis, _ := model.StdToOed25519(privGo)
	msg := []byte("hello world")
	sig := ed25519.Sign(privGo, msg)
	sig2 := oed.Sign(privOasis, msg)
	fmt.Println("sig1:", hex.EncodeToString(sig))
	fmt.Println("sig2:", hex.EncodeToString(sig2))

	suite := suites.MustFind("Ed25519")
	privateKey, publicKey, err := model.GenerateKeyPair(suite, rand.Reader)
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
	privkey2, err := model.PrivateKeyFromLibp2pBytes(bt)
	if err != nil {
		fmt.Println("privkey2, err := model.PrivateKeyFromBytes(Ed25519, bt)", err)
		return
	}
	fmt.Println("privkey2:", privkey2.String())

	fmt.Println("libp2p pub:", privkey2.GetPublic().String())

	pub, _ := publicKey.Std()
	fmt.Println("pub", pub)

	var accpunt [32]byte
	copy(accpunt[:], pub.(ed25519.PublicKey))
	b := model.P2PAddr{
		Ip: model.Ip{
			Ipv4:   inkutil.NewNone[uint32](),
			Ipv6:   inkutil.NewNone[stypes.U128](),
			Domain: inkutil.NewSome([]byte("xiaobai.asyou.me")),
		},
		Port: 1234,
		Id:   accpunt,
	}

	pubHex := hex.EncodeToString(b.Id[:])
	pub2, err := model.PublicKeyFromHex(pubHex)

	fmt.Println("pub2", pub2)
	n := &model.Node{
		ID: *pub2,
	}
	d := util.GetUrlFromIp(b.Ip)
	fmt.Println(n.PeerID())
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
