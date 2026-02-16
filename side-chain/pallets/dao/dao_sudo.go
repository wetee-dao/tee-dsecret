package dao

import (
	"errors"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func daoSetPublicJoin(caller []byte, m *model.DaoSetPublicJoin, txn *model.Txn) error {
	st, err := loadDaoState(txn)
	if err != nil || st == nil {
		return errors.New("dao not initialized")
	}
	if !isSudo(caller, st) {
		return errors.New("must call by gov/sudo")
	}
	st.PublicJoin = m.GetPublicJoin()
	state := newDaoStateState(txn)
	return state.SetPublicJoin(st.PublicJoin)
}
