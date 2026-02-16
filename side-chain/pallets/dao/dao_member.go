package dao

import (
	"errors"
	"math/big"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func daoInit(caller []byte, m *model.DaoInit, height int64, txn *model.Txn) error {
	st, err := loadDaoState(txn)
	if err != nil {
		return err
	}
	if len(st.Members) > 0 {
		return errors.New("dao already initialized")
	}
	state := newDaoStateState(txn)
	total := big.NewInt(0)
	members := make([][]byte, 0, len(m.GetInitialMembers()))
	for _, mem := range m.GetInitialMembers() {
		if len(mem.Account) == 0 {
			continue
		}
		members = append(members, mem.Account)
		bal := model.BytesToU128(mem.Balance)
		if err := state.MemberBalance.Set(txn, mem.Account, model.U128ToBytes(bal)); err != nil {
			return err
		}
		total.Add(total, bal)
	}
	st.Members = members
	st.TotalIssuance = total
	st.PublicJoin = m.GetPublicJoin()
	if len(m.GetSudoAccount()) > 0 {
		st.SudoAccount = m.SudoAccount
	}
	st.Transfer = false
	if m.GetDefaultTrack() != nil {
		trackID := uint32(0)
		if err := model.SetMappingJson(state.Tracks, txn, trackID, trackFromProto(m.GetDefaultTrack())); err != nil {
			return err
		}
		st.DefaultTrack = &trackID
	}
	st.NextProposalId = 0
	st.NextVoteId = 0
	st.NextSpendId = 0
	_ = state.SetMembers(st.Members)
	_ = state.SetTotalIssuance(st.TotalIssuance)
	_ = state.SetPublicJoin(st.PublicJoin)
	_ = state.SetSudoAccount(st.SudoAccount)
	_ = state.SetTransfer(st.Transfer)
	_ = state.SetDefaultTrack(st.DefaultTrack)
	_ = state.SetNextProposalId(0)
	_ = state.SetNextVoteId(0)
	_ = state.SetNextSpendId(0)
	return nil
}

func daoPublicJoin(caller []byte, txn *model.Txn) error {
	st, err := loadDaoState(txn)
	if err != nil {
		return err
	}
	if !st.PublicJoin {
		return errors.New("public join not allowed")
	}
	state := newDaoStateState(txn)
	b, _ := state.MemberBalance.Get(txn, caller)
	if model.BytesToU128(b).Sign() != 0 {
		return errors.New("member existed")
	}
	if err := state.MemberBalance.Set(txn, caller, model.U128ToBytes(big.NewInt(0))); err != nil {
		return err
	}
	st.Members = append(st.Members, caller)
	return state.SetMembers(st.Members)
}

func daoJoin(caller []byte, m *model.DaoJoin, txn *model.Txn) error {
	st, err := loadDaoState(txn)
	if err != nil {
		return err
	}
	if !isSudo(caller, st) {
		return errors.New("must call by gov/sudo")
	}
	newUser := m.GetNewUser()
	balance := model.BytesToU128(m.GetBalance())
	if len(newUser) == 0 || balance.Sign() == 0 {
		return errors.New("new_user and balance required")
	}
	state := newDaoStateState(txn)
	b, _ := state.MemberBalance.Get(txn, newUser)
	if model.BytesToU128(b).Sign() != 0 {
		return errors.New("member existed")
	}
	if err := state.MemberBalance.Set(txn, newUser, model.U128ToBytes(balance)); err != nil {
		return err
	}
	st.Members = append(st.Members, newUser)
	total := new(big.Int).Set(st.TotalIssuance)
	total.Add(total, balance)
	st.TotalIssuance = total
	_ = state.SetMembers(st.Members)
	return state.SetTotalIssuance(st.TotalIssuance)
}

func daoLeave(caller []byte, txn *model.Txn) error {
	st, _ := loadDaoState(txn)
	if st == nil {
		return errors.New("dao state not found")
	}
	if !memberInList(st.Members, caller) {
		return errors.New("member not existed")
	}
	state := newDaoStateState(txn)
	b, _ := state.MemberBalance.Get(txn, caller)
	if model.BytesToU128(b).Sign() != 0 {
		return errors.New("member balance not zero")
	}
	lockB, _ := state.MemberLock.Get(txn, caller)
	if model.BytesToU128(lockB).Sign() != 0 {
		return errors.New("member lock not zero")
	}
	st.Members = removeMember(st.Members, caller)
	_ = state.SetMembers(st.Members)
	_ = state.MemberBalance.Delete(txn, caller)
	_ = state.MemberLock.Delete(txn, caller)
	return nil
}

func daoLeaveWithBurn(caller []byte, txn *model.Txn) error {
	st, err := loadDaoState(txn)
	if err != nil || st == nil {
		return errors.New("dao state not found")
	}
	if !memberInList(st.Members, caller) {
		return errors.New("member not existed")
	}
	state := newDaoStateState(txn)
	b, _ := state.MemberBalance.Get(txn, caller)
	lockB, _ := state.MemberLock.Get(txn, caller)
	total := new(big.Int).Add(model.BytesToU128(b), model.BytesToU128(lockB))
	if st.TotalIssuance.Cmp(total) < 0 {
		return errors.New("low balance")
	}
	st.TotalIssuance = new(big.Int).Sub(st.TotalIssuance, total)
	st.Members = removeMember(st.Members, caller)
	_ = state.SetMembers(st.Members)
	if err := state.SetTotalIssuance(st.TotalIssuance); err != nil {
		return err
	}
	_ = state.MemberBalance.Delete(txn, caller)
	_ = state.MemberLock.Delete(txn, caller)
	return nil
}
