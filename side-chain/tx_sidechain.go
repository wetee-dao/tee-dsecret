package sidechain

import (
	"bytes"

	abcicli "github.com/cometbft/cometbft/abci/client"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

// Submit tx to sidechain
func SubmitTx(tx *model.Tx) (*abcicli.ReqRes, error) {
	return SideChainNode.Mempool().CheckTx(GetTxBytes(tx), SideChainNode.NodeInfo().ID())
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
