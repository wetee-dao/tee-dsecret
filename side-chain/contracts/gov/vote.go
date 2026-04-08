package gov

import (
	"bytes"
	"math/big"

	"github.com/wetee-dao/ink.go/util"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func (d GovQuery) Vote(id uint64) (util.Option[Vote], error) {
	v, err := d.votes.Get(d.api.GetTxn(), id)
	if err != nil {
		return util.NewNone[Vote](), err
	}
	if v == nil {
		return util.NewNone[Vote](), ErrInvalidVote
	}
	return util.NewSome(*v), nil
}

func (d GovQuery) Votes() ([]Vote, error) {
	_, votes, err := d.votes.List(d.api.GetTxn())
	if err != nil {
		return nil, err
	}
	return votes, nil
}

func (d GovMutation) SubmitVote(proposalID uint32, opinionYes bool, lockAmount model.Amount) error {
	height := d.api.GetHeight()
	caller := d.api.GetCaller()

	member, err := d.members.GetOrDefault(d.api.GetTxn(), caller, model.ZeroAmount)
	if err != nil {
		return err
	}
	if member.Int.Sign() == 0 {
		return ErrMemberNotExisted
	}

	prop, err := d.proposal(proposalID)
	if err != nil {
		return err
	}
	if prop.Status.State != ProposalOngoing {
		return ErrPropNotOngoing
	}
	track, err := d.track(prop.TrackID)
	if err != nil {
		return err
	}
	if height > prop.Deposit.Block+int64(track.MaxDeciding) {
		return ErrInvalidVoteTime
	}

	lock, err := d.memberLocks.GetOrDefault(d.api.GetTxn(), caller, model.ZeroAmount)
	if err != nil {
		return err
	}
	free := model.AmountSub(member, lock)
	if free.Int.Cmp(lockAmount.Int) < 0 {
		return ErrLowBalance
	}

	id, err := d.nextVoteIDStore.GetOrDefault(d.api.GetTxn(), uint64(0))
	if err != nil {
		return err
	}
	unlockAfter := int64(track.MinEnactmentPeriod)

	one := big.NewInt(1)
	vote := Vote{
		ID:          id,
		ProposalID:  proposalID,
		Caller:      caller,
		Pledge:      lockAmount,
		OpinionYes:  opinionYes,
		VoteWeight:  model.Amount{Int: one},
		UnlockBlock: unlockAfter,
		VoteBlock:   height,
	}

	if err := d.memberLocks.Set(d.api.GetTxn(), caller, model.AmountAdd(lock, lockAmount)); err != nil {
		return err
	}
	if err := d.votes.Set(d.api.GetTxn(), id, vote); err != nil {
		return err
	}
	return d.nextVoteIDStore.Set(d.api.GetTxn(), id+1)
}

func (d GovMutation) CancelVote(voteID uint64) error {
	caller := d.api.GetCaller()
	vote, err := d.vote(voteID)
	if err != nil {
		return err
	}
	if !bytes.Equal(vote.Caller.V, caller.V) {
		return ErrInvalidVoteUser
	}
	prop, err := d.proposal(vote.ProposalID)
	if err != nil {
		return err
	}
	if prop.Status.State != ProposalOngoing {
		return ErrPropNotOngoing
	}
	vote.Deleted = true

	lock, err := d.memberLocks.GetOrDefault(d.api.GetTxn(), caller, model.ZeroAmount)
	if err != nil {
		return err
	}
	if lock.Int.Cmp(vote.Pledge.Int) < 0 {
		return ErrLowBalance
	}

	if err := d.memberLocks.Set(d.api.GetTxn(), caller, model.AmountSub(lock, vote.Pledge)); err != nil {
		return err
	}
	return d.votes.Set(d.api.GetTxn(), vote.ID, vote)
}

func (d GovMutation) Unlock(voteID uint64) error {
	height := d.api.GetHeight()
	caller := d.api.GetCaller()

	unlocked, err := d.voteUnlocks.Get(d.api.GetTxn(), voteID)
	if err != nil {
		return err
	}
	if unlocked != nil && *unlocked {
		return ErrVoteAlreadyUnlocked
	}

	vote, err := d.vote(voteID)
	if err != nil {
		return err
	}
	if vote.Deleted {
		return ErrInvalidVoteStatus
	}
	if !bytes.Equal(vote.Caller.V, caller.V) {
		return ErrInvalidVoteUser
	}

	end, err := d.calculateProposalEndBlock(vote.ProposalID)
	if err != nil {
		return err
	}
	if height < end+vote.UnlockBlock {
		return ErrInvalidVoteUnlockTime
	}

	lock, err := d.memberLocks.GetOrDefault(d.api.GetTxn(), caller, model.ZeroAmount)
	if err != nil {
		return err
	}
	if lock.Int.Cmp(vote.Pledge.Int) < 0 {
		return ErrLowBalance
	}

	if err := d.memberLocks.Set(d.api.GetTxn(), caller, model.AmountSub(lock, vote.Pledge)); err != nil {
		return err
	}
	return d.voteUnlocks.Set(d.api.GetTxn(), vote.ID, true)
}

func (d Gov) calculateProposalStatus(prop Proposal) (bool, int64, TrackData, error) {
	track, err := d.track(prop.TrackID)
	if err != nil {
		return false, 0, TrackData{}, err
	}
	end := prop.Deposit.Block + int64(track.MaxDeciding)

	_, votes, err := d.votes.List(d.api.GetTxn())
	if err != nil {
		return false, 0, TrackData{}, err
	}

	totalSupply, err := d.totalIssuance.GetOrDefault(d.api.GetTxn(), model.ZeroAmount)
	if err != nil {
		return false, 0, TrackData{}, err
	}

	yes := big.NewInt(0)
	support := big.NewInt(0)
	all := totalSupply.Int

	for _, vote := range votes {
		if vote.ProposalID != prop.ID || vote.Deleted {
			continue
		}
		pledge := vote.Pledge.Int
		support.Add(support, pledge)
		if vote.OpinionYes {
			yes.Add(yes, pledge)
		}
	}

	if all.Sign() == 0 {
		return false, end, track, nil
	}

	depositInt := track.DecisionDeposit
	if yes.Sign() > 0 && support.Cmp(depositInt.Int) >= 0 {
		return true, end, track, nil
	}
	return false, end, track, nil
}

func (d Gov) calculateProposalEndBlock(id uint32) (int64, error) {
	height := d.api.GetHeight()
	prop, err := d.proposal(id)
	if err != nil {
		return 0, err
	}
	switch prop.Status.State {
	case ProposalRejected, ProposalApproved:
		return prop.Status.Block, nil
	case ProposalOngoing:
		confirmed, end, _, err := d.calculateProposalStatus(prop)
		if err != nil {
			return 0, err
		}
		if !confirmed && height > end {
			return end, nil
		}
	}
	return 0, ErrInvalidProposalStatus
}
