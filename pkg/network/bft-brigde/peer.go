package bftbrigde

import (
	"errors"

	"github.com/cometbft/cometbft/p2p"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func (p *BTFReactor) Send(node *model.To, message any) error {
	sendData := p2p.Envelope{}

	// set sender id
	sendData.Src = LocalPeer{id: p.Switch.NodeInfo().ID()}
	switch msg := message.(type) {
	case *model.DkgMessage:
		sendData.ChannelID = topics["dkg"].ID
		msg.To = node
		sendData.Message = msg
	case *model.BlockPartialSign:
		sendData.ChannelID = topics["block-partial-sign"].ID
		msg.To = node
		sendData.Message = msg
	default:
		return errors.New("unknown message type")
	}

	p.Receive(sendData)
	p.Switch.Broadcast(sendData)
	return nil
}

// func (p *BTFReactor) Pub(topic string, data []byte) error {
// 	channel, ok := topics[topic]
// 	if !ok {
// 		return errors.New("topic not found")
// 	}

// 	p.Switch.Broadcast(p2p.Envelope{
// 		ChannelID: channel.ID,
// 	})

// 	return nil
// }

func (p *BTFReactor) Sub(topic string, handler func(any) error) error {
	switch topic {
	case "dkg":
		p.dkgHandler = handler
	case "block-partial-sign":
		p.blockPartialSignHandler = handler
	default:
		return errors.New("topic not found")
	}
	return nil
}

// Get all available nodes
func (p *BTFReactor) AvailableNodes() []*model.PubKey {
	peers := p.Switch.Peers()
	nodes := make([]*model.PubKey, 0, peers.Size())
	for _, n := range p.nodekeys {
		peers.ForEach(func(peer p2p.Peer) {
			if peer.ID() == n.SideChainNodeID() {
				nodes = append(nodes, n)
			}
		})
	}

	return nodes
}

// Get all nodes
func (p *BTFReactor) AllNodes() []*model.PubKey {
	return p.nodekeys
}
