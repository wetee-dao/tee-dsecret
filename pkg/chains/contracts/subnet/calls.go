package subnet

import (
	"errors"
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/ink.go/util"
)

func InitSubnetContract(client *chain.ChainClient, address string) (*Subnet, error) {
	contractAddress, err := util.HexToH160(address)
	if err != nil {
		return nil, err
	}
	return &Subnet{
		ChainClient: client,
		Address:     contractAddress,
	}, nil
}

func DeploySubnetWithNew(__ink_params chain.DeployParams) (*types.H160, error) {
	return __ink_params.Client.DeployContract(
		__ink_params.Code, __ink_params.Signer, types.NewU128(*big.NewInt(0)),
		util.InkContractInput{
			Selector: "0x9bae9d5e",
			Args:     []any{},
		},
		__ink_params.Salt,
	)
}

type Subnet struct {
	ChainClient *chain.ChainClient
	Address     types.H160
}

func (c *Subnet) Client() *chain.ChainClient {
	return c.ChainClient
}

func (c *Subnet) ContractAddress() types.H160 {
	return c.Address
}

func (c *Subnet) QueryBootNodes(
	params chain.DryRunCallParams,
) (*util.Result[[]SecretNode, Error], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Result[[]SecretNode, Error]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x3fd8cc61",
			Args:     []any{},
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

func (c *Subnet) DryRunSetBootNodes(
	nodes []uint64, params chain.DryRunCallParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xe6b90091",
			Args:     []any{nodes},
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

func (c *Subnet) CallSetBootNodes(
	nodes []uint64, __ink_params chain.CallParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSetBootNodes(nodes, _param)
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
			Selector: "0xe6b90091",
			Args:     []any{nodes},
		},
	)
}

func (c *Subnet) CallOfSetBootNodesTx(
	nodes []uint64, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSetBootNodes(nodes, _param)
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
			Selector: "0xe6b90091",
			Args:     []any{nodes},
		},
	)
}

func (c *Subnet) DryRunSetRegion(
	region_id uint32, name []byte, params chain.DryRunCallParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xb6993f90",
			Args:     []any{region_id, name},
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

func (c *Subnet) CallSetRegion(
	region_id uint32, name []byte, __ink_params chain.CallParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSetRegion(region_id, name, _param)
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
			Selector: "0xb6993f90",
			Args:     []any{region_id, name},
		},
	)
}

func (c *Subnet) CallOfSetRegionTx(
	region_id uint32, name []byte, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSetRegion(region_id, name, _param)
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
			Selector: "0xb6993f90",
			Args:     []any{region_id, name},
		},
	)
}

func (c *Subnet) QueryWorker(
	id uint64, params chain.DryRunCallParams,
) (*util.Option[K8sCluster], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Option[K8sCluster]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xdfcf3455",
			Args:     []any{id},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Subnet) QueryWorkers(
	params chain.DryRunCallParams,
) (*[]Tuple_88, *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[[]Tuple_88](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xc9dfba3b",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Subnet) DryRunWorkerRegister(
	name []byte, p2p_id AccountId, ip Ip, port uint32, level byte, region_id uint32, params chain.DryRunCallParams,
) (*util.Result[uint64, Error], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Result[uint64, Error]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xb90fc981",
			Args:     []any{name, p2p_id, ip, port, level, region_id},
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

func (c *Subnet) CallWorkerRegister(
	name []byte, p2p_id AccountId, ip Ip, port uint32, level byte, region_id uint32, __ink_params chain.CallParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunWorkerRegister(name, p2p_id, ip, port, level, region_id, _param)
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
			Selector: "0xb90fc981",
			Args:     []any{name, p2p_id, ip, port, level, region_id},
		},
	)
}

func (c *Subnet) CallOfWorkerRegisterTx(
	name []byte, p2p_id AccountId, ip Ip, port uint32, level byte, region_id uint32, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunWorkerRegister(name, p2p_id, ip, port, level, region_id, _param)
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
			Selector: "0xb90fc981",
			Args:     []any{name, p2p_id, ip, port, level, region_id},
		},
	)
}

func (c *Subnet) DryRunWorkerMortgage(
	id uint64, cpu uint32, mem uint32, cvm_cpu uint32, cvm_mem uint32, disk uint32, gpu uint32, deposit types.U256, params chain.DryRunCallParams,
) (*util.Result[uint32, Error], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Result[uint32, Error]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xf70c3369",
			Args:     []any{id, cpu, mem, cvm_cpu, cvm_mem, disk, gpu, deposit},
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

func (c *Subnet) CallWorkerMortgage(
	id uint64, cpu uint32, mem uint32, cvm_cpu uint32, cvm_mem uint32, disk uint32, gpu uint32, deposit types.U256, __ink_params chain.CallParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunWorkerMortgage(id, cpu, mem, cvm_cpu, cvm_mem, disk, gpu, deposit, _param)
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
			Selector: "0xf70c3369",
			Args:     []any{id, cpu, mem, cvm_cpu, cvm_mem, disk, gpu, deposit},
		},
	)
}

func (c *Subnet) CallOfWorkerMortgageTx(
	id uint64, cpu uint32, mem uint32, cvm_cpu uint32, cvm_mem uint32, disk uint32, gpu uint32, deposit types.U256, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunWorkerMortgage(id, cpu, mem, cvm_cpu, cvm_mem, disk, gpu, deposit, _param)
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
			Selector: "0xf70c3369",
			Args:     []any{id, cpu, mem, cvm_cpu, cvm_mem, disk, gpu, deposit},
		},
	)
}

func (c *Subnet) DryRunWorkerUnmortgage(
	worker_id uint64, mortgage_id uint32, params chain.DryRunCallParams,
) (*util.Result[uint32, Error], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Result[uint32, Error]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x6d25dbe9",
			Args:     []any{worker_id, mortgage_id},
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

func (c *Subnet) CallWorkerUnmortgage(
	worker_id uint64, mortgage_id uint32, __ink_params chain.CallParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunWorkerUnmortgage(worker_id, mortgage_id, _param)
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
			Selector: "0x6d25dbe9",
			Args:     []any{worker_id, mortgage_id},
		},
	)
}

func (c *Subnet) CallOfWorkerUnmortgageTx(
	worker_id uint64, mortgage_id uint32, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunWorkerUnmortgage(worker_id, mortgage_id, _param)
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
			Selector: "0x6d25dbe9",
			Args:     []any{worker_id, mortgage_id},
		},
	)
}

func (c *Subnet) DryRunWorkerStop(
	id uint64, params chain.DryRunCallParams,
) (*util.Result[uint64, Error], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Result[uint64, Error]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xeab3ba14",
			Args:     []any{id},
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

func (c *Subnet) CallWorkerStop(
	id uint64, __ink_params chain.CallParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunWorkerStop(id, _param)
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
			Selector: "0xeab3ba14",
			Args:     []any{id},
		},
	)
}

func (c *Subnet) CallOfWorkerStopTx(
	id uint64, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunWorkerStop(id, _param)
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
			Selector: "0xeab3ba14",
			Args:     []any{id},
		},
	)
}

func (c *Subnet) QuerySecrets(
	params chain.DryRunCallParams,
) (*[]Tuple_95, *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[[]Tuple_95](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xd91d379f",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Subnet) DryRunSecretRegister(
	name []byte, validator_id AccountId, p2p_id AccountId, ip Ip, port uint32, params chain.DryRunCallParams,
) (*util.Result[uint64, Error], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Result[uint64, Error]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x55719146",
			Args:     []any{name, validator_id, p2p_id, ip, port},
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

func (c *Subnet) CallSecretRegister(
	name []byte, validator_id AccountId, p2p_id AccountId, ip Ip, port uint32, __ink_params chain.CallParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSecretRegister(name, validator_id, p2p_id, ip, port, _param)
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
			Selector: "0x55719146",
			Args:     []any{name, validator_id, p2p_id, ip, port},
		},
	)
}

func (c *Subnet) CallOfSecretRegisterTx(
	name []byte, validator_id AccountId, p2p_id AccountId, ip Ip, port uint32, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSecretRegister(name, validator_id, p2p_id, ip, port, _param)
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
			Selector: "0x55719146",
			Args:     []any{name, validator_id, p2p_id, ip, port},
		},
	)
}

func (c *Subnet) DryRunSecretDeposit(
	id uint64, deposit types.U256, params chain.DryRunCallParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x4d815cac",
			Args:     []any{id, deposit},
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

func (c *Subnet) CallSecretDeposit(
	id uint64, deposit types.U256, __ink_params chain.CallParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSecretDeposit(id, deposit, _param)
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
			Selector: "0x4d815cac",
			Args:     []any{id, deposit},
		},
	)
}

func (c *Subnet) CallOfSecretDepositTx(
	id uint64, deposit types.U256, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSecretDeposit(id, deposit, _param)
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
			Selector: "0x4d815cac",
			Args:     []any{id, deposit},
		},
	)
}

func (c *Subnet) DryRunSecretDelete(
	id uint64, params chain.DryRunCallParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x016716ab",
			Args:     []any{id},
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

func (c *Subnet) CallSecretDelete(
	id uint64, __ink_params chain.CallParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSecretDelete(id, _param)
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
			Selector: "0x016716ab",
			Args:     []any{id},
		},
	)
}

func (c *Subnet) CallOfSecretDeleteTx(
	id uint64, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSecretDelete(id, _param)
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
			Selector: "0x016716ab",
			Args:     []any{id},
		},
	)
}

func (c *Subnet) QueryValidators(
	params chain.DryRunCallParams,
) (*[]Tuple_98, *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[[]Tuple_98](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xcc64f718",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Subnet) QueryGetPendingSecrets(
	params chain.DryRunCallParams,
) (*[]Tuple_61, *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[[]Tuple_61](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x8c922079",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Subnet) DryRunValidatorJoin(
	id uint64, params chain.DryRunCallParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x3c10643c",
			Args:     []any{id},
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

func (c *Subnet) CallValidatorJoin(
	id uint64, __ink_params chain.CallParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunValidatorJoin(id, _param)
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
			Selector: "0x3c10643c",
			Args:     []any{id},
		},
	)
}

func (c *Subnet) CallOfValidatorJoinTx(
	id uint64, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunValidatorJoin(id, _param)
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
			Selector: "0x3c10643c",
			Args:     []any{id},
		},
	)
}

func (c *Subnet) DryRunValidatorDelete(
	id uint64, params chain.DryRunCallParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xdbac71a8",
			Args:     []any{id},
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

func (c *Subnet) CallValidatorDelete(
	id uint64, __ink_params chain.CallParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunValidatorDelete(id, _param)
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
			Selector: "0xdbac71a8",
			Args:     []any{id},
		},
	)
}

func (c *Subnet) CallOfValidatorDeleteTx(
	id uint64, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunValidatorDelete(id, _param)
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
			Selector: "0xdbac71a8",
			Args:     []any{id},
		},
	)
}

func (c *Subnet) QueryEpochInfo(
	params chain.DryRunCallParams,
) (*EpochInfo, *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[EpochInfo](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xfd83a947",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Subnet) DryRunSetEpochSolt(
	epoch_solt uint32, params chain.DryRunCallParams,
) (*util.NullTuple, *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.NullTuple](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x6f527ac8",
			Args:     []any{epoch_solt},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Subnet) CallSetEpochSolt(
	epoch_solt uint32, __ink_params chain.CallParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSetEpochSolt(epoch_solt, _param)
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
			Selector: "0x6f527ac8",
			Args:     []any{epoch_solt},
		},
	)
}

func (c *Subnet) CallOfSetEpochSoltTx(
	epoch_solt uint32, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSetEpochSolt(epoch_solt, _param)
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
			Selector: "0x6f527ac8",
			Args:     []any{epoch_solt},
		},
	)
}

func (c *Subnet) DryRunSetNextEpoch(
	_node_id uint64, params chain.DryRunCallParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xfc3bf295",
			Args:     []any{_node_id},
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

func (c *Subnet) CallSetNextEpoch(
	_node_id uint64, __ink_params chain.CallParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSetNextEpoch(_node_id, _param)
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
			Selector: "0xfc3bf295",
			Args:     []any{_node_id},
		},
	)
}

func (c *Subnet) CallOfSetNextEpochTx(
	_node_id uint64, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSetNextEpoch(_node_id, _param)
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
			Selector: "0xfc3bf295",
			Args:     []any{_node_id},
		},
	)
}

func (c *Subnet) QueryNextEpochValidators(
	params chain.DryRunCallParams,
) (*util.Result[[]Tuple_98, Error], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Result[[]Tuple_98, Error]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x9f8ccaab",
			Args:     []any{},
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

func (c *Subnet) DryRunSetCode(
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

func (c *Subnet) CallSetCode(
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

func (c *Subnet) CallOfSetCodeTx(
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
