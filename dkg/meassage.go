package dkg

import (
	"context"
	"encoding/json"

	rabin "go.dedis.ch/kyber/v3/share/dkg/rabin"
	"wetee.app/dsecret/types"
)

// SendDealMessage 发送Deal消息
func (dkg *DKG) SendDealMessage(ctx context.Context, node *types.Node, message *rabin.Deal) error {
	pmessage, err := types.DealToProtocol(message)
	if err != nil {
		return err
	}

	bt, err := json.Marshal(pmessage)
	if err != nil {
		return err
	}

	return dkg.Peer.Send(ctx, node, "deal", &types.Message{
		Type:    "deal",
		Payload: bt,
	})
}
