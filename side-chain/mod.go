package sidechain

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	cfg "github.com/cometbft/cometbft/config"
	cmtflags "github.com/cometbft/cometbft/libs/cli/flags"
	cmtlog "github.com/cometbft/cometbft/libs/log"
	nm "github.com/cometbft/cometbft/node"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/proxy"
	"github.com/cometbft/cometbft/version"
	"wetee.app/dsecret/chains"
	"wetee.app/dsecret/internal/dkg"
	bftbrigde "wetee.app/dsecret/internal/peer/bft-brigde"
	"wetee.app/dsecret/internal/util"
)

var SideChainNode *nm.Node

func Init(
	chainPort int,
	mainChain chains.MainChain,
	callback func(),
) (*nm.Node, *SideChain, *bftbrigde.BTFReactor, error) {
	// Get boot peers
	boots, err := mainChain.GetBootPeers()
	if err != nil {
		return nil, nil, nil, errors.New("NewSideChain error: " + err.Error())
	}

	// 创建侧链实例
	sideChain, err := NewSideChain()
	if err != nil {
		return nil, nil, nil, errors.New("NewSideChain error: " + err.Error())
	}

	addr := util.GetEnv("SIDE_CHAIN_ADDR", "0.0.0.0")
	p2pConf := cfg.DefaultP2PConfig()
	p2pConf.ListenAddress = "tcp://" + addr + ":" + fmt.Sprint(chainPort)
	p2pConf.AllowDuplicateIP = true
	p2pConf.Seeds = ""

	consensusConf := cfg.DefaultConsensusConfig()
	consensusConf.CreateEmptyBlocks = false

	rpcConf := cfg.DefaultRPCConfig()
	// rpcConf.ListenAddress = "tcp://0.0.0.0:" + fmt.Sprint(chainPort+1)
	rpcConf.ListenAddress = ""

	// init BFT node config
	config := &cfg.Config{
		BaseConfig: cfg.BaseConfig{
			Version:            version.CMTSemVer,
			Genesis:            "config/genesis.json",
			PrivValidatorKey:   "config/priv_validator_key.json",
			PrivValidatorState: "data/priv_validator_state.json",
			NodeKey:            "config/node_key.json",
			Moniker:            "WeTEE Chain",
			ProxyApp:           "tcp://127.0.0.1:" + fmt.Sprint(chainPort+1),
			ABCI:               "socket",
			LogLevel:           "error",
			LogFormat:          cfg.LogFormatPlain,
			FilterPeers:        false,
			DBBackend:          "pebbledb",
			DBPath:             "BFT",
		},
		RPC:             rpcConf,
		GRPC:            cfg.DefaultGRPCConfig(),
		P2P:             p2pConf,
		Mempool:         cfg.DefaultMempoolConfig(),
		StateSync:       cfg.DefaultStateSyncConfig(),
		BlockSync:       cfg.DefaultBlockSyncConfig(),
		Consensus:       consensusConf,
		Storage:         cfg.DefaultStorageConfig(),
		TxIndex:         cfg.DefaultTxIndexConfig(),
		Instrumentation: cfg.DefaultInstrumentationConfig(),
	}
	config.SetRoot("./chain_data")

	// init sidechain node key
	nodeKey, err := p2p.LoadNodeKey(config.NodeKeyFile())
	if err != nil {
		return nil, nil, nil, errors.New("failed to load node key: " + err.Error())
	}

	// add boot nodes

	seeds := []string{}
	for _, boot := range boots {
		if util.ToSideChainNodeID(boot.Id[:]) == nodeKey.ID() {
			continue
		}
		seeds = append(seeds, boot.SideChainUrl())
	}
	config.P2P.Seeds = strings.Join(seeds, ",")

	// load validator key
	pv := privval.LoadFilePV(
		config.PrivValidatorKeyFile(),
		config.PrivValidatorStateFile(),
	)

	// init logger
	logger := cmtlog.NewTMLogger(cmtlog.NewSyncWriter(os.Stdout))
	logger, err = cmtflags.ParseLogLevel(config.LogLevel, logger, cfg.DefaultLogLevel)
	if err != nil {
		return nil, nil, nil, errors.New("init logger error: " + err.Error())
	}

	// add DKG to chain node
	dkgReactor := bftbrigde.NewBTFReactor("DKG", mainChain)

	// init BFT node
	SideChainNode, err = nm.NewNode(
		context.Background(),
		config,
		pv,
		nodeKey,
		proxy.NewLocalClientCreator(sideChain),
		nm.DefaultGenesisDocProviderFunc(config),
		cfg.DefaultDBProvider,
		nm.DefaultMetricsProvider(config.Instrumentation),
		logger,
		nm.CustomReactors(map[string]p2p.Reactor{
			"DKG": dkgReactor,
		}),
	)
	if err != nil {
		return nil, nil, nil, errors.New("init BFT node error: " + err.Error())
	}

	callback()
	if p2pConf.Seeds != "" {
		util.LogWithRed("Boot Nodes ", p2pConf.Seeds)
	}

	return SideChainNode, sideChain, dkgReactor, err
}

func (s *SideChain) SetDKG(dkg *dkg.DKG) {
	s.dkg = dkg
}
