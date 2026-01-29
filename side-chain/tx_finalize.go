package sidechain

import (
	"bytes"

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
			LogWithTime("Empty TX:", p.Empty)
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
			err = HubSyncStep2(p.SyncTxStart, txn)
			if err != nil {
				return nil, err
			}
		case *model.Tx_SyncTxEnd: // end hub sync tx
			err = HubSyncEnd(p.SyncTxEnd, txn)
			if err != nil {
				return nil, err
			}
			// 所有节点在处理 SyncTxEnd 时统一清理 tx_index_ 储存
			deleteTxIndexStore(p.SyncTxEnd)
		case *model.Tx_SyncTxRetry: // retry hub sync tx
			if app.dkg == nil {
				// dkg 未初始化，无法处理重试
				LogWithTime("SyncTxRetry", "dkg is nil, skipping retry for txIndex:", p.SyncTxRetry)
				break
			}

			// 重新获取该交易的所有部分签名
			sigs, err := app.SigListOfTx(p.SyncTxRetry)
			if err != nil {
				return nil, errors.Wrap(err, "SyncTxRetry: failed to get signatures")
			}

			// 提取签名
			shares := make([][]byte, 0, len(sigs))
			for _, sig := range sigs {
				shares = append(shares, sig.HubSig)
			}

			// 检查是否有足够的签名
			if len(shares) < app.dkg.Threshold+1 {
				// 签名不足，等待更多签名，不处理重试
				LogWithTime("SyncTxRetry", "insufficient signatures", "txIndex:", p.SyncTxRetry, "got:", len(shares), "need:", app.dkg.Threshold+1)
			} else {
				// 重新调用 SyncToHub 提交到主链
				err = app.SyncToHub(p.SyncTxRetry, shares)
				if err != nil {
					// 再次失败时 SyncToHub 会再次 SubmitTx(SyncTxRetry)，由下一块重试
					LogWithTime("SyncTxRetry", "retry failed", "txIndex:", p.SyncTxRetry, "error:", err.Error())
				}
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
		err := app.sendPartialSign(hubCalls[0].ChainId, txIndex, hubCalls, app.ProposerAddressToNodeKey(proposer))
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
