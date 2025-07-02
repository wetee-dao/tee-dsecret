package pallets

import (
	"errors"

	"github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/app"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/types"
)

func (c *Chain) GetApp(publey []byte, id uint64) (*types.TeeApp, error) {
	if len(publey) != 32 {
		return nil, errors.New("publey length error")
	}

	var mss [32]byte
	copy(mss[:], publey)

	res, ok, err := app.GetTEEAppsLatest(c.Api.RPC.State, mss, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("GetK8sClusterAccountsLatest => not start")
	}
	return &res, nil
}

func (c *Chain) GetAppAccount(id uint64) ([]byte, error) {
	res, ok, err := app.GetAppIdAccountsLatest(c.Api.RPC.State, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("GetAppIdAccountsLatest => not ok")
	}
	return res[:], nil
}

func (c *Chain) GetAppVersionLatest(id uint64) (uint64, error) {
	res, ok, err := app.GetAppVersionLatest(c.Api.RPC.State, id)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("GetVersionLatest => not ok")
	}
	return uint64(res), nil
}

func (c *Chain) GetAppSecretEnv(id uint64) ([]byte, bool, error) {
	return app.GetSecretEnvsLatest(c.Api.RPC.State, id)
}
