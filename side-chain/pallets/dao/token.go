package dao

func (d DaoQuery) TotalSupply() ([]byte, error) {
	return d.bytesValue(d.totalIssuance)
}

func (d DaoQuery) BalanceOf(owner []byte) ([]byte, error) {
	member, err := d.member(owner)
	if err != nil || member == nil {
		return nil, err
	}
	lock, err := d.lockOf(owner)
	if err != nil {
		return nil, err
	}
	return sub(member.Balance, lock), nil
}

func (d DaoQuery) LockBalanceOf(owner []byte) ([]byte, error) {
	return d.lockOf(owner)
}

func (d DaoQuery) Allowance(owner, spender []byte) ([]byte, error) {
	return d.allowanceOf(owner, spender)
}

func (d DaoMutation) Transfer(to, value []byte) error {
	enabled, err := d.boolValue(d.transferEnabled, false)
	if err != nil {
		return err
	}
	if !enabled {
		return ErrTransferDisabled
	}
	return d.transfer(d.api.GetCaller(), to, value)
}

func (d DaoMutation) Approve(spender, value []byte) error {
	caller := d.api.GetCaller()
	return d.allowances.Set(d.api.GetTxn(), allowanceKey(caller, spender), cloneBytes(value))
}

func (d DaoMutation) TransferFrom(from, to, value []byte) error {
	enabled, err := d.boolValue(d.transferEnabled, false)
	if err != nil {
		return err
	}
	if !enabled {
		return ErrTransferDisabled
	}

	caller := d.api.GetCaller()
	allowance, err := d.allowanceOf(from, caller)
	if err != nil {
		return err
	}
	if cmp(allowance, value) < 0 {
		return ErrInsufficientAllowance
	}
	if err := d.transfer(from, to, value); err != nil {
		return err
	}
	return d.allowances.Set(d.api.GetTxn(), allowanceKey(from, caller), sub(allowance, value))
}

func (d DaoMutation) transfer(from, to, value []byte) error {
	sender, err := d.member(from)
	if err != nil {
		return err
	}
	receiver, err := d.member(to)
	if err != nil {
		return err
	}
	if sender == nil || receiver == nil {
		return ErrMemberNotExisted
	}
	free := sub(sender.Balance, mustBytes(d.lockOf(from)))
	if cmp(free, value) < 0 {
		return ErrLowBalance
	}
	sender.Balance = sub(sender.Balance, value)
	receiver.Balance = add(receiver.Balance, value)
	if err := d.members.Set(d.api.GetTxn(), from, *sender); err != nil {
		return err
	}
	return d.members.Set(d.api.GetTxn(), to, *receiver)
}
