package peer

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
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
func NewP2PNetwork(ctx context.Context, peerSecret string, boots []string, tcp, udp uint32) (*Peer, error) {
	var idht *dht.IpfsDHT
	var dhtOptions []dht.Option
	if len(boots) == 0 {
		fmt.Println("Host running as a bootsrap node")
		dhtOptions = append(dhtOptions, dht.Mode(dht.ModeServer))
	}
	dhtOptions = append(dhtOptions, dht.ProtocolPrefix("/wetee"))

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
			idht, err = dht.New(ctx, h, dhtOptions...)
			return idht, err
		}),
		libp2p.EnableNATService(),
	)

	fmt.Println("Local P2P addr: /ip4/0.0.0.0/tcp/" + fmt.Sprint(tcp) + "/p2p/" + fmt.Sprint(host.ID()))

	// 创建 gossipsub 实例
	pubsubTracer := new(pubsubTracer)
	gossipSub, err := pubsub.NewGossipSub(ctx, host, pubsub.WithEventTracer(pubsubTracer))
	if err != nil {
		return nil, fmt.Errorf("create gossipsub: %w", err)
	}

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
	privKey     crypto.PrivKey
	idht        *dht.IpfsDHT
	pubsub      *pubsub.PubSub
	topics      map[string]*pubsub.Topic
	topicsLock  sync.Mutex
	bootPeers   map[peer.ID]peer.AddrInfo
	reonnecting sync.Map
}

func (p *Peer) Send(ctx context.Context, node *types.Node, pid string, message *types.Message) error {
	var err error
	peerID := node.PeerID()
	protocolID := protocol.ConvertFromStrings([]string{pid})

	fmt.Printf("transport.Send(): peerID: %s , ProtocolID: %v \n", peerID, protocolID)
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

func (p *Peer) AddHandler(pid string, handler func(*types.Message) error) {
	streamHandler := genStream(handler)
	p.Host.SetStreamHandler(protocol.ID(pid), streamHandler)
}

func (t *Peer) RemoveHandler(pid protocol.ID) {
	t.Host.RemoveStreamHandler(pid)
}

func (p *Peer) Close() error {
	return p.Host.Close()
}

func genStream(handler func(*types.Message) error) func(network.Stream) {
	return func(stream network.Stream) {
		fmt.Printf("New call from %s", stream.Conn().RemotePeer())
		buf, err := io.ReadAll(stream)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("read stream: %s", err)
			}

			err = stream.Reset()
			if err != nil {
				fmt.Printf("reset stream: %s", err)
			}

			return
		}

		err = stream.Close()
		if err != nil {
			fmt.Printf("close stream: %s", err)
			return
		}

		data := &types.Message{}
		err = json.Unmarshal(buf, data)
		if err != nil {
			fmt.Printf("unmarshal data: %s", err)
			return
		}

		fmt.Printf("received message:  type: %s", data.Type)
		err = handler(data)
		if err != nil {
			fmt.Printf("handle data: %s", err)
			return
		}
	}
}
