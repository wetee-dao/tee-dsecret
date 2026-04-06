package dao_test

import (
	"math/big"
	"os"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/stretchr/testify/require"

	"github.com/wetee-dao/ink.go/util"
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
		model.DBINS = nil
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

// ---------------------------------------------------------------------------------------------------------------------

// daoWithGov creates a temporary DB-backed DAO with gov as the sole member and default track.
func daoWithGov(t *testing.T) (*testRuntime, func()) {
	t.Helper()
	tmp := t.TempDir()
	oldWD, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmp))
	require.NoError(t, os.MkdirAll(testDBSubdir, 0o755))

	db, err := model.NewDB()
	require.NoError(t, err)

	gov := []byte{0x01}
	members := []dao.Member{{Account: gov, Balance: big.NewInt(10_000).Bytes()}}
	track := dao.TrackData{
		Name:               "gov_track",
		PreparePeriod:      0,
		MaxDeciding:        100,
		ConfirmPeriod:      1,
		DecisionPeriod:     100,
		MinEnactmentPeriod: 0,
		DecisionDeposit:    big.NewInt(1).Bytes(),
		MaxBalance:         big.NewInt(1_000_000).Bytes(),
	}

	txn := db.NewTransaction()
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

	return rt, func() {
		_ = txn.Rollback()
		_ = db.Close()
		model.DBINS = nil
		_ = os.Chdir(oldWD)
	}
}

// TestIntegration_SetPublicJoin covers SetPublicJoin and the query side GetPublicJoin.
func TestIntegration_SetPublicJoin(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	q := dao.DaoQuery{DAO: *dao.NewDAO(rt)}
	raw, err := q.ExecQuery(&model.ContractCall{Method: "get_public_join", Args: nil})
	require.NoError(t, err)
	var pj bool
	require.NoError(t, codec.Decode(raw, &pj))
	require.False(t, pj)

	// Gov enables public join
	mut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	require.NoError(t, mut.ExecCall(&model.ContractCall{
		Method: "set_public_join",
		Args:   [][]byte{scaleBytes(t, true)},
	}))
	raw, err = q.ExecQuery(&model.ContractCall{Method: "get_public_join", Args: nil})
	require.NoError(t, err)
	require.NoError(t, codec.Decode(raw, &pj))
	require.True(t, pj)
}

// TestIntegration_Leave covers Leave: member with zero balance can leave.
func TestIntegration_Leave(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	bob := []byte{0x02}
	mut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	require.NoError(t, mut.ExecCall(&model.ContractCall{
		Method: "join",
		Args:   [][]byte{scaleBytes(t, bob), scaleBytes(t, big.NewInt(0).Bytes())},
	}))

	rt.caller = bob
	require.NoError(t, mut.ExecCall(&model.ContractCall{Method: "leave", Args: nil}))

	q := dao.DaoQuery{DAO: *dao.NewDAO(rt)}
	raw, err := q.ExecQuery(&model.ContractCall{Method: "members", Args: nil})
	require.NoError(t, err)
	var members []dao.Member
	require.NoError(t, codec.Decode(raw, &members))
	require.Len(t, members, 1) // only gov
}

// TestIntegration_Leave_FailsWhenBalanceNonZero verifies that a member with non-zero
// balance cannot leave.
func TestIntegration_Leave_FailsWhenBalanceNonZero(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	bob := []byte{0x02}
	mut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	require.NoError(t, mut.ExecCall(&model.ContractCall{
		Method: "join",
		Args:   [][]byte{scaleBytes(t, bob), scaleBytes(t, big.NewInt(100).Bytes())},
	}))

	rt.caller = bob
	err := mut.ExecCall(&model.ContractCall{Method: "leave", Args: nil})
	require.ErrorIs(t, err, dao.ErrMemberBalanceNotZero)
}

// TestIntegration_LeaveWithBurn verifies Delete removes member and burns issuance.
// Note: "delete" is not in ExecCall router (skipped in contractgen), so we call the Go method directly.
func TestIntegration_LeaveWithBurn(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	bob := []byte{0x02}
	daoMut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}

	// Join bob
	require.NoError(t, daoMut.Join(bob, big.NewInt(200).Bytes()))

	// Supply is 10200
	q := dao.DaoQuery{DAO: *dao.NewDAO(rt)}
	supply, err := q.TotalSupply()
	require.NoError(t, err)
	require.Equal(t, 0, big.NewInt(10_200).Cmp(new(big.Int).SetBytes(supply)))

	// Gov removes bob
	rt.caller = []byte{0x01}
	require.NoError(t, daoMut.Delete(bob))

	// Supply dropped by 200
	supply, err = q.TotalSupply()
	require.NoError(t, err)
	require.Equal(t, 0, big.NewInt(10_000).Cmp(new(big.Int).SetBytes(supply)))
}

// TestIntegration_DeleteByNonGov verifies that only gov can call Delete.
// Note: "delete" is not in ExecCall router, called directly.
func TestIntegration_DeleteByNonGov(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	bob := []byte{0x02}
	daoMut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	require.NoError(t, daoMut.Join(bob, big.NewInt(50).Bytes()))

	// Bob tries to delete themselves — requires gov
	rt.caller = bob
	err := daoMut.Delete(bob)
	require.ErrorIs(t, err, dao.ErrMustCallByGov)
}

// TestIntegration_AddTrack covers AddTrack and Tracks/Track queries.
func TestIntegration_AddTrack(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	mut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	newTrack := dao.TrackData{
		Name:               "new_track",
		PreparePeriod:      0,
		MaxDeciding:        50,
		ConfirmPeriod:      1,
		DecisionPeriod:     50,
		MinEnactmentPeriod: 0,
		DecisionDeposit:    big.NewInt(2).Bytes(),
		MaxBalance:         big.NewInt(500_000).Bytes(),
	}
	require.NoError(t, mut.ExecCall(&model.ContractCall{
		Method: "add_track",
		Args:   [][]byte{scaleBytes(t, newTrack)},
	}))

	q := dao.DaoQuery{DAO: *dao.NewDAO(rt)}
	raw, err := q.ExecQuery(&model.ContractCall{Method: "tracks", Args: nil})
	require.NoError(t, err)
	var tracks []dao.TrackData
	require.NoError(t, codec.Decode(raw, &tracks))
	require.Len(t, tracks, 2) // track 0 + track 1

	// Query track by id
	raw, err = q.ExecQuery(&model.ContractCall{
		Method: "track",
		Args:   [][]byte{scaleBytes(t, uint32(1))},
	})
	require.NoError(t, err)
	var opt util.Option[dao.TrackData]
	require.NoError(t, codec.Decode(raw, &opt))
	require.True(t, opt.IsSome())
	require.Equal(t, "new_track", string(opt.V.Name))

	// Non-existent track returns error (Track query API returns error for missing keys)
	_, err = q.ExecQuery(&model.ContractCall{
		Method: "track",
		Args:   [][]byte{scaleBytes(t, uint32(99))},
	})
	require.Error(t, err)
	require.ErrorIs(t, err, dao.ErrNoTrack)
}

// TestIntegration_SetDefaultTrack covers SetDefaultTrack and DefaultTrack query.
func TestIntegration_SetDefaultTrack(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	mut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	require.NoError(t, mut.ExecCall(&model.ContractCall{
		Method: "add_track",
		Args: [][]byte{scaleBytes(t, dao.TrackData{
			Name:               "t2",
			PreparePeriod:      0,
			MaxDeciding:        10,
			ConfirmPeriod:      1,
			DecisionPeriod:     10,
			MinEnactmentPeriod: 0,
			DecisionDeposit:    big.NewInt(5).Bytes(),
			MaxBalance:         big.NewInt(100_000).Bytes(),
		})},
	}))

	q := dao.DaoQuery{DAO: *dao.NewDAO(rt)}

	// Default is 0
	raw, err := q.ExecQuery(&model.ContractCall{Method: "default_track", Args: nil})
	require.NoError(t, err)
	var dt uint32
	require.NoError(t, codec.Decode(raw, &dt))
	require.Equal(t, uint32(0), dt)

	// Set to 1
	require.NoError(t, mut.ExecCall(&model.ContractCall{
		Method: "set_default_track",
		Args:   [][]byte{scaleBytes(t, uint32(1))},
	}))
	raw, err = q.ExecQuery(&model.ContractCall{Method: "default_track", Args: nil})
	require.NoError(t, err)
	require.NoError(t, codec.Decode(raw, &dt))
	require.Equal(t, uint32(1), dt)

	// Non-existent track fails
	err = mut.ExecCall(&model.ContractCall{
		Method: "set_default_track",
		Args:   [][]byte{scaleBytes(t, uint32(99))},
	})
	require.ErrorIs(t, err, dao.ErrNoTrack)
}

// TestIntegration_SubmitVote covers SubmitVote and Vote/Votes queries.
// SubmitProposal sets Pending; DepositProposal is required to advance to Ongoing before voting.
func TestIntegration_SubmitVote(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	daoMut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	daoQ := dao.DaoQuery{DAO: *dao.NewDAO(rt)}
	require.NoError(t, daoMut.SubmitProposal(dao.CallContent{Amount: big.NewInt(1).Bytes()}, 0))

	rt.height = 2
	require.NoError(t, daoMut.DepositProposal(0, big.NewInt(1).Bytes())) // advance to Ongoing
	require.NoError(t, daoMut.SubmitVote(0, true, big.NewInt(50).Bytes()))

	vote, err := daoQ.Vote(0)
	require.NoError(t, err)
	require.True(t, vote.IsSome())
	require.True(t, vote.V.OpinionYes)

	votes, err := daoQ.Votes()
	require.NoError(t, err)
	require.Len(t, votes, 1)
}

// TestIntegration_SubmitVote_FailsIfNotMember verifies that a non-member cannot vote.
func TestIntegration_SubmitVote_FailsIfNotMember(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	daoMut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	require.NoError(t, daoMut.SubmitProposal(dao.CallContent{Amount: big.NewInt(1).Bytes()}, 0))

	rt.height = 2
	rt.caller = []byte{0x99} // non-member
	err := daoMut.SubmitVote(0, true, big.NewInt(10).Bytes())
	require.ErrorIs(t, err, dao.ErrMemberNotExisted)
}

// TestIntegration_CancelVote covers CancelVote.
func TestIntegration_CancelVote(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	daoMut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	daoQ := dao.DaoQuery{DAO: *dao.NewDAO(rt)}
	require.NoError(t, daoMut.SubmitProposal(dao.CallContent{Amount: big.NewInt(1).Bytes()}, 0))

	rt.height = 2
	require.NoError(t, daoMut.DepositProposal(0, big.NewInt(1).Bytes())) // advance to Ongoing
	require.NoError(t, daoMut.SubmitVote(0, true, big.NewInt(50).Bytes()))
	require.NoError(t, daoMut.CancelVote(0))

	vote, err := daoQ.Vote(0)
	require.NoError(t, err)
	require.True(t, vote.IsSome())
	require.True(t, vote.V.Deleted)
}

// TestIntegration_Unlock covers Unlock after vote matures.
// Uses a separate track with MinEnactmentPeriod=1 so the unlock window is testable.
func TestIntegration_Unlock(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	daoMut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	daoQ := dao.DaoQuery{DAO: *dao.NewDAO(rt)}
	track := dao.TrackData{
		Name:               "test",
		PreparePeriod:      0,
		MaxDeciding:        100,
		ConfirmPeriod:      1,
		DecisionPeriod:     100,
		MinEnactmentPeriod: 1, // UnlockBlock=1, unlock at height >= end+1
		DecisionDeposit:    big.NewInt(1).Bytes(),
		MaxBalance:         big.NewInt(1_000_000).Bytes(),
	}
	require.NoError(t, daoMut.AddTrack(track))

	require.NoError(t, daoMut.SubmitProposal(dao.CallContent{Amount: big.NewInt(1).Bytes()}, 1))

	rt.height = 2
	require.NoError(t, daoMut.DepositProposal(0, big.NewInt(1).Bytes())) // advance to Ongoing
	require.NoError(t, daoMut.SubmitVote(0, true, big.NewInt(50).Bytes()))

	rt.height = 200 // past DecisionPeriod, confirmed
	require.NoError(t, daoMut.ExecProposal(0))

	// Unlock at height=103 (end=102, UnlockBlock=1): 103 >= 102+1 → succeeds
	rt.height = 103
	require.NoError(t, daoMut.Unlock(0))

	lock, err := daoQ.LockBalanceOf(rt.caller)
	require.NoError(t, err)
	require.Equal(t, 0, big.NewInt(0).Cmp(new(big.Int).SetBytes(lock)))
}

// TestIntegration_Unlock_FailsIfAlreadyUnlocked verifies double-unlock returns error.
func TestIntegration_Unlock_FailsIfAlreadyUnlocked(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	daoMut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	daoQ := dao.DaoQuery{DAO: *dao.NewDAO(rt)}
	track := dao.TrackData{
		Name:               "test",
		PreparePeriod:      0,
		MaxDeciding:        100,
		ConfirmPeriod:      1,
		DecisionPeriod:     100,
		MinEnactmentPeriod: 1, // UnlockBlock=1
		DecisionDeposit:    big.NewInt(1).Bytes(),
		MaxBalance:         big.NewInt(1_000_000).Bytes(),
	}
	require.NoError(t, daoMut.AddTrack(track))

	require.NoError(t, daoMut.SubmitProposal(dao.CallContent{Amount: big.NewInt(1).Bytes()}, 1))

	rt.height = 2
	require.NoError(t, daoMut.DepositProposal(0, big.NewInt(1).Bytes()))
	require.NoError(t, daoMut.SubmitVote(0, true, big.NewInt(50).Bytes()))

	rt.height = 200
	require.NoError(t, daoMut.ExecProposal(0))

	// First unlock at height=103 (end=102, UnlockBlock=1): 103 >= 103 → succeeds
	rt.height = 103
	require.NoError(t, daoMut.Unlock(0))

	// Second unlock fails
	err := daoMut.Unlock(0)
	require.ErrorIs(t, err, dao.ErrVoteAlreadyUnlocked)
	_ = daoQ // suppress unused
}

// TestIntegration_DepositProposal covers DepositProposal (Pending → Ongoing transition).
func TestIntegration_DepositProposal(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	mut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	require.NoError(t, mut.ExecCall(&model.ContractCall{
		Method: "submit_proposal",
		Args:   [][]byte{scaleBytes(t, dao.CallContent{Amount: big.NewInt(1).Bytes()}), scaleBytes(t, uint32(0))},
	}))

	q := dao.DaoQuery{DAO: *dao.NewDAO(rt)}
	raw, err := q.ExecQuery(&model.ContractCall{
		Method: "proposal_status",
		Args:   [][]byte{scaleBytes(t, uint32(0))},
	})
	require.NoError(t, err)
	var ps util.Option[dao.ProposalStatus]
	require.NoError(t, codec.Decode(raw, &ps))
	require.True(t, ps.IsSome())
	require.Equal(t, dao.ProposalPending, ps.V.State)

	rt.height = 2
	require.NoError(t, mut.ExecCall(&model.ContractCall{
		Method: "deposit_proposal",
		Args:   [][]byte{scaleBytes(t, uint32(0)), scaleBytes(t, big.NewInt(1).Bytes())},
	}))

	raw, err = q.ExecQuery(&model.ContractCall{
		Method: "proposal_status",
		Args:   [][]byte{scaleBytes(t, uint32(0))},
	})
	require.NoError(t, err)
	require.NoError(t, codec.Decode(raw, &ps))
	require.True(t, ps.IsSome())
	require.Equal(t, dao.ProposalOngoing, ps.V.State)
}

// TestIntegration_ExecProposal covers ExecProposal for an approved proposal.
func TestIntegration_ExecProposal(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	daoMut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	daoQ := dao.DaoQuery{DAO: *dao.NewDAO(rt)}
	track := dao.TrackData{
		Name:               "test",
		PreparePeriod:      0,
		MaxDeciding:        100,
		ConfirmPeriod:      1,
		DecisionPeriod:     100,
		MinEnactmentPeriod: 0,
		DecisionDeposit:    big.NewInt(1).Bytes(),
		MaxBalance:         big.NewInt(1_000_000).Bytes(),
	}
	require.NoError(t, daoMut.AddTrack(track))

	require.NoError(t, daoMut.SubmitProposal(dao.CallContent{Amount: big.NewInt(1).Bytes()}, 1))

	rt.height = 2
	require.NoError(t, daoMut.DepositProposal(0, big.NewInt(1).Bytes()))

	rt.height = 3
	require.NoError(t, daoMut.SubmitVote(0, true, big.NewInt(200).Bytes()))

	rt.height = 200 // > DecisionPeriod=100, confirmed
	require.NoError(t, daoMut.ExecProposal(0))

	prop, err := daoQ.Proposal(0)
	require.NoError(t, err)
	require.True(t, prop.IsSome())
	require.Equal(t, dao.ProposalApproved, prop.V.Status.State)
}

// TestIntegration_ExecProposal_NotConfirmed verifies exec fails when proposal is not approved.
func TestIntegration_ExecProposal_NotConfirmed(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	daoMut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	track := dao.TrackData{
		Name:               "test",
		PreparePeriod:      0,
		MaxDeciding:        100,
		ConfirmPeriod:      1,
		DecisionPeriod:     100,
		MinEnactmentPeriod: 0,
		DecisionDeposit:    big.NewInt(1).Bytes(),
		MaxBalance:         big.NewInt(1_000_000).Bytes(),
	}
	require.NoError(t, daoMut.AddTrack(track))

	require.NoError(t, daoMut.SubmitProposal(dao.CallContent{Amount: big.NewInt(1).Bytes()}, 1))

	rt.height = 2
	require.NoError(t, daoMut.DepositProposal(0, big.NewInt(1).Bytes()))

	// No votes → not confirmed. At height=200 > end, proposal rejected.
	rt.height = 200
	require.NoError(t, daoMut.ExecProposal(0)) // sets to Rejected

	// Already decided
	err := daoMut.ExecProposal(0)
	require.Error(t, err)
}

// TestIntegration_TransferEnabled covers Transfer after enabling transfers.
func TestIntegration_TransferEnabled(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	daoMut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	daoQ := dao.DaoQuery{DAO: *dao.NewDAO(rt)}
	alice, bob := []byte{0x0a}, []byte{0x0b}
	for _, acc := range [][]byte{alice, bob} {
		require.NoError(t, daoMut.Join(acc, big.NewInt(500).Bytes()))
	}

	require.NoError(t, daoMut.SetPublicJoin(true))

	rt.caller = alice
	require.NoError(t, daoMut.Transfer(bob, big.NewInt(100).Bytes()))

	bal, err := daoQ.BalanceOf(bob)
	require.NoError(t, err)
	require.Equal(t, 0, big.NewInt(600).Cmp(new(big.Int).SetBytes(bal)))
}

// TestIntegration_ApproveAndTransferFrom covers Approve + TransferFrom.
// Uses a dedicated track to avoid cross-test contamination.
func TestIntegration_ApproveAndTransferFrom(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	daoMut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	daoQ := dao.DaoQuery{DAO: *dao.NewDAO(rt)}
	alice, bob, charles := []byte{0x0a}, []byte{0x0b}, []byte{0x0c}
	for _, acc := range [][]byte{alice, bob, charles} {
		require.NoError(t, daoMut.Join(acc, big.NewInt(1000).Bytes()))
	}

	require.NoError(t, daoMut.SetPublicJoin(true))

	rt.caller = alice
	require.NoError(t, daoMut.Approve(charles, big.NewInt(300).Bytes()))

	allow, err := daoQ.Allowance(alice, charles)
	require.NoError(t, err)
	require.Equal(t, 0, big.NewInt(300).Cmp(new(big.Int).SetBytes(allow)))

	rt.caller = charles
	require.NoError(t, daoMut.TransferFrom(alice, bob, big.NewInt(200).Bytes()))

	abal, err := daoQ.BalanceOf(alice)
	require.NoError(t, err)
	require.Equal(t, 0, big.NewInt(800).Cmp(new(big.Int).SetBytes(abal)))

	allow, err = daoQ.Allowance(alice, charles)
	require.NoError(t, err)
	require.Equal(t, 0, big.NewInt(100).Cmp(new(big.Int).SetBytes(allow)))
}

// TestIntegration_Spend covers Spend (creates spend record + treasury proposal).
func TestIntegration_Spend(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	mut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	to := []byte{0x0b}
	require.NoError(t, mut.ExecCall(&model.ContractCall{
		Method: "spend",
		Args:   [][]byte{scaleBytes(t, to), scaleBytes(t, big.NewInt(100).Bytes()), scaleBytes(t, uint32(0))},
	}))

	q := dao.DaoQuery{DAO: *dao.NewDAO(rt)}
	raw, err := q.ExecQuery(&model.ContractCall{Method: "proposals", Args: nil})
	require.NoError(t, err)
	var props []dao.Proposal
	require.NoError(t, codec.Decode(raw, &props))
	require.Len(t, props, 1)
	require.Equal(t, dao.ProposalPending, props[0].Status.State)
}

// TestIntegration_Payout covers Payout: gov marks spend as paid.
func TestIntegration_Payout(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	mut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	require.NoError(t, mut.ExecCall(&model.ContractCall{
		Method: "spend",
		Args:   [][]byte{scaleBytes(t, []byte{0x0b}), scaleBytes(t, big.NewInt(100).Bytes()), scaleBytes(t, uint32(0))},
	}))

	require.NoError(t, mut.ExecCall(&model.ContractCall{
		Method: "payout",
		Args:   [][]byte{scaleBytes(t, uint64(0))},
	}))

	err := mut.ExecCall(&model.ContractCall{
		Method: "payout",
		Args:   [][]byte{scaleBytes(t, uint64(0))},
	})
	require.ErrorIs(t, err, dao.ErrSpendAlreadyExecuted)

	err = mut.ExecCall(&model.ContractCall{
		Method: "payout",
		Args:   [][]byte{scaleBytes(t, uint64(99))},
	})
	require.ErrorIs(t, err, dao.ErrSpendNotFound)
}

// TestIntegration_LockBalanceOf covers lock_balance_of after submitting a vote.
func TestIntegration_LockBalanceOf(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	mut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	require.NoError(t, mut.ExecCall(&model.ContractCall{
		Method: "submit_proposal",
		Args:   [][]byte{scaleBytes(t, dao.CallContent{Amount: big.NewInt(1).Bytes()}), scaleBytes(t, uint32(0))},
	}))
	rt.height = 2
	require.NoError(t, mut.ExecCall(&model.ContractCall{
		Method: "submit_vote",
		Args:   [][]byte{scaleBytes(t, uint32(0)), scaleBytes(t, true), scaleBytes(t, big.NewInt(50).Bytes())},
	}))

	q := dao.DaoQuery{DAO: *dao.NewDAO(rt)}
	raw, err := q.ExecQuery(&model.ContractCall{
		Method: "lock_balance_of",
		Args:   [][]byte{scaleBytes(t, rt.caller)},
	})
	require.NoError(t, err)
	var lock []byte
	require.NoError(t, codec.Decode(raw, &lock))
	require.Equal(t, 0, big.NewInt(50).Cmp(new(big.Int).SetBytes(lock)))
}

// TestIntegration_ProposalStatusQuery covers ProposalStatus for Pending and Ongoing.
func TestIntegration_ProposalStatusQuery(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	mut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	require.NoError(t, mut.ExecCall(&model.ContractCall{
		Method: "submit_proposal",
		Args:   [][]byte{scaleBytes(t, dao.CallContent{Amount: big.NewInt(1).Bytes()}), scaleBytes(t, uint32(0))},
	}))

	q := dao.DaoQuery{DAO: *dao.NewDAO(rt)}
	raw, err := q.ExecQuery(&model.ContractCall{
		Method: "proposal_status",
		Args:   [][]byte{scaleBytes(t, uint32(0))},
	})
	require.NoError(t, err)
	var ps util.Option[dao.ProposalStatus]
	require.NoError(t, codec.Decode(raw, &ps))
	require.True(t, ps.IsSome())
	require.Equal(t, dao.ProposalPending, ps.V.State)

	rt.height = 2
	require.NoError(t, mut.ExecCall(&model.ContractCall{
		Method: "deposit_proposal",
		Args:   [][]byte{scaleBytes(t, uint32(0)), scaleBytes(t, big.NewInt(1).Bytes())},
	}))

	raw, err = q.ExecQuery(&model.ContractCall{
		Method: "proposal_status",
		Args:   [][]byte{scaleBytes(t, uint32(0))},
	})
	require.NoError(t, err)
	require.NoError(t, codec.Decode(raw, &ps))
	require.True(t, ps.IsSome())
	require.Equal(t, dao.ProposalOngoing, ps.V.State)
}

// TestIntegration_ProposalQueryById covers single Proposal query.
func TestIntegration_ProposalQueryById(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	daoMut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	daoQ := dao.DaoQuery{DAO: *dao.NewDAO(rt)}
	require.NoError(t, daoMut.SubmitProposal(dao.CallContent{Amount: big.NewInt(1).Bytes()}, 0))

	prop, err := daoQ.Proposal(0)
	require.NoError(t, err)
	require.True(t, prop.IsSome())
	require.Equal(t, uint32(0), prop.V.ID)
	require.Equal(t, dao.ProposalPending, prop.V.Status.State)

	// Non-existent
	_, err = daoQ.Proposal(99)
	require.ErrorIs(t, err, dao.ErrInvalidProposal)
}

// TestIntegration_VoteQueryById covers Vote query and Votes list.
func TestIntegration_VoteQueryById(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	daoMut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	daoQ := dao.DaoQuery{DAO: *dao.NewDAO(rt)}
	require.NoError(t, daoMut.SubmitProposal(dao.CallContent{Amount: big.NewInt(1).Bytes()}, 0))

	rt.height = 2
	require.NoError(t, daoMut.DepositProposal(0, big.NewInt(1).Bytes())) // advance to Ongoing
	require.NoError(t, daoMut.SubmitVote(0, false, big.NewInt(30).Bytes()))

	v, err := daoQ.Vote(0)
	require.NoError(t, err)
	require.True(t, v.IsSome())
	require.False(t, v.V.OpinionYes)
	require.Equal(t, uint32(0), v.V.ProposalID)
}

// TestIntegration_MembersQuery covers Members, BalanceOf, non-member queries.
func TestIntegration_MembersQuery(t *testing.T) {
	rt, cleanup := daoWithGov(t)
	defer cleanup()

	q := dao.DaoQuery{DAO: *dao.NewDAO(rt)}
	raw, err := q.ExecQuery(&model.ContractCall{Method: "members", Args: nil})
	require.NoError(t, err)
	var members []dao.Member
	require.NoError(t, codec.Decode(raw, &members))
	require.Len(t, members, 1)

	bob := []byte{0x02}
	mut := dao.DaoMutation{DAO: *dao.NewDAO(rt)}
	require.NoError(t, mut.ExecCall(&model.ContractCall{
		Method: "join",
		Args:   [][]byte{scaleBytes(t, bob), scaleBytes(t, big.NewInt(100).Bytes())},
	}))

	raw, err = q.ExecQuery(&model.ContractCall{Method: "members", Args: nil})
	require.NoError(t, err)
	require.NoError(t, codec.Decode(raw, &members))
	require.Len(t, members, 2)

	// Balance of bob
	raw, err = q.ExecQuery(&model.ContractCall{
		Method: "balance_of",
		Args:   [][]byte{scaleBytes(t, bob)},
	})
	require.NoError(t, err)
	var bal []byte
	require.NoError(t, codec.Decode(raw, &bal))
	require.Equal(t, 0, big.NewInt(100).Cmp(new(big.Int).SetBytes(bal)))

	// Non-member balance
	raw, err = q.ExecQuery(&model.ContractCall{
		Method: "balance_of",
		Args:   [][]byte{scaleBytes(t, []byte{0xff})},
	})
	require.NoError(t, err)
	var empty []byte
	require.NoError(t, codec.Decode(raw, &empty))
	require.Nil(t, empty)

	// Lock of bob (no locks)
	raw, err = q.ExecQuery(&model.ContractCall{
		Method: "lock_balance_of",
		Args:   [][]byte{scaleBytes(t, bob)},
	})
	require.NoError(t, err)
	var lock []byte
	require.NoError(t, codec.Decode(raw, &lock))
	require.Equal(t, 0, big.NewInt(0).Cmp(new(big.Int).SetBytes(lock)))
}
