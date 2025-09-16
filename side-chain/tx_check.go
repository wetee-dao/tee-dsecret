package sidechain

import (
	"bytes"
	"fmt"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/model/protoio"
)

func (app *SideChain) checkTx(txbt []byte) uint32 {
	tx := new(model.TxBox)
	err := protoio.ReadMessage(bytes.NewBuffer(txbt), tx)
	if err != nil {
		return CodeTypeEncodingError
	}

	if len(tx.Org) == 0 {
		fmt.Println("invalid node1")
		return CodeInvalidNode
	}

	keys := app.p2p.AllNodes()
	isIn := false
	for _, key := range keys {
		if bytes.Equal(tx.Org, key.Byte()) {
			isIn = true
			break
		}
	}

	if !isIn {
		return CodeInvalidNode
	}

	return CodeTypeOK
}
