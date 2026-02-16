package dao

import (
	"errors"
	"math/big"
	"strconv"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func daoSpend(caller []byte, m *model.DaoSpend, txn *model.Txn) error {
	state := newDaoStateState(txn)
	if len(state.Members()) == 0 {
		return errors.New("dao not initialized")
	}
	b, _ := state.MemberBalance.Get(txn, caller)
	if model.BytesToU128(b).Sign() == 0 {
		return errors.New("member not existed")
	}
	spendId := state.NextSpendId()
	proposalId := state.NextProposalId()
	_ = state.SetNextSpendId(spendId + 1)
	_ = state.SetNextProposalId(proposalId + 1)
	amount := model.BytesToU128(m.GetAmount())
	_ = model.SetMappingJson(state.Spends, txn, spendId, &spendRecord{Caller: caller, To: m.GetTo(), Amount: model.NewU128(amount), Payout: false})
	input := make([]byte, 8)
	for i := 0; i < 8; i++ {
		input[i] = byte(spendId >> (i * 8))
	}
	_ = model.SetMappingJson(state.Proposals, txn, proposalId, &DaoCallContent{Amount: model.NewU128(amount), Input: input})
	_ = state.ProposalStatus.Set(txn, proposalId, []byte("pending"))
	_ = state.ProposalTrack.Set(txn, proposalId, []byte(strconv.FormatUint(uint64(m.GetTrackId()), 10)))
	_ = state.ProposalCaller.Set(txn, proposalId, caller)
	_ = state.SubmitBlock.Set(txn, proposalId, []byte(strconv.FormatInt(0, 10)))
	return nil
}

func daoPayout(caller []byte, m *model.DaoPayout, txn *model.Txn) error {
	state := newDaoStateState(txn)
	if len(state.Members()) == 0 {
		return errors.New("dao not initialized")
	}
	if !isSudo(caller, state) {
		return errors.New("must call by gov/sudo")
	}
	spendId := m.GetSpendId()
	s, _ := model.GetMappingJson[uint64, spendRecord](state.Spends, txn, spendId)
	if s == nil {
		return errors.New("spend not found")
	}
	if s.Payout {
		return errors.New("spend already executed")
	}
	toB, _ := state.MemberBalance.Get(txn, s.To)
	if model.BytesToU128(toB).Sign() != 0 {
		cur := model.BytesToU128(toB)
		_ = state.MemberBalance.Set(txn, s.To, model.U128ToBytes(new(big.Int).Add(cur, s.Amount.ToBigInt())))
	}
	s.Payout = true
	_ = model.SetMappingJson(state.Spends, txn, spendId, s)
	return nil
}
