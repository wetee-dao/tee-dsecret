package local

import (
	"context"
	"fmt"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	libp2pCrypto "github.com/libp2p/go-libp2p/core/crypto"
	types "wetee.app/dsecret/type"
	"wetee.app/dsecret/util"
)

var peers = make(map[string]*Peer)
var GlobleNodes = make([]*types.Node, 0, 100)

func NewNetwork(ctx context.Context, priv *types.PrivKey, boots []string, nodes []*types.Node, tcp, udp uint32) (*Peer, error) {
	id := priv.GetPublic().PeerID().String()

	// 创建 P2P 网络实例
	peer := &Peer{
		id:       id,
		privKey:  priv.PrivKey,
		nodes:    nodes,
		handlers: make(map[string]func(*types.Message) error),
		netHook: func() error {
			fmt.Println("network hook not implement")
			return nil
		},
	}

	peers[id] = peer

	return peer, nil
}

type Peer struct {
	id       string
	privKey  libp2pCrypto.PrivKey
	nodes    []*types.Node
	handlers map[string]func(*types.Message) error
	netHook  func() error
}

func (p *Peer) Send(ctx context.Context, node *types.Node, pid string, message *types.Message) error {
	util.LogSendmsg(">>>>>> P2P Send()", "to", node.ID.PeerID(), "-", node.ID.SS58(), "| type:", pid+"."+message.Type)
	peer := peers[node.PeerID().String()]
	if handler, ok := peer.handlers[pid]; ok {
		go handler(message)
	}
	return nil
}

func (p *Peer) AddHandler(pid string, handler func(*types.Message) error) {
	p.handlers[pid] = handler
}

func (p *Peer) RemoveHandler(pid string) {
	delete(p.handlers, pid)
}

func (p *Peer) Close() error {
	return nil
}

func (p *Peer) Start(ctx context.Context) {
	time.Sleep(time.Second * 2)

	p.nodes = GlobleNodes

	// 触发网络钩子
	p.netHook()
}

func (p *Peer) Discover(ctx context.Context) error {
	return nil
}

func (p *Peer) PeerStrID() string {
	return p.id
}

func (p *Peer) Pub(ctx context.Context, topic string, data []byte) error {
	panic("Pub not implement")
}

func (p *Peer) Sub(ctx context.Context, topic string) (*pubsub.Subscription, error) {
	panic("Sub not implement")
}

func (p *Peer) NetResetHook(hook func() error) {
	p.netHook = hook
}

func (p *Peer) NodeIds() []string {
	ns := make([]string, len(p.nodes))
	for i, n := range p.nodes {
		ns[i] = n.PeerID().String()
	}

	return ns
}

func (p *Peer) Nodes() []*types.Node {
	return p.nodes
}

func (p *Peer) Version() uint32 {
	return 1
}
