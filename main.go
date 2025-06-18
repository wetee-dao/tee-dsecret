package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/privval"
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
	chainPort := util.GetEnvInt("SIDE_CHAIN_PORT", 61001)
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
	p2pKey, err := model.PrivateKeyFromOed25519(nodeKey.PrivKey.Bytes())
	if err != nil {
		fmt.Println("Marshal PKG_PK error:", err)
		os.Exit(1)
	}

	// Init key for DKG
	validatorKey := privval.LoadFilePV(
		"./chain_data/config/priv_validator_key.json",
		"./chain_data/data/priv_validator_state.json",
	)
	dkgKey := validatorKey.Key.PrivKey

	// Init node key for Mainchain
	nodePriv, err := model.PrivateKeyFromOed25519(dkgKey.Bytes())
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

	// Init node
	node, sideChain, dkgReactor, err := sidechain.Init(chainPort, mainChain, func() {
		fmt.Println()
		util.LogWithYellow("Main Chain", chainAddr)
		util.LogWithYellow("Validator Key", nodePriv.GetPublic().SS58())
		util.LogWithYellow("P2P Key", p2pKey.GetPublic().SS58(), " ", p2pKey.GetPublic().SideChainNodeID())
	})
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
	dkgIns, err := dkg.NewDKG(nodePriv, dkgReactor, nil)
	if err != nil {
		fmt.Println("Create DKG error:", err)
		os.Exit(1)
	}
	go dkgIns.Start()
	defer dkgIns.Stop()

	// Set DKG to sideChain
	sideChain.SetDKG(dkgIns)

	// 启动 graphql 服务器
	go graph.StartServer(dkgIns, node, sideChain, gqlPort)

	// wait for stop signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
