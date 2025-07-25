package sidechain

import (
	"errors"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/wetee-dao/tee-dsecret/pkg/chains"
	"github.com/wetee-dao/tee-dsecret/pkg/dkg"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/util"
)

// Submit sync tx to polkadot hub
func (s *SideChain) SyncToHub(txIndex int64, sigs [][]byte) error {
	call, err := model.GetCodec[types.Call](GLOABL_STATE, "tx_index"+fmt.Sprint(txIndex))
	if err != nil || call == nil {
		util.LogWithRed("Sync to polkadot hub", "error: call not found call data")
		return errors.New("sync to polkadot hub error: call not found call data")
	}

	// Aggregate signature
	util.LogWithGray("Sync to polkadot hub", "sync sigs = ", len(sigs))
	signer := dkg.NewDssSigner(s.dkg)
	signer.SetSigs(sigs)

	// submit sync tx to polkadot hub
	client := chains.MainChain.GetClient()
	err = client.SignAndSubmit(signer, *call, false, 0)
	if err != nil {
		util.LogWithRed("Sync to polkadot hub", "error => ", err.Error())
		fmt.Println("                    ", " SS58 => ", s.dkg.DkgPubKey.SS58())
		fmt.Println("                    ", " SYNC at batch tx ", txIndex)
	} else {
		util.LogWithGreen("Sync to polkadot hub", "success at batch tx ", fmt.Sprint(txIndex))
	}

	// stop sync transaction
	SubmitTx(&model.Tx{
		Payload: &model.Tx_SyncTxEnd{
			SyncTxEnd: txIndex,
		},
	})

	return err
}

// sync transaction index
var SyncTxIndexKey = "sync_transaction"

type AsyncBatchState struct {
	Going int64
	Done  int64
}

// check sync is running
func IsSyncRuning() bool {
	tx, err := model.GetJson[AsyncBatchState](GLOABL_STATE, SyncTxIndexKey)
	if err != nil {
		return true
	}
	if tx == nil {
		tx = &AsyncBatchState{
			Going: 0,
			Done:  0,
		}
	}

	return false
	// return tx.Going > tx.Done
}

// sync transaction step1
func SyncStep1() ([]byte, error) {
	tx, err := model.GetJson[AsyncBatchState](GLOABL_STATE, SyncTxIndexKey)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		tx = &AsyncBatchState{
			Going: 0,
			Done:  0,
		}
	}

	// if tx.Going > tx.Done {
	// 	return nil, errors.New("SyncStep1 one transaction is runing")
	// }

	return GetTxBytes(&model.Tx{
		Payload: &model.Tx_SyncTxStart{
			SyncTxStart: tx.Going + 1,
		},
	}), nil
}

// sync transaction step2
func SyncStep2(i int64) error {
	tx, err := model.GetJson[AsyncBatchState](GLOABL_STATE, SyncTxIndexKey)
	if err != nil {
		return err
	}

	if tx.Going > tx.Done {
		return errors.New("SyncStep2 one transaction is runing")
	}

	tx.Going = i
	tx.Done = i - 1
	return model.SetJson(GLOABL_STATE, SyncTxIndexKey, tx)
}

// sync transaction step3
func SyncEnd(i int64) error {
	tx, err := model.GetJson[AsyncBatchState](GLOABL_STATE, SyncTxIndexKey)
	if err != nil {
		return err
	}

	if i != tx.Going {
		// util.LogWithRed("SyncEnd", "i is not equal to tx.Going")
	}

	tx.Done = tx.Going
	return model.SetJson(GLOABL_STATE, SyncTxIndexKey, tx)
}
