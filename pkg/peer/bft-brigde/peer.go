package bftbrigde

import (
	"errors"

	"github.com/cometbft/cometbft/p2p"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func (p *BTFReactor) Send(node model.PubKey, topic string, message any) error {
	sendData := p2p.Envelope{}
	switch topic {
	case "dkg":
		sendData.ChannelID = topics["dkg"].ID
		sendData.Message = message.(*model.DkgMessage)
	case "block-partial-sign":
		sendData.ChannelID = topics["block-partial-sign"].ID
		sendData.Message = message.(*model.BlockPartialSign)
	}

	peers := p.Switch.Peers()
	peers.ForEach(func(p p2p.Peer) {
		if node.SideChainNodeID() == p.ID() {
			// util.LogError("P2P Send To", node.SS58(), topic+"."+message.Type)
			p.Send(sendData)
		}
	})

	return nil
}

func (p *BTFReactor) Pub(topic string, data []byte) error {
	channel, ok := topics[topic]
	if !ok {
		return errors.New("topic not found")
	}

	p.Switch.Broadcast(p2p.Envelope{
		ChannelID: channel.ID,
	})

	return nil
}

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

func (p *BTFReactor) PeerID() string {
	return ""
}

func (p *BTFReactor) Nodes() []*model.PubKey {
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

func (p *BTFReactor) Nodekeys() []*model.PubKey {
	return p.nodekeys
}

func (p *BTFReactor) LinkToNetwork() {

}
