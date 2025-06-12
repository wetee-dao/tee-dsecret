package sidechain

import (
	"context"
	"errors"
	"os"

	cfg "github.com/cometbft/cometbft/config"
	cmtflags "github.com/cometbft/cometbft/libs/cli/flags"
	cmtlog "github.com/cometbft/cometbft/libs/log"
	nm "github.com/cometbft/cometbft/node"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/proxy"
	"github.com/cometbft/cometbft/version"
)

func InitNode() (*nm.Node, error) {
	// 创建侧链实例
	sideChain, err := NewSideChain()
	if err != nil {
		return nil, errors.New("NewSideChain error: " + err.Error())
	}

	// init BFT node config
	config := &cfg.Config{
		BaseConfig: cfg.BaseConfig{
			Version:            version.CMTSemVer,
			Genesis:            "config/genesis.json",
			PrivValidatorKey:   "config/priv_validator_key.json",
			PrivValidatorState: "data/priv_validator_state.json",
			NodeKey:            "config/node_key.json",
			Moniker:            "WeTEE Chain",
			ProxyApp:           "tcp://127.0.0.1:26658",
			ABCI:               "socket",
			LogLevel:           cfg.DefaultLogLevel,
			LogFormat:          cfg.LogFormatPlain,
			FilterPeers:        false,
			DBBackend:          "pebbledb",
			DBPath:             "BFT",
		},
		RPC:             cfg.DefaultRPCConfig(),
		GRPC:            cfg.DefaultGRPCConfig(),
		P2P:             cfg.DefaultP2PConfig(),
		Mempool:         cfg.DefaultMempoolConfig(),
		StateSync:       cfg.DefaultStateSyncConfig(),
		BlockSync:       cfg.DefaultBlockSyncConfig(),
		Consensus:       cfg.DefaultConsensusConfig(),
		Storage:         cfg.DefaultStorageConfig(),
		TxIndex:         cfg.DefaultTxIndexConfig(),
		Instrumentation: cfg.DefaultInstrumentationConfig(),
	}
	config.SetRoot("./chain_data")

	// init sidechain node key
	nodeKey, err := p2p.LoadNodeKey(config.NodeKeyFile())
	if err != nil {
		return nil, errors.New("failed to load node key: " + err.Error())
	}

	// load validator key
	pv := privval.LoadFilePV(
		config.PrivValidatorKeyFile(),
		config.PrivValidatorStateFile(),
	)

	// init logger
	logger := cmtlog.NewTMLogger(cmtlog.NewSyncWriter(os.Stdout))
	logger, err = cmtflags.ParseLogLevel(config.LogLevel, logger, cfg.DefaultLogLevel)
	if err != nil {
		return nil, errors.New("init logger error: " + err.Error())
	}

	// init BFT node
	return nm.NewNode(
		context.Background(),
		config,
		pv,
		nodeKey,
		proxy.NewLocalClientCreator(sideChain),
		nm.DefaultGenesisDocProviderFunc(config),
		cfg.DefaultDBProvider,
		nm.DefaultMetricsProvider(config.Instrumentation),
		logger,
	)
}
