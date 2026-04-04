package dao

import (
	"bytes"
	"math/big"

	"github.com/wetee-dao/ink.go/util"
)

func (d DaoQuery) Vote(id uint64) (util.Option[Vote], error) {
	v, err := d.votes.Get(d.api.GetTxn(), id)
	if err != nil {
		return util.NewNone[Vote](), err
	}
	if v == nil {
		return util.NewNone[Vote](), ErrInvalidVote
	}
	return util.NewSome(*v), nil
}

func (d DaoQuery) Votes() ([]Vote, error) {
	return d.votes.List(d.api.GetTxn())
}

func (d DaoMutation) SubmitVote(proposalID uint32, opinionYes bool, lockAmount []byte) error {
	height := d.api.GetHeight()
	caller := d.api.GetCaller()
	member, err := d.member(caller)
	if err != nil {
		return err
	}
	if member == nil {
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
	free := sub(member.Balance, mustBytes(d.lockOf(caller)))
	if cmp(free, lockAmount) < 0 {
		return ErrLowBalance
	}

	id, err := d.nextVoteID()
	if err != nil {
		return err
	}
	lock, err := d.lockOf(caller)
	if err != nil {
		return err
	}
	vote := Vote{
		ID:          id,
		ProposalID:  proposalID,
		Caller:      cloneBytes(caller),
		Pledge:      cloneBytes(lockAmount),
		OpinionYes:  opinionYes,
		VoteWeight:  encodeAmount(big.NewInt(1)),
		UnlockBlock: 1,
		VoteBlock:   height,
	}
	if err := d.memberLocks.Set(d.api.GetTxn(), caller, add(lock, lockAmount)); err != nil {
		return err
	}
	if err := d.votes.Set(d.api.GetTxn(), id, vote); err != nil {
		return err
	}
	return d.nextVoteIDStore.Set(d.api.GetTxn(), singletonKey, id+1)
}

func (d DaoMutation) CancelVote(voteID uint64) error {
	caller := d.api.GetCaller()
	vote, err := d.vote(voteID)
	if err != nil {
		return err
	}
	if !bytes.Equal(vote.Caller, caller) {
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
	lock, err := d.lockOf(caller)
	if err != nil {
		return err
	}
	if cmp(lock, vote.Pledge) < 0 {
		return ErrLowBalance
	}
	if err := d.memberLocks.Set(d.api.GetTxn(), caller, sub(lock, vote.Pledge)); err != nil {
		return err
	}
	return d.votes.Set(d.api.GetTxn(), vote.ID, vote)
}

func (d DaoMutation) Unlock(voteID uint64) error {
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
	if !bytes.Equal(vote.Caller, caller) {
		return ErrInvalidVoteUser
	}
	end, err := d.calculateProposalEndBlock(vote.ProposalID)
	if err != nil {
		return err
	}
	if height < end+vote.UnlockBlock {
		return ErrInvalidVoteUnlockTime
	}
	lock, err := d.lockOf(caller)
	if err != nil {
		return err
	}
	if cmp(lock, vote.Pledge) < 0 {
		return ErrLowBalance
	}
	if err := d.memberLocks.Set(d.api.GetTxn(), caller, sub(lock, vote.Pledge)); err != nil {
		return err
	}
	return d.voteUnlocks.Set(d.api.GetTxn(), vote.ID, true)
}
