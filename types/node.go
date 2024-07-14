package types

import (
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p/core/peer"
)

type Node struct {
	ID   string `json:"id"`
	Addr string `json:"addr"`
}

func (n *Node) PeerID() peer.ID {
	pk, err := PublicKeyFromHex("08011220" + n.ID)
	if err != nil {
		fmt.Println("Node types.PublicKeyFromHex error:", err)
		os.Exit(1)
	}
	peerID, err := peer.IDFromPublicKey(pk)
	if err != nil {
		fmt.Println("Node peer.IDFromPublicKey error:", err)
		os.Exit(1)
	}
	return peer.ID(peerID)
}
