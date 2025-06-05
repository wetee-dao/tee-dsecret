package module

import (
	chain "github.com/wetee-dao/go-sdk"
	"github.com/wetee-dao/go-sdk/pallet/gpu"
	"github.com/wetee-dao/go-sdk/pallet/types"

	"errors"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
)

// Worker
type GpuApp struct {
	Client *chain.ChainClient
	Signer *signature.KeyringPair
}

func (w *GpuApp) GetApp(publey []byte, id uint64) (*types.GpuApp, error) {
	if len(publey) != 32 {
		return nil, errors.New("publey length error")
	}

	var mss [32]byte
	copy(mss[:], publey)

	res, ok, err := gpu.GetGPUAppsLatest(w.Client.Api.RPC.State, mss, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("GetK8sClusterAccountsLatest => not start")
	}
	return &res, nil
}

func (w *GpuApp) GetAccount(id uint64) ([]byte, error) {
	res, ok, err := gpu.GetAppIdAccountsLatest(w.Client.Api.RPC.State, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("GetAppIdAccountsLatest => not ok")
	}
	return res[:], nil
}

func (w *GpuApp) GetVersionLatest(id uint64) (uint64, error) {
	res, ok, err := gpu.GetAppVersionLatest(w.Client.Api.RPC.State, id)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("GetVersionLatest => not ok")
	}
	return uint64(res), nil
}

func (w *GpuApp) GetSecretEnv(id uint64) ([]byte, bool, error) {
	return gpu.GetSecretEnvsLatest(w.Client.Api.RPC.State, id)
}
