package sidechain

// import (
// 	"time"

// 	"wetee.app/dsecret/chains"
// 	"wetee.app/dsecret/internal/model"
// )

// var (
// 	lastSyncTime uint32 = 0
// 	epoch        uint32 = 0
// 	validators          = []*model.Validator{}
// 	bridgeCalls         = []string{}
// )

// func SyncMainChainState() {
// 	// 6s 同步一次主链状态
// 	now := uint32(time.Now().Unix())
// 	if now-lastSyncTime < 6 {
// 		return
// 	}

// 	var err error
// 	validators, _, err = chains.MainChain.GetNodes()
// 	if err != nil {
// 		return
// 	}

// }
