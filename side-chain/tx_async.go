package sidechain

import (
	"errors"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/wetee-dao/tee-dsecret/pkg/chains"
	"github.com/wetee-dao/tee-dsecret/pkg/dkg"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/util"
)

type AsyncTransaction struct {
	Going int64
	Done  int64
}

func (s SideChain) SyncTxIsGoing() bool {
	tx, err := model.GetJson[AsyncTransaction]("G", "back_transaction")
	if err != nil {
		return true
	}
	if tx == nil {
		tx = &AsyncTransaction{
			Going: 0,
			Done:  0,
		}
	}
	return tx.Going > tx.Done
}

func (s *SideChain) TrySyncTxStart() ([]byte, error) {
	tx, err := model.GetJson[AsyncTransaction]("G", "back_transaction")
	if err != nil {
		return nil, err
	}
	if tx == nil {
		tx = &AsyncTransaction{
			Going: 0,
			Done:  0,
		}
	}

	if tx.Going > tx.Done {
		return nil, errors.New("one transaction is runing")
	}

	return GetTxBytes(&model.Tx{
		Payload: &model.Tx_SyncTxStart{
			SyncTxStart: tx.Going + 1,
		},
	}), nil
}

func (s *SideChain) SyncTxStart(i int64) error {
	tx, err := model.GetJson[AsyncTransaction]("G", "back_transaction")
	if err != nil {
		return err
	}

	if tx.Going > tx.Done {
		return errors.New("one transaction is runing")
	}

	tx.Going = i
	return nil
}

func (s *SideChain) SyncTxEnd(i int64) error {
	tx, err := model.GetJson[AsyncTransaction]("G", "back_transaction")
	if err != nil {
		return err
	}

	if tx.Going == tx.Done || i != tx.Going {
		return errors.New("no transaction is runing")
	}

	tx.Done = tx.Going
	return nil
}

func (s *SideChain) SyncToHub(txIndex int64, sigs [][]byte) error {
	bt, err := model.GetKey("G", "tx_index"+fmt.Sprint(txIndex))
	if err != nil {
		return err
	}
	call := new(types.Call)
	codec.Decode(bt, call)

	util.LogWithGray("Sync to PolkadotHub", "call with sigs => ", len(sigs))
	signer := dkg.NewDssSigner(s.dkg)
	signer.SetSigs(sigs)

	client := chains.MainChain.GetClient()
	err = client.SignAndSubmit(signer, *call, false, 0)
	if err != nil {
		util.LogWithRed("Sync to PolkadotHub error", err.Error())
	} else {
		util.LogWithGreen("Sync to PolkadotHub success tx => ", fmt.Sprint(txIndex))
	}

	return err
}
