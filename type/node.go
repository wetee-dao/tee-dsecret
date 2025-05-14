package types

import "github.com/libp2p/go-libp2p/core/peer"

type Node struct {
	ID   PubKey `json:"id"`
	Type uint8  `json:"type"` // 0: worker, 1: dsecret
}

// 计算 PeerID
func (n *Node) PeerID() peer.ID {
	return n.ID.PeerID()
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
