package dao

func (d DaoMutation) Spend(to, amount []byte, trackID uint32) error {
	caller := d.api.GetCaller()
	member, err := d.member(caller)
	if err != nil {
		return err
	}
	if member == nil {
		return ErrMemberNotExisted
	}
	id, err := d.nextSpendID()
	if err != nil {
		return err
	}

	record := Spend{
		ID:      id,
		Caller:  cloneBytes(caller),
		To:      cloneBytes(to),
		Amount:  cloneBytes(amount),
		TrackID: trackID,
	}

	if err := d.spends.Set(d.api.GetTxn(), id, record); err != nil {
		return err
	}

	if err := d.nextSpendIDStore.Set(d.api.GetTxn(), singletonKey, id+1); err != nil {
		return err
	}

	return d.SubmitProposal(CallContent{
		Input:  encodeUint64(id),
		Amount: cloneBytes(amount),
	}, trackID)
}

func (d DaoMutation) Payout(spendID uint64) error {
	if err := d.ensureGov(); err != nil {
		return err
	}

	spend, err := d.spends.Get(d.api.GetTxn(), spendID)
	if err != nil {
		return err
	}
	if spend == nil {
		return ErrSpendNotFound
	}
	if spend.Payout {
		return ErrSpendAlreadyExecuted
	}
	spend.Payout = true

	return d.spends.Set(d.api.GetTxn(), spend.ID, *spend)
}
