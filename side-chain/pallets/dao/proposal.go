package dao

import (
	"bytes"
	"fmt"

	"github.com/wetee-dao/ink.go/util"
)

func (d DaoQuery) Proposal(id uint32) (util.Option[Proposal], error) {
	prop, err := d.proposals.Get(d.api.GetTxn(), id)
	if err != nil {
		return util.NewNone[Proposal](), err
	}
	if prop == nil {
		return util.NewNone[Proposal](), ErrInvalidProposal
	}
	return util.NewSome(*prop), nil
}

func (d DaoQuery) Proposals() ([]Proposal, error) {
	return d.proposals.List(d.api.GetTxn())
}

func (d DaoQuery) ProposalStatus(id uint32) (util.Option[ProposalStatus], error) {
	height := d.api.GetHeight()
	prop, err := d.proposal(id)
	if err != nil {
		return util.NewNone[ProposalStatus](), err
	}
	if prop.Status.State != ProposalOngoing {
		return util.NewSome(prop.Status), nil
	}

	confirmed, end, _, err := d.calculateProposalStatus(prop)
	if err != nil {
		return util.NewNone[ProposalStatus](), err
	}
	if !confirmed && height > end {
		status := ProposalStatus{State: ProposalRejected, Block: end}
		return util.NewSome(status), nil
	}
	if confirmed {
		status := ProposalStatus{State: ProposalApproved, Block: 0}
		return util.NewSome(status), nil
	}
	status := ProposalStatus{State: ProposalOngoing, Block: 0}
	return util.NewSome(status), nil
}

func (d DaoMutation) SubmitProposal(call CallContent, trackID uint32) error {
	height := d.api.GetHeight()
	caller := d.api.GetCaller()
	member, err := d.member(caller)
	if err != nil {
		return err
	}
	if member == nil {
		return ErrMemberNotExisted
	}
	track, err := d.tracks.Get(d.api.GetTxn(), trackID)
	if err != nil {
		return err
	}
	if track == nil {
		return ErrNoTrack
	}
	if cmp(call.Amount, track.MaxBalance) > 0 {
		return fmt.Errorf("max balance overflow")
	}

	id, err := d.nextProposalID()
	if err != nil {
		return err
	}
	record := Proposal{
		ID:          id,
		Call:        util.NewSome(call),
		TrackID:     trackID,
		Caller:      cloneBytes(caller),
		Status:      ProposalStatus{State: ProposalPending},
		SubmitBlock: height,
		Deposit: ProposalDeposit{
			Depositor: cloneBytes(caller),
			Amount:    nil,
			Block:     height,
		},
	}
	if err := d.proposals.Set(d.api.GetTxn(), id, record); err != nil {
		return err
	}
	return d.nextProposalIDStore.Set(d.api.GetTxn(), singletonKey, id+1)
}

func (d DaoMutation) CancelProposal(proposalID uint32) error {
	caller := d.api.GetCaller()
	prop, err := d.proposal(proposalID)
	if err != nil {
		return err
	}
	if prop.Status.State != ProposalPending {
		return ErrInvalidProposalStatus
	}
	if !bytes.Equal(prop.Caller, caller) {
		return ErrInvalidProposalCaller
	}
	prop.Status = ProposalStatus{State: ProposalCanceled}
	return d.proposals.Set(d.api.GetTxn(), prop.ID, prop)
}

func (d DaoMutation) DepositProposal(proposalID uint32, amount []byte) error {
	height := d.api.GetHeight()
	caller := d.api.GetCaller()
	prop, err := d.proposal(proposalID)
	if err != nil {
		return err
	}
	if prop.Status.State != ProposalPending {
		return ErrInvalidProposalStatus
	}
	track, err := d.track(prop.TrackID)
	if err != nil {
		return err
	}
	if height < prop.Deposit.Block+int64(track.PreparePeriod) {
		return ErrInvalidDepositTime
	}
	if cmp(amount, track.DecisionDeposit) < 0 {
		return ErrInvalidDeposit
	}
	prop.Deposit = ProposalDeposit{Depositor: cloneBytes(caller), Amount: cloneBytes(amount), Block: height}
	prop.Status = ProposalStatus{State: ProposalOngoing}
	return d.proposals.Set(d.api.GetTxn(), prop.ID, prop)
}

func (d DaoMutation) ExecProposal(proposalID uint32) error {
	height := d.api.GetHeight()
	prop, err := d.proposal(proposalID)
	if err != nil {
		return err
	}
	if prop.Status.State != ProposalOngoing {
		return ErrPropNotOngoing
	}
	confirmed, end, track, err := d.calculateProposalStatus(prop)
	if err != nil {
		return err
	}
	if !confirmed {
		if height > end {
			prop.Status = ProposalStatus{State: ProposalRejected, Block: end}
			return d.proposals.Set(d.api.GetTxn(), prop.ID, prop)
		}
		return ErrProposalNotConfirmed
	}
	if height < end+int64(track.DecisionPeriod) {
		return ErrProposalInDecision
	}
	prop.Status = ProposalStatus{State: ProposalApproved, Block: height}
	prop.DecisionBlock = height
	return d.proposals.Set(d.api.GetTxn(), prop.ID, prop)
}
