package dao

import (
	"errors"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func daoSetPublicJoin(caller []byte, m *model.DaoSetPublicJoin, txn *model.Txn) error {
	state := newDaoStateState(txn)
	if len(state.Members()) == 0 {
		return errors.New("dao not initialized")
	}
	if !isSudo(caller, state) {
		return errors.New("must call by gov/sudo")
	}
	return state.SetPublicJoin(m.GetPublicJoin())
}
