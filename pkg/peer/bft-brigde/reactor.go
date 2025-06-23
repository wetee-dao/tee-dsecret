package bftbrigde

import (
	"errors"

	"github.com/cometbft/cometbft/libs/service"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/p2p/conn"
	"github.com/cometbft/cometbft/types"

	"github.com/wetee-dao/tee-dsecret/chains"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/util"
)

var MaxMsgSize = types.MaxBlockSizeBytes
var topics = map[string]p2p.ChannelDescriptor{
	"dkg": p2p.ChannelDescriptor{
		ID:                  byte(0xFF),
		Priority:            10000,
		SendQueueCapacity:   1000,
		RecvBufferCapacity:  50 * 4096,
		RecvMessageCapacity: MaxMsgSize,
		MessageType:         &DkgMessage{},
	},
}

type BTFReactor struct {
	service.BaseService
	Switch *p2p.Switch

	nodekeys   []*model.PubKey
	dkgHandler func(*model.Message) error
	callDkg    func(string) error
}

func NewBTFReactor(name string, main chains.Chain) *BTFReactor {
	r := &BTFReactor{}
	r.BaseService = *service.NewBaseService(nil, name, r)

	return r
}

// 实现 Service 接口的 OnStart 生命周期钩子
func (r *BTFReactor) OnStart() error {
	nodeInfo := r.Switch.NodeInfo()
	address, _ := nodeInfo.NetAddress()
	util.LogError("Local Address ", address.String())
	r.PrintPeers("P2P OnStart")

	// 启动协程、初始化资源
	return nil
}

// 实现 OnStop 生命周期钩子
func (r *BTFReactor) OnStop() {
	util.LogError("P2P OnStop")
}

func (r *BTFReactor) OnReset() error {
	util.LogError("P2P OnReset")
	return nil
}

func (dr *BTFReactor) SetSwitch(sw *p2p.Switch) {
	dr.Switch = sw
}

func (*BTFReactor) GetChannels() []*conn.ChannelDescriptor {
	channels := make([]*conn.ChannelDescriptor, 0, len(topics))
	for _, c := range topics {
		channels = append(channels, &c)
	}
	return channels
}

func (r *BTFReactor) AddPeer(p2p.Peer) {
	r.PrintPeers("P2P AddPeer")
}

func (r *BTFReactor) RemovePeer(p2p.Peer, any) {
	r.PrintPeers("P2P DelPeer")
}

func (r *BTFReactor) Receive(e p2p.Envelope) {
	switch msg := e.Message.(type) {
	case *DkgMessage:
		pub, err := r.GetPubkeyFromPeerID(e.Src.ID())
		if err == nil {
			// util.LogWithCyan("P2P Receive From", pub.SS58(), "dkg."+msg.Type)
			r.dkgHandler(&model.Message{
				MsgID:   msg.MsgId,
				Payload: msg.Payload,
				Type:    msg.Type,
				OrgId:   pub.String(),
			})
		} else {
			util.LogError("P2P PubkeyFromPeerID", "Receive unknown node", e.Src.ID())
		}
	default:
		util.LogError("P2P Receive", "Receive error", "msg", msg)
	}
}

func (*BTFReactor) InitPeer(peer p2p.Peer) p2p.Peer {
	return peer
}

func (r *BTFReactor) GetPubkeyFromPeerID(peer p2p.ID) (*model.PubKey, error) {
	for _, node := range r.nodekeys {
		if node.SideChainNodeID() == peer {
			return node, nil
		}
	}

	return nil, errors.New("not found")
}

func (r *BTFReactor) PrintPeers(event string) {
	if chains.MainChain == nil {
		return
	}
	_, pubkeys, err := chains.MainChain.GetNodes()
	if err == nil {
		r.nodekeys = pubkeys
	} else {
		util.LogError(event, "GetNodes", err.Error())
	}

	outbound, inbound, dialing := r.Switch.NumPeers()
	util.LogError(event, "Peers outbound=>", outbound, "inbound=>", inbound, "dialing=>", dialing, " || nodekeys=>", len(r.nodekeys))
	// r.Switch.Peers().ForEach(func(peer p2p.Peer) {
	// 	fmt.Println("             ", peer.ID())
	// })
}
