package gov

import (
	"errors"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/ink.go/util"
)

func InitGovContract(client *chain.ChainClient, address string) (*Gov, error) {
	contractAddress, err := util.HexToH160(address)
	if err != nil {
		return nil, err
	}
	return &Gov{
		ChainClient: client,
		Address:     contractAddress,
	}, nil
}

type Gov struct {
	ChainClient *chain.ChainClient
	Address     types.H160
}

func (c *Gov) Client() *chain.ChainClient {
	return c.ChainClient
}

func (c *Gov) ContractAddress() types.H160 {
	return c.Address
}

func (c *Gov) DryRunCancelVote(
	vote_i_d uint64, __ink_params chain.DryRunParams,
) (*Error, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "cancel_vote")
	}
	v, gas, err := chain.DryRunInk[Error](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x06991e89",
			Args:     []any{vote_i_d},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) ExecCancelVote(
	vote_i_d uint64, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunCancelVote(vote_i_d, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x06991e89",
			Args:     []any{vote_i_d},
		},
		__ink_params,
	)
}

func (c *Gov) CallOfCancelVote(
	vote_i_d uint64, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunCancelVote(vote_i_d, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x06991e89",
			Args:     []any{vote_i_d},
		},
	)
}

func (c *Gov) DryRunSetPublicJoin(
	public_join bool, __ink_params chain.DryRunParams,
) (*Error, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "set_public_join")
	}
	v, gas, err := chain.DryRunInk[Error](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x13647515",
			Args:     []any{public_join},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) ExecSetPublicJoin(
	public_join bool, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSetPublicJoin(public_join, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x13647515",
			Args:     []any{public_join},
		},
		__ink_params,
	)
}

func (c *Gov) CallOfSetPublicJoin(
	public_join bool, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunSetPublicJoin(public_join, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x13647515",
			Args:     []any{public_join},
		},
	)
}

func (c *Gov) DryRunApprove(
	spender UniAddr, value types.U256, __ink_params chain.DryRunParams,
) (*Error, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "approve")
	}
	v, gas, err := chain.DryRunInk[Error](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x13960eb8",
			Args:     []any{spender, value},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) ExecApprove(
	spender UniAddr, value types.U256, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunApprove(spender, value, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x13960eb8",
			Args:     []any{spender, value},
		},
		__ink_params,
	)
}

func (c *Gov) CallOfApprove(
	spender UniAddr, value types.U256, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunApprove(spender, value, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x13960eb8",
			Args:     []any{spender, value},
		},
	)
}

func (c *Gov) DryRunSubmitProposal(
	call CallContent, track_i_d uint32, __ink_params chain.DryRunParams,
) (*Error, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "submit_proposal")
	}
	v, gas, err := chain.DryRunInk[Error](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x1c0471c4",
			Args:     []any{call, track_i_d},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) ExecSubmitProposal(
	call CallContent, track_i_d uint32, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSubmitProposal(call, track_i_d, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x1c0471c4",
			Args:     []any{call, track_i_d},
		},
		__ink_params,
	)
}

func (c *Gov) CallOfSubmitProposal(
	call CallContent, track_i_d uint32, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunSubmitProposal(call, track_i_d, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x1c0471c4",
			Args:     []any{call, track_i_d},
		},
	)
}

func (c *Gov) QueryAllowance(
	owner UniAddr, spender UniAddr, __ink_params chain.DryRunParams,
) (*types.U256, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "allowance")
	}
	v, gas, err := chain.DryRunInk[types.U256](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x1d623327",
			Args:     []any{owner, spender},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) DryRunLeave(
	__ink_params chain.DryRunParams,
) (*Error, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "leave")
	}
	v, gas, err := chain.DryRunInk[Error](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x1fac3dd3",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) ExecLeave(
	__ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunLeave(_param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x1fac3dd3",
			Args:     []any{},
		},
		__ink_params,
	)
}

func (c *Gov) CallOfLeave(
	__ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunLeave(__ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x1fac3dd3",
			Args:     []any{},
		},
	)
}

func (c *Gov) DryRunSetDefaultTrack(
	track_i_d uint32, __ink_params chain.DryRunParams,
) (*Error, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "set_default_track")
	}
	v, gas, err := chain.DryRunInk[Error](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x2987d4ba",
			Args:     []any{track_i_d},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) ExecSetDefaultTrack(
	track_i_d uint32, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSetDefaultTrack(track_i_d, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x2987d4ba",
			Args:     []any{track_i_d},
		},
		__ink_params,
	)
}

func (c *Gov) CallOfSetDefaultTrack(
	track_i_d uint32, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunSetDefaultTrack(track_i_d, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x2987d4ba",
			Args:     []any{track_i_d},
		},
	)
}

func (c *Gov) DryRunUnlock(
	vote_i_d uint64, __ink_params chain.DryRunParams,
) (*Error, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "unlock")
	}
	v, gas, err := chain.DryRunInk[Error](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x2ab81dd1",
			Args:     []any{vote_i_d},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) ExecUnlock(
	vote_i_d uint64, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunUnlock(vote_i_d, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x2ab81dd1",
			Args:     []any{vote_i_d},
		},
		__ink_params,
	)
}

func (c *Gov) CallOfUnlock(
	vote_i_d uint64, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunUnlock(vote_i_d, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x2ab81dd1",
			Args:     []any{vote_i_d},
		},
	)
}

func (c *Gov) QueryVotes(
	__ink_params chain.DryRunParams,
) (*[]Vote, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "votes")
	}
	v, gas, err := chain.DryRunInk[[]Vote](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x340f2cf4",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) QueryDefaultTrack(
	__ink_params chain.DryRunParams,
) (*uint32, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "default_track")
	}
	v, gas, err := chain.DryRunInk[uint32](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x3626e9bf",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) QueryTracks(
	__ink_params chain.DryRunParams,
) (*[]TrackData, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "tracks")
	}
	v, gas, err := chain.DryRunInk[[]TrackData](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x3c05be89",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) QueryLockBalanceOf(
	owner UniAddr, __ink_params chain.DryRunParams,
) (*types.U256, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "lock_balance_of")
	}
	v, gas, err := chain.DryRunInk[types.U256](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x3cd4966d",
			Args:     []any{owner},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) DryRunExecProposal(
	proposal_i_d uint32, __ink_params chain.DryRunParams,
) (*Error, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "exec_proposal")
	}
	v, gas, err := chain.DryRunInk[Error](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x3e44340e",
			Args:     []any{proposal_i_d},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) ExecExecProposal(
	proposal_i_d uint32, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunExecProposal(proposal_i_d, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x3e44340e",
			Args:     []any{proposal_i_d},
		},
		__ink_params,
	)
}

func (c *Gov) CallOfExecProposal(
	proposal_i_d uint32, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunExecProposal(proposal_i_d, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x3e44340e",
			Args:     []any{proposal_i_d},
		},
	)
}

func (c *Gov) DryRunCancelProposal(
	proposal_i_d uint32, __ink_params chain.DryRunParams,
) (*Error, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "cancel_proposal")
	}
	v, gas, err := chain.DryRunInk[Error](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x46c421db",
			Args:     []any{proposal_i_d},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) ExecCancelProposal(
	proposal_i_d uint32, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunCancelProposal(proposal_i_d, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x46c421db",
			Args:     []any{proposal_i_d},
		},
		__ink_params,
	)
}

func (c *Gov) CallOfCancelProposal(
	proposal_i_d uint32, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunCancelProposal(proposal_i_d, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x46c421db",
			Args:     []any{proposal_i_d},
		},
	)
}

func (c *Gov) QueryVote(
	id uint64, __ink_params chain.DryRunParams,
) (*util.Option[Vote], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "vote")
	}
	v, gas, err := chain.DryRunInk[util.Option[Vote]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x4d151cf1",
			Args:     []any{id},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) QueryMembers(
	__ink_params chain.DryRunParams,
) (*[]Member, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "members")
	}
	v, gas, err := chain.DryRunInk[[]Member](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x53550f05",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) DryRunInit(
	__ink_params chain.DryRunParams,
) (*Error, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "init")
	}
	v, gas, err := chain.DryRunInk[Error](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x53e39003",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) ExecInit(
	__ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunInit(_param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x53e39003",
			Args:     []any{},
		},
		__ink_params,
	)
}

func (c *Gov) CallOfInit(
	__ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunInit(__ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x53e39003",
			Args:     []any{},
		},
	)
}

func (c *Gov) DryRunPublicJoin(
	__ink_params chain.DryRunParams,
) (*Error, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "public_join")
	}
	v, gas, err := chain.DryRunInk[Error](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x6d6f0772",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) ExecPublicJoin(
	__ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunPublicJoin(_param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x6d6f0772",
			Args:     []any{},
		},
		__ink_params,
	)
}

func (c *Gov) CallOfPublicJoin(
	__ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunPublicJoin(__ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x6d6f0772",
			Args:     []any{},
		},
	)
}

func (c *Gov) QueryGetPublicJoin(
	__ink_params chain.DryRunParams,
) (*bool, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "get_public_join")
	}
	v, gas, err := chain.DryRunInk[bool](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x76f20724",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) QueryTotalSupply(
	__ink_params chain.DryRunParams,
) (*types.U256, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "total_supply")
	}
	v, gas, err := chain.DryRunInk[types.U256](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x771771f7",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) QueryProposal(
	id uint32, __ink_params chain.DryRunParams,
) (*util.Option[Proposal], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "proposal")
	}
	v, gas, err := chain.DryRunInk[util.Option[Proposal]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x799fe2dd",
			Args:     []any{id},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) DryRunSpend(
	to UniAddr, amount types.U256, track_i_d uint32, __ink_params chain.DryRunParams,
) (*Error, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "spend")
	}
	v, gas, err := chain.DryRunInk[Error](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x872c40e9",
			Args:     []any{to, amount, track_i_d},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) ExecSpend(
	to UniAddr, amount types.U256, track_i_d uint32, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSpend(to, amount, track_i_d, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x872c40e9",
			Args:     []any{to, amount, track_i_d},
		},
		__ink_params,
	)
}

func (c *Gov) CallOfSpend(
	to UniAddr, amount types.U256, track_i_d uint32, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunSpend(to, amount, track_i_d, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x872c40e9",
			Args:     []any{to, amount, track_i_d},
		},
	)
}

func (c *Gov) QueryBalanceOf(
	owner UniAddr, __ink_params chain.DryRunParams,
) (*types.U256, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "balance_of")
	}
	v, gas, err := chain.DryRunInk[types.U256](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x87e0ff32",
			Args:     []any{owner},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) DryRunLeaveWithBurn(
	__ink_params chain.DryRunParams,
) (*Error, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "leave_with_burn")
	}
	v, gas, err := chain.DryRunInk[Error](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x8d6e0f8c",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) ExecLeaveWithBurn(
	__ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunLeaveWithBurn(_param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x8d6e0f8c",
			Args:     []any{},
		},
		__ink_params,
	)
}

func (c *Gov) CallOfLeaveWithBurn(
	__ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunLeaveWithBurn(__ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x8d6e0f8c",
			Args:     []any{},
		},
	)
}

func (c *Gov) DryRunPayout(
	spend_i_d uint64, __ink_params chain.DryRunParams,
) (*Error, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "payout")
	}
	v, gas, err := chain.DryRunInk[Error](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xa4f0196c",
			Args:     []any{spend_i_d},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) ExecPayout(
	spend_i_d uint64, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunPayout(spend_i_d, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xa4f0196c",
			Args:     []any{spend_i_d},
		},
		__ink_params,
	)
}

func (c *Gov) CallOfPayout(
	spend_i_d uint64, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunPayout(spend_i_d, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xa4f0196c",
			Args:     []any{spend_i_d},
		},
	)
}

func (c *Gov) DryRunTransfer(
	to UniAddr, value types.U256, __ink_params chain.DryRunParams,
) (*Error, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "transfer")
	}
	v, gas, err := chain.DryRunInk[Error](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xada06462",
			Args:     []any{to, value},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) ExecTransfer(
	to UniAddr, value types.U256, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunTransfer(to, value, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xada06462",
			Args:     []any{to, value},
		},
		__ink_params,
	)
}

func (c *Gov) CallOfTransfer(
	to UniAddr, value types.U256, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunTransfer(to, value, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xada06462",
			Args:     []any{to, value},
		},
	)
}

func (c *Gov) QueryTrack(
	id uint32, __ink_params chain.DryRunParams,
) (*util.Option[TrackData], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "track")
	}
	v, gas, err := chain.DryRunInk[util.Option[TrackData]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xb9fd54c4",
			Args:     []any{id},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) DryRunSubmitVote(
	proposal_i_d uint32, opinion_yes bool, lock_amount types.U256, __ink_params chain.DryRunParams,
) (*Error, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "submit_vote")
	}
	v, gas, err := chain.DryRunInk[Error](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xbb5ca128",
			Args:     []any{proposal_i_d, opinion_yes, lock_amount},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) ExecSubmitVote(
	proposal_i_d uint32, opinion_yes bool, lock_amount types.U256, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSubmitVote(proposal_i_d, opinion_yes, lock_amount, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xbb5ca128",
			Args:     []any{proposal_i_d, opinion_yes, lock_amount},
		},
		__ink_params,
	)
}

func (c *Gov) CallOfSubmitVote(
	proposal_i_d uint32, opinion_yes bool, lock_amount types.U256, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunSubmitVote(proposal_i_d, opinion_yes, lock_amount, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xbb5ca128",
			Args:     []any{proposal_i_d, opinion_yes, lock_amount},
		},
	)
}

func (c *Gov) DryRunJoin(
	new_user UniAddr, balance types.U256, __ink_params chain.DryRunParams,
) (*Error, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "join")
	}
	v, gas, err := chain.DryRunInk[Error](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xc6c83af9",
			Args:     []any{new_user, balance},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) ExecJoin(
	new_user UniAddr, balance types.U256, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunJoin(new_user, balance, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xc6c83af9",
			Args:     []any{new_user, balance},
		},
		__ink_params,
	)
}

func (c *Gov) CallOfJoin(
	new_user UniAddr, balance types.U256, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunJoin(new_user, balance, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xc6c83af9",
			Args:     []any{new_user, balance},
		},
	)
}

func (c *Gov) DryRunAddTrack(
	track TrackData, __ink_params chain.DryRunParams,
) (*Error, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "add_track")
	}
	v, gas, err := chain.DryRunInk[Error](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xdf60a515",
			Args:     []any{track},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) ExecAddTrack(
	track TrackData, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunAddTrack(track, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xdf60a515",
			Args:     []any{track},
		},
		__ink_params,
	)
}

func (c *Gov) CallOfAddTrack(
	track TrackData, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunAddTrack(track, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xdf60a515",
			Args:     []any{track},
		},
	)
}

func (c *Gov) DryRunDepositProposal(
	proposal_i_d uint32, amount types.U256, __ink_params chain.DryRunParams,
) (*Error, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "deposit_proposal")
	}
	v, gas, err := chain.DryRunInk[Error](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xea58b31d",
			Args:     []any{proposal_i_d, amount},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) ExecDepositProposal(
	proposal_i_d uint32, amount types.U256, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunDepositProposal(proposal_i_d, amount, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xea58b31d",
			Args:     []any{proposal_i_d, amount},
		},
		__ink_params,
	)
}

func (c *Gov) CallOfDepositProposal(
	proposal_i_d uint32, amount types.U256, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunDepositProposal(proposal_i_d, amount, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xea58b31d",
			Args:     []any{proposal_i_d, amount},
		},
	)
}

func (c *Gov) DryRunTransferFrom(
	from UniAddr, to UniAddr, amount types.U256, __ink_params chain.DryRunParams,
) (*Error, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "transfer_from")
	}
	v, gas, err := chain.DryRunInk[Error](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xee3bc685",
			Args:     []any{from, to, amount},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) ExecTransferFrom(
	from UniAddr, to UniAddr, amount types.U256, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunTransferFrom(from, to, amount, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xee3bc685",
			Args:     []any{from, to, amount},
		},
		__ink_params,
	)
}

func (c *Gov) CallOfTransferFrom(
	from UniAddr, to UniAddr, amount types.U256, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunTransferFrom(from, to, amount, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xee3bc685",
			Args:     []any{from, to, amount},
		},
	)
}

func (c *Gov) QueryProposalStatus(
	id uint32, __ink_params chain.DryRunParams,
) (*util.Option[ProposalStatus], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "proposal_status")
	}
	v, gas, err := chain.DryRunInk[util.Option[ProposalStatus]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xfd88907b",
			Args:     []any{id},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Gov) QueryProposals(
	__ink_params chain.DryRunParams,
) (*[]Proposal, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "proposals")
	}
	v, gas, err := chain.DryRunInk[[]Proposal](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xfe26fd17",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}
