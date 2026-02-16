package sidechain

import (
	"bytes"
	"fmt"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/model/protoio"
)

func (app *SideChain) checkTx(txbt []byte) uint32 {
	txbox := new(model.TxBox)
	err := protoio.ReadMessage(bytes.NewBuffer(txbt), txbox)
	if err != nil {
		return CodeTypeEncodingError
	}

	innerTx := new(model.Tx)
	err = protoio.ReadMessage(bytes.NewBuffer(txbox.Tx), innerTx)
	if err != nil {
		return CodeTypeEncodingError
	}

	// 验证交易签名
	if err := model.VerifyTxSigner(innerTx); err != nil {
		return CodeTypeInvalidTxFormat
	}

	if len(txbox.Org) == 0 {
		fmt.Println("invalid node1")
		return CodeInvalidNode
	}

	keys := app.p2p.AllNodes()
	isIn := false
	for _, key := range keys {
		if bytes.Equal(txbox.Org, key.Byte()) {
			isIn = true
			break
		}
	}

	if !isIn {
		return CodeInvalidNode
	}

	return CodeTypeOK
}
