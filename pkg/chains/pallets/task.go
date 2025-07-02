package pallets

import (
	"errors"

	"github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/task"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/types"
)

func (c *Chain) GetTaskAccount(id uint64) ([]byte, error) {
	res, ok, err := task.GetTaskIdAccountsLatest(c.Api.RPC.State, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("GetAppIdAccountsLatest => not ok")
	}
	return res[:], nil
}

func (c *Chain) GetTaskVersionLatest(id uint64) (uint64, error) {
	res, ok, err := task.GetTaskVersionLatest(c.Api.RPC.State, id)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("GetVersionLatest => not ok")
	}
	return uint64(res), nil
}

func (c *Chain) GetTask(publey []byte, id uint64) (*types.TeeTask, error) {
	if len(publey) != 32 {
		return nil, errors.New("publey length error")
	}

	var mss [32]byte
	copy(mss[:], publey)

	res, ok, err := task.GetTEETasksLatest(c.Api.RPC.State, mss, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("GetK8sClusterAccountsLatest => not ok")
	}
	return &res, nil
}

func (c *Chain) GetTaskSecretEnv(id uint64) ([]byte, bool, error) {
	return task.GetSecretEnvsLatest(c.Api.RPC.State, id)
}
