package bftbrigde

import (
	"errors"

	"github.com/cometbft/cometbft/libs/service"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/p2p/conn"
	"github.com/cometbft/cometbft/types"

	"github.com/wetee-dao/tee-dsecret/pkg/chains"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/util"
)

var MaxMsgSize = types.MaxBlockSizeBytes
var topics = map[string]p2p.ChannelDescriptor{
	"dkg": { // dkg msg
		ID:                  255,
		Priority:            10000,
		SendQueueCapacity:   1000,
		RecvBufferCapacity:  50 * 4096,
		RecvMessageCapacity: MaxMsgSize,
		MessageType:         &model.DkgMessage{},
	},
	"block-partial-sign": { // block partial sign msg
		ID:                  254,
		Priority:            10000,
		SendQueueCapacity:   1000,
		RecvBufferCapacity:  50 * 4096,
		RecvMessageCapacity: MaxMsgSize,
		MessageType:         &model.BlockPartialSign{},
	},
}

type BTFReactor struct {
	service.BaseService
	Switch *p2p.Switch

	id                      *model.PubKey
	validators              []*model.PubKey
	nodekeys                []*model.PubKey
	dkgHandler              func(any) error
	blockPartialSignHandler func(any) error
}

func NewBTFReactor(name string) *BTFReactor {
	r := &BTFReactor{}
	r.BaseService = *service.NewBaseService(nil, name, r)

	return r
}

// 实现 Service 接口的 OnStart 生命周期钩子
func (r *BTFReactor) OnStart() error {
	nodeInfo := r.Switch.NodeInfo()
	address, _ := nodeInfo.NetAddress()
	util.LogWithYellow("Local Address ", address.String())
	r.PrintPeers("P2P OnStart")

	return nil
}

// 实现 OnStop 生命周期钩子
func (r *BTFReactor) OnStop() {
	util.LogWithYellow("P2P OnStop")
}

func (r *BTFReactor) OnReset() error {
	util.LogWithYellow("P2P OnReset")
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
	case *model.DkgMessage:
		if !msg.To.Check(r.id) {
			return
		}

		if r.dkgHandler == nil {
			util.LogWithRed("P2P Receive", "dkgHandler not set")
			return
		}

		pub, err := r.GetPubkeyFromPeerID(e.Src.ID())
		if err != nil {
			util.LogWithRed("P2P PubkeyFromPeerID", "Receive unknown node", e.Src.ID())
		}

		msg.From = pub.String()
		r.dkgHandler(msg)
		return
	case *model.BlockPartialSign:
		if !msg.To.Check(r.id) {
			return
		}

		if r.blockPartialSignHandler == nil {
			util.LogWithRed("P2P Receive", "blockPartialSignHandler not set")
			return
		}

		pub, err := r.GetPubkeyFromPeerID(e.Src.ID())
		if err != nil {
			util.LogWithRed("P2P PubkeyFromPeerID", "Receive unknown node", e.Src.ID())
		}

		msg.From = pub.String()
		r.blockPartialSignHandler(msg)
	default:
		util.LogWithRed("P2P Receive", "Receive error", "msg", msg)
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

	// get from main chain
	validatorWrap, pubkeys, err := chains.MainChain.GetNodes()
	if err != nil {
		util.LogWithRed(event, "GetNodes", err.Error())
	}

	// save self nodekey
	for _, v := range pubkeys {
		if v.SideChainNodeID() == r.Switch.NodeInfo().ID() {
			r.id = v
		}
	}

	// save all nodekeys and validators
	r.nodekeys = pubkeys
	validators := make([]*model.PubKey, len(validatorWrap))
	for i, v := range validatorWrap {
		validators[i] = &v.P2pId
	}
	r.validators = validators

	outbound, inbound, dialing := r.Switch.NumPeers()
	util.LogWithYellow(event, "Peers outbound=>", outbound, "inbound=>", inbound, "dialing=>", dialing, " || nodekeys=>", len(r.nodekeys))
}
