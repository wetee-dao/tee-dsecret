package dkg

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"wetee.app/dsecret/peer"
	"wetee.app/dsecret/peer/local"
	"wetee.app/dsecret/store"
	types "wetee.app/dsecret/type"
)

var peerSecret = []string{
	"080112407512939e37970b04c2b9a6060b16654473cf0721b71f8e56126ee314cbd0a7e9fe9125a81688d932ea792e9722c777f7696117363f86a107ab9d3681a8c922c8",
	"08011240dc72ceaf44e1e382a2ee4bbf47d6eabcca460740d5c28dd7bc097db90700594636922ff5c13eded54d1ba710bf043eed221112eb726da532149e1e3619fb9336",
	"08011240ebecd85320c3a05a1c170738c321a58d6560be9a1fbb323028f24aa874dece755a10bfc8865e3015c7496a3831d0e5e3abaadf964116d5ea5d8ada52135286a1",
}

func TestNetwork(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := store.InitDB("")
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

	peers := make([]peer.Peer, len(nodes))
	for i, s := range peerSecret {
		nodeSecret, _ := types.PrivateKeyFromLibp2pHex(s)

		// 启动 P2P 网络
		peer, err := local.NewNetwork(ctx, nodeSecret, []string{}, nodes, uint32(0), uint32(0))
		require.NoErrorf(t, err, "failed peer.NewNetwork")

		peers[i] = peer
	}

	for i, s := range peerSecret {
		nodeSecret, _ := types.PrivateKeyFromLibp2pHex(s)

		// 创建 DKG 实例
		dkg, err := NewRabinDKG(nodeSecret, peers[i])
		require.NoErrorf(t, err, "failed NewRabinDKG")

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

	time.Sleep(time.Second * 10)
}
