package sidechain

import (
	"wetee.app/dsecret/chains"
	"wetee.app/dsecret/internal/model"
)

var stateEpoch uint32 = 0
var stateLastEpochBlock uint32 = 0

func (app *SideChain) CheckEpoch() {
	epoch, lastEpochBlock, now, err := chains.ChainIns.GetEpoch()
	if err != nil {
		return
	}

	// util.LogWithGreen("SideChain CheckEpoch", "epoch:", epoch, "lastEpochBlock:", lastEpochBlock, "now:", now)
	if epoch > stateEpoch || now-lastEpochBlock >= 72000 || epoch <= 1 {
		validators, err := chains.ChainIns.GetValidatorList()
		if err != nil {
			return
		}

		// util.LogWithGreen("SideChain GetValidatorList", "validators:", validators)
		app.dkg.TryConsensus(model.ConsensusMsg{
			Validators: validators,
			Epoch:      epoch,
		})
	}

	return
}
