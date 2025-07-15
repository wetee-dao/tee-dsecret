package sidechain

import (
	"bytes"
	"errors"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/model/protoio"
)

func (app *SideChain) FinalizeTx(txs [][]byte, txn *model.Txn, height int64, proposer []byte) ([]*abci.ExecTxResult, error) {
	res := []*abci.ExecTxResult{}
	hubCalls := make([]*model.HubCall, 0, len(txs))
	var txIndex int64 = 0

	for _, txbt := range txs {
		txbox := new(model.TxBox)
		err := protoio.ReadMessage(bytes.NewBuffer(txbt), txbox)
		if err != nil {
			return nil, err
		}

		tx := new(model.Tx)
		err = protoio.ReadMessage(bytes.NewBuffer(txbox.Tx), tx)
		if err != nil {
			return nil, err
		}

		switch p := tx.Payload.(type) {
		case *model.Tx_EpochStart:
			app.SetEpochStatus(p.EpochStart) // set epoch last time
		case *model.Tx_EpochEnd:
			app.calcValidatorUpdates(p.EpochEnd) // calc validator updates
			err = app.SetEpoch(p.EpochEnd, txn)  // set epoch and validators
			if err != nil {
				return nil, err
			}
			app.SetEpochStatus(0)
		case *model.Tx_Bridge:
			break
		case *model.Tx_SyncTxStart:
			txIndex = p.SyncTxStart
			app.SyncTxStart(p.SyncTxStart)
		case *model.Tx_SyncTxEnd:
			app.SyncTxEnd(p.SyncTxEnd)
		case *model.Tx_HubCall:
			hubCalls = append(hubCalls, p.HubCall)
		case *model.Tx_Test:
			break
		default:
			return nil, errors.New("invalid tx type")
		}

		res = append(res, &abci.ExecTxResult{Code: uint32(abci.CodeTypeOK)})
	}

	if txIndex > 0 && len(hubCalls) > 0 && app.dkg != nil {
		app.sendPartialSign(txIndex, hubCalls, app.ProposerAddressToNodeKey(proposer))
	}
	return res, nil
}
