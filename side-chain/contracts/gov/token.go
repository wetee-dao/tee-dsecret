package gov

import (
	"fmt"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func (d GovQuery) TotalSupply() (model.Amount, error) {
	return d.totalIssuance.GetOrDefault(d.api.GetTxn(), model.ZeroAmount)
}

func (d GovQuery) BalanceOf(owner model.UniAddr) (model.Amount, error) {
	member, err := d.members.GetOrDefault(d.api.GetTxn(), owner, model.ZeroAmount)
	if err != nil {
		return model.ZeroAmount, err
	}
	lock, err := d.memberLocks.GetOrDefault(d.api.GetTxn(), owner, model.ZeroAmount)
	if err != nil {
		return model.ZeroAmount, err
	}

	fmt.Println("member:", member.String(), "lock:", lock.String())

	return model.AmountSub(member, lock), nil
}

func (d GovQuery) LockBalanceOf(owner model.UniAddr) (model.Amount, error) {
	return d.memberLocks.GetOrDefault(d.api.GetTxn(), owner, model.ZeroAmount)
}

func (d GovQuery) Allowance(owner, spender model.UniAddr) (model.Amount, error) {
	return d.allowances.GetOrDefault(d.api.GetTxn(), allowanceKey(owner, spender), model.ZeroAmount)
}

func (d GovMutation) Transfer(to model.UniAddr, value model.Amount) error {
	enabled, err := d.transferEnabled.GetOrDefault(d.api.GetTxn(), false)
	if err != nil {
		return err
	}
	if !enabled {
		return ErrTransferDisabled
	}
	return d.transfer(d.api.GetCaller(), to, value)
}

func (d GovMutation) Approve(spender model.UniAddr, value model.Amount) error {
	caller := d.api.GetCaller()
	return d.allowances.Set(d.api.GetTxn(), allowanceKey(caller, spender), value)
}

func (d GovMutation) TransferFrom(from, to model.UniAddr, amount model.Amount) error {
	enabled, err := d.transferEnabled.GetOrDefault(d.api.GetTxn(), false)
	if err != nil {
		return err
	}
	if !enabled {
		return ErrTransferDisabled
	}

	caller := d.api.GetCaller()
	allowance, err := d.allowances.GetOrDefault(d.api.GetTxn(), allowanceKey(from, caller), model.ZeroAmount)
	if err != nil {
		return err
	}
	if allowance.Int.Cmp(amount.Int) < 0 {
		return ErrInsufficientAllowance
	}
	if err := d.transfer(from, to, amount); err != nil {
		return err
	}
	return d.allowances.Set(d.api.GetTxn(), allowanceKey(from, caller), model.AmountSub(allowance, amount))
}

func (d GovMutation) transfer(from, to model.UniAddr, value model.Amount) error {
	sender, err := d.members.GetOrDefault(d.api.GetTxn(), from, model.ZeroAmount)
	if err != nil {
		return err
	}
	receiver, err := d.members.GetOrDefault(d.api.GetTxn(), to, model.ZeroAmount)
	if err != nil {
		return err
	}
	if sender.Int.Sign() == 0 || receiver.Int.Sign() == 0 {
		return ErrMemberNotExisted
	}

	lock, err := d.memberLocks.GetOrDefault(d.api.GetTxn(), from, model.ZeroAmount)
	if err != nil {
		return err
	}

	free := model.AmountSub(sender, lock)
	if free.Int.Cmp(value.Int) < 0 {
		return ErrLowBalance
	}

	sender = model.AmountSub(sender, value)
	receiver = model.AmountAdd(receiver, value)

	if err := d.members.Set(d.api.GetTxn(), from, sender); err != nil {
		return err
	}
	return d.members.Set(d.api.GetTxn(), to, receiver)
}
