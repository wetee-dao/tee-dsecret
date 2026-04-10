package gov

import (
	"fmt"
	"math/big"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func (d GovQuery) Members() ([]Member, error) {
	keys, balances, err := d.members.List(d.api.GetTxn())
	if err != nil {
		return nil, err
	}

	out := make([]Member, len(keys))
	for i := range keys {
		out[i] = Member{Account: keys[i], Balance: balances[i]}
	}
	return out, err
}

func (d GovQuery) GetPublicJoin() (bool, error) {
	return d.publicJoin.GetOrDefault(d.api.GetTxn(), true)
}

func (d GovMutation) PublicJoin() error {
	fmt.Println("PublicJoin1")
	enabled, err := d.publicJoin.GetOrDefault(d.api.GetTxn(), true)
	if err != nil {
		return err
	}

	if !enabled {
		return ErrPublicJoinNotAllowed
	}

	fmt.Println("PublicJoin")
	return d.addMember(d.api.GetCaller(), model.ZeroAmount)
}

func (d GovMutation) SetPublicJoin(publicJoin bool) error {
	if err := d.ensureGov(); err != nil {
		return err
	}
	return d.publicJoin.Set(d.api.GetTxn(), publicJoin)
}

func (d GovMutation) Join(newUser model.UniAddr, balance model.Amount) error {
	if err := d.ensureGov(); err != nil {
		return err
	}

	return d.addMember(newUser, balance)
}

func (d GovMutation) Leave() error {
	caller := d.api.GetCaller()
	member, err := d.members.GetOrDefault(d.api.GetTxn(), caller, model.ZeroAmount)
	if err != nil {
		return err
	}

	lock, err := d.memberLocks.GetOrDefault(d.api.GetTxn(), caller, model.ZeroAmount)
	if err != nil {
		return err
	}

	if member.Cmp(big.NewInt(0)) >= 0 || lock.Cmp(big.NewInt(0)) >= 0 {
		return ErrMemberBalanceNotZero
	}
	return d.removeMember(caller)
}

func (d GovMutation) LeaveWithBurn() error {
	caller := d.api.GetCaller()
	return d.deleteMemberAndBurn(caller, false)
}

func (d GovMutation) Delete(account model.UniAddr) error {
	return d.deleteMemberAndBurn(account, true)
}

func (d Gov) addMember(account model.UniAddr, balance model.Amount) error {
	if len(account.V) == 0 {
		return ErrMemberNotExisted
	}
	member, err := d.members.Get(d.api.GetTxn(), account)
	if err != nil {
		return err
	}
	if member != nil {
		return ErrMemberExisted
	}

	if err := d.members.Set(d.api.GetTxn(), account, balance); err != nil {
		return err
	}

	total, err := d.totalIssuance.GetOrDefault(d.api.GetTxn(), model.ZeroAmount)
	if err != nil {
		return err
	}
	return d.totalIssuance.Set(d.api.GetTxn(), model.AmountAdd(total, balance))
}

func (d Gov) deleteMemberAndBurn(account model.UniAddr, requireGov bool) error {
	if requireGov {
		if err := d.ensureGov(); err != nil {
			return err
		}
	}

	member, err := d.members.GetOrDefault(d.api.GetTxn(), account, model.ZeroAmount)
	if err != nil {
		return err
	}

	lock, err := d.memberLocks.GetOrDefault(d.api.GetTxn(), account, model.ZeroAmount)
	if err != nil {
		return err
	}

	amount := model.AmountAdd(member, lock)
	total, err := d.totalIssuance.GetOrDefault(d.api.GetTxn(), model.ZeroAmount)
	if err != nil {
		return err
	}

	if err := d.totalIssuance.Set(d.api.GetTxn(), model.AmountSub(total, amount)); err != nil {
		return err
	}
	return d.removeMember(account)
}

func (d Gov) removeMember(account model.UniAddr) error {
	if err := d.members.Delete(d.api.GetTxn(), account); err != nil {
		return err
	}
	return d.memberLocks.Delete(d.api.GetTxn(), account)
}
