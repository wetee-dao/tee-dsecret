package sidechain

import (
	"bytes"
	"errors"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	abcicli "github.com/cometbft/cometbft/abci/client"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/wetee-dao/tee-dsecret/pkg/chains"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/util"
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

func TEECallToHubCall(tcall *model.TeeCall, dkgKey types.AccountID) (*types.Call, error) {
	callWrap := tcall.Tx
	switch tx := callWrap.(type) {
	case *model.TeeCall_PodStart:
		pod := tx.PodStart

		pod_key, err := types.NewAccountID(tcall.Caller)
		if err != nil {
			return nil, errors.New("get pod key error")
		}

		call, err := chains.MainChain.TxCallOfStartPod(pod.Id, *pod_key, dkgKey)
		if err != nil {
			util.LogError("TxCallOfStartPod", err)
			return nil, err
		}

		return call, nil
	case *model.TeeCall_PodMint:
		pod := tx.PodMint
		call, err := chains.MainChain.TxCallOfMintPod(pod.Id, types.NewH256(pod.ReportHash), dkgKey)
		if err != nil {
			util.LogError("TxCallOfStartPod", err)
			return nil, err
		}

		return call, nil
	case *model.TeeCall_UploadSecret:
		upload := tx.UploadSecret
		call, err := chains.MainChain.TxCallOfUploadSecret(types.NewH160(upload.User), upload.Index, dkgKey)
		if err != nil {
			util.LogError("TxCallOfUploadSecret", err)
			return nil, err
		}

		return call, nil
	case *model.TeeCall_InitDisk:
		init := tx.InitDisk
		call, err := chains.MainChain.TxCallOfInitDisk(types.NewH160(init.User), init.Index, types.NewH256(init.Hash), dkgKey)
		if err != nil {
			util.LogError("TxCallOfInitDisk", err)
			return nil, err
		}

		return call, nil
		// case *model.TeeCall_BridgeCall:
	}
	return nil, errors.New("invalid tee call")
}
