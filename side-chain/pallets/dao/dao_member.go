package dao

import (
	"errors"
	"math/big"
	"strconv"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func daoInit(caller []byte, m *model.DaoInit, height int64, txn *model.Txn) error {
	state := newDaoStateState(txn)
	if len(state.Members()) > 0 {
		return errors.New("dao already initialized")
	}
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
	_ = state.SetMembers(members)
	_ = state.SetTotalIssuance(total)
	_ = state.SetPublicJoin(m.GetPublicJoin())
	if len(m.GetSudoAccount()) > 0 {
		_ = state.SetSudoAccount(m.SudoAccount)
	}
	_ = state.SetTransfer(false)
	if m.GetDefaultTrack() != nil {
		trackID := uint32(0)
		if err := model.SetMappingJson(state.Tracks, txn, trackID, trackFromProto(m.GetDefaultTrack())); err != nil {
			return err
		}
		_ = state.SetDefaultTrack(&trackID)
	}
	_ = state.SetNextProposalId(0)
	_ = state.SetNextVoteId(0)
	_ = state.SetNextSpendId(0)
	return nil
}

func daoPublicJoin(caller []byte, txn *model.Txn) error {
	state := newDaoStateState(txn)
	if !state.PublicJoin() {
		return errors.New("public join not allowed")
	}
	b, _ := state.MemberBalance.Get(txn, caller)
	if model.BytesToU128(b).Sign() != 0 {
		return errors.New("member existed")
	}
	if err := state.MemberBalance.Set(txn, caller, model.U128ToBytes(big.NewInt(0))); err != nil {
		return err
	}
	members := append(state.Members(), caller)
	return state.SetMembers(members)
}

func daoJoin(caller []byte, m *model.DaoJoin, txn *model.Txn) error {
	state := newDaoStateState(txn)
	if !isSudo(caller, state) {
		return errors.New("must call by gov/sudo")
	}
	newUser := m.GetNewUser()
	balance := model.BytesToU128(m.GetBalance())
	if len(newUser) == 0 || balance.Sign() == 0 {
		return errors.New("new_user and balance required")
	}
	b, _ := state.MemberBalance.Get(txn, newUser)
	if model.BytesToU128(b).Sign() != 0 {
		return errors.New("member existed")
	}
	if err := state.MemberBalance.Set(txn, newUser, model.U128ToBytes(balance)); err != nil {
		return err
	}
	members := append(state.Members(), newUser)
	_ = state.SetMembers(members)
	total := state.TotalIssuance()
	if total == nil {
		total = big.NewInt(0)
	}
	total = new(big.Int).Add(total, balance)
	return state.SetTotalIssuance(total)
}

func daoLeave(caller []byte, txn *model.Txn) error {
	state := newDaoStateState(txn)
	members := state.Members()
	if !memberInList(members, caller) {
		return errors.New("member not existed")
	}
	b, _ := state.MemberBalance.Get(txn, caller)
	if model.BytesToU128(b).Sign() != 0 {
		return errors.New("member balance not zero")
	}
	lockB, _ := state.MemberLock.Get(txn, caller)
	if model.BytesToU128(lockB).Sign() != 0 {
		return errors.New("member lock not zero")
	}
	_ = state.SetMembers(removeMember(members, caller))
	_ = state.MemberBalance.Delete(txn, caller)
	_ = state.MemberLock.Delete(txn, caller)
	return nil
}

func daoLeaveWithBurn(caller []byte, txn *model.Txn) error {
	state := newDaoStateState(txn)
	members := state.Members()
	if !memberInList(members, caller) {
		return errors.New("member not existed")
	}
	b, _ := state.MemberBalance.Get(txn, caller)
	lockB, _ := state.MemberLock.Get(txn, caller)
	total := new(big.Int).Add(model.BytesToU128(b), model.BytesToU128(lockB))
	issuance := state.TotalIssuance()
	if issuance == nil || issuance.Cmp(total) < 0 {
		return errors.New("low balance")
	}
	_ = state.SetTotalIssuance(new(big.Int).Sub(issuance, total))
	_ = state.SetMembers(removeMember(members, caller))
	_ = state.MemberBalance.Delete(txn, caller)
	_ = state.MemberLock.Delete(txn, caller)
	return nil
}

// MemberList 返回所有成员地址列表（对应 ink list()）。
func MemberList() [][]byte {
	v, _ := model.GetJson[[][]byte]("dao", stateKeydaoStateMembers)
	if v == nil {
		return nil
	}
	return *v
}

// GetPublicJoin 返回是否允许公开加入（对应 ink get_public_join()）。
func GetPublicJoin() bool {
	b, _ := model.GetKey("dao", stateKeydaoStatePublicJoin)
	if len(b) == 0 {
		return false
	}
	v, _ := strconv.ParseBool(string(b))
	return v
}

func memberInList(members [][]byte, k []byte) bool {
	for _, m := range members {
		if bytesEqual(m, k) {
			return true
		}
	}
	return false
}

func removeMember(members [][]byte, k []byte) [][]byte {
	out := make([][]byte, 0, len(members))
	for _, m := range members {
		if !bytesEqual(m, k) {
			out = append(out, m)
		}
	}
	return out
}
