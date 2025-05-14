package main

import (
	"context"
	"fmt"
	"os"

	"wetee.app/dsecret/graph"
	"wetee.app/dsecret/internal/chain"
	"wetee.app/dsecret/internal/dkg"
	"wetee.app/dsecret/internal/peer/p2p"
	"wetee.app/dsecret/internal/store"
	types "wetee.app/dsecret/type"
	"wetee.app/dsecret/util"
)

var DefaultChainUrl string = "ws://wetee-node.worker-addon.svc.cluster.local:9944"

func main() {
	// 获取环境变量
	peerSecret := util.GetEnv("PEER_PK", "")
	tcpPort := util.GetEnvInt("TCP_PORT", 61000)
	udpPort := util.GetEnvInt("UDP_PORT", 61000)
	chainAddr := util.GetEnv("CHAIN_ADDR", DefaultChainUrl)
	password := util.GetEnv("PASSWORD", "")

	// 初始化数据库
	err := store.InitDB(password)
	if err != nil {
		fmt.Println("Init db error:", err)
		os.Exit(1)
	}

	// 初始化加密套件
	nodeSecret, err := types.PrivateKeyFromLibp2pHex(peerSecret)
	if err != nil {
		fmt.Println("Marshal PKG_PK error:", err)
		os.Exit(1)
	}

	// 链接区块链
	err = chain.InitChain(chainAddr, nodeSecret)
	if err != nil {
		fmt.Println("Connect to chain error:", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动 P2P 网络
	peer, err := p2p.NewP2PNetwork(ctx, nodeSecret, uint32(tcpPort), uint32(udpPort))
	if err != nil {
		fmt.Println("Start P2P peer error:", err)
		os.Exit(1)
	}

	// 创建 DKG 实例
	dkgIns, err := dkg.NewRabinDKG(nodeSecret, peer)
	if err != nil {
		fmt.Println("Create DKG error:", err)
		os.Exit(1)
	}

	// 启动节点
	go peer.Start(ctx)

	// 运行 DKG 协议
	if err := dkgIns.Start(ctx, nil); err != nil {
		fmt.Println("Start DKG error:", err)
		os.Exit(1)
	}

	graph.StartServer(dkgIns)
}
