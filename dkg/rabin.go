package dkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/share"
	rabin "go.dedis.ch/kyber/v3/share/dkg/rabin"
	"go.dedis.ch/kyber/v3/suites"
	p2peer "wetee.app/dsecret/peer"
	"wetee.app/dsecret/store"
	"wetee.app/dsecret/types"
)

// DKG 代表 Rabin DKG 协议的实例。
type DKG struct {
	mu sync.Mutex
	// Host 是 P2P 网络主机。
	Peer *p2peer.Peer
	// Suite 是加密套件。
	Suite suites.Suite
	// NodeSecret 是长期的私钥。
	NodeSecret kyber.Scalar
	// Participants 是参与者的公钥列表。
	Participants []kyber.Point
	// Peer 是 P2P 网络主机。
	Nodes []*types.Node
	// Threshold 是密钥重建所需的最小份额数量。
	Threshold int
	// Shares 是当前节点持有的密钥份额。
	Shares map[peer.ID]*share.PriShare
	// DistKeyGenerator
	DistKeyGenerator *rabin.DistKeyGenerator
	// DistPubKey globle public key
	DkgPubKey kyber.Point // DKG group Public key
	// DistKeyShare is the node private share
	DkgKeyShare types.DistKeyShare // DKG node private share
	// preRecerve is the channel to receive SendEncryptedSecretRequest
	preRecerve map[string]chan *share.PubShare
}

// NewRabinDKG 创建一个新的 Rabin DKG 实例。
func NewRabinDKG(NodeSecret *types.PrivKey, nodes []*types.Node, threshold int, p *p2peer.Peer) (*DKG, error) {
	// 检查参数。
	if len(nodes) < threshold {
		return nil, errors.New("阈值必须小于参与者数量")
	}

	suite := suites.MustFind("Ed25519")
	// 获取节点公钥列表。
	participants := make([]kyber.Point, len(nodes))
	for i, n := range nodes {
		pk, err := types.PublicKeyFromHex(n.ID)
		if err != nil {
			fmt.Println("解析 PKG_PUBS 失败:", err)
			os.Exit(1)
		}
		_, err = peer.IDFromPublicKey(pk)
		if err != nil {
			fmt.Println("IDFromPublicKey 失败:", err)
			os.Exit(1)
		}
		participants[i] = pk.Point()
	}

	// 创建 DKG 对象。
	dkg := &DKG{
		Suite:        suite,
		Participants: participants,
		Threshold:    threshold,
		Shares:       make(map[peer.ID]*share.PriShare),
		NodeSecret:   NodeSecret.Scalar(),
		Peer:         p,
		Nodes:        nodes,
		preRecerve:   make(map[string]chan *share.PubShare),
	}

	// 存储 DKG 对象
	// dkg.restore()
	return dkg, nil
}

// Start 启动 Rabin DKG 协议。
func (dkg *DKG) Start(ctx context.Context) error {
	var err error

	// initialize vss dealer
	dkg.DistKeyGenerator, err = rabin.NewDistKeyGenerator(dkg.Suite, dkg.NodeSecret, dkg.Participants, dkg.Threshold)
	if err != nil {
		return fmt.Errorf("初始化 VSS 协议失败: %w", err)
	}

	// 获取当前节点的
	deals, err := dkg.DistKeyGenerator.Deals()
	if err != nil {
		return fmt.Errorf("生成密钥份额失败: %w", err)
	}

	// Add 请求回调 handler
	dkg.Peer.AddHandler("dkg", dkg.HandleMessage)

	for {
		err = dkg.Peer.Discover(ctx)
		if err != nil {
			fmt.Println("Discover error:", err)
		}

		if len(dkg.Peer.Network().Peers())+1 <= dkg.Threshold {
			time.Sleep(time.Second * 10)
			fmt.Println("")
		} else {
			break
		}
	}

	for i, deal := range deals {
		err = dkg.SendDealMessage(ctx, dkg.Nodes[i], deal)
		if err != nil {
			fmt.Println("Send error:", err)
		}
	}

	return nil
}

func (d *DKG) Share() types.DistKeyShare {
	return d.DkgKeyShare
}

func (dkg *DKG) ID() int {
	pub := dkg.Suite.Point().Mul(dkg.NodeSecret, nil)
	var index int = -1
	for i, p := range dkg.Participants {
		if p.Equal(pub) {
			index = i
			break
		}
	}
	return index
}

func (dkg *DKG) restore() error {
	v, err := store.GetKey("G", "dkg")
	if err != nil {
		return fmt.Errorf("get dkg: %w", err)
	}
	if v != nil {
		err = json.Unmarshal(v, dkg)
		if err != nil {
			return fmt.Errorf("unmarshal dkg: %w", err)
		}
	}
	return nil
}

func (dkg *DKG) store() error {
	payload, err := json.Marshal(dkg)
	if err != nil {
		return fmt.Errorf("marshal dkg: %w", err)
	}
	return store.SetKey("G", "dkg", payload)
}
