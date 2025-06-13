package local

import (
	"context"
	"fmt"
	"time"

	libp2pCrypto "github.com/libp2p/go-libp2p/core/crypto"
	"go.dedis.ch/kyber/v4"
	"wetee.app/dsecret/internal/model"
	"wetee.app/dsecret/internal/util"
)

var (
	peers       = make(map[string]*Peer)
	GlobleNodes = make([]*model.Node, 0, 100)
	Version     uint32
	Commits     []kyber.Point
)

func NewNetwork(ctx context.Context, priv *model.PrivKey, boots []string, nodes []*model.Node, tcp, udp uint32) (*Peer, error) {
	id := priv.GetPublic().PeerID().String()

	// 创建 P2P 网络实例
	peer := &Peer{
		id:       id,
		privKey:  priv.PrivKey,
		nodes:    nodes,
		handlers: make(map[string]func(*model.Message) error),
		netHook: func([]kyber.Point) error {
			fmt.Println("::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::: netHook not found error")
			return nil
		},
		version: 1,
	}

	peers[id] = peer

	return peer, nil
}

type Peer struct {
	id       string
	privKey  libp2pCrypto.PrivKey
	nodes    []*model.Node
	handlers map[string]func(*model.Message) error
	netHook  func([]kyber.Point) error
	version  uint32
}

func (p *Peer) Send(node *model.Node, pid string, message *model.Message) error {
	util.LogSendmsg(">>>>>> P2P Send()", "to", node.ID.PeerID(), "-", node.ID.SS58(), "| type:", pid+"."+message.Type)
	peer := peers[node.PeerID().String()]
	if handler, ok := peer.handlers[pid]; ok {
		go handler(message)
	} else {
		fmt.Println("handler not found for pid: ", node.ID.PeerID())
	}
	return nil
}
func (p *Peer) Close() error {
	return nil
}

func (p *Peer) GoStart() {
	p.nodes = GlobleNodes
	p.version = Version

	time.Sleep(500 * time.Millisecond)

	// 触发网络钩子
	err := p.netHook(Commits)
	if err != nil {
		fmt.Println("netHook error: ", err)
	}
}

func (p *Peer) PeerStrID() string {
	return p.id
}

func (p *Peer) Pub(topic string, data []byte) error {
	panic("Pub not implement")
}

func (p *Peer) Sub(topic string, handler func(*model.Message) error) error {
	p.handlers[topic] = handler
	return nil
}

func (p *Peer) NetResetHook(hook func([]kyber.Point) error) {
	p.netHook = hook
}

func (p *Peer) NodeIds() []string {
	ns := make([]string, len(p.nodes))
	for i, n := range p.nodes {
		ns[i] = n.PeerID().String()
	}

	return ns
}

func (p *Peer) Nodes() []*model.Node {
	return p.nodes
}

func (p *Peer) Version() uint32 {
	return p.version
}
