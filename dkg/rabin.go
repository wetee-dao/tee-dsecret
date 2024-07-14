package dkg

import (
	"context"
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
	p2p "wetee.app/dsecret/peer"
	"wetee.app/dsecret/types"
	"wetee.app/dsecret/util"
)

// DKG 代表 Rabin DKG 协议的实例。
type DKG struct {
	mu sync.Mutex
	// Host 是 P2P 网络主机。
	Peer *p2p.Peer
	// Suite 是加密套件。
	Suite suites.Suite
	// NodeSecret 是长期的私钥。
	NodeSecret kyber.Scalar
	// rabin dkg internal private polynomial (f)
	FPoly *share.PriPoly
	// rabin dkg internal private polynimial (g)
	GPoly *share.PriPoly
	// Participants 是参与者的公钥列表。
	Participants []kyber.Point
	// Peer 是 P2P 网络主机。
	Peers []*types.Node
	// Threshold 是密钥重建所需的最小份额数量。
	Threshold int
	// Shares 是当前节点持有的密钥份额。
	Shares map[peer.ID]*share.PriShare
	// DistKeyGenerator
	DistKeyGenerator *rabin.DistKeyGenerator
}

// NewRabinDKG 创建一个新的 Rabin DKG 实例。
func NewRabinDKG(suite suites.Suite, NodeSecret *types.PrivKey, nodes []*types.Node, threshold int) (*DKG, error) {
	// 检查参数。
	if len(nodes) < threshold {
		return nil, errors.New("阈值必须小于参与者数量")
	}

	participants := make([]kyber.Point, len(nodes))
	for i, n := range nodes {
		pk, err := types.PublicKeyFromHex("08011220" + n.ID)
		if err != nil {
			fmt.Println("解析 PKG_PUBS 失败:", err)
			os.Exit(1)
		}
		peerID, err := peer.IDFromPublicKey(pk)
		fmt.Println("peerID:", peerID)
		participants[i] = pk.Point()
	}

	// 创建 DKG 对象。
	dkg := &DKG{
		Suite:        suite,
		Participants: participants,
		Threshold:    threshold,
		Shares:       make(map[peer.ID]*share.PriShare),
		NodeSecret:   NodeSecret.Scalar(),
		Peers:        nodes,
	}

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
	for index, deal := range deals {
		fmt.Printf("参与者 %d 的 Deal: %+v\n", index, deal)
	}

	time.Sleep(15 * time.Second)
	for i, deal := range deals {
		dkg.SendDealMessage(ctx, dkg.Peers[i], deal)
	}

	// Add 请求回调 handler
	dkg.Peer.AddHandler("deal", dkg.HandleMessage)

	// Listen to deal
	deal, err := dkg.Peer.Sub(ctx, "deal")
	if err != nil {
		fmt.Println("peer.Receive error:", err)
		os.Exit(1)
	}
	// Listen to deal response
	response, err := dkg.Peer.Sub(ctx, "response")
	if err != nil {
		fmt.Println("peer.Receive error:", err)
		os.Exit(1)
	}
	// Listen to secret commits
	secretCommits, err := dkg.Peer.Sub(ctx, "secret_commits")
	if err != nil {
		fmt.Println("peer.Receive error:", err)
		os.Exit(1)
	}

	go func() {
		for {
			msg, err := deal.Next(context.Background())
			if err != nil {
				fmt.Println("接收消息失败:", err)
				continue
			}
			err = dkg.HandleDeal(msg.Data)
			if err != nil {
				util.LogWithRed("HandleDeal error:", err)
			}
		}
	}()

	go func() {
		for {
			msg, err := response.Next(context.Background())
			if err != nil {
				fmt.Println("接收消息失败:", err)
				continue
			}
			err = dkg.HandleDealResp(msg.Data)
			if err != nil {
				util.LogWithRed("HandleDealResp error:", err)
			}
		}
	}()

	for {
		msg, err := secretCommits.Next(context.Background())
		if err != nil {
			fmt.Println("接收消息失败:", err)
			continue
		}
		err = dkg.HandleSecretCommits(msg.Data)
		if err != nil {
			util.LogWithRed("HandleDeal error:", err)
		}
	}

	return nil
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
