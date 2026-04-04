package dao

// 本文件覆盖 dao_gen 中 ExecCall / ExecQuery 的主要成功与失败路径。
// 未通过合约暴露的方法（如 DAO.Delete）不在此测。

import (
	"math/big"
	"os"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/stretchr/testify/require"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

// fullTestRT 与 integration_test 中 testRuntime 相同，供本包完整场景使用。
type fullTestRT struct {
	height int64
	caller []byte
	txn    *model.Txn
}

func (r *fullTestRT) GetHeight() int64   { return r.height }
func (r *fullTestRT) GetTxn() *model.Txn { return r.txn }
func (r *fullTestRT) GetCaller() []byte  { return r.caller }

func fullScale(t *testing.T, v any) []byte {
	t.Helper()
	b, err := codec.Encode(v)
	require.NoError(t, err)
	return b
}

func fullChdirDB(t *testing.T) func() {
	t.Helper()
	tmp := t.TempDir()
	old, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmp))
	require.NoError(t, os.MkdirAll("chain_data/wetee", 0o755))
	db, err := model.NewDB()
	require.NoError(t, err)
	return func() {
		_ = db.Close()
		_ = os.Chdir(old)
	}
}

func setTransferEnabled(m *DaoMutation, on bool) error {
	return m.transferEnabled.Set(m.api.GetTxn(), singletonKey, on)
}

func defaultTrackBig() TrackData {
	return TrackData{
		Name:               "gov",
		PreparePeriod:      0,
		MaxDeciding:        50,
		ConfirmPeriod:      1,
		DecisionPeriod:     2,
		MinEnactmentPeriod: 0,
		DecisionDeposit:    big.NewInt(10).Bytes(),
		MaxBalance:         big.NewInt(1_000_000).Bytes(),
	}
}

func exec(m *DaoMutation, method string, args [][]byte) error {
	return m.ExecCall(&model.ContractCall{Method: method, Args: args})
}

func query(q *DaoQuery, method string, args [][]byte) ([]byte, error) {
	return q.ExecQuery(&model.ContractCall{Method: method, Args: args})
}

// --- ExecQuery：参数缺失与全方法 smoke ---

func TestExecQuery_MissingArgs(t *testing.T) {
	defer fullChdirDB(t)()
	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()
	rt := &fullTestRT{txn: txn}
	q := DaoQuery{DAO: *NewDAO(rt)}

	_, err := query(&q, "balance_of", nil)
	require.Error(t, err)

	_, err = query(&q, "lock_balance_of", [][]byte{})
	require.Error(t, err)

	_, err = query(&q, "allowance", [][]byte{fullScale(t, []byte{1})})
	require.Error(t, err)

	_, err = query(&q, "track", nil)
	require.Error(t, err)

	_, err = query(&q, "proposal", nil)
	require.Error(t, err)

	_, err = query(&q, "proposal_status", nil)
	require.Error(t, err)

	_, err = query(&q, "vote", nil)
	require.Error(t, err)
}

func TestExecQuery_AllSmokeAfterInit(t *testing.T) {
	defer fullChdirDB(t)()
	gov := []byte{0x01}
	members := []Member{{Account: gov, Balance: big.NewInt(5000).Bytes()}}
	tr := defaultTrackBig()

	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()
	rt := &fullTestRT{height: 1, caller: gov, txn: txn}
	mut := DaoMutation{DAO: *NewDAO(rt)}
	require.NoError(t, exec(&mut, "init", [][]byte{
		fullScale(t, members),
		fullScale(t, true),
		fullScale(t, gov),
		fullScale(t, tr),
	}))

	q := DaoQuery{DAO: *NewDAO(rt)}

	raw, err := query(&q, "members", nil)
	require.NoError(t, err)
	var ms []Member
	require.NoError(t, codec.Decode(raw, &ms))

	raw, err = query(&q, "public_join", nil)
	require.NoError(t, err)
	var pj bool
	require.NoError(t, codec.Decode(raw, &pj))
	require.True(t, pj)

	raw, err = query(&q, "total_supply", nil)
	require.NoError(t, err)
	var ts []byte
	require.NoError(t, codec.Decode(raw, &ts))

	raw, err = query(&q, "balance_of", [][]byte{fullScale(t, gov)})
	require.NoError(t, err)

	raw, err = query(&q, "lock_balance_of", [][]byte{fullScale(t, gov)})
	require.NoError(t, err)

	raw, err = query(&q, "allowance", [][]byte{fullScale(t, gov), fullScale(t, []byte{0x02})})
	require.NoError(t, err)

	raw, err = query(&q, "tracks", nil)
	require.NoError(t, err)
	var tracks []TrackData
	require.NoError(t, codec.Decode(raw, &tracks))
	require.NotEmpty(t, tracks)

	raw, err = query(&q, "track", [][]byte{fullScale(t, uint32(0))})
	require.NoError(t, err)
	require.NotEmpty(t, raw)

	raw, err = query(&q, "default_track", nil)
	require.NoError(t, err)
	var dt uint32
	require.NoError(t, codec.Decode(raw, &dt))
	require.Equal(t, uint32(0), dt)

	_, err = query(&q, "proposal", [][]byte{fullScale(t, uint32(0))})
	require.Error(t, err) // 尚无提案

	_, err = query(&q, "proposals", nil)
	require.NoError(t, err)

	_, err = query(&q, "proposal_status", [][]byte{fullScale(t, uint32(0))})
	require.Error(t, err)

	_, err = query(&q, "vote", [][]byte{fullScale(t, uint64(0))})
	require.Error(t, err)

	_, err = query(&q, "votes", nil)
	require.NoError(t, err)
}

// --- 成员与公开加入 ---

func TestScenario_PublicJoinSuccess(t *testing.T) {
	defer fullChdirDB(t)()
	gov := []byte{0x01}
	alice := []byte{0xaa}
	members := []Member{{Account: gov, Balance: big.NewInt(100).Bytes()}}

	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()
	rt := &fullTestRT{height: 1, caller: gov, txn: txn}
	mut := DaoMutation{DAO: *NewDAO(rt)}
	require.NoError(t, exec(&mut, "init", [][]byte{
		fullScale(t, members),
		fullScale(t, true),
		fullScale(t, gov),
	}))

	rt.caller = alice
	require.NoError(t, exec(&mut, "public_join", nil))

	q := DaoQuery{DAO: *NewDAO(rt)}
	raw, err := query(&q, "balance_of", [][]byte{fullScale(t, alice)})
	require.NoError(t, err)
	var bal []byte
	require.NoError(t, codec.Decode(raw, &bal))
	require.True(t, isZero(bal))
}

func TestScenario_SetPublicJoin_NotGov(t *testing.T) {
	defer fullChdirDB(t)()
	gov := []byte{0x01}
	alice := []byte{0xaa}
	members := []Member{{Account: gov, Balance: big.NewInt(10).Bytes()}}

	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()
	rt := &fullTestRT{height: 1, caller: alice, txn: txn}
	mut := DaoMutation{DAO: *NewDAO(rt)}
	require.NoError(t, exec(&mut, "init", [][]byte{
		fullScale(t, members),
		fullScale(t, false),
		fullScale(t, gov),
	}))

	err := exec(&mut, "set_public_join", [][]byte{fullScale(t, true)})
	require.ErrorIs(t, err, ErrMustCallByGov)
}

func TestScenario_Join_NotGov(t *testing.T) {
	defer fullChdirDB(t)()
	gov := []byte{0x01}
	alice := []byte{0xaa}
	members := []Member{{Account: gov, Balance: big.NewInt(10).Bytes()}}

	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()
	rt := &fullTestRT{height: 1, caller: alice, txn: txn}
	mut := DaoMutation{DAO: *NewDAO(rt)}
	require.NoError(t, exec(&mut, "init", [][]byte{
		fullScale(t, members),
		fullScale(t, false),
		fullScale(t, gov),
	}))

	err := exec(&mut, "join", [][]byte{fullScale(t, []byte{0xbb}), fullScale(t, big.NewInt(1).Bytes())})
	require.ErrorIs(t, err, ErrMustCallByGov)
}

func TestScenario_LeaveAndLeaveWithBurn(t *testing.T) {
	defer fullChdirDB(t)()
	gov := []byte{0x01}
	bob := []byte{0xbb}
	members := []Member{{Account: gov, Balance: big.NewInt(1000).Bytes()}}

	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()
	rt := &fullTestRT{height: 1, caller: gov, txn: txn}
	mut := DaoMutation{DAO: *NewDAO(rt)}
	require.NoError(t, exec(&mut, "init", [][]byte{
		fullScale(t, members),
		fullScale(t, false),
		fullScale(t, gov),
	}))
	require.NoError(t, exec(&mut, "join", [][]byte{
		fullScale(t, bob),
		fullScale(t, []byte{}),
	}))

	rt.caller = bob
	require.NoError(t, exec(&mut, "leave", nil))

	rt.caller = gov
	require.NoError(t, exec(&mut, "join", [][]byte{
		fullScale(t, bob),
		fullScale(t, big.NewInt(100).Bytes()),
	}))
	rt.caller = bob
	require.ErrorIs(t, exec(&mut, "leave", nil), ErrMemberBalanceNotZero)

	require.NoError(t, exec(&mut, "leave_with_burn", nil))
}

// --- 代币：转账 / 授权 ---

func TestScenario_TransferApproveTransferFrom(t *testing.T) {
	defer fullChdirDB(t)()
	gov := []byte{0x01}
	alice := []byte{0x0a}
	bob := []byte{0x0b}
	members := []Member{
		{Account: alice, Balance: big.NewInt(1000).Bytes()},
		{Account: bob, Balance: big.NewInt(100).Bytes()},
	}

	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()
	rt := &fullTestRT{height: 1, caller: alice, txn: txn}
	mut := DaoMutation{DAO: *NewDAO(rt)}
	require.NoError(t, exec(&mut, "init", [][]byte{
		fullScale(t, members),
		fullScale(t, false),
		fullScale(t, gov),
	}))
	require.NoError(t, setTransferEnabled(&mut, true))

	require.NoError(t, exec(&mut, "transfer", [][]byte{
		fullScale(t, bob),
		fullScale(t, big.NewInt(50).Bytes()),
	}))

	rt.caller = alice
	require.NoError(t, exec(&mut, "approve", [][]byte{
		fullScale(t, bob),
		fullScale(t, big.NewInt(30).Bytes()),
	}))

	q := DaoQuery{DAO: *NewDAO(rt)}
	raw, err := query(&q, "allowance", [][]byte{fullScale(t, alice), fullScale(t, bob)})
	require.NoError(t, err)
	var alw []byte
	require.NoError(t, codec.Decode(raw, &alw))
	require.Equal(t, 0, big.NewInt(30).Cmp(new(big.Int).SetBytes(alw)))

	rt.caller = bob
	require.NoError(t, exec(&mut, "transfer_from", [][]byte{
		fullScale(t, alice),
		fullScale(t, bob),
		fullScale(t, big.NewInt(10).Bytes()),
	}))
}

func TestScenario_Transfer_LowBalance(t *testing.T) {
	defer fullChdirDB(t)()
	gov := []byte{0x01}
	alice := []byte{0x0a}
	bob := []byte{0x0b}
	members := []Member{
		{Account: alice, Balance: big.NewInt(5).Bytes()},
		{Account: bob, Balance: big.NewInt(1).Bytes()},
	}
	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()
	rt := &fullTestRT{height: 1, caller: alice, txn: txn}
	mut := DaoMutation{DAO: *NewDAO(rt)}
	require.NoError(t, exec(&mut, "init", [][]byte{
		fullScale(t, members),
		fullScale(t, false),
		fullScale(t, gov),
	}))
	require.NoError(t, setTransferEnabled(&mut, true))
	require.ErrorIs(t, exec(&mut, "transfer", [][]byte{
		fullScale(t, bob),
		fullScale(t, big.NewInt(100).Bytes()),
	}), ErrLowBalance)
}

// --- 轨道治理 ---

func TestScenario_AddTrack_SetDefaultTrack(t *testing.T) {
	defer fullChdirDB(t)()
	gov := []byte{0x01}
	members := []Member{{Account: gov, Balance: big.NewInt(100).Bytes()}}
	tr0 := defaultTrackBig()

	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()
	rt := &fullTestRT{height: 1, caller: gov, txn: txn}
	mut := DaoMutation{DAO: *NewDAO(rt)}
	require.NoError(t, exec(&mut, "init", [][]byte{
		fullScale(t, members),
		fullScale(t, false),
		fullScale(t, gov),
		fullScale(t, tr0),
	}))

	tr1 := tr0
	tr1.Name = "second"
	require.NoError(t, exec(&mut, "add_track", [][]byte{fullScale(t, tr1)}))
	require.NoError(t, exec(&mut, "set_default_track", [][]byte{fullScale(t, uint32(1))}))

	q := DaoQuery{DAO: *NewDAO(rt)}
	raw, err := query(&q, "default_track", nil)
	require.NoError(t, err)
	var dt uint32
	require.NoError(t, codec.Decode(raw, &dt))
	require.Equal(t, uint32(1), dt)
}

// --- 提案生命周期：存款错误、投票、执行拒绝/通过 ---

func TestScenario_Deposit_InvalidTimeAndAmount(t *testing.T) {
	defer fullChdirDB(t)()
	gov := []byte{0x01}
	members := []Member{{Account: gov, Balance: big.NewInt(10_000).Bytes()}}
	tr := defaultTrackBig()
	tr.PreparePeriod = 10

	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()
	rt := &fullTestRT{height: 1, caller: gov, txn: txn}
	mut := DaoMutation{DAO: *NewDAO(rt)}
	require.NoError(t, exec(&mut, "init", [][]byte{
		fullScale(t, members),
		fullScale(t, false),
		fullScale(t, gov),
		fullScale(t, tr),
	}))
	require.NoError(t, exec(&mut, "submit_proposal", [][]byte{
		fullScale(t, CallContent{Amount: big.NewInt(100).Bytes()}),
		fullScale(t, uint32(0)),
	}))

	rt.height = 5
	require.ErrorIs(t, exec(&mut, "deposit_proposal", [][]byte{
		fullScale(t, uint32(0)),
		fullScale(t, big.NewInt(10).Bytes()),
	}), ErrInvalidDepositTime)

	rt.height = 20
	require.ErrorIs(t, exec(&mut, "deposit_proposal", [][]byte{
		fullScale(t, uint32(0)),
		fullScale(t, big.NewInt(5).Bytes()),
	}), ErrInvalidDeposit)
}

func TestScenario_SubmitVote_NotOngoing(t *testing.T) {
	defer fullChdirDB(t)()
	gov := []byte{0x01}
	members := []Member{{Account: gov, Balance: big.NewInt(10_000).Bytes()}}
	tr := defaultTrackBig()

	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()
	rt := &fullTestRT{height: 1, caller: gov, txn: txn}
	mut := DaoMutation{DAO: *NewDAO(rt)}
	require.NoError(t, exec(&mut, "init", [][]byte{
		fullScale(t, members),
		fullScale(t, false),
		fullScale(t, gov),
		fullScale(t, tr),
	}))
	require.NoError(t, exec(&mut, "submit_proposal", [][]byte{
		fullScale(t, CallContent{Amount: big.NewInt(100).Bytes()}),
		fullScale(t, uint32(0)),
	}))

	require.ErrorIs(t, exec(&mut, "submit_vote", [][]byte{
		fullScale(t, uint32(0)),
		fullScale(t, true),
		fullScale(t, big.NewInt(10).Bytes()),
	}), ErrPropNotOngoing)
}

func TestScenario_ExecProposal_Reject(t *testing.T) {
	defer fullChdirDB(t)()
	gov := []byte{0x01}
	members := []Member{{Account: gov, Balance: big.NewInt(10_000).Bytes()}}
	tr := defaultTrackBig()
	tr.MaxDeciding = 20

	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()
	rt := &fullTestRT{height: 10, caller: gov, txn: txn}
	mut := DaoMutation{DAO: *NewDAO(rt)}
	require.NoError(t, exec(&mut, "init", [][]byte{
		fullScale(t, members),
		fullScale(t, false),
		fullScale(t, gov),
		fullScale(t, tr),
	}))
	require.NoError(t, exec(&mut, "submit_proposal", [][]byte{
		fullScale(t, CallContent{Amount: big.NewInt(50).Bytes()}),
		fullScale(t, uint32(0)),
	}))
	require.NoError(t, exec(&mut, "deposit_proposal", [][]byte{
		fullScale(t, uint32(0)),
		fullScale(t, big.NewInt(10).Bytes()),
	}))

	depositH := rt.height
	end := depositH + int64(tr.MaxDeciding)
	rt.height = end + 5
	require.NoError(t, exec(&mut, "exec_proposal", [][]byte{fullScale(t, uint32(0))}))

	q := DaoQuery{DAO: *NewDAO(rt)}
	raw, err := query(&q, "proposals", nil)
	require.NoError(t, err)
	var props []Proposal
	require.NoError(t, codec.Decode(raw, &props))
	require.Equal(t, ProposalRejected, props[0].Status.State)
}

func TestScenario_ExecProposal_Approve(t *testing.T) {
	defer fullChdirDB(t)()
	gov := []byte{0x01}
	members := []Member{{Account: gov, Balance: big.NewInt(100_000).Bytes()}}
	tr := defaultTrackBig()
	tr.MaxDeciding = 30
	tr.DecisionPeriod = 3

	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()
	rt := &fullTestRT{height: 100, caller: gov, txn: txn}
	mut := DaoMutation{DAO: *NewDAO(rt)}
	require.NoError(t, exec(&mut, "init", [][]byte{
		fullScale(t, members),
		fullScale(t, false),
		fullScale(t, gov),
		fullScale(t, tr),
	}))
	require.NoError(t, exec(&mut, "submit_proposal", [][]byte{
		fullScale(t, CallContent{Amount: big.NewInt(500).Bytes()}),
		fullScale(t, uint32(0)),
	}))
	depositH := rt.height
	require.NoError(t, exec(&mut, "deposit_proposal", [][]byte{
		fullScale(t, uint32(0)),
		fullScale(t, big.NewInt(10).Bytes()),
	}))

	require.NoError(t, exec(&mut, "submit_vote", [][]byte{
		fullScale(t, uint32(0)),
		fullScale(t, true),
		fullScale(t, big.NewInt(10).Bytes()),
	}))

	end := depositH + int64(tr.MaxDeciding)
	rt.height = end + int64(tr.DecisionPeriod) + 1

	require.NoError(t, exec(&mut, "exec_proposal", [][]byte{fullScale(t, uint32(0))}))

	q := DaoQuery{DAO: *NewDAO(rt)}
	raw, err := query(&q, "proposals", nil)
	require.NoError(t, err)
	var props []Proposal
	require.NoError(t, codec.Decode(raw, &props))
	require.Equal(t, ProposalApproved, props[0].Status.State)
}

func TestScenario_CancelVote(t *testing.T) {
	defer fullChdirDB(t)()
	gov := []byte{0x01}
	members := []Member{{Account: gov, Balance: big.NewInt(10_000).Bytes()}}
	tr := defaultTrackBig()

	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()
	rt := &fullTestRT{height: 1, caller: gov, txn: txn}
	mut := DaoMutation{DAO: *NewDAO(rt)}
	require.NoError(t, exec(&mut, "init", [][]byte{
		fullScale(t, members),
		fullScale(t, false),
		fullScale(t, gov),
		fullScale(t, tr),
	}))
	require.NoError(t, exec(&mut, "submit_proposal", [][]byte{
		fullScale(t, CallContent{Amount: big.NewInt(10).Bytes()}),
		fullScale(t, uint32(0)),
	}))
	require.NoError(t, exec(&mut, "deposit_proposal", [][]byte{
		fullScale(t, uint32(0)),
		fullScale(t, big.NewInt(10).Bytes()),
	}))
	require.NoError(t, exec(&mut, "submit_vote", [][]byte{
		fullScale(t, uint32(0)),
		fullScale(t, true),
		fullScale(t, big.NewInt(10).Bytes()),
	}))
	require.NoError(t, exec(&mut, "cancel_vote", [][]byte{fullScale(t, uint64(0))}))
}

func TestScenario_Unlock_AfterReject(t *testing.T) {
	defer fullChdirDB(t)()
	gov := []byte{0x01}
	members := []Member{{Account: gov, Balance: big.NewInt(10_000).Bytes()}}
	tr := defaultTrackBig()
	tr.MaxDeciding = 15

	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()
	rt := &fullTestRT{height: 5, caller: gov, txn: txn}
	mut := DaoMutation{DAO: *NewDAO(rt)}
	require.NoError(t, exec(&mut, "init", [][]byte{
		fullScale(t, members),
		fullScale(t, false),
		fullScale(t, gov),
		fullScale(t, tr),
	}))
	require.NoError(t, exec(&mut, "submit_proposal", [][]byte{
		fullScale(t, CallContent{Amount: big.NewInt(20).Bytes()}),
		fullScale(t, uint32(0)),
	}))
	require.NoError(t, exec(&mut, "deposit_proposal", [][]byte{
		fullScale(t, uint32(0)),
		fullScale(t, big.NewInt(10).Bytes()),
	}))
	require.NoError(t, exec(&mut, "submit_vote", [][]byte{
		fullScale(t, uint32(0)),
		fullScale(t, true),
		fullScale(t, big.NewInt(10).Bytes()),
	}))

	depositH := int64(5)
	end := depositH + int64(tr.MaxDeciding)
	rt.height = end + 2
	require.NoError(t, exec(&mut, "exec_proposal", [][]byte{fullScale(t, uint32(0))}))

	rt.height = end + 3
	require.NoError(t, exec(&mut, "unlock", [][]byte{fullScale(t, uint64(0))}))
}

// --- Treasury spend / payout ---

func TestScenario_SpendAndPayout(t *testing.T) {
	defer fullChdirDB(t)()
	gov := []byte{0x01}
	bob := []byte{0xbb}
	members := []Member{{Account: gov, Balance: big.NewInt(50_000).Bytes()}}
	tr := defaultTrackBig()

	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()
	rt := &fullTestRT{height: 1, caller: gov, txn: txn}
	mut := DaoMutation{DAO: *NewDAO(rt)}
	require.NoError(t, exec(&mut, "init", [][]byte{
		fullScale(t, members),
		fullScale(t, false),
		fullScale(t, gov),
		fullScale(t, tr),
	}))

	require.NoError(t, exec(&mut, "spend", [][]byte{
		fullScale(t, bob),
		fullScale(t, big.NewInt(100).Bytes()),
		fullScale(t, uint32(0)),
	}))

	rt.caller = bob
	require.ErrorIs(t, exec(&mut, "payout", [][]byte{fullScale(t, uint64(0))}), ErrMustCallByGov)

	rt.caller = gov
	require.NoError(t, exec(&mut, "payout", [][]byte{fullScale(t, uint64(0))}))
	require.ErrorIs(t, exec(&mut, "payout", [][]byte{fullScale(t, uint64(0))}), ErrSpendAlreadyExecuted)
}

func TestScenario_SubmitProposal_NoSuchTrack(t *testing.T) {
	defer fullChdirDB(t)()
	gov := []byte{0x01}
	members := []Member{{Account: gov, Balance: big.NewInt(1000).Bytes()}}
	tr := defaultTrackBig()

	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()
	rt := &fullTestRT{height: 1, caller: gov, txn: txn}
	mut := DaoMutation{DAO: *NewDAO(rt)}
	require.NoError(t, exec(&mut, "init", [][]byte{
		fullScale(t, members),
		fullScale(t, false),
		fullScale(t, gov),
		fullScale(t, tr),
	}))
	require.ErrorIs(t, exec(&mut, "submit_proposal", [][]byte{
		fullScale(t, CallContent{Amount: big.NewInt(1).Bytes()}),
		fullScale(t, uint32(99)),
	}), ErrNoTrack)
}

func TestScenario_ExecProposal_NotConfirmedOrInDecision(t *testing.T) {
	defer fullChdirDB(t)()
	gov := []byte{0x01}
	members := []Member{{Account: gov, Balance: big.NewInt(50_000).Bytes()}}
	tr := defaultTrackBig()
	tr.MaxDeciding = 100

	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()
	rt := &fullTestRT{height: 1, caller: gov, txn: txn}
	mut := DaoMutation{DAO: *NewDAO(rt)}
	require.NoError(t, exec(&mut, "init", [][]byte{
		fullScale(t, members),
		fullScale(t, false),
		fullScale(t, gov),
		fullScale(t, tr),
	}))
	require.NoError(t, exec(&mut, "submit_proposal", [][]byte{
		fullScale(t, CallContent{Amount: big.NewInt(10).Bytes()}),
		fullScale(t, uint32(0)),
	}))
	require.NoError(t, exec(&mut, "deposit_proposal", [][]byte{
		fullScale(t, uint32(0)),
		fullScale(t, big.NewInt(10).Bytes()),
	}))
	// 投票质押不足 DecisionDeposit，confirmed 为 false，height 未过 end → ErrProposalNotConfirmed
	require.NoError(t, exec(&mut, "submit_vote", [][]byte{
		fullScale(t, uint32(0)),
		fullScale(t, true),
		fullScale(t, big.NewInt(5).Bytes()),
	}))
	rt.height = 5
	require.ErrorIs(t, exec(&mut, "exec_proposal", [][]byte{fullScale(t, uint32(0))}), ErrProposalNotConfirmed)

	// 足够质押达成 confirmed，但高度仍在决策窗口内 → ErrProposalInDecision
	require.NoError(t, exec(&mut, "cancel_vote", [][]byte{fullScale(t, uint64(0))}))
	require.NoError(t, exec(&mut, "submit_vote", [][]byte{
		fullScale(t, uint32(0)),
		fullScale(t, true),
		fullScale(t, big.NewInt(10).Bytes()),
	}))
	rt.height = 10
	require.ErrorIs(t, exec(&mut, "exec_proposal", [][]byte{fullScale(t, uint32(0))}), ErrProposalInDecision)
}

func TestScenario_TransferFrom_InsufficientAllowance(t *testing.T) {
	defer fullChdirDB(t)()
	gov := []byte{0x01}
	alice := []byte{0x0a}
	bob := []byte{0x0b}
	members := []Member{
		{Account: alice, Balance: big.NewInt(500).Bytes()},
		{Account: bob, Balance: big.NewInt(50).Bytes()},
	}
	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()
	rt := &fullTestRT{height: 1, caller: bob, txn: txn}
	mut := DaoMutation{DAO: *NewDAO(rt)}
	require.NoError(t, exec(&mut, "init", [][]byte{
		fullScale(t, members),
		fullScale(t, false),
		fullScale(t, gov),
	}))
	require.NoError(t, setTransferEnabled(&mut, true))
	require.ErrorIs(t, exec(&mut, "transfer_from", [][]byte{
		fullScale(t, alice),
		fullScale(t, bob),
		fullScale(t, big.NewInt(10).Bytes()),
	}), ErrInsufficientAllowance)
}

func TestScenario_Payout_NotFound(t *testing.T) {
	defer fullChdirDB(t)()
	gov := []byte{0x01}
	members := []Member{{Account: gov, Balance: big.NewInt(100).Bytes()}}
	tr := defaultTrackBig()

	db := model.DBINS
	txn := db.NewTransaction()
	defer func() { _ = txn.Rollback() }()
	rt := &fullTestRT{height: 1, caller: gov, txn: txn}
	mut := DaoMutation{DAO: *NewDAO(rt)}
	require.NoError(t, exec(&mut, "init", [][]byte{
		fullScale(t, members),
		fullScale(t, false),
		fullScale(t, gov),
		fullScale(t, tr),
	}))
	require.ErrorIs(t, exec(&mut, "payout", [][]byte{fullScale(t, uint64(99))}), ErrSpendNotFound)
}
