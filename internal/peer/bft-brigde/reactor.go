package bftbrigde

import (
	bcproto "github.com/cometbft/cometbft/api/cometbft/blocksync/v1"
	"github.com/cometbft/cometbft/libs/service"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/p2p/conn"
	"github.com/cometbft/cometbft/types"
	"go.dedis.ch/kyber/v4"

	"wetee.app/dsecret/chains"
	"wetee.app/dsecret/internal/model"
	"wetee.app/dsecret/internal/util"
)

var MaxMsgSize = types.MaxBlockSizeBytes
var topics = map[string]p2p.ChannelDescriptor{
	"dkg": p2p.ChannelDescriptor{
		ID:                  byte(0xFF),
		Priority:            5,
		SendQueueCapacity:   1000,
		RecvBufferCapacity:  50 * 4096,
		RecvMessageCapacity: MaxMsgSize,
		MessageType:         &DkgMessage{},
	},
}

type BTFReactor struct {
	service.BaseService
	Switch *p2p.Switch

	mainChain chains.MainChain
	callbacks map[string][]func(*model.Message) error
	netHook   func([]kyber.Point) error
}

func NewBTFReactor(name string, main chains.MainChain) *BTFReactor {
	r := &BTFReactor{
		callbacks: map[string][]func(*model.Message) error{},
	}
	r.BaseService = *service.NewBaseService(nil, name, r)

	return r

}

// 实现 Service 接口的 OnStart 生命周期钩子
func (r *BTFReactor) OnStart() error {
	nodeInfo := r.Switch.NodeInfo()
	address, _ := nodeInfo.NetAddress()
	util.LogError("Local Node", address.String())

	r.PrintPeers("BTF OnStart")

	// 启动协程、初始化资源
	return nil
}

// 实现 OnStop 生命周期钩子
func (r *BTFReactor) OnStop() {
	util.LogError("BTF OnStop")
}

func (r *BTFReactor) OnReset() error {
	util.LogError("BTF OnReset")
	return nil
}

func (dr *BTFReactor) SetSwitch(sw *p2p.Switch) {
	dr.Switch = sw
}

func (*BTFReactor) GetChannels() []*conn.ChannelDescriptor {
	return []*p2p.ChannelDescriptor{
		{
			ID:                  byte(0xFF),
			Priority:            5,
			SendQueueCapacity:   1000,
			RecvBufferCapacity:  50 * 4096,
			RecvMessageCapacity: MaxMsgSize,
			MessageType:         &bcproto.Message{},
		},
	}
}

func (r *BTFReactor) AddPeer(p2p.Peer) {
	r.PrintPeers("BTF AddPeer")
}

func (r *BTFReactor) RemovePeer(p2p.Peer, any) {
	r.PrintPeers("BTF RemovePeer")
}

func (*BTFReactor) Receive(e p2p.Envelope) {
	switch msg := e.Message.(type) {
	case *DkgMessage:
		util.LogError("Receive tee msg", "e.Src.ID()", msg)
	default:
		util.LogError("Receive error", "msg", msg)
	}
}

func (*BTFReactor) InitPeer(peer p2p.Peer) p2p.Peer {
	return peer
}

func (r *BTFReactor) PrintPeers(event string) {
	outbound, inbound, dialing := r.Switch.NumPeers()
	util.LogError(event, "Peers outbound=>", outbound, "inbound=>", inbound, "dialing=>", dialing)
}
