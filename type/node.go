package types

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
)

type Node struct {
	ID   string `json:"id"`
	Type uint8  `json:"type"` // 0: worker, 1: dsecret
}

// 计算 PeerID
func (n *Node) PeerID() peer.ID {
	pk, err := PublicKeyFromLibp2pHex(n.ID)
	if err != nil {
		fmt.Println("Node types.PublicKeyFromHex error:", err)
		return peer.ID("")
	}
	peerID, err := peer.IDFromPublicKey(pk)
	if err != nil {
		fmt.Println("Node peer.IDFromPublicKey error:", err)
		return peer.ID("")
	}
	return peer.ID(peerID)
}

// func NodeFromString(id string) (*Node, error) {
// 	bt, err := b58.Decode(id)
// 	if err != nil {
// 		return nil, errors.New("b58.Decode error: " + err.Error())
// 	}

// 	pid := peer.ID(string(id))
// 	pub, err := pid.ExtractPublicKey()
// 	if err != nil {
// 		return nil, err
// 	}

// 	bt, err = libp2pCrypto.MarshalPublicKey(pub)
// 	if err != nil {
// 		return nil, err
// 	}
// 	id = hex.EncodeToString(bt)

// 	return &Node{
// 		ID:   id,
// 		Type: 0,
// 	}, nil
// }
