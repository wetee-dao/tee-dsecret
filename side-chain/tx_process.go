package sidechain

import (
	"bytes"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/model/protoio"
)

// Process tx
func (app *SideChain) ProcessTx(txs [][]byte) abci.ProcessProposalStatus {
	for _, txbt := range txs {
		txbox := new(model.TxBox)
		err := protoio.ReadMessage(bytes.NewBuffer(txbt), txbox)
		if err != nil {
			return abci.PROCESS_PROPOSAL_STATUS_REJECT
		}

		tx := new(model.Tx)
		err = protoio.ReadMessage(bytes.NewBuffer(txbox.Tx), tx)
		if err != nil {
			return abci.PROCESS_PROPOSAL_STATUS_REJECT
		}

		switch tx.Payload.(type) {
		case *model.Tx_EpochStart:
		case *model.Tx_EpochEnd:
		case *model.Tx_SyncTxStart:
		case *model.Tx_SyncTxEnd:
		case *model.Tx_SyncTxRetry:
		case *model.Tx_Empty:
		case *model.Tx_HubCall:
		case *model.Tx_DaoCall:
		default:
			fmt.Println("Payload is not set")
		}
	}

	return abci.PROCESS_PROPOSAL_STATUS_ACCEPT
}
