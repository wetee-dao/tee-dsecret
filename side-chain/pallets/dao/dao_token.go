package dao

import (
	"errors"
	"math/big"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func daoTransfer(caller []byte, m *model.DaoTransfer, txn *model.Txn) error {
	st, _ := loadDaoState(txn)
	if st != nil && !st.Transfer {
		return errors.New("transfer disable")
	}
	state := newDaoStateState(txn)
	to := m.GetTo()
	value := model.BytesToU128(m.GetValue())
	cb, _ := state.MemberBalance.Get(txn, caller)
	tb, _ := state.MemberBalance.Get(txn, to)
	if model.BytesToU128(cb).Sign() == 0 || model.BytesToU128(tb).Sign() == 0 {
		return errors.New("member not existed")
	}
	lockB, _ := state.MemberLock.Get(txn, caller)
	free := new(big.Int).Sub(model.BytesToU128(cb), model.BytesToU128(lockB))
	if free.Cmp(value) < 0 {
		return errors.New("low balance")
	}
	_ = state.MemberBalance.Set(txn, caller, model.U128ToBytes(new(big.Int).Sub(model.BytesToU128(cb), value)))
	_ = state.MemberBalance.Set(txn, to, model.U128ToBytes(new(big.Int).Add(model.BytesToU128(tb), value)))
	return nil
}

func daoApprove(caller []byte, m *model.DaoApprove, txn *model.Txn) error {
	state := newDaoStateState(txn)
	_ = state.Allowance.Set(txn, addrKey(caller)+"_"+addrKey(m.GetSpender()), model.U128ToBytes(model.BytesToU128(m.GetValue())))
	return nil
}

func daoTransferFrom(caller []byte, m *model.DaoTransferFrom, txn *model.Txn) error {
	st, _ := loadDaoState(txn)
	if st != nil && !st.Transfer {
		return errors.New("transfer disable")
	}
	state := newDaoStateState(txn)
	from := m.GetFrom()
	to := m.GetTo()
	value := model.BytesToU128(m.GetValue())
	fb, _ := state.MemberBalance.Get(txn, from)
	tb, _ := state.MemberBalance.Get(txn, to)
	if model.BytesToU128(fb).Sign() == 0 || model.BytesToU128(tb).Sign() == 0 {
		return errors.New("member not existed")
	}
	allowanceKey := addrKey(from) + "_" + addrKey(caller)
	allowB, _ := state.Allowance.Get(txn, allowanceKey)
	allow := model.BytesToU128(allowB)
	if allow.Cmp(value) < 0 {
		return errors.New("insufficient allowance")
	}
	lockB, _ := state.MemberLock.Get(txn, from)
	free := new(big.Int).Sub(model.BytesToU128(fb), model.BytesToU128(lockB))
	if free.Cmp(value) < 0 {
		return errors.New("low balance")
	}
	_ = state.MemberBalance.Set(txn, from, model.U128ToBytes(new(big.Int).Sub(model.BytesToU128(fb), value)))
	_ = state.MemberBalance.Set(txn, to, model.U128ToBytes(new(big.Int).Add(model.BytesToU128(tb), value)))
	_ = state.Allowance.Set(txn, allowanceKey, model.U128ToBytes(new(big.Int).Sub(allow, value)))
	return nil
}
