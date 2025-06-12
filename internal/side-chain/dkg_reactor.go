package sidechain

import (
	"github.com/cometbft/cometbft/libs/service"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/p2p/conn"
)

type DKGReactor struct {
	service.BaseService // Provides Start, Stop, .Quit
	Switch              *p2p.Switch
}

func NewBaseReactor(name string, impl p2p.Reactor) *DKGReactor {
	return &DKGReactor{
		BaseService: *service.NewBaseService(nil, name, impl),
		Switch:      nil,
	}
}

func (dr *DKGReactor) SetSwitch(sw *p2p.Switch) {
	dr.Switch = sw
}
func (*DKGReactor) GetChannels() []*conn.ChannelDescriptor { return nil }
func (*DKGReactor) AddPeer(p2p.Peer)                       {}
func (*DKGReactor) RemovePeer(p2p.Peer, any)               {}
func (*DKGReactor) Receive(p2p.Envelope)                   {}
func (*DKGReactor) InitPeer(peer p2p.Peer) p2p.Peer        { return peer }
