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
		case *model.Tx_Epoch:
			app.calcValidatorUpdates(p.Epoch) // calc validator updates
			err = app.SetEpoch(p.Epoch, txn)  // set epoch and validators
			if err != nil {
				return nil, err
			}
		case *model.Tx_Bridge:
			break
		case *model.Tx_Test:
			break
		default:
			return nil, errors.New("invalid tx type")
		}
		res = append(res, &abci.ExecTxResult{Code: uint32(abci.CodeTypeOK)})
	}

	return res, nil
}
