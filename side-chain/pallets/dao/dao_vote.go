package dao

import (
	"errors"
	"math/big"
	"strconv"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func daoSubmitVote(caller []byte, m *model.DaoSubmitVote, height int64, txn *model.Txn) error {
	state := newDaoStateState(txn)
	b, _ := state.MemberBalance.Get(txn, caller)
	bal := model.BytesToU128(b)
	if bal.Sign() == 0 {
		return errors.New("member not existed")
	}
	lockB, _ := state.MemberLock.Get(txn, caller)
	lock := model.BytesToU128(lockB)
	free := new(big.Int).Sub(bal, lock)
	lockAmount := model.BytesToU128(m.GetLockAmount())
	if free.Cmp(lockAmount) < 0 {
		return errors.New("low balance")
	}
	proposalId := m.GetProposalId()
	stB, _ := state.ProposalStatus.Get(txn, proposalId)
	if string(stB) != "ongoing" {
		return errors.New("proposal not ongoing")
	}
	sb, _ := state.SubmitBlock.Get(txn, proposalId)
	submitBlock := int64(0)
	if len(sb) > 0 {
		submitBlock, _ = strconv.ParseInt(string(sb), 10, 64)
	}
	trB, _ := state.ProposalTrack.Get(txn, proposalId)
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
		Pledge:      model.NewU128(lockAmount),
		OpinionYes:  m.GetOpinionYes(),
		CallId:      proposalId,
		Caller:      caller,
		VoteBlock:   height,
		UnlockBlock: int64(1),
		Deleted:     false,
	}
	_ = model.SetMappingJson(state.Votes, txn, vid, v)
	_ = state.VoteOfMember.Set(txn, caller, []byte(strconv.FormatUint(vid, 10)))
	curLock, _ := state.MemberLock.Get(txn, caller)
	_ = state.MemberLock.Set(txn, caller, model.U128ToBytes(new(big.Int).Add(model.BytesToU128(curLock), lockAmount)))
	return nil
}

func daoCancelVote(caller []byte, m *model.DaoCancelVote, height int64, txn *model.Txn) error {
	state := newDaoStateState(txn)
	voteId := m.GetVoteId()
	v, _ := model.GetMappingJson[uint64, voteInfo](state.Votes, txn, voteId)
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
	_ = model.SetMappingJson(state.Votes, txn, voteId, v)
	curLock, _ := state.MemberLock.Get(txn, caller)
	_ = state.MemberLock.Set(txn, caller, model.U128ToBytes(new(big.Int).Sub(model.BytesToU128(curLock), v.Pledge.ToBigInt())))
	return nil
}

func daoUnlock(caller []byte, m *model.DaoUnlock, height int64, txn *model.Txn) error {
	state := newDaoStateState(txn)
	voteId := m.GetVoteId()
	unlockB, _ := state.Unlock.Get(txn, voteId)
	if len(unlockB) > 0 {
		return errors.New("vote already unlocked")
	}
	v, _ := model.GetMappingJson[uint64, voteInfo](state.Votes, txn, voteId)
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
	_ = state.Unlock.Set(txn, voteId, []byte("1"))
	curLock, _ := state.MemberLock.Get(txn, caller)
	_ = state.MemberLock.Set(txn, caller, model.U128ToBytes(new(big.Int).Sub(model.BytesToU128(curLock), v.Pledge.ToBigInt())))
	return nil
}

// VoteList 返回某提案下的投票列表（对应 ink vote_list(proposal_id)）。
func VoteList(proposalId uint32) []*voteInfo {
	list, _, _ := model.GetJsonList[voteInfo]("dao", stateKeydaoStateVotes)
	if list == nil {
		return nil
	}
	var out []*voteInfo
	for _, v := range list {
		if v != nil && v.CallId == proposalId && !v.Deleted {
			out = append(out, v)
		}
	}
	return out
}

// Vote 按 vote_id 返回投票信息（对应 ink vote(vote_id)）。
func Vote(voteId uint64) *voteInfo {
	key := stateKeydaoStateVotes + strconv.FormatUint(voteId, 10)
	val, _ := model.GetJson[voteInfo]("dao", key)
	return val
}
