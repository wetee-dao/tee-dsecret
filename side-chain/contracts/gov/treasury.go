package gov

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func (d GovMutation) Spend(to model.UniAddr, amount model.Amount, trackID uint32) error {
	caller := d.api.GetCaller()

	member, err := d.members.GetOrDefault(d.api.GetTxn(), caller, model.ZeroAmount)
	if err != nil {
		return err
	}
	if member.Int.Sign() == 0 {
		return ErrMemberNotExisted
	}

	id, err := d.nextSpendIDStore.GetOrDefault(d.api.GetTxn(), uint64(0))
	if err != nil {
		return err
	}

	record := Spend{
		ID:      id,
		Caller:  caller,
		To:      to,
		Amount:  amount,
		TrackID: trackID,
	}

	if err := d.spends.Set(d.api.GetTxn(), id, record); err != nil {
		return err
	}

	if err := d.nextSpendIDStore.Set(d.api.GetTxn(), id+1); err != nil {
		return err
	}

	idbt, _ := codec.Encode(id)
	return d.SubmitProposal(CallContent{
		Contract: []byte("gov"),
		Selector: model.MethodToSelector("Payout"),
		Args:     [][]byte{idbt},
		Amount:   amount,
	}, trackID)
}

func (d GovMutation) Payout(spendID uint64) error {
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
