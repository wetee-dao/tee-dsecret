package module

import (
	chain "github.com/wetee-dao/go-sdk"
	"github.com/wetee-dao/go-sdk/pallet/task"
	"github.com/wetee-dao/go-sdk/pallet/types"

	"errors"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
)

// Worker
type Task struct {
	Client *chain.ChainClient
	Signer *signature.KeyringPair
}

func (w *Task) GetAccount(id uint64) ([]byte, error) {
	res, ok, err := task.GetTaskIdAccountsLatest(w.Client.Api.RPC.State, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("GetAppIdAccountsLatest => not ok")
	}
	return res[:], nil
}

func (w *Task) GetVersionLatest(id uint64) (uint64, error) {
	res, ok, err := task.GetTaskVersionLatest(w.Client.Api.RPC.State, id)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("GetVersionLatest => not ok")
	}
	return uint64(res), nil
}

func (w *Task) GetTask(publey []byte, id uint64) (*types.TeeTask, error) {
	if len(publey) != 32 {
		return nil, errors.New("publey length error")
	}

	var mss [32]byte
	copy(mss[:], publey)

	res, ok, err := task.GetTEETasksLatest(w.Client.Api.RPC.State, mss, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("GetK8sClusterAccountsLatest => not ok")
	}
	return &res, nil
}

func (w *Task) GetSecretEnv(id uint64) ([]byte, bool, error) {
	return task.GetSecretEnvsLatest(w.Client.Api.RPC.State, id)
}
