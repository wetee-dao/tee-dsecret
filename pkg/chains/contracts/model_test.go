package contracts

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/ink.go/util"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/contracts/cloud"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/contracts/subnet"
)

func TestCloud(t *testing.T) {
	client, err := chain.ClientInit("ws://127.0.0.1:9944", true)
	if err != nil {
		panic(err)
	}

	pk, err := chain.Sr25519PairFromSecret("//Alice", 42)
	if err != nil {
		util.LogWithPurple("Sr25519PairFromSecret", err)
		panic(err)
	}

	cloudIns, err := cloud.InitCloudContract(client, cloudAddress)
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

func TestSetSubnetSolt(t *testing.T) {
	client, err := chain.ClientInit("ws://127.0.0.1:9944", true)
	if err != nil {
		panic(err)
	}

	pk, err := chain.Sr25519PairFromSecret("//Alice", 42)
	if err != nil {
		util.LogWithPurple("Sr25519PairFromSecret", err)
		panic(err)
	}

	subnetIns, err := subnet.InitSubnetContract(client, subnetAddress)
	if err != nil {
		util.LogWithPurple("InitCloudContract", err)
		panic(err)
	}

	subnetIns.ExecSetEpochSolt(100, chain.ExecParams{
		Signer:    &pk,
		PayAmount: types.NewU128(*big.NewInt(0)),
	})
}
