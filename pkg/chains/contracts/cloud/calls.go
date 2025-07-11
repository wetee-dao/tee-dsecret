package cloud

import (
	"errors"
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/ink.go/util"
)

func InitCloudContract(client *chain.ChainClient, address string) (*Cloud, error) {
	contractAddress, err := util.HexToH160(address)
	if err != nil {
		return nil, err
	}
	return &Cloud{
		ChainClient: client,
		Address:     contractAddress,
	}, nil
}

func DeployCloudWithNew(subnet_addr types.H160, pod_contract_code_hash types.H256, __ink_params chain.DeployParams) (*types.H160, error) {
	return __ink_params.Client.DeployContract(
		__ink_params.Code, __ink_params.Signer, types.NewU128(*big.NewInt(0)),
		util.InkContractInput{
			Selector: "0x9bae9d5e",
			Args:     []any{subnet_addr, pod_contract_code_hash},
		},
		__ink_params.Salt,
	)
}

type Cloud struct {
	ChainClient *chain.ChainClient
	Address     types.H160
}

func (c *Cloud) Client() *chain.ChainClient {
	return c.ChainClient
}

func (c *Cloud) ContractAddress() types.H160 {
	return c.Address
}

func (c *Cloud) DryRunSetPodContract(
	pod_contract types.H256, params chain.DryRunCallParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xeebfb380",
			Args:     []any{pod_contract},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	if v != nil && v.IsErr {
		return nil, nil, errors.New("Contract Reverted: " + v.E.Error())
	}

	return v, gas, nil
}

func (c *Cloud) CallSetPodContract(
	pod_contract types.H256, __ink_params chain.CallParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSetPodContract(pod_contract, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		__ink_params.Signer,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xeebfb380",
			Args:     []any{pod_contract},
		},
	)
}

func (c *Cloud) CallOfSetPodContractTx(
	pod_contract types.H256, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSetPodContract(pod_contract, _param)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.Signer,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xeebfb380",
			Args:     []any{pod_contract},
		},
	)
}

func (c *Cloud) QuerySubnetAddress(
	params chain.DryRunCallParams,
) (*types.H160, *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[types.H160](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x241d1854",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) DryRunCreateUserPod(
	name []byte, pod_type PodType, tee_type TEEType, containers []Container, region_id uint32, level byte, worker_id uint64, params chain.DryRunCallParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x1035bf95",
			Args:     []any{name, pod_type, tee_type, containers, region_id, level, worker_id},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	if v != nil && v.IsErr {
		return nil, nil, errors.New("Contract Reverted: " + v.E.Error())
	}

	return v, gas, nil
}

func (c *Cloud) CallCreateUserPod(
	name []byte, pod_type PodType, tee_type TEEType, containers []Container, region_id uint32, level byte, worker_id uint64, __ink_params chain.CallParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunCreateUserPod(name, pod_type, tee_type, containers, region_id, level, worker_id, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		__ink_params.Signer,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x1035bf95",
			Args:     []any{name, pod_type, tee_type, containers, region_id, level, worker_id},
		},
	)
}

func (c *Cloud) CallOfCreateUserPodTx(
	name []byte, pod_type PodType, tee_type TEEType, containers []Container, region_id uint32, level byte, worker_id uint64, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunCreateUserPod(name, pod_type, tee_type, containers, region_id, level, worker_id, _param)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.Signer,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x1035bf95",
			Args:     []any{name, pod_type, tee_type, containers, region_id, level, worker_id},
		},
	)
}

func (c *Cloud) QueryPodLen(
	params chain.DryRunCallParams,
) (*uint64, *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[uint64](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xaf63d0e1",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) QueryPods(
	page uint64, size uint64, params chain.DryRunCallParams,
) (*[]Tuple_87, *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[[]Tuple_87](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xba743fed",
			Args:     []any{page, size},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) QueryUserPodLen(
	params chain.DryRunCallParams,
) (*uint32, *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[uint32](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x31385138",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) QueryUserPods(
	page uint32, size uint32, params chain.DryRunCallParams,
) (*[]Tuple_87, *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[[]Tuple_87](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x2ba5c5d5",
			Args:     []any{page, size},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) QueryWorkerPodsVersion(
	worker_id uint64, params chain.DryRunCallParams,
) (*[]Tuple_91, *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[[]Tuple_91](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x56d09cd0",
			Args:     []any{worker_id},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) QueryWorkerPods(
	worker_id uint64, page uint32, size uint32, params chain.DryRunCallParams,
) (*[]Tuple_87, *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[[]Tuple_87](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xd2d1cf5e",
			Args:     []any{worker_id, page, size},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) QueryPod(
	pod_id uint64, params chain.DryRunCallParams,
) (*util.Option[Tuple_94], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Option[Tuple_94]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xb431f434",
			Args:     []any{pod_id},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) QueryPodsByIds(
	pod_ids []uint64, params chain.DryRunCallParams,
) (*[]Tuple_98, *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[[]Tuple_98](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x711ca8a1",
			Args:     []any{pod_ids},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) QueryWorkerPodLen(
	worker_id uint64, params chain.DryRunCallParams,
) (*uint32, *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[uint32](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x2fced50e",
			Args:     []any{worker_id},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) DryRunSetCode(
	code_hash types.H256, params chain.DryRunCallParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x694fb50f",
			Args:     []any{code_hash},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	if v != nil && v.IsErr {
		return nil, nil, errors.New("Contract Reverted: " + v.E.Error())
	}

	return v, gas, nil
}

func (c *Cloud) CallSetCode(
	code_hash types.H256, __ink_params chain.CallParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSetCode(code_hash, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		__ink_params.Signer,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x694fb50f",
			Args:     []any{code_hash},
		},
	)
}

func (c *Cloud) CallOfSetCodeTx(
	code_hash types.H256, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSetCode(code_hash, _param)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.Signer,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x694fb50f",
			Args:     []any{code_hash},
		},
	)
}
