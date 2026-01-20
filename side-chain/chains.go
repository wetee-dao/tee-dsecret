package sidechain

import (
	"errors"
	"fmt"

	"github.com/cockroachdb/pebble"
	"github.com/wetee-dao/tee-dsecret/pkg/chains"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

type ChainState struct {
	Chains map[uint32]*model.ChainConfig
	chains map[uint32]*chains.ChainApi
}

var chainsKey = "chains"

func (s *SideChain) loadChains() error {
	// 如果 chains map 未初始化，先初始化
	if s.chains == nil {
		s.chains = make(map[uint32]*chains.ChainApi)
	}

	// 从数据库加载 ChainConfig map
	chainConfigs, err := model.GetJson[map[uint32]*model.ChainConfig]("", chainsKey)
	if err != nil && !errors.Is(err, pebble.ErrNotFound) {
		return fmt.Errorf("failed to load chains from database: %w", err)
	}

	// 如果没有配置，直接返回
	if chainConfigs == nil || len(*chainConfigs) == 0 {
		return nil
	}

	// 检查 DKG 是否已设置（需要私钥来连接链）
	if s.dkg == nil || s.dkg.Signer == nil {
		return errors.New("DKG not initialized, cannot connect to chains")
	}

	// 遍历所有 ChainConfig，连接到每个区块链
	for chainId, config := range *chainConfigs {
		if config == nil {
			continue
		}

		// 连接到区块链
		chainApi, err := chains.ConnectChain(config.Urls, s.dkg.Signer)
		if err != nil {
			return fmt.Errorf("failed to connect to chain %d: %w", chainId, err)
		}

		// 存储连接结果
		s.chains[chainId] = &chainApi
	}

	return nil
}

func (s *SideChain) getChain(chainId uint32) (chains.ChainApi, error) {
	chain, ok := s.chains[chainId]
	if !ok {
		return nil, errors.New("chain not found")
	}
	return *chain, nil
}
