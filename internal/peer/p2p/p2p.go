package p2p

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	libp2pCrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/event"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	"go.dedis.ch/kyber/v4"

	"wetee.app/dsecret/internal/chain"
	types "wetee.app/dsecret/type"
	"wetee.app/dsecret/util"
)

// NewP2PNetwork 创建一个新的 P2P 网络实例
func NewP2PNetwork(ctx context.Context, priv *types.PrivKey, tcp, udp uint32) (*Peer, error) {
	nodes, boots, version, err := GetChainNodes()
	if err != nil {
		return nil, err
	}

	var idht *dht.IpfsDHT
	var dhtOptions []dht.Option

	// 判断是否是种子节点
	var peerId = priv.GetPublic().PeerID()
	isBoot := false
	for _, b := range boots {
		if strings.Index(b, peerId.String()) > -1 {
			isBoot = true
		}
	}
	if isBoot {
		dhtOptions = append(dhtOptions, dht.Mode(dht.ModeServer))
	}

	// 创建连接筛选器
	gater := newConnectionGater(nodes)
	dhtOptions = append(dhtOptions, dht.RoutingTableFilter(gater.chainRoutingTableFilter))
	dhtOptions = append(dhtOptions, dht.ProtocolPrefix("/wetee"))

	// 创建连接管理器
	connmgr, err := connmgr.NewConnManager(
		100,                                  // Lowwater
		400,                                  // HighWater,
		connmgr.WithGracePeriod(time.Minute), // 1 minute grace period
	)

	// 创建 P2P 网络主机
	host, err := libp2p.New(
		libp2p.Identity(priv.PrivKey),
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/"+fmt.Sprint(tcp),         // TCP endpoint
			"/ip4/0.0.0.0/udp/"+fmt.Sprint(udp)+"/quic", // UDP endpoint for the QUIC transport
		),
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		libp2p.Security(noise.ID, noise.New),
		libp2p.DefaultTransports,
		libp2p.ConnectionManager(connmgr),
		libp2p.NATPortMap(),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			var err error
			idht, err = dht.New(ctx, h, dhtOptions...)
			return idht, err
		}),
		libp2p.EnableNATService(),
		libp2p.ConnectionGater(gater),
	)

	fmt.Println("Local P2P addr: /ip4/0.0.0.0/tcp/"+fmt.Sprint(tcp)+"/p2p/"+fmt.Sprint(host.ID()), " --- ", priv.GetPublic().SS58())

	// 创建 gossipsub 实例
	pubsubTracer := new(pubsubTracer)
	gossipSub, err := pubsub.NewGossipSub(ctx, host, pubsub.WithEventTracer(pubsubTracer))
	if err != nil {
		return nil, fmt.Errorf("create gossipsub: %w", err)
	}

	// 创建 boot peers
	bootPeers := make(map[peer.ID]peer.AddrInfo)
	for _, b := range boots {
		if b == "" {
			continue
		}
		peerInfo, err := peer.AddrInfoFromString(b)
		if err != nil {
			return nil, fmt.Errorf("decode boot peer: %w", err)
		}
		bootPeers[peerInfo.ID] = *peerInfo
	}

	// 创建 P2P 网络实例
	peer := &Peer{
		Host:      host,
		privKey:   priv.PrivKey,
		idht:      idht,
		pubsub:    gossipSub,
		topics:    make(map[string]*pubsub.Topic),
		bootPeers: bootPeers,
		gater:     gater,
		nodes:     nodes,
		version:   version,
		netHook: func([]kyber.Point) error {
			fmt.Println("network hook not implement")
			return nil
		},
	}

	return peer, nil
}

// Peer P2P 网络实例
type Peer struct {
	host.Host
	privKey     libp2pCrypto.PrivKey
	idht        *dht.IpfsDHT
	pubsub      *pubsub.PubSub
	topics      map[string]*pubsub.Topic
	topicsLock  sync.Mutex
	bootPeers   map[peer.ID]peer.AddrInfo
	gater       *ChainConnectionGater
	reonnecting sync.Map
	nodes       []*types.Node
	netHook     func([]kyber.Point) error
	version     uint32
}

func (p *Peer) PeerStrID() string {
	return p.ID().String()
}

func (p *Peer) NodeIds() []string {
	peers := p.Network().Peers()
	nodes := make([]string, len(peers))
	for i, peer := range peers {
		nodes[i] = peer.String()
	}
	return nodes
}

func (p *Peer) Nodes() []*types.Node {
	return p.nodes
}

func (p *Peer) Version() uint32 {
	return p.version
}

// Send 发送消息
func (p *Peer) Send(ctx context.Context, node *types.Node, pid string, message *types.Message) error {
	var err error
	peerID := node.PeerID()
	protocolID := protocol.ConvertFromStrings([]string{pid})

	fmt.Println(p.Network().Peers())

	util.LogSendmsg(">>>>>> P2P Send()", "to", peerID, "-", node.ID.SS58(), "| type:", pid+"."+message.Type)
	var stream network.Stream
	newStream := func() error {
		stream, err = p.Host.NewStream(ctx, peer.ID(peerID), protocolID...)
		return err
	}

	// 生成消息串流
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 10 * time.Second
	bctx := backoff.WithContext(b, ctx)
	err = backoff.Retry(newStream, bctx)
	if err != nil {
		return fmt.Errorf("Host.NewStream error: %v", err)
	}
	defer stream.Close()

	buf, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	_, err = stream.Write(buf)
	if err != nil {
		return fmt.Errorf("write stream: %w", err)
	}

	return nil
}

// AddHandler 添加消息处理器
func (p *Peer) AddHandler(pid string, handler func(*types.Message) error) {
	streamHandler := genStream(handler)
	p.Host.SetStreamHandler(protocol.ID(pid), streamHandler)
}

// RemoveHandler 移除消息处理器
func (t *Peer) RemoveHandler(pid string) {
	t.Host.RemoveStreamHandler(protocol.ID(pid))
}

func (p *Peer) NetResetHook(hook func([]kyber.Point) error) {
	p.netHook = hook
}

func (p *Peer) Start(ctx context.Context) {
	for _, peer := range p.bootPeers {
		if err := p.Connect(ctx, peer); err != nil {
			fmt.Println("Can't connect to peer:", peer, err)
		} else {
			fmt.Println("Connected to bootstrap node:", peer)
		}
	}

	go func() {
		for {
			nodes, _, version, err := GetChainNodes()
			if err == nil && p.version != version {
				// 重新加载节点
				p.version = version
				p.nodes = nodes
				p.gater.Nodes = nodes

				// 触发网络钩子
				p.netHook([]kyber.Point{})
			}

			p.Discover(ctx)
			fmt.Println("Peer len:", len(p.Network().Peers()))
			time.Sleep(time.Second * 30)
		}
	}()

	subCh, err := p.EventBus().Subscribe(new(event.EvtPeerConnectednessChanged))
	if err != nil {
		fmt.Printf("Error subscribing to peer connectedness changes: %s \n", err)
	}
	defer subCh.Close()

	for {
		select {
		case ev, ok := <-subCh.Out():
			if !ok {
				return
			}

			evt := ev.(event.EvtPeerConnectednessChanged)
			if evt.Connectedness != network.NotConnected {
				continue
			}

			if _, ok := p.bootPeers[evt.Peer]; !ok {
				continue
			}

			paddr := p.bootPeers[evt.Peer]
			go p.reconnectToPeer(ctx, paddr)
		case <-ctx.Done():
			return
		}
	}
}

// Close 关闭 P2P 网络实例
func (p *Peer) Close() error {
	return p.Host.Close()
}

func GetChainNodes() ([]*types.Node, []string, uint32, error) {
	// Get session index
	// version, err := session.GetCurrentIndexLatest(chain.ChainIns.Api.RPC.State)
	// if err != nil {
	// 	fmt.Println("Get session index error:", err)
	// 	return nil, nil, 0, err
	// }

	var version uint32 = 1

	// Get boot peers from chain
	bootPeers, err := chain.ChainIns.GetBootPeers()
	if err != nil {
		fmt.Println("Get node list error:", err)
		return nil, nil, 0, err
	}

	// get node list from chain
	_, _, nodes, err := chain.ChainIns.GetNodes()
	if err != nil {
		fmt.Println("Get node list error:", err)
		return nil, nil, 0, err
	}

	// 计算 p2p 地址
	boots := make([]string, 0, len(bootPeers))
	for _, b := range bootPeers {
		var gopub ed25519.PublicKey = b.Id[:]
		pub, _ := types.PubKeyFromStdPubKey(gopub)
		n := &types.Node{
			ID: *pub,
		}
		d := util.GetUrlFromIp(b.Ip)
		url := d + "/tcp/" + fmt.Sprint(b.Port) + "/p2p/" + n.PeerID().String()
		boots = append(boots, url)
	}

	return nodes, boots, version, nil
}
