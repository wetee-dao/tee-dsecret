package dao

import (
	"errors"
	"math/big"
	"strconv"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func daoSubmitProposal(caller []byte, p *DaoCallPayload, height int64, txn *model.Txn) error {
	state := newDaoStateState(txn)
	b, _ := state.MemberBalance.Get(txn, caller)
	if model.BytesToU128(b).Sign() == 0 {
		return errors.New("member not existed")
	}
	st, err := loadDaoState(txn)
	if err != nil || st == nil {
		return errors.New("dao not initialized")
	}
	if p.Call == nil {
		return errors.New("call required")
	}
	if int(p.TrackId) >= len(getDaoTracks(txn)) {
		return errors.New("no track or invalid track")
	}
	track, _ := model.GetMappingJson[uint32, DaoTrackData](state.Tracks, txn, p.TrackId)
	if track != nil && p.Call.Amount.Sign() != 0 {
		if p.Call.Amount.ToBigInt().Cmp(track.MaxBalance.ToBigInt()) > 0 {
			return errors.New("max balance overflow")
		}
	}
	proposalId := state.NextProposalId()
	_ = state.SetNextProposalId(proposalId + 1)
	_ = model.SetMappingJson(state.Proposals, txn, proposalId, p.Call)
	_ = state.ProposalStatus.Set(txn, proposalId, []byte("pending"))
	_ = state.ProposalTrack.Set(txn, proposalId, []byte(strconv.FormatUint(uint64(p.TrackId), 10)))
	_ = state.ProposalCaller.Set(txn, proposalId, caller)
	_ = state.SubmitBlock.Set(txn, proposalId, []byte(strconv.FormatInt(height, 10)))
	return nil
}

func daoDepositProposal(caller []byte, p *DaoCallPayload, height int64, txn *model.Txn) error {
	state := newDaoStateState(txn)
	b, _ := state.ProposalStatus.Get(txn, p.ProposalId)
	if string(b) != "pending" {
		return errors.New("invalid proposal status")
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
	if track == nil {
		return errors.New("no track")
	}
	if height < submitBlock+int64(track.PreparePeriod) {
		return errors.New("invalid deposit time")
	}
	depositNeed := track.DecisionDeposit.ToBigInt()
	if p.Amount.ToBigInt().Cmp(depositNeed) < 0 {
		return errors.New("invalid deposit")
	}
	_ = model.SetMappingJson(state.Deposits, txn, p.ProposalId, &proposalDeposit{Depositor: caller, Amount: p.Amount, Block: height})
	_ = state.ProposalStatus.Set(txn, p.ProposalId, []byte("ongoing"))
	return nil
}

func daoProposalEndBlock(txn *model.Txn, proposalId uint32) int64 {
	state := newDaoStateState(txn)
	trB, _ := state.ProposalTrack.Get(txn, proposalId)
	trackId := uint32(0)
	if len(trB) > 0 {
		u, _ := strconv.ParseUint(string(trB), 10, 32)
		trackId = uint32(u)
	}
	track, _ := model.GetMappingJson[uint32, DaoTrackData](state.Tracks, txn, trackId)
	if track == nil {
		return -1
	}
	dep, _ := model.GetMappingJson[uint32, proposalDeposit](state.Deposits, txn, proposalId)
	if dep == nil {
		return -1
	}
	return dep.Block + int64(track.MaxDeciding)
}

func daoExecProposal(caller []byte, p *DaoCallPayload, height int64, txn *model.Txn) error {
	state := newDaoStateState(txn)
	stB, _ := state.ProposalStatus.Get(txn, p.ProposalId)
	if string(stB) != "ongoing" {
		return errors.New("proposal not ongoing")
	}
	endBlock := daoProposalEndBlock(txn, p.ProposalId)
	trB, _ := state.ProposalTrack.Get(txn, p.ProposalId)
	trackId := uint32(0)
	if len(trB) > 0 {
		u, _ := strconv.ParseUint(string(trB), 10, 32)
		trackId = uint32(u)
	}
	track, _ := model.GetMappingJson[uint32, DaoTrackData](state.Tracks, txn, trackId)
	if track == nil {
		return errors.New("no track")
	}
	if height <= endBlock+int64(track.ConfirmPeriod) {
		return errors.New("proposal in decision")
	}
	call, _ := model.GetMappingJson[uint32, DaoCallContent](state.Proposals, txn, p.ProposalId)
	if call == nil {
		return errors.New("invalid proposal")
	}
	_ = state.ProposalStatus.Set(txn, p.ProposalId, []byte("approved"))
	if len(call.Selector) >= 4 && len(call.Input) >= 8 {
		var spendId uint64
		for i := 0; i < 8; i++ {
			spendId |= uint64(call.Input[i]) << (i * 8)
		}
		s, _ := model.GetMappingJson[uint64, spendRecord](state.Spends, txn, spendId)
		if s != nil && !s.Payout {
			toB, _ := state.MemberBalance.Get(txn, s.To)
			if model.BytesToU128(toB).Sign() != 0 {
				cur := model.BytesToU128(toB)
				_ = state.MemberBalance.Set(txn, s.To, model.U128ToBytes(new(big.Int).Add(cur, s.Amount.ToBigInt())))
			}
			s.Payout = true
			_ = model.SetMappingJson(state.Spends, txn, spendId, s)
		}
	}
	return nil
}

func daoCancelProposal(caller []byte, p *DaoCallPayload, height int64, txn *model.Txn) error {
	state := newDaoStateState(txn)
	statusB, _ := state.ProposalStatus.Get(txn, p.ProposalId)
	if string(statusB) != "pending" {
		return errors.New("invalid proposal status")
	}
	callerB, _ := state.ProposalCaller.Get(txn, p.ProposalId)
	if !bytesEqual(callerB, caller) {
		return errors.New("invalid proposal caller")
	}
	_ = state.ProposalStatus.Set(txn, p.ProposalId, []byte("canceled"))
	return nil
}
