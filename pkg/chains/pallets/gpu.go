package pallets

import (
	"errors"

	"github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/gpu"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/types"
)

func (c *Chain) GetGpuApp(publey []byte, id uint64) (*types.GpuApp, error) {
	if len(publey) != 32 {
		return nil, errors.New("publey length error")
	}

	var mss [32]byte
	copy(mss[:], publey)

	res, ok, err := gpu.GetGPUAppsLatest(c.Api().RPC.State, mss, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("GetK8sClusterAccountsLatest => not start")
	}
	return &res, nil
}

func (c *Chain) GetGpuAccount(id uint64) ([]byte, error) {
	res, ok, err := gpu.GetAppIdAccountsLatest(c.Api().RPC.State, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("GetAppIdAccountsLatest => not ok")
	}
	return res[:], nil
}

func (c *Chain) GetGpuVersionLatest(id uint64) (uint64, error) {
	res, ok, err := gpu.GetAppVersionLatest(c.Api().RPC.State, id)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("GetVersionLatest => not ok")
	}
	return uint64(res), nil
}

func (c *Chain) GetGpuSecretEnv(id uint64) ([]byte, bool, error) {
	return gpu.GetSecretEnvsLatest(c.Api().RPC.State, id)
}
