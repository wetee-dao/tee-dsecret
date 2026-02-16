package dao

import (
	"errors"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func daoSetPublicJoin(caller []byte, p *DaoCallPayload, txn *model.Txn) error {
	st, err := loadDaoState(txn)
	if err != nil || st == nil {
		return errors.New("dao not initialized")
	}
	if !isSudo(caller, st) {
		return errors.New("must call by gov/sudo")
	}
	if p.PublicJoin != nil {
		st.PublicJoin = *p.PublicJoin
		state := newDaoStateState(txn)
		return state.SetPublicJoin(st.PublicJoin)
	}
	return errors.New("public_join required")
}
