package local

import (
	"crypto/ed25519"
	"fmt"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"go.dedis.ch/kyber/v4"
)

var (
	peers = make(map[string]*Peer)
)

func NewNetwork(priv *model.PrivKey, boots []string, nodes []*model.PubKey, tcp, udp uint32) (*Peer, error) {
	id := priv.GetPublic().String()

	// 创建 P2P 网络实例
	peer := &Peer{
		id:       id,
		privKey:  priv.PrivateKey,
		nodes:    nodes,
		handlers: make(map[string]func(any) error),
		callBack: func(ty string) error {
			fmt.Println("::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::: netHook not found error")
			return nil
		},
		version: 1,
	}

	peers[id] = peer

	return peer, nil
}

type Peer struct {
	id         string
	privKey    ed25519.PrivateKey
	nodes      []*model.PubKey
	handlers   map[string]func(any) error
	callBack   func(string) error
	version    uint32
	PreCommits []kyber.Point
}

func (p *Peer) Send(node model.PubKey, topic string, message any) error {
	// util.LogSendmsg(">>>>>> P2P Send()", "to", node.String(), "-", node.SS58(), "| type:", topic+"."+message.Type)
	peer := peers[node.String()]
	if handler, ok := peer.handlers[topic]; ok {
		go handler(message)
	} else {
		fmt.Println("handler not found for topic: ", topic, "node", node)
	}
	return nil
}

func (p *Peer) LinkToNetwork() {

}

func (p *Peer) PeerID() string {
	return p.id
}

func (p *Peer) Pub(topic string, data []byte) error {
	panic("Pub not implement")
}

func (p *Peer) Sub(topic string, handler func(any) error) error {
	p.handlers[topic] = handler
	return nil
}

func (p *Peer) Nodes() []*model.PubKey {
	return p.nodes
}

func (p *Peer) SetNodes(nodes []*model.PubKey) {
	p.nodes = nodes
}
