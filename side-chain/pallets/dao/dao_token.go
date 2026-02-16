package dao

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func daoTransfer(caller []byte, m *model.DaoTransfer, txn *model.Txn) error {
	state := newDaoStateState(txn)
	if !state.Transfer() {
		return errors.New("transfer disable")
	}
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
	state := newDaoStateState(txn)
	if !state.Transfer() {
		return errors.New("transfer disable")
	}
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

// Name 返回代币名称（对应 ink name()，当前实现返回空）。
func Name() []byte { return nil }

// Symbol 返回代币符号（对应 ink symbol()，当前实现返回空）。
func Symbol() []byte { return nil }

// Decimals 返回小数位数（对应 ink decimals()，12）。
func Decimals() uint8 { return 12 }

// TotalSupply 返回代币总供应量（对应 ink total_supply()）。
func TotalSupply() *big.Int {
	b, _ := model.GetKey("dao", stateKeydaoStateTotalIssuance)
	return model.BytesToU128(b)
}

// BalanceOf 返回账户余额（对应 ink balance_of(owner)）。
func BalanceOf(owner []byte) *big.Int {
	if len(owner) == 0 {
		return big.NewInt(0)
	}
	key := stateKeydaoStateMemberBalance + hex.EncodeToString(owner)
	b, _ := model.GetKey("dao", key)
	return model.BytesToU128(b)
}

// Allowance 返回 owner 授权给 spender 的额度（对应 ink allowance(owner, spender)）。
func Allowance(owner, spender []byte) *big.Int {
	if len(owner) == 0 || len(spender) == 0 {
		return big.NewInt(0)
	}
	key := stateKeydaoStateAllowance + hex.EncodeToString(owner) + "_" + hex.EncodeToString(spender)
	b, _ := model.GetKey("dao", key)
	return model.BytesToU128(b)
}

// LockBalanceOf 返回账户锁定余额（对应 ink lock_balance_of(owner)）。
func LockBalanceOf(owner []byte) *big.Int {
	if len(owner) == 0 {
		return big.NewInt(0)
	}
	key := stateKeydaoStateMemberLock + hex.EncodeToString(owner)
	b, _ := model.GetKey("dao", key)
	return model.BytesToU128(b)
}
