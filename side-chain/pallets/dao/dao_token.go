package dao

import (
	"errors"
	"math/big"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func daoTransfer(caller []byte, p *DaoCallPayload, txn *model.Txn) error {
	st, _ := loadDaoState(txn)
	if st != nil && !st.Transfer {
		return errors.New("transfer disable")
	}
	state := newDaoStateState(txn)
	cb, _ := state.MemberBalance.Get(txn, caller)
	tb, _ := state.MemberBalance.Get(txn, p.To)
	if model.BytesToU128(cb).Sign() == 0 || model.BytesToU128(tb).Sign() == 0 {
		return errors.New("member not existed")
	}
	lockB, _ := state.MemberLock.Get(txn, caller)
	free := new(big.Int).Sub(model.BytesToU128(cb), model.BytesToU128(lockB))
	if free.Cmp(p.Value.ToBigInt()) < 0 {
		return errors.New("low balance")
	}
	_ = state.MemberBalance.Set(txn, caller, model.U128ToBytes(new(big.Int).Sub(model.BytesToU128(cb), p.Value.ToBigInt())))
	_ = state.MemberBalance.Set(txn, p.To, model.U128ToBytes(new(big.Int).Add(model.BytesToU128(tb), p.Value.ToBigInt())))
	return nil
}

func daoApprove(caller []byte, p *DaoCallPayload, txn *model.Txn) error {
	state := newDaoStateState(txn)
	_ = state.Allowance.Set(txn, addrKey(caller)+"_"+addrKey(p.Spender), model.U128ToBytes(p.Value.ToBigInt()))
	return nil
}

func daoTransferFrom(caller []byte, p *DaoCallPayload, txn *model.Txn) error {
	st, _ := loadDaoState(txn)
	if st != nil && !st.Transfer {
		return errors.New("transfer disable")
	}
	state := newDaoStateState(txn)
	fb, _ := state.MemberBalance.Get(txn, p.From)
	tb, _ := state.MemberBalance.Get(txn, p.To)
	if model.BytesToU128(fb).Sign() == 0 || model.BytesToU128(tb).Sign() == 0 {
		return errors.New("member not existed")
	}
	allowanceKey := addrKey(p.From) + "_" + addrKey(caller)
	allowB, _ := state.Allowance.Get(txn, allowanceKey)
	allow := model.BytesToU128(allowB)
	if allow.Cmp(p.Value.ToBigInt()) < 0 {
		return errors.New("insufficient allowance")
	}
	lockB, _ := state.MemberLock.Get(txn, p.From)
	free := new(big.Int).Sub(model.BytesToU128(fb), model.BytesToU128(lockB))
	if free.Cmp(p.Value.ToBigInt()) < 0 {
		return errors.New("low balance")
	}
	_ = state.MemberBalance.Set(txn, p.From, model.U128ToBytes(new(big.Int).Sub(model.BytesToU128(fb), p.Value.ToBigInt())))
	_ = state.MemberBalance.Set(txn, p.To, model.U128ToBytes(new(big.Int).Add(model.BytesToU128(tb), p.Value.ToBigInt())))
	_ = state.Allowance.Set(txn, allowanceKey, model.U128ToBytes(new(big.Int).Sub(allow, p.Value.ToBigInt())))
	return nil
}
