package dao

import (
	"errors"
	"math/big"
	"strconv"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func daoSubmitVote(caller []byte, p *DaoCallPayload, height int64, txn *model.Txn) error {
	state := newDaoStateState(txn)
	b, _ := state.MemberBalance.Get(txn, caller)
	bal := model.BytesToU128(b)
	if bal.Sign() == 0 {
		return errors.New("member not existed")
	}
	lockB, _ := state.MemberLock.Get(txn, caller)
	lock := model.BytesToU128(lockB)
	free := new(big.Int).Sub(bal, lock)
	if free.Cmp(p.LockAmount.ToBigInt()) < 0 {
		return errors.New("low balance")
	}
	stB, _ := state.ProposalStatus.Get(txn, p.ProposalId)
	if string(stB) != "ongoing" {
		return errors.New("proposal not ongoing")
	}
	sb, _ := state.SubmitBlock.Get(txn, p.ProposalId)
	submitBlock := int64(0)
	if len(sb) > 0 {
		submitBlock, _ = strconv.ParseInt(string(sb), 10, 64)
	}
	trB, _ := state.ProposalTrack.Get(txn, p.ProposalId)
	trackId := uint32(0)
	if len(trB) > 0 {
		u, _ := strconv.ParseUint(string(trB), 10, 32)
		trackId = uint32(u)
	}
	track, _ := model.GetMappingJson[uint32, DaoTrackData](state.Tracks, txn, trackId)
	if track != nil && height > submitBlock+int64(track.MaxDeciding) {
		return errors.New("invalid vote time")
	}
	vid := state.NextVoteId()
	_ = state.SetNextVoteId(vid + 1)
	v := &voteInfo{
		Pledge:      p.LockAmount,
		OpinionYes:  p.OpinionYes,
		CallId:      p.ProposalId,
		Caller:      caller,
		VoteBlock:   height,
		UnlockBlock: int64(1),
		Deleted:     false,
	}
	_ = model.SetMappingJson(state.Votes, txn, vid, v)
	_ = state.VoteOfMember.Set(txn, caller, []byte(strconv.FormatUint(vid, 10)))
	curLock, _ := state.MemberLock.Get(txn, caller)
	_ = state.MemberLock.Set(txn, caller, model.U128ToBytes(new(big.Int).Add(model.BytesToU128(curLock), p.LockAmount.ToBigInt())))
	return nil
}

func daoCancelVote(caller []byte, p *DaoCallPayload, height int64, txn *model.Txn) error {
	state := newDaoStateState(txn)
	v, _ := model.GetMappingJson[uint64, voteInfo](state.Votes, txn, p.VoteId)
	if v == nil {
		return errors.New("invalid vote")
	}
	if !bytesEqual(v.Caller, caller) {
		return errors.New("invalid vote user")
	}
	stB, _ := state.ProposalStatus.Get(txn, v.CallId)
	if string(stB) != "ongoing" {
		return errors.New("proposal not ongoing")
	}
	v.Deleted = true
	_ = model.SetMappingJson(state.Votes, txn, p.VoteId, v)
	curLock, _ := state.MemberLock.Get(txn, caller)
	_ = state.MemberLock.Set(txn, caller, model.U128ToBytes(new(big.Int).Sub(model.BytesToU128(curLock), v.Pledge.ToBigInt())))
	return nil
}

func daoUnlock(caller []byte, p *DaoCallPayload, height int64, txn *model.Txn) error {
	state := newDaoStateState(txn)
	unlockB, _ := state.Unlock.Get(txn, p.VoteId)
	if len(unlockB) > 0 {
		return errors.New("vote already unlocked")
	}
	v, _ := model.GetMappingJson[uint64, voteInfo](state.Votes, txn, p.VoteId)
	if v == nil || v.Deleted {
		return errors.New("invalid vote")
	}
	if !bytesEqual(v.Caller, caller) {
		return errors.New("invalid vote user")
	}
	endBlock := daoProposalEndBlock(txn, v.CallId)
	if endBlock < 0 {
		return errors.New("invalid proposal")
	}
	if height < endBlock+v.UnlockBlock {
		return errors.New("invalid vote unlock time")
	}
	_ = state.Unlock.Set(txn, p.VoteId, []byte("1"))
	curLock, _ := state.MemberLock.Get(txn, caller)
	_ = state.MemberLock.Set(txn, caller, model.U128ToBytes(new(big.Int).Sub(model.BytesToU128(curLock), v.Pledge.ToBigInt())))
	return nil
}
