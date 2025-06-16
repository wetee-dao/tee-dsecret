package bftbrigde

import (
	"errors"

	"github.com/cometbft/cometbft/p2p"
	"wetee.app/dsecret/internal/model"
)

func (p *BTFReactor) Send(node model.PubKey, topic string, message *model.Message) error {
	channel, ok := topics[topic]
	if !ok {
		return errors.New("topic not found")
	}

	peers := p.Switch.Peers()
	peers.ForEach(func(p p2p.Peer) {
		if node.SideChainNodeID() == p.ID() {
			p.Send(p2p.Envelope{
				ChannelID: channel.ID,
				Message: &DkgMessage{
					Type:    message.Type,
					MsgId:   message.MsgID,
					Payload: message.Payload,
				},
			})
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

func (p *BTFReactor) Sub(topic string, handler func(*model.Message) error) error {
	p.messageHandlers[topic] = handler
	return nil
}

func (p *BTFReactor) PeerID() string {
	return ""
}

func (p *BTFReactor) SetNetworkChangeBack(hook func(string) error) {
	p.callDkg = hook
}

func (p *BTFReactor) Nodes() []*model.PubKey {
	return p.nodekeys
}

func (p *BTFReactor) LinkToNetwork() {

}
