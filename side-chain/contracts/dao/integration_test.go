package dao_test

import (
	"math/big"
	"os"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/stretchr/testify/require"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/side-chain/contracts/dao"
)

const testDBSubdir = "chain_data/wetee"

// testRuntime implements model.ContractApi against a single pebble transaction.
type testRuntime struct {
	height int64
	caller []byte
	txn    *model.Txn
}

func (r *testRuntime) GetHeight() int64   { return r.height }
func (r *testRuntime) GetTxn() *model.Txn { return r.txn }
func (r *testRuntime) GetCaller() []byte  { return r.caller }

func scaleBytes(t *testing.T, v any) []byte {
	t.Helper()
	b, err := codec.Encode(v)
	require.NoError(t, err)
	return b
}

// chdirTempDB opens model DB under t.TempDir() so tests do not touch the repo's chain_data.
func chdirTempDB(t *testing.T) (cleanup func()) {
	t.Helper()
	tmp := t.TempDir()
	oldWD, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmp))
	require.NoError(t, os.MkdirAll(testDBSubdir, 0o755))

	db, err := model.NewDB()
	require.NoError(t, err)

	return func() {
		_ = db.Close()
		_ = os.Chdir(oldWD)
	}
}

func TestIntegration_InitAndQueryMembers(t *testing.T) {
	defer chdirTempDB(t)()

	gov := []byte{0x01, 0x02, 0x03}
	members := []dao.Member{
		{Account: []byte{0x0a}, Balance: big.NewInt(1000).Bytes()},
	}
	track := dao.TrackData{
		Name:               "t",
		PreparePeriod:      0,
		MaxDeciding:        1000,
		ConfirmPeriod:      1,
		DecisionPeriod:     1,
		MinEnactmentPeriod: 0,
		DecisionDeposit:    big.NewInt(1).Bytes(),
		MaxBalance:         big.NewInt(1_000_000).Bytes(),
	}

	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()

	rt := &testRuntime{height: 1, caller: gov, txn: txn}
	mut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}

	args := [][]byte{
		scaleBytes(t, members),
		scaleBytes(t, true),
		scaleBytes(t, gov),
		scaleBytes(t, track),
	}
	err := mut.ExecCall(&model.ContractCall{
		Contract: "dao",
		Method:   "init",
		Args:     args,
	})
	require.NoError(t, err)

	q := dao.DaoQuery{DAO: *dao.NewDAO(rt)}
	raw, err := q.ExecQuery(&model.ContractCall{Method: "members", Args: nil})
	require.NoError(t, err)
	var out []dao.Member
	require.NoError(t, codec.Decode(raw, &out))
	require.Len(t, out, 1)
	require.Equal(t, members[0].Account, out[0].Account)
	require.Equal(t, 0, big.NewInt(1000).Cmp(new(big.Int).SetBytes(out[0].Balance)))

	ts, err := q.ExecQuery(&model.ContractCall{Method: "total_supply", Args: nil})
	require.NoError(t, err)
	var supply []byte
	require.NoError(t, codec.Decode(ts, &supply))
	require.Equal(t, 0, big.NewInt(1000).Cmp(new(big.Int).SetBytes(supply)))

	bal, err := q.ExecQuery(&model.ContractCall{
		Method: "balance_of",
		Args:   [][]byte{scaleBytes(t, []byte{0x0a})},
	})
	require.NoError(t, err)
	var b []byte
	require.NoError(t, codec.Decode(bal, &b))
	require.Equal(t, 0, big.NewInt(1000).Cmp(new(big.Int).SetBytes(b)))
}

func TestIntegration_PublicJoinDeniedWhenDisabled(t *testing.T) {
	defer chdirTempDB(t)()

	gov := []byte{0x01}
	members := []dao.Member{{Account: gov, Balance: big.NewInt(10).Bytes()}}

	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()

	rt := &testRuntime{height: 1, caller: gov, txn: txn}
	mut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	err := mut.ExecCall(&model.ContractCall{
		Method: "init",
		Args: [][]byte{
			scaleBytes(t, members),
			scaleBytes(t, false), // public join off
			scaleBytes(t, gov),
		},
	})
	require.NoError(t, err)

	rt.caller = []byte{0x99} // non-gov user
	err = mut.ExecCall(&model.ContractCall{Method: "public_join", Args: nil})
	require.ErrorIs(t, err, dao.ErrPublicJoinNotAllowed)
}

func TestIntegration_JoinByGov(t *testing.T) {
	defer chdirTempDB(t)()

	gov := []byte{0x01}
	bob := []byte{0x02}
	members := []dao.Member{{Account: gov, Balance: big.NewInt(100).Bytes()}}

	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()

	rt := &testRuntime{height: 1, caller: gov, txn: txn}
	mut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	require.NoError(t, mut.ExecCall(&model.ContractCall{
		Method: "init",
		Args: [][]byte{
			scaleBytes(t, members),
			scaleBytes(t, false),
			scaleBytes(t, gov),
		},
	}))

	err := mut.ExecCall(&model.ContractCall{
		Method: "join",
		Args: [][]byte{
			scaleBytes(t, bob),
			scaleBytes(t, big.NewInt(50).Bytes()),
		},
	})
	require.NoError(t, err)

	q := dao.DaoQuery{DAO: *dao.NewDAO(rt)}
	raw, err := q.ExecQuery(&model.ContractCall{
		Method: "balance_of",
		Args:   [][]byte{scaleBytes(t, bob)},
	})
	require.NoError(t, err)
	var bal []byte
	require.NoError(t, codec.Decode(raw, &bal))
	require.Equal(t, 0, big.NewInt(50).Cmp(new(big.Int).SetBytes(bal)))
}

func TestIntegration_TransferDisabledByDefault(t *testing.T) {
	defer chdirTempDB(t)()

	gov := []byte{0x01}
	alice := []byte{0x0a}
	bob := []byte{0x0b}
	members := []dao.Member{
		{Account: alice, Balance: big.NewInt(100).Bytes()},
		{Account: bob, Balance: big.NewInt(10).Bytes()},
	}

	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()

	rt := &testRuntime{height: 1, caller: alice, txn: txn}
	mut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	require.NoError(t, mut.ExecCall(&model.ContractCall{
		Method: "init",
		Args: [][]byte{
			scaleBytes(t, members),
			scaleBytes(t, false),
			scaleBytes(t, gov),
		},
	}))

	err := mut.ExecCall(&model.ContractCall{
		Method: "transfer",
		Args: [][]byte{
			scaleBytes(t, bob),
			scaleBytes(t, big.NewInt(5).Bytes()),
		},
	})
	require.ErrorIs(t, err, dao.ErrTransferDisabled)
}

func TestIntegration_CancelProposal(t *testing.T) {
	defer chdirTempDB(t)()

	gov := []byte{0x01}
	members := []dao.Member{{Account: gov, Balance: big.NewInt(10_000).Bytes()}}
	track := dao.TrackData{
		Name:            "t",
		MaxDeciding:     100,
		DecisionDeposit: big.NewInt(1).Bytes(),
		MaxBalance:      big.NewInt(5000).Bytes(),
	}

	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()

	rt := &testRuntime{height: 1, caller: gov, txn: txn}
	mut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	require.NoError(t, mut.ExecCall(&model.ContractCall{
		Method: "init",
		Args: [][]byte{
			scaleBytes(t, members),
			scaleBytes(t, false),
			scaleBytes(t, gov),
			scaleBytes(t, track),
		},
	}))

	call := dao.CallContent{Amount: big.NewInt(10).Bytes()}
	require.NoError(t, mut.ExecCall(&model.ContractCall{
		Method: "submit_proposal",
		Args: [][]byte{
			scaleBytes(t, call),
			scaleBytes(t, uint32(0)),
		},
	}))

	require.NoError(t, mut.ExecCall(&model.ContractCall{
		Method: "cancel_proposal",
		Args:   [][]byte{scaleBytes(t, uint32(0))},
	}))

	q := dao.DaoQuery{DAO: *dao.NewDAO(rt)}
	raw, err := q.ExecQuery(&model.ContractCall{Method: "proposals", Args: nil})
	require.NoError(t, err)
	var props []dao.Proposal
	require.NoError(t, codec.Decode(raw, &props))
	require.Len(t, props, 1)
	require.Equal(t, dao.ProposalCanceled, props[0].Status.State)
}
