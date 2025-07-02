package contracts

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/ink.go/util"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/contracts/subnet"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func TestInitNetwork(t *testing.T) {
	contractAddress, err := util.HexToH160("0xC2A11E61acC3Bc9598150Fd3086Ea88f8B5c1377")
	if err != nil {
		util.LogWithPurple("HexToH160", err)
		t.Fatal(err)
	}

	pk, err := chain.Sr25519PairFromSecret("//Alice", 42)
	if err != nil {
		util.LogWithPurple("Sr25519PairFromSecret", err)
		t.Fatal(err)
	}

	client, err := chain.ClientInit("ws://127.0.0.1:9944", true)
	if err != nil {
		t.Fatal(err)
	}

	contract := subnet.Subnet{
		ChainClient: client,
		Address:     contractAddress,
	}

	v1, _ := model.PubKeyFromSS58("5CdERUzLMFh5D8RB82bd6t4nuqKJLdNr6ZQ9NAsoQqVMyz5B")
	p1, _ := model.PubKeyFromSS58("5CAG6XhZY5Q3seRa4BwDhSQGFHqoA4H2m3GJKew7xArJwcNJ")
	_, gas, err := contract.DryRunSecretRegister(
		[]byte("node0"),
		v1.AccountID(),
		p1.AccountID(),
		subnet.Ip{
			Ipv4:   util.NewSome[uint32](2130706433),
			Ipv6:   util.NewNone[types.U128](),
			Domain: util.NewNone[[]byte](),
		},
		31000,
		chain.DefaultParamWithOragin(types.AccountID(pk.AccountID())),
	)
	if err != nil {
		t.Fatal(err)
	}

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
			Signer:              &pk,
			PayAmount:           types.NewU128(*big.NewInt(0)),
			GasLimit:            gas.GasRequired,
			StorageDepositLimit: gas.StorageDeposit,
		},
	)
	fmt.Println("worker register result:", err)

	v2, _ := model.PubKeyFromSS58("5Fk6tyXKk9HmATcSvtcEjMHsyfn2e49H76qP72yFXzUU4ws6")
	p2, _ := model.PubKeyFromSS58("5GuRb3N6Qraej2S3kQNX33UMnk47saYTAH4EBGzPiuqG8kni")
	_, gas, err = contract.DryRunSecretRegister(
		[]byte("node1"),
		v2.AccountID(),
		p2.AccountID(),
		subnet.Ip{
			Ipv4:   util.NewSome[uint32](2130706433),
			Ipv6:   util.NewNone[types.U128](),
			Domain: util.NewNone[[]byte](),
		},
		41000,
		chain.DefaultParamWithOragin(types.AccountID(pk.AccountID())),
	)
	if err != nil {
		t.Fatal(err)
	}
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
			Signer:              &pk,
			PayAmount:           types.NewU128(*big.NewInt(0)),
			GasLimit:            gas.GasRequired,
			StorageDepositLimit: gas.StorageDeposit,
		},
	)
	fmt.Println("worker register result:", err)

	v3, _ := model.PubKeyFromSS58("5CK7kDvy6svMswxifABZAu8GFrcAvEw1z9nt7Wuuvh8YMzx1")
	p3, _ := model.PubKeyFromSS58("5FgmV7fM5yAyZK5DfbAv3x9CrSBcnNt3Zykbxs9S9HHrvbeG")
	_, gas, err = contract.DryRunSecretRegister(
		[]byte("node2"),
		v3.AccountID(),
		p3.AccountID(),
		subnet.Ip{
			Ipv4:   util.NewSome[uint32](2130706433),
			Ipv6:   util.NewNone[types.U128](),
			Domain: util.NewNone[[]byte](),
		},
		51000,
		chain.DefaultParamWithOragin(types.AccountID(pk.AccountID())),
	)
	if err != nil {
		t.Fatal(err)
	}
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
			Signer:              &pk,
			PayAmount:           types.NewU128(*big.NewInt(0)),
			GasLimit:            gas.GasRequired,
			StorageDepositLimit: gas.StorageDeposit,
		},
	)
	fmt.Println("worker register result:", err)

	contract.CallSetBootNodes([]uint64{0, 1, 2}, chain.CallParams{
		Signer:              &pk,
		PayAmount:           types.NewU128(*big.NewInt(0)),
		GasLimit:            gas.GasRequired,
		StorageDepositLimit: gas.StorageDeposit,
	})

	contract.CallValidatorJoin(1, chain.CallParams{
		Signer:              &pk,
		PayAmount:           types.NewU128(*big.NewInt(0)),
		GasLimit:            gas.GasRequired,
		StorageDepositLimit: gas.StorageDeposit,
	})

	contract.CallValidatorJoin(2, chain.CallParams{
		Signer:              &pk,
		PayAmount:           types.NewU128(*big.NewInt(0)),
		GasLimit:            gas.GasRequired,
		StorageDepositLimit: gas.StorageDeposit,
	})

}
