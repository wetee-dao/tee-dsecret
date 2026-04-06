package dao

func (d DaoQuery) Members() ([]Member, error) {
	// return d.members.List(d.api.GetTxn())
	return []Member{}, nil
}

func (d DaoQuery) GetPublicJoin() (bool, error) {
	return d.boolValue(d.publicJoin, false)
}

func (d DaoMutation) PublicJoin() error {
	enabled, err := d.boolValue(d.publicJoin, false)
	if err != nil {
		return err
	}
	if !enabled {
		return ErrPublicJoinNotAllowed
	}
	return d.addMember(d.api.GetCaller(), nil)
}

func (d DaoMutation) SetPublicJoin(publicJoin bool) error {
	if err := d.ensureGov(); err != nil {
		return err
	}
	return d.publicJoin.Set(d.api.GetTxn(), singletonKey, publicJoin)
}

func (d DaoMutation) Join(newUser, balance []byte) error {
	if err := d.ensureGov(); err != nil {
		return err
	}
	return d.addMember(newUser, balance)
}

func (d DaoMutation) Leave() error {
	caller := d.api.GetCaller()
	member, err := d.member(caller)
	if err != nil {
		return err
	}
	if member == nil {
		return ErrMemberNotExisted
	}
	lock, err := d.lockOf(caller)
	if err != nil {
		return err
	}
	if !isZero(member.Balance) || !isZero(lock) {
		return ErrMemberBalanceNotZero
	}
	return d.removeMember(caller)
}

func (d DaoMutation) LeaveWithBurn() error {
	caller := d.api.GetCaller()
	return d.deleteMemberAndBurn(caller, false)
}

func (d DaoMutation) Delete(account []byte) error {
	return d.deleteMemberAndBurn(account, true)
}

func (d DAO) addMember(account, balance []byte) error {
	if len(account) == 0 {
		return ErrMemberNotExisted
	}
	member, err := d.member(account)
	if err != nil {
		return err
	}
	if member != nil {
		return ErrMemberExisted
	}
	record := Member{Account: cloneBytes(account), Balance: cloneBytes(balance)}
	if err := d.members.Set(d.api.GetTxn(), account, record); err != nil {
		return err
	}
	total, err := d.bytesValue(d.totalIssuance)
	if err != nil {
		return err
	}
	return d.totalIssuance.Set(d.api.GetTxn(), singletonKey, add(total, balance))
}

func (d DAO) removeMember(account []byte) error {
	if err := d.members.Delete(d.api.GetTxn(), account); err != nil {
		return err
	}
	return d.memberLocks.Delete(d.api.GetTxn(), account)
}

func (d DAO) deleteMemberAndBurn(account []byte, requireGov bool) error {
	if requireGov {
		if err := d.ensureGov(); err != nil {
			return err
		}
	}
	member, err := d.member(account)
	if err != nil {
		return err
	}
	if member == nil {
		return ErrMemberNotExisted
	}
	lock, err := d.lockOf(account)
	if err != nil {
		return err
	}
	amount := add(member.Balance, lock)
	total, err := d.bytesValue(d.totalIssuance)
	if err != nil {
		return err
	}
	if cmp(total, amount) < 0 {
		return ErrLowBalance
	}
	if err := d.totalIssuance.Set(d.api.GetTxn(), singletonKey, sub(total, amount)); err != nil {
		return err
	}
	return d.removeMember(account)
}
