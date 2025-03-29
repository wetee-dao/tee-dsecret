package dkg

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"go.dedis.ch/kyber/v4"
	"go.dedis.ch/kyber/v4/share"
	pedersen "go.dedis.ch/kyber/v4/share/dkg/pedersen"
	"go.dedis.ch/kyber/v4/sign/schnorr"
	"go.dedis.ch/kyber/v4/suites"
	p2peer "wetee.app/dsecret/peer"
	"wetee.app/dsecret/store"
	types "wetee.app/dsecret/type"
)

// DKG 代表 Rabin DKG 协议的实例
type DKG struct {
	// 操作互斥锁
	mu sync.RWMutex
	// Host 是 P2P 网络主机
	Peer p2peer.Peer
	// Suite 是加密套件
	Suite suites.Suite
	// NodeSecret 是长期的私钥
	NodeSecret kyber.Scalar
	// Signer 是用于签名的私钥
	Signer *types.PrivKey
	// AllNodes 是所有节点的集合
	AllNodes []*types.Node
	// Peer 是 P2P 网络主机
	DkgNodes []*types.Node
	// Threshold 是密钥重建所需的最小份额数量
	Threshold int

	// DistKeyGenerator
	DistKeyGenerator *pedersen.DistKeyGenerator
	// DistPubKey globle public key
	DkgPubKey kyber.Point

	// cache the deal, response, justification, result
	deals     map[string]*pedersen.DealBundle
	responses map[string]*pedersen.ResponseBundle
	justifs   []*pedersen.JustificationBundle
	results   *pedersen.Result

	// Shares 是当前节点持有的密钥份额
	Shares map[peer.ID]*share.PriShare
	// DistKeyShare is the node private share
	DkgKeyShare types.DistKeyShare

	// preRecerve is the channel to receive SendEncryptedSecretRequest
	preRecerve map[string]chan any
	// 未初始化状态 => 0 | 初始化成功 => 1
	status uint8
}

// NewRabinDKG 创建一个新的 Rabin DKG 实例
func NewRabinDKG(NodeSecret *types.PrivKey, p p2peer.Peer) (*DKG, error) {
	nodes := p.Nodes()

	// 获取节点公钥列表
	dkgNodes := make([]*types.Node, 0, len(nodes))
	for _, n := range nodes {
		// 过滤不是dkg节点
		if n.Type != 1 {
			continue
		}

		dkgNodes = append(dkgNodes, n)
	}

	// 获取密钥重建所需的最小份额数量
	threshold := len(dkgNodes) * 2 / 3

	// 创建 DKG 对象
	dkg := &DKG{
		Suite:      suites.MustFind("Ed25519"),
		NodeSecret: NodeSecret.Scalar(),
		Signer:     NodeSecret,
		Threshold:  threshold,
		Shares:     make(map[peer.ID]*share.PriShare),
		Peer:       p,
		AllNodes:   nodes,
		DkgNodes:   dkgNodes,
		preRecerve: make(map[string]chan any),
		deals:      make(map[string]*pedersen.DealBundle),
		responses:  make(map[string]*pedersen.ResponseBundle),
	}

	// 添加网络节点变化回调
	p.NetResetHook(dkg.ReShare)

	// 复原 DKG 对象
	dkg.reStore()

	return dkg, nil
}

// Start 启动 Rabin DKG 协议
func (dkg *DKG) Start(ctx context.Context, log pedersen.Logger) error {
	// Add 请求回调 handler
	dkg.Peer.AddHandler("dkg", dkg.HandleDkg)
	dkg.Peer.AddHandler("worker", dkg.HandleWorker)

	if flag.Lookup("test.v") == nil {
		go dkg.HandleSecretSave(ctx)
	}

	// 如果已经初始化，则直接返回
	if dkg.status == 1 {
		return nil
	}

	// dkg 节点列表
	nodes := make([]pedersen.Node, 0, len(dkg.DkgNodes))
	for i, p := range dkg.DkgNodes {
		nodes = append(nodes, pedersen.Node{
			Index:  uint32(i),
			Public: p.ID.Point(),
		})
	}

	// 初始化协议配置
	conf := pedersen.Config{
		Suite:     dkg.Suite,
		NewNodes:  nodes,
		Threshold: dkg.Threshold,
		Auth:      schnorr.NewScheme(dkg.Suite),
		FastSync:  true,
		Longterm:  dkg.NodeSecret,
		Nonce:     Version2Nonce(dkg.Peer.Version()),
		Log:       log,
	}

	// initialize dealer
	var err error
	dkg.DistKeyGenerator, err = pedersen.NewDistKeyHandler(&conf)
	if err != nil {
		return fmt.Errorf("Failed to initialize DKG protocol: %w", err)
	}

	// 等待节点连接
	for {
		if dkg.nodeLen()+1 < len(dkg.DkgNodes) {
			time.Sleep(time.Second * 10)
			fmt.Println("Number of nodes:", dkg.nodeLen(), " len(dkg.DkgNodes) ", len(dkg.DkgNodes))
			fmt.Println("The number of nodes is insufficient, please wait for more nodes to join")
		} else {
			break
		}
	}

	// 获取当前节点的协议
	deal, err := dkg.DistKeyGenerator.Deals()
	if err != nil {
		return fmt.Errorf("Failed to generate key shares: %w", err)
	}

	// 开启节点共识
	for _, node := range dkg.DkgNodes {
		err = dkg.SendDealMessage(ctx, node, deal)
		if err != nil {
			fmt.Println("Send error:", err)
		}
	}

	// 等待节点完成重组
	for {
		if dkg.DkgKeyShare.PriShare != nil {
			break
		}
		fmt.Println("Wait for the DKG network to complete reorganization")
		time.Sleep(time.Second * 5)
	}
	fmt.Println("The DKG protocol has been successfully initiated")

	return nil
}

func (dkg *DKG) ReShare() error {
	fmt.Println("ReShare")
	return nil
}

func (dkg *DKG) nodeLen() int {
	var len int
	peers := dkg.Peer.NodeIds()
	for _, p := range peers {
		for _, node := range dkg.DkgNodes {
			if p == node.PeerID().String() {
				len = len + 1
			}
		}
	}
	return len
}

func (d *DKG) Share() types.DistKeyShare {
	return d.DkgKeyShare
}

func (dkg *DKG) ID() int {
	pub := dkg.Suite.Point().Mul(dkg.NodeSecret, nil)
	var index int = -1
	for i, p := range dkg.DkgNodes {
		if p.ID.Point().String() == pub.String() {
			index = i
			break
		}
	}
	return index
}

func (dkg *DKG) reStore() error {
	v, err := store.GetKey("G", "dkg-"+dkg.Signer.GetPublic().SS58())
	if err != nil {
		return fmt.Errorf("get dkg: %w", err)
	}

	d, err := types.DistKeyShareFromProtocol(dkg.Suite, v)
	if err != nil {
		return fmt.Errorf("unmarshal dkg: %w", err)
	}
	dkg.DkgKeyShare = *d
	dkg.status = 1

	return nil
}

func (dkg *DKG) saveStore() error {
	payload, err := types.DistKeyShareToProtocol(&dkg.DkgKeyShare)
	if err != nil {
		return fmt.Errorf("marshal dkg: %w", err)
	}

	return store.SetKey("G", "dkg-"+dkg.Signer.GetPublic().SS58(), payload)
}

func (dkg *DKG) SendToNode(ctx context.Context, node *types.Node, pid string, message *types.Message) error {
	if node == nil {
		fmt.Println("node is nil")
		return errors.New("node is nil")
	}

	message.OrgId = dkg.Peer.PeerStrID()
	if dkg.Peer.PeerStrID() == node.PeerID().String() {
		switch pid {
		case "dkg":
			go dkg.HandleDkg(message)
			return nil
		case "worker":
			go dkg.HandleWorker(message)
			return nil
		default:
			return errors.New("invalid pid")
		}
	}
	return dkg.Peer.Send(ctx, node, pid, message)
}

func (dkg *DKG) GetNode(nodeId string) *types.Node {
	for _, node := range dkg.AllNodes {
		if node.PeerID().String() == nodeId {
			return node
		}
	}

	return nil
}

func Version2Nonce(v uint32) []byte {
	var nonce [pedersen.NonceLength]byte
	var version = fmt.Append(nil, v)

	copy(nonce[:], version)

	return nonce[:]
}
