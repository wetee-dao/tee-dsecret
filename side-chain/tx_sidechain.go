package sidechain

import (
	"bytes"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/mempool"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

// Submit tx to sidechain
func SubmitTx(tx *model.Tx) error {
	return SideChainNode.Mempool().CheckTx(GetTxBytes(tx), func(r *abci.ResponseCheckTx) {}, mempool.TxInfo{
		SenderP2PID: SideChainNode.NodeInfo().ID(),
	})
}

// Get tx bytes
func GetTxBytes(tx *model.Tx) []byte {
	buf := new(bytes.Buffer)
	abci.WriteMessage(tx, buf)

	org := P2PKey.Byte()
	txbox := model.TxBox{
		Org: org,
		Tx:  buf.Bytes(),
	}

	boxbuf := new(bytes.Buffer)
	abci.WriteMessage(&txbox, boxbuf)

	return boxbuf.Bytes()
}
