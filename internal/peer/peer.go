package peer

import (
	"wetee.app/dsecret/internal/model"
)

type Peer interface {
	Send(node_id model.PubKey, pid string, message *model.Message) error
	Pub(topic string, data []byte) error
	Sub(topic string, handler func(*model.Message) error) error

	LinkToNetwork()

	// Stop() error
	// PeerID() string

	Nodes() []*model.PubKey

	// Epoch() uint32
	SetNetworkChangeBack(hook func(back_type string) error)
}
