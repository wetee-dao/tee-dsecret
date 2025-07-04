package subnet

import (
	"errors"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/ink.go/util"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

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
			Selector: "set_code",
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
			Selector: "set_code",
			Args:     []any{code_hash},
		},
	)
}

func (c *Subnet) TxCallOfSetCode(
	code_hash types.H256, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSetCode(code_hash, _param)
	if err != nil {
		return nil, err
	}
	return chain.TxCall(
		c,
		__ink_params.Signer,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "set_code",
			Args:     []any{code_hash},
		},
	)
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
			Selector: "boot_nodes",
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
			Selector: "set_boot_nodes",
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
			Selector: "set_boot_nodes",
			Args:     []any{nodes},
		},
	)
}

func (c *Subnet) TxCallOfSetBootNodes(
	nodes []uint64, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSetBootNodes(nodes, _param)
	if err != nil {
		return nil, err
	}
	return chain.TxCall(
		c,
		__ink_params.Signer,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "set_boot_nodes",
			Args:     []any{nodes},
		},
	)
}

func (c *Subnet) QueryWorkers(
	params chain.DryRunCallParams,
) (*[]Tuple_70, *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[[]Tuple_70](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "workers",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Subnet) DryRunWorkerRegister(
	name []byte, p2p_id AccountId, ip Ip, port uint32, level byte, params chain.DryRunCallParams,
) (*util.Result[uint64, Error], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Result[uint64, Error]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "worker_register",
			Args:     []any{name, p2p_id, ip, port, level},
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
	name []byte, p2p_id AccountId, ip Ip, port uint32, level byte, __ink_params chain.CallParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunWorkerRegister(name, p2p_id, ip, port, level, _param)
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
			Selector: "worker_register",
			Args:     []any{name, p2p_id, ip, port, level},
		},
	)
}

func (c *Subnet) TxCallOfWorkerRegister(
	name []byte, p2p_id AccountId, ip Ip, port uint32, level byte, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunWorkerRegister(name, p2p_id, ip, port, level, _param)
	if err != nil {
		return nil, err
	}
	return chain.TxCall(
		c,
		__ink_params.Signer,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "worker_register",
			Args:     []any{name, p2p_id, ip, port, level},
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
			Selector: "worker_mortgage",
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
			Selector: "worker_mortgage",
			Args:     []any{id, cpu, mem, cvm_cpu, cvm_mem, disk, gpu, deposit},
		},
	)
}

func (c *Subnet) TxCallOfWorkerMortgage(
	id uint64, cpu uint32, mem uint32, cvm_cpu uint32, cvm_mem uint32, disk uint32, gpu uint32, deposit types.U256, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunWorkerMortgage(id, cpu, mem, cvm_cpu, cvm_mem, disk, gpu, deposit, _param)
	if err != nil {
		return nil, err
	}
	return chain.TxCall(
		c,
		__ink_params.Signer,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "worker_mortgage",
			Args:     []any{id, cpu, mem, cvm_cpu, cvm_mem, disk, gpu, deposit},
		},
	)
}

func (c *Subnet) DryRunWorkerUnmortgage(
	id uint64, mortgage_id uint32, params chain.DryRunCallParams,
) (*util.Result[uint32, Error], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Result[uint32, Error]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "worker_unmortgage",
			Args:     []any{id, mortgage_id},
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
	id uint64, mortgage_id uint32, __ink_params chain.CallParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunWorkerUnmortgage(id, mortgage_id, _param)
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
			Selector: "worker_unmortgage",
			Args:     []any{id, mortgage_id},
		},
	)
}

func (c *Subnet) TxCallOfWorkerUnmortgage(
	id uint64, mortgage_id uint32, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunWorkerUnmortgage(id, mortgage_id, _param)
	if err != nil {
		return nil, err
	}
	return chain.TxCall(
		c,
		__ink_params.Signer,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "worker_unmortgage",
			Args:     []any{id, mortgage_id},
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
			Selector: "worker_stop",
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
			Selector: "worker_stop",
			Args:     []any{id},
		},
	)
}

func (c *Subnet) TxCallOfWorkerStop(
	id uint64, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunWorkerStop(id, _param)
	if err != nil {
		return nil, err
	}
	return chain.TxCall(
		c,
		__ink_params.Signer,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "worker_stop",
			Args:     []any{id},
		},
	)
}

func (c *Subnet) QuerySecrets(
	params chain.DryRunCallParams,
) (*[]Tuple_77, *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[[]Tuple_77](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "secrets",
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
			Selector: "secret_register",
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
			Selector: "secret_register",
			Args:     []any{name, validator_id, p2p_id, ip, port},
		},
	)
}

func (c *Subnet) TxCallOfSecretRegister(
	name []byte, validator_id AccountId, p2p_id AccountId, ip Ip, port uint32, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSecretRegister(name, validator_id, p2p_id, ip, port, _param)
	if err != nil {
		return nil, err
	}
	return chain.TxCall(
		c,
		__ink_params.Signer,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "secret_register",
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
			Selector: "secret_deposit",
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
			Selector: "secret_deposit",
			Args:     []any{id, deposit},
		},
	)
}

func (c *Subnet) TxCallOfSecretDeposit(
	id uint64, deposit types.U256, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSecretDeposit(id, deposit, _param)
	if err != nil {
		return nil, err
	}
	return chain.TxCall(
		c,
		__ink_params.Signer,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "secret_deposit",
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
			Selector: "secret_delete",
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
			Selector: "secret_delete",
			Args:     []any{id},
		},
	)
}

func (c *Subnet) TxCallOfSecretDelete(
	id uint64, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSecretDelete(id, _param)
	if err != nil {
		return nil, err
	}
	return chain.TxCall(
		c,
		__ink_params.Signer,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "secret_delete",
			Args:     []any{id},
		},
	)
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
			Selector: "validator_join",
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
			Selector: "validator_join",
			Args:     []any{id},
		},
	)
}

func (c *Subnet) TxCallOfValidatorJoin(
	id uint64, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunValidatorJoin(id, _param)
	if err != nil {
		return nil, err
	}
	return chain.TxCall(
		c,
		__ink_params.Signer,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "validator_join",
			Args:     []any{id},
		},
	)
}

func (c *Subnet) QueryGetPendingSecrets(
	params chain.DryRunCallParams,
) (*[]Tuple_45, *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[[]Tuple_45](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "get_pending_secrets",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
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
			Selector: "validator_delete",
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
			Selector: "validator_delete",
			Args:     []any{id},
		},
	)
}

func (c *Subnet) TxCallOfValidatorDelete(
	id uint64, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunValidatorDelete(id, _param)
	if err != nil {
		return nil, err
	}
	return chain.TxCall(
		c,
		__ink_params.Signer,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "validator_delete",
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
			Selector: "epoch_info",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Subnet) QueryValidators(
	params chain.DryRunCallParams,
) (*[]Tuple_83, *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[[]Tuple_83](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "validators",
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
			Selector: "set_epoch_solt",
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
			Selector: "set_epoch_solt",
			Args:     []any{epoch_solt},
		},
	)
}

func (c *Subnet) TxCallOfSetEpochSolt(
	epoch_solt uint32, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSetEpochSolt(epoch_solt, _param)
	if err != nil {
		return nil, err
	}
	return chain.TxCall(
		c,
		__ink_params.Signer,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "set_epoch_solt",
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
			Selector: "set_next_epoch",
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
			Selector: "set_next_epoch",
			Args:     []any{_node_id},
		},
	)
}

func (c *Subnet) TxCallOfSetNextEpoch(
	_node_id uint64, __ink_params chain.CallParams,
) (*types.Call, error) {
	pubkey := model.PubKeyFromByte(__ink_params.Signer.Public())
	fmt.Println(pubkey.SS58())
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSetNextEpoch(_node_id, _param)
	if err != nil {
		return nil, err
	}
	return chain.TxCall(
		c,
		__ink_params.Signer,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "set_next_epoch",
			Args:     []any{_node_id},
		},
	)
}

func (c *Subnet) QueryNextEpochValidators(
	params chain.DryRunCallParams,
) (*util.Result[[]Tuple_83, Error], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Result[[]Tuple_83, Error]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "next_epoch_validators",
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
