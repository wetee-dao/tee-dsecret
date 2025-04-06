package peer

import (
	"context"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"go.dedis.ch/kyber/v4"
	types "wetee.app/dsecret/type"
)

// func NewNetwork(ctx context.Context, ty string, priv *types.PrivKey, boots []string, nodes []*types.Node, tcp, udp uint32) (Peer, error) {
// 	if ty == "P2P" {
// 		return p2p.NewP2PNetwork(ctx, priv, boots, nodes, tcp, udp)
// 	} else if ty == "LOCAL" {
// 		return local.NewNetwork(ctx, priv, boots, nodes, tcp, udp)
// 	}
// 	return nil, errors.New("network 404")
// }

type Peer interface {
	Send(ctx context.Context, node *types.Node, pid string, message *types.Message) error
	AddHandler(pid string, handler func(*types.Message) error)
	RemoveHandler(pid string)
	Close() error
	Start(ctx context.Context)
	Discover(ctx context.Context) error
	PeerStrID() string
	Pub(ctx context.Context, topic string, data []byte) error
	Sub(ctx context.Context, topic string) (*pubsub.Subscription, error)
	NodeIds() []string
	Nodes() []*types.Node
	Version() uint32
	NetResetHook(hook func([]kyber.Point) error)
}
