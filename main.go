package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cometbft/cometbft/p2p"
	chain "wetee.app/dsecret/chains"
	"wetee.app/dsecret/graph"
	"wetee.app/dsecret/internal/dkg"
	"wetee.app/dsecret/internal/model"
	"wetee.app/dsecret/internal/util"
	sidechain "wetee.app/dsecret/side-chain"
)

var DefaultChainUrl string = "ws://wetee-node.worker-addon.svc.cluster.local:9944"

func main() {
	// 获取环境变量
	gqlPort := util.GetEnvInt("GQL_PORT", 61000)
	chainAddr := util.GetEnv("CHAIN_ADDR", DefaultChainUrl)
	chainPort := util.GetEnvInt("CHAIN_PORT", 61001)
	// password := util.GetEnv("PASSWORD", "")

	// Init app db
	db, err := model.NewDB()
	if err != nil {
		fmt.Println("Init db error:", err)
		os.Exit(1)
	}
	defer db.Close()

	// init sidechain node key
	nodeKey, err := p2p.LoadNodeKey("./chain_data/config/node_key.json")
	if err != nil {
		fmt.Println("failed to load node key:", err)
		os.Exit(1)
	}

	// Init node key for Mainchain
	nodePriv, err := model.PrivateKeyFromOed25519(nodeKey.PrivKey.Bytes())
	if err != nil {
		fmt.Println("Marshal PKG_PK error:", err)
		os.Exit(1)
	}

	// Link to polkadot
	mainChain, err := chain.ConnectMainChain(chainAddr, nodePriv)
	if err != nil {
		fmt.Println("Connect to chain error:", err)
		os.Exit(1)
	}

	// Get boot peers
	boots, err := mainChain.GetBootPeers()
	if err != nil {
		fmt.Println("Get boot peers error:", err)
		os.Exit(1)
	}

	// Init node
	node, _, dkgReactor, err := sidechain.InitNode(chainPort, boots)
	if err != nil {
		log.Fatalf("failed to init node: %v", err)
	}

	// Start BFT node
	if err := node.Start(); err != nil {
		log.Fatalf("failed to start BFT node: %v", err)
	}
	defer func() {
		_ = node.Stop()
		node.Wait()
	}()

	// Create DKG
	dkgIns, err := dkg.NewDKG(nodePriv, dkgReactor)
	if err != nil {
		fmt.Println("Create DKG error:", err)
		os.Exit(1)
	}

	// // 运行 DKG 协议
	// if err := dkgIns.Start(ctx, nil); err != nil {
	// 	fmt.Println("Start DKG error:", err)
	// 	os.Exit(1)
	// }

	// 启动 graphql 服务器
	go graph.StartServer(dkgIns, gqlPort)

	// wait for stop signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
