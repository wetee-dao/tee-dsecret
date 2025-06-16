package sidechain

import (
	"wetee.app/dsecret/chains"
)

var stateEpoch uint32 = 0
var stateLastEpochBlock uint32 = 0

func (app *SideChain) CheckEpoch() {
	epoch, lastEpochBlock, now, err := chains.ChainIns.GetEpoch()
	if err != nil {
		return
	}

	if epoch > stateEpoch || now-lastEpochBlock >= 72000 || epoch == 0 {
		// validators, err := chains.ChainIns.GetValidatorList()
		// if err != nil {
		// 	return
		// }

		// app.dkg.StartConsensus([]kyber.Point{}, validators, epoch)
	}

	return
}
