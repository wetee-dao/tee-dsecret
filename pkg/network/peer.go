package peer

import (
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

type Peer interface {
	Send(to *model.To, message any) error
	Sub(topic string, handler func(any) error) error

	AvailableNodes() []*model.PubKey
	AllNodes() []*model.PubKey
}
