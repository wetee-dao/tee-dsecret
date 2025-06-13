package peer

import (
	"go.dedis.ch/kyber/v4"
	"wetee.app/dsecret/internal/model"
)

type Peer interface {
	Send(node *model.Node, pid string, message *model.Message) error
	Pub(topic string, data []byte) error
	Sub(topic string, handler func(*model.Message) error) error

	GoStart()
	Close() error

	PeerStrID() string
	NodeIds() []string
	Nodes() []*model.Node

	Version() uint32
	NetResetHook(hook func([]kyber.Point) error)
}
