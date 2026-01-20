package contracts

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/ink.go/util"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/ink/cloud"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/ink/subnet"
)

func TestCloud(t *testing.T) {
	client, err := chain.InitClient([]string{TestChainUrl}, true)
	if err != nil {
		panic(err)
	}

	pk, err := chain.Sr25519PairFromSecret("//Alice", 42)
	if err != nil {
		util.LogWithPurple("Sr25519PairFromSecret", err)
		panic(err)
	}

	cloudIns, err := cloud.InitCloudContract(client, CloudAddress)
	if err != nil {
		util.LogWithPurple("InitCloudContract", err)
		panic(err)
	}

	subnet, _, err := cloudIns.QuerySubnetAddress(chain.DefaultParamWithOrigin(pk.AccountID()))
	if err != nil {
		util.LogWithPurple("QueryPodLen", err)
		panic(err)
	}
	fmt.Println(subnet.Hex())
}

func TestQueryCloudStage(t *testing.T) {
	client, err := chain.InitClient([]string{TestChainUrl}, true)
	if err != nil {
		panic(err)
	}

	pk, err := chain.Sr25519PairFromSecret("//Alice", 42)
	if err != nil {
		util.LogWithPurple("Sr25519PairFromSecret", err)
		panic(err)
	}

	cloudIns, err := cloud.InitCloudContract(client, CloudAddress)
	if err != nil {
		util.LogWithPurple("InitCloudContract", err)
		panic(err)
	}

	stage, _, err := cloudIns.QueryMintInterval(chain.DefaultParamWithOrigin(pk.AccountID()))
	if err != nil {
		util.LogWithPurple("QueryPodLen", err)
		panic(err)
	}
	fmt.Println(*stage)
}

func TestSetCloudStage(t *testing.T) {
	client, err := chain.InitClient([]string{TestChainUrl}, true)
	if err != nil {
		panic(err)
	}

	pk, err := chain.Sr25519PairFromSecret("//Alice", 42)
	if err != nil {
		util.LogWithPurple("Sr25519PairFromSecret", err)
		panic(err)
	}

	cloudIns, err := cloud.InitCloudContract(client, CloudAddress)
	if err != nil {
		util.LogWithPurple("InitCloudContract", err)
		panic(err)
	}

	err = cloudIns.ExecSetMintInterval(200, chain.ExecParams{
		Signer:    &pk,
		PayAmount: types.NewU128(*big.NewInt(0)),
	})
	if err != nil {
		util.LogWithPurple("ExecSetMintInterval", err)
	}
}

func TestSetSubnetSolt(t *testing.T) {
	client, err := chain.InitClient([]string{TestChainUrl}, true)
	if err != nil {
		panic(err)
	}

	pk, err := chain.Sr25519PairFromSecret("//Alice", 42)
	if err != nil {
		util.LogWithPurple("Sr25519PairFromSecret", err)
		panic(err)
	}

	subnetIns, err := subnet.InitSubnetContract(client, SubnetAddress)
	if err != nil {
		util.LogWithPurple("InitCloudContract", err)
		panic(err)
	}

	subnetIns.ExecSetEpochSolt(1000, chain.ExecParams{
		Signer:    &pk,
		PayAmount: types.NewU128(*big.NewInt(0)),
	})
}

func TestCloudUpdate(t *testing.T) {
	client, err := chain.InitClient([]string{TestChainUrl}, true)
	if err != nil {
		panic(err)
	}

	pk, err := chain.Sr25519PairFromSecret("//Alice", 42)
	if err != nil {
		util.LogWithPurple("Sr25519PairFromSecret", err)
		panic(err)
	}

	/// init pod
	cloudData, err := os.ReadFile("../../../hack/contract_cache/cloud.polkavm")
	if err != nil {
		util.LogWithPurple("read file error", err)
		panic(err)
	}

	code, err := client.UploadInkCode(cloudData, &pk)
	if err != nil {
		util.LogWithPurple("UploadInkCode", err)
		panic(err)
	}

	cloudIns, err := cloud.InitCloudContract(client, CloudAddress)
	if err != nil {
		util.LogWithPurple("InitCloudContract", err)
		panic(err)
	}

	err = cloudIns.ExecSetCode(*code, chain.ExecParams{
		Signer:    &pk,
		PayAmount: types.NewU128(*big.NewInt(0)),
	})

	if err != nil {
		util.LogWithPurple("ExecSetCode", err)
	}
}

func TestSubnetUpdate(t *testing.T) {
	client, err := chain.InitClient([]string{TestChainUrl}, true)
	if err != nil {
		panic(err)
	}

	pk, err := chain.Sr25519PairFromSecret("//Alice", 42)
	if err != nil {
		util.LogWithPurple("Sr25519PairFromSecret", err)
		panic(err)
	}

	/// init pod
	netData, err := os.ReadFile("../../../hack/contract_cache/subnet.polkavm")
	if err != nil {
		util.LogWithPurple("read file error", err)
		panic(err)
	}

	netCode, err := client.UploadInkCode(netData, &pk)
	if err != nil {
		util.LogWithPurple("UploadInkCode", err)
		panic(err)
	}

	fmt.Println("cloudAddress: ", CloudAddress)

	subnetIns, err := subnet.InitSubnetContract(client, SubnetAddress)
	if err != nil {
		util.LogWithPurple("InitCloudContract", err)
		panic(err)
	}

	err = subnetIns.ExecSetCode(*netCode, chain.ExecParams{
		Signer:    &pk,
		PayAmount: types.NewU128(*big.NewInt(0)),
	})

	if err != nil {
		util.LogWithPurple("subnet ExecSetCode", err)
	}
}

func TestWorkerUpdate(t *testing.T) {
	client, err := chain.InitClient([]string{TestChainUrl}, true)
	if err != nil {
		panic(err)
	}

	pk, err := chain.Sr25519PairFromSecret("//Alice", 42)
	if err != nil {
		util.LogWithPurple("Sr25519PairFromSecret", err)
		panic(err)
	}

	subnetIns, err := subnet.InitSubnetContract(client, SubnetAddress)
	if err != nil {
		util.LogWithPurple("InitSubnetContract", err)
		panic(err)
	}

	err = subnetIns.ExecSecretUpdate(0, []byte("v0"), subnet.Ip{
		Ipv4:   util.NewSome[uint32](3232263885),
		Ipv6:   util.NewNone[types.U128](),
		Domain: util.NewNone[[]byte](),
	}, 31000, chain.ExecParams{
		Signer:    &pk,
		PayAmount: types.NewU128(*big.NewInt(0)),
	})
	fmt.Println(err)

	subnetIns.ExecSecretUpdate(1, []byte("v1"), subnet.Ip{
		Ipv4:   util.NewSome[uint32](3232263885),
		Ipv6:   util.NewNone[types.U128](),
		Domain: util.NewNone[[]byte](),
	}, 41000, chain.ExecParams{
		Signer:    &pk,
		PayAmount: types.NewU128(*big.NewInt(0)),
	})

	subnetIns.ExecSecretUpdate(2, []byte("v2"), subnet.Ip{
		Ipv4:   util.NewSome[uint32](3232263885),
		Ipv6:   util.NewNone[types.U128](),
		Domain: util.NewNone[[]byte](),
	}, 51000, chain.ExecParams{
		Signer:    &pk,
		PayAmount: types.NewU128(*big.NewInt(0)),
	})
}
