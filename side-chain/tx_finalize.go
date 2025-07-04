package sidechain

import (
	"bytes"
	"errors"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/model/protoio"
)

func (app *SideChain) FinalizeTx(txs [][]byte, txn *model.Txn) ([]*abci.ExecTxResult, error) {
	res := []*abci.ExecTxResult{}
	for _, txbt := range txs {
		tx := new(model.Tx)
		err := protoio.ReadMessage(bytes.NewBuffer(txbt), tx)
		if err != nil {
			return nil, err
		}

		switch p := tx.Payload.(type) {
		case *model.Tx_EpochStatus:
			app.SetEpochStatus(p.EpochStatus) // set epoch last time
			res = append(res, &abci.ExecTxResult{Code: uint32(abci.CodeTypeOK)})
		case *model.Tx_Epoch:
			app.calcValidatorUpdates(p.Epoch) // calc validator updates
			err = app.SetEpoch(p.Epoch, txn)  // set epoch and validators
			if err != nil {
				return nil, err
			}
			res = append(res, &abci.ExecTxResult{Code: uint32(abci.CodeTypeOK)})
		case *model.Tx_Bridge:
			res = append(res, &abci.ExecTxResult{Code: uint32(abci.CodeTypeOK)})
		case *model.Tx_Test:
			res = append(res, &abci.ExecTxResult{Code: uint32(abci.CodeTypeOK)})
		default:
			return nil, errors.New("invalid tx type")
		}
	}

	return res, nil
}
