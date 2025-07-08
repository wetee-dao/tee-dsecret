package pod

import (
	"errors"
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/ink.go/util"
)

func InitPodContract(client *chain.ChainClient, address string) (*Pod, error) {
	contractAddress, err := util.HexToH160(address)
	if err != nil {
		return nil, err
	}
	return &Pod{
		ChainClient: client,
		Address:     contractAddress,
	}, nil
}

func DeployPodWithNew(id uint64, owner types.H160, __ink_params chain.DeployParams) (*types.H160, error) {
	return __ink_params.Client.DeployContract(
		__ink_params.Code, __ink_params.Signer, types.NewU128(*big.NewInt(0)),
		util.InkContractInput{
			Selector: "0x9bae9d5e",
			Args:     []any{id, owner},
		},
		__ink_params.Salt,
	)
}

type Pod struct {
	ChainClient *chain.ChainClient
	Address     types.H160
}

func (c *Pod) Client() *chain.ChainClient {
	return c.ChainClient
}

func (c *Pod) ContractAddress() types.H160 {
	return c.Address
}

func (c *Pod) DryRunCloud(
	params chain.DryRunCallParams,
) (*types.H160, *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[types.H160](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xb24fd0f6",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Pod) CallCloud(
	__ink_params chain.CallParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunCloud(_param)
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
			Selector: "0xb24fd0f6",
			Args:     []any{},
		},
	)
}

func (c *Pod) CallOfCloudTx(
	__ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunCloud(_param)
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
			Selector: "0xb24fd0f6",
			Args:     []any{},
		},
	)
}

func (c *Pod) DryRunApprove(
	value util.Option[types.U256], params chain.DryRunCallParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x681266a0",
			Args:     []any{value},
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

func (c *Pod) CallApprove(
	value util.Option[types.U256], __ink_params chain.CallParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunApprove(value, _param)
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
			Selector: "0x681266a0",
			Args:     []any{value},
		},
	)
}

func (c *Pod) CallOfApproveTx(
	value util.Option[types.U256], __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunApprove(value, _param)
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
			Selector: "0x681266a0",
			Args:     []any{value},
		},
	)
}

func (c *Pod) DryRunPayForWoker(
	worker types.H160, amount types.U256, params chain.DryRunCallParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xd51e3b30",
			Args:     []any{worker, amount},
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

func (c *Pod) CallPayForWoker(
	worker types.H160, amount types.U256, __ink_params chain.CallParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunPayForWoker(worker, amount, _param)
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
			Selector: "0xd51e3b30",
			Args:     []any{worker, amount},
		},
	)
}

func (c *Pod) CallOfPayForWokerTx(
	worker types.H160, amount types.U256, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunPayForWoker(worker, amount, _param)
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
			Selector: "0xd51e3b30",
			Args:     []any{worker, amount},
		},
	)
}

func (c *Pod) DryRunCharge(
	params chain.DryRunCallParams,
) (*util.NullTuple, *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.NullTuple](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x1906ffe6",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Pod) CallCharge(
	__ink_params chain.CallParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunCharge(_param)
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
			Selector: "0x1906ffe6",
			Args:     []any{},
		},
	)
}

func (c *Pod) CallOfChargeTx(
	__ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunCharge(_param)
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
			Selector: "0x1906ffe6",
			Args:     []any{},
		},
	)
}

func (c *Pod) DryRunWithdraw(
	amount types.U256, params chain.DryRunCallParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		params.Origin,
		params.PayAmount,
		params.GasLimit,
		params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x410fcc9d",
			Args:     []any{amount},
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

func (c *Pod) CallWithdraw(
	amount types.U256, __ink_params chain.CallParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunWithdraw(amount, _param)
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
			Selector: "0x410fcc9d",
			Args:     []any{amount},
		},
	)
}

func (c *Pod) CallOfWithdrawTx(
	amount types.U256, __ink_params chain.CallParams,
) (*types.Call, error) {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunWithdraw(amount, _param)
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
			Selector: "0x410fcc9d",
			Args:     []any{amount},
		},
	)
}

func (c *Pod) DryRunSetCode(
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

func (c *Pod) CallSetCode(
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

func (c *Pod) CallOfSetCodeTx(
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
