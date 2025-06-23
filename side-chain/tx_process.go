package sidechain

import (
	"bytes"
	"fmt"

	abcicli "github.com/cometbft/cometbft/abci/client"
	"github.com/cometbft/cometbft/abci/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/model/protoio"
)

func SubmitTx(tx *model.Tx) (*abcicli.ReqRes, error) {
	buf := new(bytes.Buffer)
	err := types.WriteMessage(tx, buf)
	if err != nil {
		return nil, err
	}

	return SideChainNode.Mempool().CheckTx(buf.Bytes(), SideChainNode.NodeInfo().ID())
}

func (app *SideChain) ProcessTx(txs [][]byte, txn *model.Txn) abci.ProcessProposalStatus {
	for _, txbt := range txs {
		tx := new(model.Tx)
		err := protoio.ReadMessage(bytes.NewBuffer(txbt), tx)
		if err != nil {
			return abci.PROCESS_PROPOSAL_STATUS_REJECT
		}

		switch tx.Payload.(type) {
		case *model.Tx_Epoch:
			break
		case *model.Tx_Bridge:
			break
		case *model.Tx_Test:
			break
		default:
			fmt.Println("Payload is not set")
		}
	}

	return abci.PROCESS_PROPOSAL_STATUS_ACCEPT
}
