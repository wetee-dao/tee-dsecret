package dao

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDaoQuery_Vote_NotFound(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	q := DaoQuery{DAO: *NewDAO(rt)}
	result, err := q.Vote(999)
	// Vote returns ErrInvalidVote when not found, not None
	require.ErrorIs(t, err, ErrInvalidVote)
	require.True(t, result.IsNone())
}

func TestDaoQuery_Votes_Empty(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	q := DaoQuery{DAO: *NewDAO(rt)}
	votes, err := q.Votes()
	require.NoError(t, err)
	require.Empty(t, votes)
}

func TestDaoMutation_SubmitVote_NotMember(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	track := TrackData{
		Name:            "test",
		PreparePeriod:    0,
		MaxDeciding:      100,
		ConfirmPeriod:    20,
		DecisionPeriod:   30,
		DecisionDeposit:  big.NewInt(50).Bytes(),
		MaxBalance:       big.NewInt(1000).Bytes(),
	}
	_ = m.Init([]Member{}, false, []byte{1}, &track)

	rt.caller = []byte{2}
	err := m.SubmitVote(0, true, big.NewInt(10).Bytes())
	require.ErrorIs(t, err, ErrMemberNotExisted)
}

func TestDaoMutation_SubmitVote_ProposalNotOngoing(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	members := []Member{{Account: []byte{2}, Balance: big.NewInt(100).Bytes()}}
	track := TrackData{
		Name:            "test",
		PreparePeriod:    0,
		MaxDeciding:      100,
		ConfirmPeriod:    20,
		DecisionPeriod:   30,
		DecisionDeposit:  big.NewInt(50).Bytes(),
		MaxBalance:       big.NewInt(1000).Bytes(),
	}
	_ = m.Init(members, false, sudo, &track)

	rt.caller = []byte{2}
	_ = m.SubmitProposal(CallContent{Amount: big.NewInt(10).Bytes()}, 0)

	err := m.SubmitVote(0, true, big.NewInt(10).Bytes()) // proposal still pending
	require.ErrorIs(t, err, ErrPropNotOngoing)
}

func TestDaoMutation_SubmitVote_LowBalance(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	members := []Member{{Account: []byte{2}, Balance: big.NewInt(100).Bytes()}}
	track := TrackData{
		Name:            "test",
		PreparePeriod:    0,
		MaxDeciding:      100,
		ConfirmPeriod:    20,
		DecisionPeriod:   30,
		DecisionDeposit:  big.NewInt(50).Bytes(),
		MaxBalance:       big.NewInt(1000).Bytes(),
	}
	_ = m.Init(members, false, sudo, &track)

	rt.caller = []byte{2}
	_ = m.SubmitProposal(CallContent{Amount: big.NewInt(10).Bytes()}, 0)
	rt.height = 1 // advance height to pass prepare period
	_ = m.DepositProposal(0, big.NewInt(50).Bytes())

	err := m.SubmitVote(0, true, big.NewInt(200).Bytes()) // more than balance
	require.ErrorIs(t, err, ErrLowBalance)
}

func TestDaoMutation_SubmitVote_Success(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	members := []Member{{Account: []byte{2}, Balance: big.NewInt(100).Bytes()}}
	track := TrackData{
		Name:            "test",
		PreparePeriod:    0,
		MaxDeciding:      100,
		ConfirmPeriod:    20,
		DecisionPeriod:   30,
		DecisionDeposit:  big.NewInt(50).Bytes(),
		MaxBalance:       big.NewInt(1000).Bytes(),
	}
	_ = m.Init(members, false, sudo, &track)

	rt.caller = []byte{2}
	_ = m.SubmitProposal(CallContent{Amount: big.NewInt(10).Bytes()}, 0)
	rt.height = 1
	_ = m.DepositProposal(0, big.NewInt(50).Bytes())

	err := m.SubmitVote(0, true, big.NewInt(10).Bytes())
	require.NoError(t, err)

	q := DaoQuery{DAO: m.DAO}
	votes, err := q.Votes()
	require.NoError(t, err)
	require.Len(t, votes, 1)
	require.True(t, votes[0].OpinionYes)
}

func TestDaoMutation_CancelVote_NotOwner(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	members := []Member{
		{Account: []byte{2}, Balance: big.NewInt(100).Bytes()},
		{Account: []byte{3}, Balance: big.NewInt(100).Bytes()},
	}
	track := TrackData{
		Name:            "test",
		PreparePeriod:    0,
		MaxDeciding:      100,
		ConfirmPeriod:    20,
		DecisionPeriod:   30,
		DecisionDeposit:  big.NewInt(50).Bytes(),
		MaxBalance:       big.NewInt(1000).Bytes(),
	}
	_ = m.Init(members, false, sudo, &track)

	rt.caller = []byte{2}
	_ = m.SubmitProposal(CallContent{Amount: big.NewInt(10).Bytes()}, 0)
	rt.height = 1
	_ = m.DepositProposal(0, big.NewInt(50).Bytes())
	_ = m.SubmitVote(0, true, big.NewInt(10).Bytes())

	rt.caller = []byte{3} // different caller
	err := m.CancelVote(0)
	require.ErrorIs(t, err, ErrInvalidVoteUser)
}

func TestDaoMutation_CancelVote_Success(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	members := []Member{{Account: []byte{2}, Balance: big.NewInt(100).Bytes()}}
	track := TrackData{
		Name:            "test",
		PreparePeriod:    0,
		MaxDeciding:      100,
		ConfirmPeriod:    20,
		DecisionPeriod:   30,
		DecisionDeposit:  big.NewInt(50).Bytes(),
		MaxBalance:       big.NewInt(1000).Bytes(),
	}
	_ = m.Init(members, false, sudo, &track)

	rt.caller = []byte{2}
	_ = m.SubmitProposal(CallContent{Amount: big.NewInt(10).Bytes()}, 0)
	rt.height = 1
	_ = m.DepositProposal(0, big.NewInt(50).Bytes())
	_ = m.SubmitVote(0, true, big.NewInt(10).Bytes())

	err := m.CancelVote(0)
	require.NoError(t, err)

	q := DaoQuery{DAO: m.DAO}
	result, err := q.Vote(0)
	require.NoError(t, err)
	require.True(t, result.IsSome())
	v, err := result.UnWrap()
	require.NoError(t, err)
	require.True(t, v.Deleted)
}

func TestDaoMutation_Unlock_AlreadyUnlocked(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	members := []Member{{Account: []byte{2}, Balance: big.NewInt(100).Bytes()}}
	track := TrackData{
		Name:               "test",
		PreparePeriod:       0,
		MaxDeciding:         100,
		ConfirmPeriod:       20,
		DecisionPeriod:      30,
		MinEnactmentPeriod:  10,
		DecisionDeposit:     big.NewInt(50).Bytes(),
		MaxBalance:          big.NewInt(1000).Bytes(),
	}
	_ = m.Init(members, false, sudo, &track)

	rt.caller = []byte{2}
	_ = m.SubmitProposal(CallContent{Amount: big.NewInt(10).Bytes()}, 0)
	rt.height = 1
	_ = m.DepositProposal(0, big.NewInt(50).Bytes())
	_ = m.SubmitVote(0, true, big.NewInt(10).Bytes())

	// Mark as unlocked
	_ = m.voteUnlocks.Set(rt.txn, 0, true)

	err := m.Unlock(0)
	require.ErrorIs(t, err, ErrVoteAlreadyUnlocked)
}

func TestDaoMutation_Unlock_InvalidVoteStatus(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	members := []Member{{Account: []byte{2}, Balance: big.NewInt(100).Bytes()}}
	track := TrackData{
		Name:               "test",
		PreparePeriod:       0,
		MaxDeciding:         100,
		ConfirmPeriod:       20,
		DecisionPeriod:      30,
		MinEnactmentPeriod:  10,
		DecisionDeposit:     big.NewInt(50).Bytes(),
		MaxBalance:          big.NewInt(1000).Bytes(),
	}
	_ = m.Init(members, false, sudo, &track)

	rt.caller = []byte{2}
	_ = m.SubmitProposal(CallContent{Amount: big.NewInt(10).Bytes()}, 0)
	rt.height = 1
	_ = m.DepositProposal(0, big.NewInt(50).Bytes())
	_ = m.SubmitVote(0, true, big.NewInt(10).Bytes())
	_ = m.CancelVote(0) // deleted

	err := m.Unlock(0)
	require.ErrorIs(t, err, ErrInvalidVoteStatus)
}

func TestDaoMutation_Unlock_InvalidUser(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	members := []Member{
		{Account: []byte{2}, Balance: big.NewInt(100).Bytes()},
		{Account: []byte{3}, Balance: big.NewInt(100).Bytes()},
	}
	track := TrackData{
		Name:               "test",
		PreparePeriod:       0,
		MaxDeciding:         100,
		ConfirmPeriod:       20,
		DecisionPeriod:      30,
		MinEnactmentPeriod:  10,
		DecisionDeposit:     big.NewInt(50).Bytes(),
		MaxBalance:          big.NewInt(1000).Bytes(),
	}
	_ = m.Init(members, false, sudo, &track)

	rt.caller = []byte{2}
	_ = m.SubmitProposal(CallContent{Amount: big.NewInt(10).Bytes()}, 0)
	rt.height = 1
	_ = m.DepositProposal(0, big.NewInt(50).Bytes())
	_ = m.SubmitVote(0, true, big.NewInt(10).Bytes())

	rt.caller = []byte{3} // different caller
	err := m.Unlock(0)
	require.ErrorIs(t, err, ErrInvalidVoteUser)
}

func TestDaoMutation_Unlock_Success(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	members := []Member{{Account: []byte{2}, Balance: big.NewInt(100).Bytes()}}
	track := TrackData{
		Name:               "test",
		PreparePeriod:       0,
		MaxDeciding:         100,
		ConfirmPeriod:       20,
		DecisionPeriod:      30,
		MinEnactmentPeriod:  10,
		DecisionDeposit:     big.NewInt(50).Bytes(),
		MaxBalance:          big.NewInt(1000).Bytes(),
	}
	_ = m.Init(members, false, sudo, &track)

	rt.caller = []byte{2}
	_ = m.SubmitProposal(CallContent{Amount: big.NewInt(10).Bytes()}, 0)
	rt.height = 1
	_ = m.DepositProposal(0, big.NewInt(50).Bytes())
	_ = m.SubmitVote(0, true, big.NewInt(10).Bytes())

	// Advance height past MaxDeciding so proposal will be rejected
	// Deposit.Block + MaxDeciding = 1 + 100 = 101
	// Need height > 101 to reject
	rt.height = 200
	err := m.ExecProposal(0) // This will reject the proposal
	require.NoError(t, err)

	// Now proposal is rejected, we can unlock
	rt.height = 300
	rt.caller = []byte{2}
	err = m.Unlock(0)
	require.NoError(t, err)

	// Verify vote is unlocked
	unlocked, err := m.voteUnlocks.Get(rt.txn, 0)
	require.NoError(t, err)
	require.NotNil(t, unlocked)
	require.True(t, *unlocked)
}

func TestDaoMutation_SubmitVote_MultipleVotes(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	members := []Member{
		{Account: []byte{2}, Balance: big.NewInt(100).Bytes()},
		{Account: []byte{3}, Balance: big.NewInt(100).Bytes()},
	}
	track := TrackData{
		Name:            "test",
		PreparePeriod:    0,
		MaxDeciding:      100,
		ConfirmPeriod:    20,
		DecisionPeriod:   30,
		DecisionDeposit:  big.NewInt(50).Bytes(),
		MaxBalance:       big.NewInt(1000).Bytes(),
	}
	_ = m.Init(members, false, sudo, &track)

	rt.caller = []byte{2}
	_ = m.SubmitProposal(CallContent{Amount: big.NewInt(10).Bytes()}, 0)
	rt.height = 1
	_ = m.DepositProposal(0, big.NewInt(50).Bytes())

	rt.caller = []byte{2}
	_ = m.SubmitVote(0, true, big.NewInt(10).Bytes())

	rt.caller = []byte{3}
	err := m.SubmitVote(0, false, big.NewInt(20).Bytes())
	require.NoError(t, err)

	q := DaoQuery{DAO: m.DAO}
	votes, err := q.Votes()
	require.NoError(t, err)
	require.Len(t, votes, 2)
}
