package contracts

//go:generate go-ink-gen -json subnet.json
//go:generate go-ink-gen -json cloud.json

import (
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/ink.go/util"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/contracts/cloud"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/contracts/subnet"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/revive"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

// Contract
type Contract struct {
	*chain.ChainClient
	signer *chain.Signer
	subnet *subnet.Subnet
	cloud  *cloud.Cloud
}

const subnetAddress = "0xd4b3f313be7143e117a021ef649e5324a2fd30fa"
const cloudAddress = "0xda3d54a1f5321c873b053f164e2ced96c614788e"

func GetCloudAddress() string {
	return cloudAddress
}

func NewContract(url string, pk *model.PrivKey) (*Contract, error) {
	client, err := chain.ClientInit(url, false)
	if err != nil {
		return nil, err
	}

	subnet, err := subnet.InitSubnetContract(client, subnetAddress)
	if err != nil {
		util.LogWithPurple("InitSubnetContract", err)
		return nil, err
	}

	cloud, err := cloud.InitCloudContract(client, cloudAddress)
	if err != nil {
		util.LogWithPurple("InitCloudContract", err)
		return nil, err
	}

	p, err := pk.ToSigner()
	if err != nil {
		return nil, err
	}

	util.LogWithYellow("Mainchain Key", pk.GetPublic().SS58())
	h160 := pk.GetPublic().H160()

	// check account is mapaccount in revive
	_, isSome, err := revive.GetOriginalAccountLatest(client.Api.RPC.State, h160)
	if err != nil {
		util.LogWithPurple("GetOriginalAccountLatest", err)
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
		cloud:       cloud,
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
	result, _, err := c.subnet.QueryBootNodes(chain.DefaultParamWithOrigin(types.AccountID(c.signer.AccountID())))
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
	workers, _, err := c.subnet.QueryWorkers(1, 5000, chain.DefaultParamWithOrigin(types.AccountID(c.signer.AccountID())))
	if err != nil {
		return nil, nil, err
	}
	dsecrets, _, err := c.subnet.QuerySecrets(chain.DefaultParamWithOrigin(types.AccountID(c.signer.AccountID())))
	if err != nil {
		return nil, nil, err
	}

	nodes := make([]*model.PubKey, 0, len(*workers)+len(*dsecrets))
	validators := make([]*model.Validator, 0, len(*dsecrets))

	for _, v := range *workers {
		nodes = append(nodes, model.PubKeyFromByte(v.F1.P2pId[:]))
	}
	for _, v := range *dsecrets {
		if !v.F1.TerminalBlock.IsNone() {
			continue
		}
		nodes = append(nodes, model.PubKeyFromByte(v.F1.P2pId[:]))
		validators = append(validators, &model.Validator{
			NodeID:      v.F0,
			P2pId:       *model.PubKeyFromByte(v.F1.P2pId[:]),
			ValidatorId: *model.PubKeyFromByte(v.F1.ValidatorId[:]),
		})
	}

	return validators, nodes, nil
}

func (c *Contract) GetValidatorList() ([]*model.Validator, error) {
	dsecrets, _, err := c.subnet.QueryValidators(chain.DefaultParamWithOrigin(types.AccountID(c.signer.AccountID())))
	if err != nil {
		return nil, err
	}

	validators := make([]*model.Validator, 0, len(*dsecrets))
	for _, v := range *dsecrets {
		validators = append(validators, &model.Validator{
			NodeID:      uint64(v.F0),
			P2pId:       *model.PubKeyFromByte(v.F1.P2pId[:]),
			ValidatorId: *model.PubKeyFromByte(v.F1.ValidatorId[:]),
		})
	}

	return validators, nil
}

// epoch
func (c *Contract) GetEpoch() (uint32, uint32, uint32, uint32, types.H160, error) {
	d, _, err := c.subnet.QueryEpochInfo(chain.DefaultParamWithOrigin(types.AccountID(c.signer.AccountID())))
	if err != nil {
		return 0, 0, 0, 0, types.H160{}, err
	}

	address := types.H160{}
	if !d.SideChainPub.IsNone() {
		address = d.SideChainPub.V
	}

	return d.Epoch, d.EpochSolt, d.LastEpochBlock, d.Now, address, nil
}

// go to new epoch
func (c *Contract) SetNewEpoch(nodeId uint64) error {
	return c.subnet.CallSetNextEpoch(nodeId, chain.CallParams{
		Signer:    c.signer,
		PayAmount: types.NewU128(*big.NewInt(0)),
	})
}

func (c *Contract) TxCallOfSetNextEpoch(nodeId uint64, signer chain.SignerType) (*types.Call, error) {
	return c.subnet.CallOfSetNextEpochTx(nodeId, chain.CallParams{
		Signer:    signer,
		PayAmount: types.NewU128(*big.NewInt(0)),
	})
}

func (c *Contract) GetNextEpochValidatorList() ([]*model.Validator, error) {
	dsecrets, _, err := c.subnet.QueryNextEpochValidators(chain.DefaultParamWithOrigin(types.AccountID(c.signer.AccountID())))
	if err != nil {
		return nil, err
	}

	validators := make([]*model.Validator, 0, len(dsecrets.V))
	for _, v := range dsecrets.V {
		validators = append(validators, &model.Validator{
			NodeID:      v.F0,
			P2pId:       *model.PubKeyFromByte(v.F1.P2pId[:]),
			ValidatorId: *model.PubKeyFromByte(v.F1.ValidatorId[:]),
		})
	}

	return validators, nil
}

func ConvertTEEType(t cloud.TEEType) model.TEEType {
	return model.TEEType{
		SGX: t.SGX,
		CVM: t.CVM,
	}
}
