package bftbrigde

import (
	"errors"

	"github.com/cometbft/cometbft/p2p"
	"go.dedis.ch/kyber/v4"
	"wetee.app/dsecret/internal/model"
)

func (p *BTFReactor) Send(node *model.Node, topic string, message *model.Message) error {
	channel, ok := topics[topic]
	if !ok {
		return errors.New("topic not found")
	}
	peers := p.Switch.Peers()
	peers.ForEach(func(p p2p.Peer) {
		if string(node.PeerID()) == string(p.ID()) {
			p.Send(p2p.Envelope{
				ChannelID: channel.ID,
				Message: &DkgMessage{
					MsgId:   message.MsgID,
					Payload: message.Payload,
				},
			})
		}
	})
	return nil
}

func (p *BTFReactor) Close() error {
	return nil
}

func (p *BTFReactor) GoStart() {

}

func (p *BTFReactor) PeerStrID() string {
	return ""
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
	calls := p.callbacks[topic]
	if calls == nil {
		calls = []func(*model.Message) error{}
	}

	p.callbacks[topic] = append(calls, handler)
	return nil
}

func (p *BTFReactor) NetResetHook(hook func([]kyber.Point) error) {
	p.netHook = hook
}

func (p *BTFReactor) NodeIds() []string {
	return []string{}
}

func (p *BTFReactor) Nodes() []*model.Node {
	return []*model.Node{}
}

func (p *BTFReactor) Version() uint32 {
	return 1
}
