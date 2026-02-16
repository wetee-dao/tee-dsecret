package dao

import (
	"errors"
	"math/big"
	"strconv"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func daoSubmitProposal(caller []byte, m *model.DaoSubmitProposal, height int64, txn *model.Txn) error {
	state := newDaoStateState(txn)
	if len(state.Members()) == 0 {
		return errors.New("dao not initialized")
	}
	b, _ := state.MemberBalance.Get(txn, caller)
	if model.BytesToU128(b).Sign() == 0 {
		return errors.New("member not existed")
	}
	if m.GetCall() == nil {
		return errors.New("call required")
	}
	trackId := m.GetTrackId()
	if int(trackId) >= len(getDaoTracks(txn)) {
		return errors.New("no track or invalid track")
	}
	track, _ := model.GetMappingJson[uint32, DaoTrackData](state.Tracks, txn, trackId)
	callAmount := model.BytesToU128(m.Call.Amount)
	if track != nil && callAmount.Sign() != 0 {
		if callAmount.Cmp(track.MaxBalance.ToBigInt()) > 0 {
			return errors.New("max balance overflow")
		}
	}
	proposalId := state.NextProposalId()
	_ = state.SetNextProposalId(proposalId + 1)
	_ = model.SetMappingJson(state.Proposals, txn, proposalId, callContentFromProto(m.Call))
	_ = state.ProposalStatus.Set(txn, proposalId, []byte("pending"))
	_ = state.ProposalTrack.Set(txn, proposalId, []byte(strconv.FormatUint(uint64(trackId), 10)))
	_ = state.ProposalCaller.Set(txn, proposalId, caller)
	_ = state.SubmitBlock.Set(txn, proposalId, []byte(strconv.FormatInt(height, 10)))
	return nil
}

func daoDepositProposal(caller []byte, m *model.DaoDepositProposal, height int64, txn *model.Txn) error {
	state := newDaoStateState(txn)
	proposalId := m.GetProposalId()
	b, _ := state.ProposalStatus.Get(txn, proposalId)
	if string(b) != "pending" {
		return errors.New("invalid proposal status")
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
	if track == nil {
		return errors.New("no track")
	}
	if height < submitBlock+int64(track.PreparePeriod) {
		return errors.New("invalid deposit time")
	}
	depositNeed := track.DecisionDeposit.ToBigInt()
	amount := model.BytesToU128(m.GetAmount())
	if amount.Cmp(depositNeed) < 0 {
		return errors.New("invalid deposit")
	}
	_ = model.SetMappingJson(state.Deposits, txn, proposalId, &proposalDeposit{Depositor: caller, Amount: model.NewU128(amount), Block: height})
	_ = state.ProposalStatus.Set(txn, proposalId, []byte("ongoing"))
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

func daoExecProposal(caller []byte, m *model.DaoExecProposal, height int64, txn *model.Txn) error {
	state := newDaoStateState(txn)
	proposalId := m.GetProposalId()
	stB, _ := state.ProposalStatus.Get(txn, proposalId)
	if string(stB) != "ongoing" {
		return errors.New("proposal not ongoing")
	}
	endBlock := daoProposalEndBlock(txn, proposalId)
	trB, _ := state.ProposalTrack.Get(txn, proposalId)
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
	call, _ := model.GetMappingJson[uint32, DaoCallContent](state.Proposals, txn, proposalId)
	if call == nil {
		return errors.New("invalid proposal")
	}
	_ = state.ProposalStatus.Set(txn, proposalId, []byte("approved"))
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

func daoCancelProposal(caller []byte, m *model.DaoCancelProposal, height int64, txn *model.Txn) error {
	state := newDaoStateState(txn)
	proposalId := m.GetProposalId()
	statusB, _ := state.ProposalStatus.Get(txn, proposalId)
	if string(statusB) != "pending" {
		return errors.New("invalid proposal status")
	}
	callerB, _ := state.ProposalCaller.Get(txn, proposalId)
	if !bytesEqual(callerB, caller) {
		return errors.New("invalid proposal caller")
	}
	_ = state.ProposalStatus.Set(txn, proposalId, []byte("canceled"))
	return nil
}

// Proposals 分页返回提案内容列表（对应 ink proposals(page, size)）。
func Proposals(page, size uint32) []*DaoCallContent {
	return proposalsFromDB(page, size)
}

// Proposal 按 ID 返回提案内容（对应 ink proposal(id)）。
func Proposal(id uint32) *DaoCallContent {
	key := stateKeydaoStateProposals + strconv.FormatUint(uint64(id), 10)
	c, _ := model.GetJson[DaoCallContent]("dao", key)
	return c
}

// ProposalStatus 返回提案状态（对应 ink proposal_status(proposal_id)）。
func ProposalStatus(proposalId uint32) string {
	key := stateKeydaoStateProposalStatus + strconv.FormatUint(uint64(proposalId), 10)
	b, _ := model.GetKey("dao", key)
	return string(b)
}

func proposalsFromDB(page, size uint32) []*DaoCallContent {
	if page < 1 {
		page = 1
	}
	start := (page - 1) * size
	var out []*DaoCallContent
	for i := uint32(0); i < size; i++ {
		id := start + i
		c := Proposal(id)
		if c == nil {
			break
		}
		out = append(out, c)
	}
	return out
}
