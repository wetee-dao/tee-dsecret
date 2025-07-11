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

	/// init pod
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

	/// init subnet
	subnetAddress := DeploySubnetContract(client, pk)

	/// init cloud
	cloudCode, err := os.ReadFile("./contract_cache/cloud.polkavm")
	if err != nil {
		util.LogWithPurple("read file error", err)
		panic(err)
	}

	salt := genSalt()
	cloudAddress, err := cloud.DeployCloudWithNew(*subnetAddress, *podCode, chain.DeployParams{
		Client: client,
		Signer: &pk,
		Code:   util.InkCode{Upload: &cloudCode},
		Salt:   util.NewSome(salt),
	})

	if err != nil {
		util.LogWithPurple("DeployContract", err)
		panic(err)
	}

	InitSubnet(client, pk, subnetAddress.Hex())
	InitWorker(client, pk, subnetAddress.Hex())
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
	_call := chain.CallParams{
		Signer:    &pk,
		PayAmount: types.NewU128(*big.NewInt(0)),
	}
	subnetContract, err := subnet.InitSubnetContract(client, subnetAddress)
	if err != nil {
		panic(err)
	}

	v1, _ := model.PubKeyFromSS58("5CdERUzLMFh5D8RB82bd6t4nuqKJLdNr6ZQ9NAsoQqVMyz5B")
	p1, _ := model.PubKeyFromSS58("5CAG6XhZY5Q3seRa4BwDhSQGFHqoA4H2m3GJKew7xArJwcNJ")
	err = subnetContract.CallSecretRegister(
		[]byte("node0"),
		v1.AccountID(),
		p1.AccountID(),
		subnet.Ip{
			Ipv4:   util.NewSome[uint32](2130706433),
			Ipv6:   util.NewNone[types.U128](),
			Domain: util.NewNone[[]byte](),
		},
		31000,
		_call,
	)
	fmt.Println("node0 register result:", err)

	v2, _ := model.PubKeyFromSS58("5Fk6tyXKk9HmATcSvtcEjMHsyfn2e49H76qP72yFXzUU4ws6")
	p2, _ := model.PubKeyFromSS58("5GuRb3N6Qraej2S3kQNX33UMnk47saYTAH4EBGzPiuqG8kni")
	err = subnetContract.CallSecretRegister(
		[]byte("node1"),
		v2.AccountID(),
		p2.AccountID(),
		subnet.Ip{
			Ipv4:   util.NewSome[uint32](2130706433),
			Ipv6:   util.NewNone[types.U128](),
			Domain: util.NewNone[[]byte](),
		},
		41000,
		_call,
	)
	fmt.Println("node1 register result:", err)

	v3, _ := model.PubKeyFromSS58("5CK7kDvy6svMswxifABZAu8GFrcAvEw1z9nt7Wuuvh8YMzx1")
	p3, _ := model.PubKeyFromSS58("5FgmV7fM5yAyZK5DfbAv3x9CrSBcnNt3Zykbxs9S9HHrvbeG")
	err = subnetContract.CallSecretRegister(
		[]byte("node2"),
		v3.AccountID(),
		p3.AccountID(),
		subnet.Ip{
			Ipv4:   util.NewSome[uint32](2130706433),
			Ipv6:   util.NewNone[types.U128](),
			Domain: util.NewNone[[]byte](),
		},
		51000,
		_call,
	)
	fmt.Println("node2 register result:", err)

	subnetContract.CallSetBootNodes([]uint64{0, 1, 2}, _call)
	subnetContract.CallValidatorJoin(1, _call)
	subnetContract.CallValidatorJoin(2, _call)
}

func InitWorker(client *chain.ChainClient, pk chain.Signer, subnetAddress string) {
	_call := chain.CallParams{
		Signer:    &pk,
		PayAmount: types.NewU128(*big.NewInt(0)),
	}
	subnetContract, err := subnet.InitSubnetContract(client, subnetAddress)
	if err != nil {
		panic(err)
	}

	subnetContract.CallSetRegion(0, []byte("defalut"), _call)

	pubkey, _ := model.PubKeyFromSS58("5GSBfdb3PxME3XM4JrkFKAgHH77ADDWXUx6o8KGVmavLnZ44")
	subnetContract.CallWorkerRegister([]byte("worker0"), pubkey.AccountID(), subnet.Ip{
		Ipv4:   util.NewNone[uint32](),
		Ipv6:   util.NewNone[types.U128](),
		Domain: util.NewSome([]byte("xiaobai.asyou.me")),
	}, 10000, 1, 0, _call)

	// subnetContract.CallWorkerMortgage(
	// 	0,
	// 	10000, 10000,
	// 	0, 0,
	// 	1000000,
	// 	0,
	// 	types.NewU256(*big.NewInt(10000000)),
	// 	_call,
	// )
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
