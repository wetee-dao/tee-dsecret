package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/ink.go/util"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/contracts/cloud"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/contracts/subnet"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func main() {
	client, err := chain.ClientInit("ws://127.0.0.1:9944", true)
	if err != nil {
		panic(err)
	}

	pk, err := chain.Sr25519PairFromSecret("//Alice", 42)
	if err != nil {
		util.LogWithPurple("Sr25519PairFromSecret", err)
		panic(err)
	}

	podData, err := os.ReadFile("./contract_cache/pod.polkavm")
	if err != nil {
		util.LogWithPurple("read file error", err)
		panic(err)
	}

	podCode, err := client.UploadInkCode(podData, &pk)
	if err != nil {
		util.LogWithPurple("UploadInkCode", err)
		panic(err)
	}

	fmt.Println(podCode.Hex())
	subnetAddress := DeploySubnetContract(client, pk)

	data, err := os.ReadFile("./contract_cache/subnet.polkavm")
	if err != nil {
		util.LogWithPurple("read file error", err)
		panic(err)
	}

	salt := genSalt()
	cloudAddress, err := cloud.DeployCloudWithNew(*subnetAddress, *podCode, chain.DeployParams{
		Client: client,
		Signer: &pk,
		Code:   util.InkCode{Upload: &data},
		Salt:   util.NewSome(salt),
	})

	if err != nil {
		util.LogWithPurple("DeployContract", err)
		panic(err)
	}

	fmt.Println("subnet address ======> ", subnetAddress.Hex())
	InitSubnet(client, pk, subnetAddress.Hex())
	fmt.Println("subnet address ======> ", subnetAddress.Hex())
	fmt.Println("cloud  address ======> ", cloudAddress.Hex())
}

func DeploySubnetContract(client *chain.ChainClient, pk chain.Signer) *types.H160 {
	data, err := os.ReadFile("./contract_cache/subnet.polkavm")
	if err != nil {
		util.LogWithPurple("read file error", err)
		panic(err)
	}

	salt := genSalt()
	res, err := subnet.DeploySubnetWithNew(chain.DeployParams{
		Client: client,
		Signer: &pk,
		Code:   util.InkCode{Upload: &data},
		Salt:   util.NewSome(salt),
	})

	if err != nil {
		util.LogWithPurple("DeployContract", err)
		panic(err)
	}

	return res
}

func InitSubnet(client *chain.ChainClient, pk chain.Signer, subnetAddress string) {
	contract, err := subnet.InitSubnetContract(client, subnetAddress)
	if err != nil {
		panic(err)
	}

	v1, _ := model.PubKeyFromSS58("5CdERUzLMFh5D8RB82bd6t4nuqKJLdNr6ZQ9NAsoQqVMyz5B")
	p1, _ := model.PubKeyFromSS58("5CAG6XhZY5Q3seRa4BwDhSQGFHqoA4H2m3GJKew7xArJwcNJ")
	err = contract.CallSecretRegister(
		[]byte("node0"),
		v1.AccountID(),
		p1.AccountID(),
		subnet.Ip{
			Ipv4:   util.NewSome[uint32](2130706433),
			Ipv6:   util.NewNone[types.U128](),
			Domain: util.NewNone[[]byte](),
		},
		31000,
		chain.CallParams{
			Signer:    &pk,
			PayAmount: types.NewU128(*big.NewInt(0)),
		},
	)
	fmt.Println("worker register result:", err)

	v2, _ := model.PubKeyFromSS58("5Fk6tyXKk9HmATcSvtcEjMHsyfn2e49H76qP72yFXzUU4ws6")
	p2, _ := model.PubKeyFromSS58("5GuRb3N6Qraej2S3kQNX33UMnk47saYTAH4EBGzPiuqG8kni")
	err = contract.CallSecretRegister(
		[]byte("node1"),
		v2.AccountID(),
		p2.AccountID(),
		subnet.Ip{
			Ipv4:   util.NewSome[uint32](2130706433),
			Ipv6:   util.NewNone[types.U128](),
			Domain: util.NewNone[[]byte](),
		},
		41000,
		chain.CallParams{
			Signer:    &pk,
			PayAmount: types.NewU128(*big.NewInt(0)),
		},
	)
	fmt.Println("worker register result:", err)

	v3, _ := model.PubKeyFromSS58("5CK7kDvy6svMswxifABZAu8GFrcAvEw1z9nt7Wuuvh8YMzx1")
	p3, _ := model.PubKeyFromSS58("5FgmV7fM5yAyZK5DfbAv3x9CrSBcnNt3Zykbxs9S9HHrvbeG")
	err = contract.CallSecretRegister(
		[]byte("node2"),
		v3.AccountID(),
		p3.AccountID(),
		subnet.Ip{
			Ipv4:   util.NewSome[uint32](2130706433),
			Ipv6:   util.NewNone[types.U128](),
			Domain: util.NewNone[[]byte](),
		},
		51000,
		chain.CallParams{
			Signer:    &pk,
			PayAmount: types.NewU128(*big.NewInt(0)),
		},
	)
	fmt.Println("worker register result:", err)

	contract.CallSetBootNodes([]uint64{0, 1, 2}, chain.CallParams{
		Signer:    &pk,
		PayAmount: types.NewU128(*big.NewInt(0)),
	})

	contract.CallValidatorJoin(1, chain.CallParams{
		Signer:    &pk,
		PayAmount: types.NewU128(*big.NewInt(0)),
	})

	contract.CallValidatorJoin(2, chain.CallParams{
		Signer:    &pk,
		PayAmount: types.NewU128(*big.NewInt(0)),
	})
}

func genSalt() [32]byte {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	randomBytes := [32]byte{}
	copy(randomBytes[:], bytes)

	return randomBytes
}
