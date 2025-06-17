package dkg

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	inkUtil "github.com/wetee-dao/ink.go/util"
	"wetee.app/dsecret/internal/model"
	"wetee.app/dsecret/internal/peer/local"
	"wetee.app/dsecret/internal/util"
)

var peerSecret = []string{
	"7512939e37970b04c2b9a6060b16654473cf0721b71f8e56126ee314cbd0a7e9fe9125a81688d932ea792e9722c777f7696117363f86a107ab9d3681a8c922c8",
	"dc72ceaf44e1e382a2ee4bbf47d6eabcca460740d5c28dd7bc097db90700594636922ff5c13eded54d1ba710bf043eed221112eb726da532149e1e3619fb9336",
	"ebecd85320c3a05a1c170738c321a58d6560be9a1fbb323028f24aa874dece755a10bfc8865e3015c7496a3831d0e5e3abaadf964116d5ea5d8ada52135286a1",
}

var newPeerSecret = []string{
	"2d1f5379cc0475c169e7d5ee5e53437c0b6fe60dc4ea822c61d7db9dfe5e33b8931f5771269844909f1e73bd51e9da82eb8a7b2def0313417003f978d9010eb4",
	"66dfdf585852c8aa4b111c15b69d4b40ad6db7cfabd7ae0b25384c950f5018cb2b4c0a71a9d97d6eb9bfcb29acc0d71df41015c5c5cc697d9471728a2a508a8a",
}

func TestNetwork(t *testing.T) {
	os.RemoveAll("./chain_data")

	db, err := model.NewDB()
	if err != nil {
		require.NoErrorf(t, err, "failed store.InitDB")
		os.Exit(1)
	}
	defer db.Close()

	nodes := []*model.PubKey{}
	validators := []*model.Validator{}
	for _, s := range peerSecret {
		nodeSecret, _ := model.PrivateKeyFromHex(s)
		nodes = append(nodes, nodeSecret.GetPublic())
		validators = append(validators, &model.Validator{
			ValidatorId: *nodeSecret.GetPublic(),
			P2pId:       *nodeSecret.GetPublic(),
		})
	}

	peers := make([]*local.Peer, 0, len(nodes))
	for _, s := range peerSecret {
		nodeSecret, _ := model.PrivateKeyFromHex(s)

		peer, err := local.NewNetwork(nodeSecret, []string{}, nodes, uint32(0), uint32(0))
		require.NoErrorf(t, err, "failed peer.NewNetwork")

		peers = append(peers, peer)
	}

	dkgs := make([]*DKG, 0, len(nodes))
	for i, s := range peerSecret {
		nodeSecret, _ := model.PrivateKeyFromHex(s)

		dkg, err := NewDKG(nodeSecret, peers[i], inkUtil.NewSome(validators), Logger{
			NodeTag: "NODE " + fmt.Sprint(i),
		})
		require.NoErrorf(t, err, "failed NewDKG")
		go dkg.Start()

		dkgs = append(dkgs, dkg)
	}

	dkgs[0].TryConsensus(model.ConsensusMsg{
		Validators: validators,
		Epoch:      1,
	})
	time.Sleep(time.Second * 1)

	for _, d := range dkgs {
		util.LogWithYellow("V0 |||", d.DkgKeyShare.PriShare.String())
	}

	util.LogWithGreen("----------------------------------------------------------------------------------------------------")

	for _, s := range newPeerSecret {
		nodeSecret, _ := model.PrivateKeyFromHex(s)

		// new peer
		nodes = append(nodes, nodeSecret.GetPublic())
		validators = append(validators, &model.Validator{
			ValidatorId: *nodeSecret.GetPublic(),
			P2pId:       *nodeSecret.GetPublic(),
		})
	}

	for _, s := range newPeerSecret {
		nodeSecret, _ := model.PrivateKeyFromHex(s)

		// 启动 P2P 网络
		peer, err := local.NewNetwork(nodeSecret, []string{}, nodes, uint32(0), uint32(0))
		require.NoErrorf(t, err, "failed peer.NewNetwork")

		peers = append(peers, peer)
	}

	// set nodes
	for _, peer := range peers {
		peer.SetNodes(nodes)
	}

	for i, s := range newPeerSecret {
		nodeSecret, _ := model.PrivateKeyFromHex(s)

		// 创建 DKG 实例
		dkg, err := NewDKG(nodeSecret, peers[3+i], inkUtil.NewSome(validators), Logger{
			NodeTag: "NODE " + fmt.Sprint(i),
		})
		require.NoErrorf(t, err, "failed NewDKG")
		go dkg.Start()

		dkgs = append(dkgs, dkg)
	}

	dkgs[0].TryConsensus(model.ConsensusMsg{
		Validators: validators,
		Epoch:      2,
	})
	time.Sleep(time.Second * 1)

	for _, d := range dkgs {
		if d.DkgKeyShare != nil {
			util.LogWithCyan("V1 |||", d.DkgKeyShare.PriShare.String())
		}
	}

	util.LogWithGreen("----------------------------------------------------------------------------------------------------")

	dkgs[0].TryConsensus(model.ConsensusMsg{
		Validators: validators,
		Epoch:      3,
	})
	time.Sleep(time.Second * 1)

	for _, d := range dkgs {
		if d.DkgKeyShare.PriShare.PriShare != nil {
			util.LogWithCyan("V2 |||", d.DkgKeyShare.PriShare.String())
		}
	}
}

// func TestSave(t *testing.T) {
// 	os.RemoveAll("./chain_data")

// 	db, err := model.NewDB()
// 	if err != nil {
// 		require.NoErrorf(t, err, "failed store.InitDB")
// 		os.Exit(1)
// 	}

// 	defer db.Close()

// 	nodeSecret, _ := model.PrivateKeyFromHex(peerSecret[0])
// 	peer, err := local.NewNetwork(nodeSecret, []string{}, []*model.PubKey{}, uint32(0), uint32(0))
// 	require.NoErrorf(t, err, "failed peer.NewNetwork")
// 	dkg, err := NewDKG(nodeSecret, peer, inkUtil.NewNone[[]*model.Validator](), Logger{
// 		NodeIndex: 0,
// 	})

// 	dkg.Epoch = 100
// 	go dkg.saveStore()

// 	time.Sleep(time.Second)

// 	dkg.Epoch = 1

// 	err = dkg.reStore()
// 	require.NoErrorf(t, err, "failed reStore")

// 	if dkg.Epoch != 100 {
// 		t.Errorf("dkg.Epoch != 100")
// 	}
// }
