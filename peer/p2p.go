package peer

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	"wetee.app/dsecret/types"
)

// NewP2PNetwork 创建一个新的 P2P 网络实例。
func NewP2PNetwork(ctx context.Context, peerSecret string, bootPeers []string, tcp, udp uint32) (*Peer, error) {
	var idht *dht.IpfsDHT

	// 解码私钥HEX
	buf, err := hex.DecodeString(peerSecret)
	if err != nil {
		return nil, fmt.Errorf("decode peer secret: %w", err)
	}

	// 解码私钥
	priv, err := crypto.UnmarshalPrivateKey(buf)
	if err != nil {
		return nil, fmt.Errorf("decode private key: %w", err)
	}

	// 创建连接管理器
	connmgr, err := connmgr.NewConnManager(
		100,                                  // Lowwater
		400,                                  // HighWater,
		connmgr.WithGracePeriod(time.Minute), // 1 minute grace period
	)

	// 创建 P2P 网络主机。
	host, err := libp2p.New(
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/"+fmt.Sprint(tcp), // regular tcp connections
			// "/ip4/0.0.0.0/udp/"+fmt.Sprint(udp)+"/quic", // a UDP endpoint for the QUIC transport
		),
		// support TLS connections
		// libp2p.Security(libp2ptls.ID, libp2ptls.New),
		libp2p.Security(noise.ID, noise.New),
		libp2p.DefaultTransports,
		libp2p.ConnectionManager(connmgr),
		libp2p.NATPortMap(),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			var err error
			idht, err = dht.New(context.Background(), h)
			return idht, err
		}),
		libp2p.EnableNATService(),
	)

	fmt.Println("P2P Addr: /ip4/0.0.0.0/tcp/" + fmt.Sprint(tcp) + "/p2p/" + fmt.Sprint(host.ID()))

	// 创建 gossipsub 实例
	pubsubTracer := new(pubsubTracer)
	gossipSub, err := pubsub.NewGossipSub(ctx, host, pubsub.WithEventTracer(pubsubTracer))
	if err != nil {
		return nil, fmt.Errorf("create gossipsub: %w", err)
	}

	peer := &Peer{
		Host:      host,
		privKey:   priv,
		idht:      idht,
		pubsub:    gossipSub,
		topics:    make(map[string]*pubsub.Topic),
		bootPeers: bootPeers,
	}

	return peer, nil
}

type Peer struct {
	host.Host
	privKey    crypto.PrivKey
	idht       *dht.IpfsDHT
	pubsub     *pubsub.PubSub
	topics     map[string]*pubsub.Topic
	topicsLock sync.Mutex
	bootPeers  []string
}

func (p *Peer) Start(ctx context.Context) {
	var wg sync.WaitGroup

	for _, peerAddr := range p.bootPeers {
		if peerAddr == "" {
			continue
		}
		pi, err := peer.AddrInfoFromString(peerAddr)
		if err != nil {
			fmt.Println("Can't parse peer addr info string: ", pi, err)
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := p.Connect(ctx, *pi); err != nil {
				fmt.Println("Can't connect to peer: ", pi, err)
			} else {
				fmt.Println("Connected to bootstrap node: ", pi)
			}
		}()
	}

	wg.Wait()
}

func (p *Peer) Receive(ctx context.Context, topic string) (*pubsub.Subscription, error) {

	t, err := p.join(topic)
	if err != nil {
		return nil, fmt.Errorf("join topic: %w", err)
	}

	sub, err := t.Subscribe()
	if err != nil {
		return nil, fmt.Errorf("subscribe topic: %w", err)
	}

	return sub, nil
}

func (p *Peer) Send(ctx context.Context, node *types.Node, message *types.Message) error {
	var err error
	peerID := node.PeerID()
	protocolID := protocol.ConvertFromStrings([]string{message.Type})

	fmt.Printf("transport.Send(): peerID:%s, ProtocolID:%v", peerID, protocolID)
	var stream network.Stream
	newStream := func() error {
		stream, err = p.Host.NewStream(ctx, peerID, protocolID...)
		return err
	}
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 10 * time.Second
	bctx := backoff.WithContext(b, ctx)

	err = backoff.Retry(newStream, bctx)
	if err != nil {
		return fmt.Errorf("new stream: %v", err)
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

func (p *Peer) Broadcast(ctx context.Context, topic string, data []byte) error {

	t, err := p.join(topic)
	if err != nil {
		return fmt.Errorf("join topic: %w", err)
	}

	err = t.Publish(ctx, data)
	if err != nil {
		return fmt.Errorf("publish topic: %w", err)
	}

	return nil
}

func (p *Peer) join(topic string) (*pubsub.Topic, error) {

	p.topicsLock.Lock()
	defer p.topicsLock.Unlock()

	t, exists := p.topics[topic]
	if exists {
		return t, nil
	}

	t, err := p.pubsub.Join(topic)
	if err != nil {
		return nil, fmt.Errorf("join topic: %w", err)
	}

	p.topics[topic] = t

	return t, nil
}

func (p *Peer) Close() error {
	return p.Host.Close()
}
