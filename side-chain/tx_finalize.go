package sidechain

import (
	"bytes"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"wetee.app/dsecret/internal/model"
	"wetee.app/dsecret/internal/model/protoio"
	"wetee.app/dsecret/internal/util"
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
		case *model.Tx_Epoch:
			// calc validator updates
			app.calcValidatorUpdates(p.Epoch)

			// set epoch and validators
			err = app.SetEpoch(p.Epoch, txn)
			if err != nil {
				return nil, err
			}

			util.LogWithPurple("SideChain FinalizeTx", "New Epoch", p.Epoch.Epoch)
			res = append(res, &abci.ExecTxResult{Code: uint32(abci.CodeTypeOK)})
		case *model.Tx_Bridge:
			res = append(res, &abci.ExecTxResult{Code: uint32(abci.CodeTypeOK)})
		case *model.Tx_Test:
			res = append(res, &abci.ExecTxResult{Code: uint32(abci.CodeTypeOK)})
		default:
			fmt.Println("Payload is unknown")
		}
	}

	return res, nil
}
