package sidechain

import (
	"bytes"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/pkg/errors"
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
		case *model.Tx_Empty:
			fmt.Println("Empty TX:", p.Empty)
		case *model.Tx_EpochStart: // set epoch last time
			err := app.SetEpochStatus(p.EpochStart)
			if err != nil {
				return nil, err
			}
		case *model.Tx_EpochEnd:
			app.calcValidatorUpdates(p.EpochEnd) // calc validator updates
			err = app.SetEpoch(p.EpochEnd, txn)  // set epoch and validators
			if err != nil {
				return nil, err
			}
			err = app.SetEpochStatus(0)
			if err != nil {
				return nil, err
			}
		case *model.Tx_SyncTxStart: // start hub sync tx
			txIndex = p.SyncTxStart
			err = SyncStep2(p.SyncTxStart, txn)
			if err != nil {
				return nil, err
			}
		case *model.Tx_SyncTxEnd: // end hub sync tx
			err = SyncEnd(p.SyncTxEnd, txn)
			if err != nil {
				return nil, err
			}
		case *model.Tx_HubCall: // add hub call
			err := app.finalizeHubCall(p.HubCall, txn)
			if err != nil {
				return nil, err
			}
			hubCalls = append(hubCalls, p.HubCall)
		default:
			return nil, errors.New("invalid tx type")
		}

		res = append(res, &abci.ExecTxResult{Code: uint32(abci.CodeTypeOK)})
	}

	// if hub tx, send partial sign
	if txIndex > 0 && len(hubCalls) > 0 && app.dkg != nil {
		err := app.sendPartialSign(txIndex, hubCalls, app.ProposerAddressToNodeKey(proposer))
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (app *SideChain) finalizeHubCall(hub *model.HubCall, txn *model.Txn) error {
	for _, callWrap := range hub.Call {
		switch tx := callWrap.Tx.(type) {
		case *model.TeeCall_PodStart:
			continue
		case *model.TeeCall_PodMint:
			continue
		case *model.TeeCall_UploadSecret:
			upload := tx.UploadSecret
			user := types.H160(upload.User)
			err := app.SaveSecret(user, upload.Index, upload.Data, txn)
			if err != nil {
				return errors.Wrap(err, "finalizeHubCall SaveSecret")
			}
		case *model.TeeCall_InitDisk:
			initDisk := tx.InitDisk
			user := types.H160(initDisk.User)
			err := app.SaveDiskKey(user, initDisk.Index, initDisk.Data, txn)
			if err != nil {
				return errors.Wrap(err, "finalizeHubCall InitDisk")
			}
		default:
			return errors.New("finalizeHubCall invalid tx type")
		}
	}
	return nil
}
