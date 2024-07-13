package dkg

import (
	"context"
	"encoding/json"

	rabin "go.dedis.ch/kyber/v3/share/dkg/rabin"
	"wetee.app/dsecret/types"
)

// BroadcastMessage 广播消息给指定参与者。
func (dkg *DKG) BroadcastMessage(message *rabin.Deal) error {
	pmessage, err := types.DealToProtocol(message)
	if err != nil {
		return err
	}

	bt, err := json.Marshal(pmessage)
	if err != nil {
		return err
	}
	return dkg.Peer.Send(context.Background(), "deal", bt)
}
