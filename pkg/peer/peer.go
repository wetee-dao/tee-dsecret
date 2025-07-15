package peer

import (
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

type Peer interface {
	Send(node_id model.PubKey, pid string, message any) error
	Pub(topic string, data []byte) error
	Sub(topic string, handler func(any) error) error

	LinkToNetwork()

	// Stop() error
	// PeerID() string

	Nodes() []*model.PubKey
}
