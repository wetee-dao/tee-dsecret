package module

import (
	chain "github.com/wetee-dao/go-sdk"
	"github.com/wetee-dao/go-sdk/pallet/app"
	"github.com/wetee-dao/go-sdk/pallet/types"

	"errors"
)

// Worker
type App struct {
	Client *chain.ChainClient
	Signer *chain.Signer
}

func (w *App) GetApp(publey []byte, id uint64) (*types.TeeApp, error) {
	if len(publey) != 32 {
		return nil, errors.New("publey length error")
	}

	var mss [32]byte
	copy(mss[:], publey)

	res, ok, err := app.GetTEEAppsLatest(w.Client.Api.RPC.State, mss, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("GetK8sClusterAccountsLatest => not start")
	}
	return &res, nil
}

func (w *App) GetAccount(id uint64) ([]byte, error) {
	res, ok, err := app.GetAppIdAccountsLatest(w.Client.Api.RPC.State, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("GetAppIdAccountsLatest => not ok")
	}
	return res[:], nil
}

func (w *App) GetVersionLatest(id uint64) (uint64, error) {
	res, ok, err := app.GetAppVersionLatest(w.Client.Api.RPC.State, id)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("GetVersionLatest => not ok")
	}
	return uint64(res), nil
}

func (w *App) GetSecretEnv(id uint64) ([]byte, bool, error) {
	return app.GetSecretEnvsLatest(w.Client.Api.RPC.State, id)
}
