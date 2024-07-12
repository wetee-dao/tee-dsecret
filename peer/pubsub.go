package peer

import (
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	pb "github.com/libp2p/go-libp2p-pubsub/pb"
	"github.com/libp2p/go-libp2p/core/peer"
)

var _ pubsub.EventTracer = (*pubsubTracer)(nil)

type pubsubTracer struct{}

func (p *pubsubTracer) Trace(evt *pb.TraceEvent) {
	// log.Debugf("PUBSUB EVENT TRACE: %s", evt.Type)
	switch evt.Type.String() {
	case pb.TraceEvent_DELIVER_MESSAGE.String():
		pid := peer.ID(string(evt.DeliverMessage.ReceivedFrom))
		fmt.Println("pubsub.tracer: event type ", evt.Type, " from ", pid, " on topic ", *(evt.DeliverMessage.Topic))
	case pb.TraceEvent_PUBLISH_MESSAGE.String():
		fmt.Println("pubsub.tracer: event type ", evt.Type, " on topic", *(evt.PublishMessage.Topic))
	}
}
