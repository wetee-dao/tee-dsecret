package sidechain

import (
	"errors"
	"fmt"
	"sort"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/wetee-dao/tee-dsecret/pkg/chains"
	"github.com/wetee-dao/tee-dsecret/pkg/dkg"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/util"
)

// sendPartialSign sends the partial signatures to the main chain.
func (s *SideChain) sendPartialSign(tx_index int64, hubs []*model.HubCall, proposer *model.PubKey) {
	if len(hubs) == 0 {
		return
	}

	indexCalls := make([]*model.IndexCall, 0, len(hubs))
	for _, hub := range hubs {
		if hub == nil {
			util.LogWithRed("sendPartialSign", "hub is nil")
			continue
		}
		indexCalls = append(indexCalls, hub.Call...)
	}

	sort.Slice(indexCalls, func(i, j int) bool {
		return indexCalls[i].Index < indexCalls[j].Index
	})

	calls := make([]types.Call, 0, len(indexCalls))
	for _, bt := range indexCalls {
		c := new(types.Call)
		codec.Decode(bt.Call, c)
		calls = append(calls, *c)
	}

	client := chains.MainChain.GetClient()
	call, err := client.BatchCall("batch_all", calls)
	if err != nil {
		util.LogWithRed("sendPartialSign", "BatchCall error", err)
		return
	}

	signer := dkg.NewDssSigner(s.dkg)
	sig, err := client.PartialSign(signer, *call)
	if err != nil {
		util.LogWithRed("sendPartialSign", "PartialSign error", err)
		return
	}

	psig := &model.BlockPartialSign{
		OrgId:   s.dkg.P2PId().String(),
		HubSig:  sig,
		TxIndex: tx_index,
	}

	err = model.SetCodec("G", "tx_index"+fmt.Sprint(tx_index), *call)
	if err != nil {
		util.LogWithRed("sendPartialSign", "SetKey error", err)
		return
	}

	err = s.p2p.Send(*proposer, "block-partial-sign", psig)
	if err != nil {
		util.LogWithRed("sendPartialSign", "Send error", err)
	}
}

// HandlePartialSign handles the block partial sign messages received via P2P.
func (s *SideChain) revPartialSign(msgBox any) error {
	if s.txCh == nil {
		return errors.New("txCh is nil")
	}

	s.txCh.Push(msgBox.(*model.BlockPartialSign))
	return nil
}

// handle block partial sign message
func (s *SideChain) handlePartialSign(msg *model.BlockPartialSign) error {
	err := s.SavePartialSig(msg.OrgId, msg)
	if err != nil {
		util.LogWithRed("HandlePartialSign", "SaveSig error", err)
		return err
	}

	sigs, err := s.SigListOfTx(msg.TxIndex)
	if err != nil {
		util.LogWithRed("HandlePartialSign", "GetSigList error", err)
		return err
	}

	if len(sigs) < s.dkg.Threshold {
		util.LogWithGray("HandlePartialSign", "sigs", len(sigs), "threshold", s.dkg.Threshold)
		return nil
	}

	if len(sigs) != s.dkg.Threshold {
		util.LogWithGray("HandlePartialSign sigs => ", len(sigs), ">", s.dkg.Threshold)
		return nil
	}

	shares := make([][]byte, 0, len(sigs))
	for _, sig := range sigs {
		shares = append(shares, sig.HubSig)
	}

	err = s.SyncToHub(msg.TxIndex, shares)
	if err != nil {
		return err
	}

	return nil
}

const PartialSigPrefix = "partial_sig_"

func (s *SideChain) SavePartialSig(user_id string, msg *model.BlockPartialSign) error {
	bt, _ := msg.Marshal()
	return model.SetKey("G", PartialSigPrefix+fmt.Sprint(msg.TxIndex)+"_"+user_id, bt)
}

func (s *SideChain) SigListOfTx(txIndex int64) ([]*model.BlockPartialSign, error) {
	bts, err := model.GetList("G", PartialSigPrefix+fmt.Sprint(txIndex)+"_", 1, 10000)
	if err != nil {
		return nil, err
	}

	sigs := make([]*model.BlockPartialSign, 0, len(bts))
	for _, bt := range bts {
		msg := new(model.BlockPartialSign)
		err := msg.Unmarshal(bt)
		if err != nil {
			util.LogWithRed("GetSig", "Unmarshal error", err)
			continue
		}

		if msg.TxIndex == txIndex {
			sigs = append(sigs, msg)
		}
	}

	return sigs, nil
}
