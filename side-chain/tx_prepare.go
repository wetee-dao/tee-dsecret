package sidechain

import (
	"bytes"
	"time"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/model/protoio"
	"github.com/wetee-dao/tee-dsecret/pkg/util"
)

func (s *SideChain) PrepareTx(txs [][]byte, finaltx *[][]byte, height int64, addMainChainTx bool) {
	hubtx := make([][]byte, 0, 50)
	for _, txbt := range txs {
		txbox := new(model.TxBox)
		err := protoio.ReadMessage(bytes.NewBuffer(txbt), txbox)
		if err != nil {
			continue
		}

		tx := new(model.Tx)
		err = protoio.ReadMessage(bytes.NewBuffer(txbox.Tx), tx)
		if err != nil {
			continue
		}

		switch tx.Payload.(type) {
		case *model.Tx_EpochStart:
			*finaltx = append(*finaltx, txbt)
		case *model.Tx_EpochEnd:
			*finaltx = append(*finaltx, txbt)
		case *model.Tx_SyncTxStart:
			*finaltx = append(*finaltx, txbt)
		case *model.Tx_SyncTxEnd:
			*finaltx = append(*finaltx, txbt)
		case *model.Tx_HubCall:
			if addMainChainTx {
				hubtx = append(hubtx, txbt)
			}
		default:
			break
		}
	}

	if len(hubtx) > 0 {
		tx, err := SyncStep1()
		if err != nil {
			util.LogWithRed("PrepareTx", "TryTxStart err:", err)
			time.Sleep(time.Second * 2)
			return
		}

		if s.dkg.AvailableNodeLen() < s.dkg.Threshold+1 {
			util.LogWithRed("PrepareTx", "exapect validator node:", s.dkg.Threshold+1, ", got:", s.dkg.AvailableNodeLen())
			time.Sleep(time.Second * 2)
			return
		}

		*finaltx = append(*finaltx, tx)
		*finaltx = append(*finaltx, hubtx...)
		s.DeleteSigOfTx(height)
	}
}
