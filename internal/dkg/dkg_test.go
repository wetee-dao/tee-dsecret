package dkg

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"wetee.app/dsecret/internal/model"
	"wetee.app/dsecret/internal/peer"
	"wetee.app/dsecret/internal/peer/local"
	types "wetee.app/dsecret/type"
)

var peerSecret = []string{
	"080112407512939e37970b04c2b9a6060b16654473cf0721b71f8e56126ee314cbd0a7e9fe9125a81688d932ea792e9722c777f7696117363f86a107ab9d3681a8c922c8",
	"08011240dc72ceaf44e1e382a2ee4bbf47d6eabcca460740d5c28dd7bc097db90700594636922ff5c13eded54d1ba710bf043eed221112eb726da532149e1e3619fb9336",
	"08011240ebecd85320c3a05a1c170738c321a58d6560be9a1fbb323028f24aa874dece755a10bfc8865e3015c7496a3831d0e5e3abaadf964116d5ea5d8ada52135286a1",
}

var newPeerSecret = []string{
	"080112402d1f5379cc0475c169e7d5ee5e53437c0b6fe60dc4ea822c61d7db9dfe5e33b8931f5771269844909f1e73bd51e9da82eb8a7b2def0313417003f978d9010eb4",
	"0801124066dfdf585852c8aa4b111c15b69d4b40ad6db7cfabd7ae0b25384c950f5018cb2b4c0a71a9d97d6eb9bfcb29acc0d71df41015c5c5cc697d9471728a2a508a8a",
}

func TestNetwork(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := model.NewDB()
	if err != nil {
		require.NoErrorf(t, err, "failed store.InitDB")
		os.Exit(1)
	}

	nodes := []*types.Node{}
	for _, s := range peerSecret {
		nodeSecret, _ := types.PrivateKeyFromLibp2pHex(s)
		nodes = append(nodes, &types.Node{
			ID:   *nodeSecret.GetPublic(),
			Type: 1,
		})
	}

	peers := make([]peer.Peer, 0, len(nodes))
	for _, s := range peerSecret {
		nodeSecret, _ := types.PrivateKeyFromLibp2pHex(s)

		// 启动 P2P 网络
		peer, err := local.NewNetwork(ctx, nodeSecret, []string{}, nodes, uint32(0), uint32(0))
		require.NoErrorf(t, err, "failed peer.NewNetwork")

		peers = append(peers, peer)
	}

	dkgs := make([]*DKG, 0, len(nodes))
	for i, s := range peerSecret {
		nodeSecret, _ := types.PrivateKeyFromLibp2pHex(s)

		// 创建 DKG 实例
		dkg, err := NewRabinDKG(nodeSecret, peers[i])
		require.NoErrorf(t, err, "failed NewRabinDKG")
		dkgs = append(dkgs, dkg)

		if i == 2 {
			dkg.Start(ctx, Logger{
				NodeIndex: i,
			})
		} else {
			go dkg.Start(ctx, Logger{
				NodeIndex: i,
			})
		}
	}

	time.Sleep(time.Second)
	for _, d := range dkgs {
		fmt.Println("old nodes ", d.DkgPubKey.String())
		fmt.Println("old nodes ", d.DkgKeyShare.PriShare.String())
		fmt.Println("old nodes ", d.DkgKeyShare.Commits)
	}

	fmt.Println("///////////////////////////////////////////////////////////////////////////")
	local.Commits = dkgs[0].DkgKeyShare.Commits

	for _, s := range newPeerSecret {
		nodeSecret, _ := types.PrivateKeyFromLibp2pHex(s)
		n := &types.Node{
			ID:   *nodeSecret.GetPublic(),
			Type: 1,
		}
		nodes = append(nodes, n)
	}

	for _, s := range newPeerSecret {
		nodeSecret, _ := types.PrivateKeyFromLibp2pHex(s)

		// 启动 P2P 网络
		peer, err := local.NewNetwork(ctx, nodeSecret, []string{}, nodes, uint32(0), uint32(0))
		require.NoErrorf(t, err, "failed peer.NewNetwork")

		peers = append(peers, peer)
	}

	for i, s := range newPeerSecret {
		nodeSecret, _ := types.PrivateKeyFromLibp2pHex(s)

		// 创建 DKG 实例
		dkg, err := NewRabinDKG(nodeSecret, peers[3+i])
		require.NoErrorf(t, err, "failed NewRabinDKG")

		dkgs = append(dkgs, dkg)
	}

	/// reshare
	local.Version = 2
	local.GlobleNodes = nodes
	for i, p := range peers {
		if i == len(peers)-1 {
			p.Start(ctx)
		} else {
			go p.Start(ctx)
		}
	}

	time.Sleep(time.Second * 10)

	for _, d := range dkgs {
		fmt.Println("new nodes ", d.DkgPubKey.String())
		fmt.Println("new nodes ", d.DkgKeyShare.PriShare.String())
		fmt.Println("new nodes ", d.DkgKeyShare.Commits)
	}
}
