package contracts

//go:generate go-ink-gen -json subnet.json

import (
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/ink.go/util"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/contracts/subnet"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/revive"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

// Contract
type Contract struct {
	*chain.ChainClient
	signer *chain.Signer
	subnet subnet.Subnet
}

func NewContract(url string, pk *model.PrivKey, subNetAddress types.H160) (*Contract, error) {
	client, err := chain.ClientInit(url, false)
	if err != nil {
		return nil, err
	}

	p, err := pk.ToSigner()
	if err != nil {
		return nil, err
	}

	util.LogWithYellow("Mainchain Key", pk.GetPublic().SS58())
	h160 := pk.GetPublic().H160()

	subnet := subnet.Subnet{
		ChainClient: client,
		Address:     subNetAddress,
	}

	// check account is mapaccount in revive
	_, isSome, err := revive.GetOriginalAccountLatest(client.Api.RPC.State, h160)
	if err != nil {
		return nil, err
	}
	if !isSome {
		runtimeCall := revive.MakeMapAccountCall()
		call, err := (runtimeCall).AsCall()
		if err != nil {
			return nil, err
		}

		err = client.SignAndSubmit(p, call, true)
		if err != nil {
			return nil, err
		}
	}

	return &Contract{
		ChainClient: client,
		signer:      p,
		subnet:      subnet,
	}, nil
}

func (c *Contract) GetClient() *chain.ChainClient {
	return c.ChainClient
}

func (c *Contract) GetSignerAddress() string {
	return c.signer.SS58Address(42)
}

// nodes
func (c *Contract) GetBootPeers() ([]model.P2PAddr, error) {
	result, _, err := c.subnet.QueryBootNodes(chain.DefaultParamWithOragin(types.AccountID(c.signer.AccountID())))
	if err != nil {
		return nil, err
	}

	nodes := result.V
	boots := make([]model.P2PAddr, 0, len(nodes))
	for _, n := range nodes {
		boots = append(boots, model.P2PAddr{
			Id: n.P2pId,
			Ip: model.Ip{
				Ipv4:   n.Ip.Ipv4,
				Ipv6:   n.Ip.Ipv6,
				Domain: n.Ip.Domain,
			},
			Port: uint16(n.Port),
		})
	}

	return boots, nil
}

func (c *Contract) GetNodes() ([]*model.Validator, []*model.PubKey, error) {
	workers, _, err := c.subnet.QueryWorkers(chain.DefaultParamWithOragin(types.AccountID(c.signer.AccountID())))
	if err != nil {
		return nil, nil, err
	}
	dsecrets, _, err := c.subnet.QuerySecrets(chain.DefaultParamWithOragin(types.AccountID(c.signer.AccountID())))
	if err != nil {
		return nil, nil, err
	}

	nodes := make([]*model.PubKey, 0, len(*workers)+len(*dsecrets))
	validators := make([]*model.Validator, 0, len(*dsecrets))

	for _, v := range *workers {
		nodes = append(nodes, model.PubKeyFromByte(v.F1.P2pId[:]))
	}
	for _, v := range *dsecrets {
		if !v.F1.TerminalBlock.IsNone {
			continue
		}
		nodes = append(nodes, model.PubKeyFromByte(v.F1.P2pId[:]))
		validators = append(validators, &model.Validator{
			P2pId:       *model.PubKeyFromByte(v.F1.P2pId[:]),
			ValidatorId: *model.PubKeyFromByte(v.F1.ValidatorId[:]),
		})
	}

	return validators, nodes, nil
}

func (c *Contract) GetValidatorList() ([]*model.Validator, error) {
	dsecrets, _, err := c.subnet.QueryValidators(chain.DefaultParamWithOragin(types.AccountID(c.signer.AccountID())))
	if err != nil {
		return nil, err
	}

	validators := make([]*model.Validator, 0, len(*dsecrets))
	for _, v := range *dsecrets {
		validators = append(validators, &model.Validator{
			P2pId:       *model.PubKeyFromByte(v.F0.P2pId[:]),
			ValidatorId: *model.PubKeyFromByte(v.F0.ValidatorId[:]),
		})
	}

	return validators, nil
}

// epoch
func (c *Contract) GetEpoch() (uint32, uint32, uint32, uint32, [32]byte, error) {
	d, _, err := c.subnet.QueryEpochInfo(chain.DefaultParamWithOragin(types.AccountID(c.signer.AccountID())))
	if err != nil {
		return 0, 0, 0, 0, [32]byte{}, err
	}

	return d.Epoch, d.EpochSolt, d.LastEpochBlock, d.Now, d.SideChainPub, nil
}

// go to new epoch
func (c *Contract) SetNewEpoch(new_key [32]byte, sig [64]byte) error {
	_, gas, err := c.subnet.DryRunSetNextEpoch(
		new_key,
		sig,
		chain.DefaultParamWithOragin(types.AccountID(c.signer.AccountID())),
	)

	if err != nil {
		return err
	}

	return c.subnet.CallSetNextEpoch(new_key, sig, chain.CallParams{
		Signer:              c.signer,
		PayAmount:           types.NewU128(*big.NewInt(0)),
		GasLimit:            gas.GasRequired,
		StorageDepositLimit: gas.StorageDeposit,
	})
}

func (c *Contract) GetNextEpochValidatorList() ([]*model.Validator, error) {
	dsecrets, _, err := c.subnet.QueryNextEpochValidators(chain.DefaultParamWithOragin(types.AccountID(c.signer.AccountID())))
	if err != nil {
		return nil, err
	}

	validators := make([]*model.Validator, 0, len(dsecrets.V))
	for _, v := range dsecrets.V {
		validators = append(validators, &model.Validator{
			P2pId:       *model.PubKeyFromByte(v.F0.P2pId[:]),
			ValidatorId: *model.PubKeyFromByte(v.F0.ValidatorId[:]),
		})
	}

	return validators, nil
}
