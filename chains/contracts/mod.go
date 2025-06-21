package contracts

//go:generate go-ink-gen -json subnet.json

import (
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/ink.go/util"
	"wetee.app/dsecret/chains/contracts/subnet"

	"wetee.app/dsecret/internal/model"
)

var defaultPrams = chain.DryRunCallParams{
	Amount:              types.NewU128(*big.NewInt(0)),
	GasLimit:            util.NewNone[types.Weight](),
	StorageDepositLimit: util.NewNone[types.U128](),
}

func queryPramsWithOragin(origin types.AccountID) chain.DryRunCallParams {
	return chain.DryRunCallParams{
		Origin:              origin,
		Amount:              types.NewU128(*big.NewInt(0)),
		GasLimit:            util.NewNone[types.Weight](),
		StorageDepositLimit: util.NewNone[types.U128](),
	}
}

// Contract
type Contract struct {
	*chain.ChainClient
	signer *chain.Signer
	subnet subnet.Subnet
}

func NewContract(client *chain.ChainClient, signer *chain.Signer, subNetAddress types.H160) *Contract {
	subnet := subnet.Subnet{
		ChainClient: client,
		Address:     subNetAddress,
	}

	return &Contract{
		ChainClient: client,
		signer:      signer,
		subnet:      subnet,
	}
}

func (c *Contract) GetSignerAddress() string {
	return c.signer.SS58Address(42)
}

// nodes
func (c *Contract) GetBootPeers() ([]model.P2PAddr, error) {
	result, err := c.subnet.QueryBootNodes(queryPramsWithOragin(types.AccountID(c.signer.AccountID())))
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
	workers, err := c.subnet.QueryWorkers(queryPramsWithOragin(types.AccountID(c.signer.AccountID())))
	if err != nil {
		return nil, nil, err
	}
	dsecrets, err := c.subnet.QuerySecrets(queryPramsWithOragin(types.AccountID(c.signer.AccountID())))
	if err != nil {
		return nil, nil, err
	}

	nodes := make([]*model.PubKey, 0, len(workers.V)+len(dsecrets.V))
	validators := make([]*model.Validator, 0, len(dsecrets.V))

	for _, v := range workers.V {
		nodes = append(nodes, model.PubKeyFromByte(v.P2pId[:]))
	}
	for _, v := range dsecrets.V {
		nodes = append(nodes, model.PubKeyFromByte(v.P2pId[:]))
		validators = append(validators, &model.Validator{
			P2pId:       *model.PubKeyFromByte(v.P2pId[:]),
			ValidatorId: *model.PubKeyFromByte(v.ValidatorId[:]),
		})
	}

	return validators, nodes, nil
}

func (c *Contract) GetValidatorList() ([]*model.Validator, error) {
	dsecrets, err := c.subnet.QuerySecrets(queryPramsWithOragin(types.AccountID(c.signer.AccountID())))
	if err != nil {
		return nil, err
	}

	validators := make([]*model.Validator, 0, len(dsecrets.V))
	for _, v := range dsecrets.V {
		validators = append(validators, &model.Validator{
			P2pId:       *model.PubKeyFromByte(v.P2pId[:]),
			ValidatorId: *model.PubKeyFromByte(v.ValidatorId[:]),
		})
	}

	return validators, nil
}

// epoch
func (c *Contract) GetEpoch() (uint32, uint32, uint32, error) {
	d, err := c.subnet.QueryEpoch(queryPramsWithOragin(types.AccountID(c.signer.AccountID())))
	if err != nil {
		return 0, 0, 0, err
	}

	return d.F0, d.F1, d.F2, nil
}
