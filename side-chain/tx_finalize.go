package sidechain

import (
	"bytes"
	"fmt"

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
			// set epoch last time
			app.SetEpochStatus(p.EpochStatus)
			res = append(res, &abci.ExecTxResult{Code: uint32(abci.CodeTypeOK)})
		case *model.Tx_Epoch:
			// calc validator updates
			app.calcValidatorUpdates(p.Epoch)
			// set epoch and validators
			err = app.SetEpoch(p.Epoch, txn)
			if err != nil {
				return nil, err
			}
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
