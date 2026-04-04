package dao_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/side-chain/contracts/dao"
)

// stubTxn satisfies ContractApi for error-path tests that never touch storage.
type stubTxn struct{}

func (stubTxn) GetHeight() int64   { return 0 }
func (stubTxn) GetTxn() *model.Txn { return nil }
func (stubTxn) GetCaller() []byte { return nil }

func TestExecCall_MissingCaller(t *testing.T) {
	m := dao.DaoMutation{DAO: *dao.NewDAO(stubTxn{})}
	err := m.ExecCall(&model.ContractCall{Method: "init", Args: nil})
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing caller")
}

func TestExecCall_UnknownMethod(t *testing.T) {
	m := dao.DaoMutation{DAO: *dao.NewDAO(&fixedCaller{caller: []byte{1}})}
	err := m.ExecCall(&model.ContractCall{Method: "no_such_method", Args: nil})
	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown method")
}

type fixedCaller struct {
	stubTxn
	caller []byte
}

func (f *fixedCaller) GetCaller() []byte { return f.caller }

func TestExecQuery_EmptyMethod(t *testing.T) {
	q := dao.DaoQuery{DAO: *dao.NewDAO(stubTxn{})}
	_, err := q.ExecQuery(&model.ContractCall{Method: "", Args: nil})
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty query method")
}

func TestExecQuery_UnknownMethod(t *testing.T) {
	q := dao.DaoQuery{DAO: *dao.NewDAO(stubTxn{})}
	_, err := q.ExecQuery(&model.ContractCall{Method: "no_such_query", Args: nil})
	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown query method")
}

func TestExecCall_InitTooFewArgs(t *testing.T) {
	m := dao.DaoMutation{DAO: *dao.NewDAO(&fixedCaller{caller: []byte{1}})}
	err := m.ExecCall(&model.ContractCall{Method: "init", Args: [][]byte{}})
	require.Error(t, err)
	require.Contains(t, err.Error(), "expects at least")
}

// TestExecCall_TableInsufficientArgs 覆盖所有需要参数的 mutation 在参数不足时的路由错误（不访问存储）。
func TestExecCall_TableInsufficientArgs(t *testing.T) {
	m := dao.DaoMutation{DAO: *dao.NewDAO(&fixedCaller{caller: []byte{1}})}
	cases := []struct {
		method string
		args   [][]byte
	}{
		{"join", nil},
		{"join", [][]byte{{1}}},
		{"submit_proposal", nil},
		{"submit_proposal", [][]byte{{1}}},
		{"deposit_proposal", nil},
		{"deposit_proposal", [][]byte{{1}}},
		{"submit_vote", nil},
		{"submit_vote", [][]byte{{1}, {1}}},
		{"cancel_vote", nil},
		{"unlock", nil},
		{"exec_proposal", nil},
		{"cancel_proposal", nil},
		{"transfer", nil},
		{"transfer", [][]byte{{1}}},
		{"approve", nil},
		{"approve", [][]byte{{1}}},
		{"transfer_from", nil},
		{"transfer_from", [][]byte{{1}, {2}}},
		{"spend", nil},
		{"spend", [][]byte{{1}, {2}}},
		{"payout", nil},
		{"set_public_join", nil},
		{"add_track", nil},
		{"set_default_track", nil},
	}
	for _, tc := range cases {
		t.Run(tc.method, func(t *testing.T) {
			err := m.ExecCall(&model.ContractCall{Method: tc.method, Args: tc.args})
			require.Error(t, err)
			require.Contains(t, err.Error(), "expects at least")
		})
	}
}
