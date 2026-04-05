package dao

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDaoQuery_Proposal_NotFound(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	q := DaoQuery{DAO: *NewDAO(rt)}
	_, err := q.Proposal(999)
	require.NoError(t, err) // Returns None, not error
}

func TestDaoQuery_Proposals_Empty(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	q := DaoQuery{DAO: *NewDAO(rt)}
	proposals, err := q.Proposals()
	require.NoError(t, err)
	require.Empty(t, proposals)
}

func TestDaoQuery_ProposalStatus_NotFound(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	q := DaoQuery{DAO: *NewDAO(rt)}
	_, err := q.ProposalStatus(999)
	require.ErrorIs(t, err, ErrInvalidProposal)
}

func TestDaoMutation_SubmitProposal_NotMember(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	rt.caller = []byte{1}

	track := TrackData{
		Name:            "test",
		PreparePeriod:    10,
		MaxDeciding:      100,
		ConfirmPeriod:    20,
		DecisionPeriod:   30,
		DecisionDeposit:  big.NewInt(50).Bytes(),
		MaxBalance:       big.NewInt(1000).Bytes(),
	}
	_ = m.Init([]Member{}, false, []byte{1}, &track)

	err := m.SubmitProposal(CallContent{Amount: big.NewInt(10).Bytes()}, 0)
	require.ErrorIs(t, err, ErrMemberNotExisted)
}

func TestDaoMutation_SubmitProposal_NoTrack(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	members := []Member{{Account: []byte{2}, Balance: big.NewInt(100).Bytes()}}
	_ = m.Init(members, false, sudo, nil)

	rt.caller = []byte{2}
	err := m.SubmitProposal(CallContent{Amount: big.NewInt(10).Bytes()}, 0)
	require.ErrorIs(t, err, ErrNoTrack)
}

func TestDaoMutation_SubmitProposal_Success(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	members := []Member{{Account: []byte{2}, Balance: big.NewInt(100).Bytes()}}
	track := TrackData{
		Name:            "test",
		PreparePeriod:    10,
		MaxDeciding:      100,
		ConfirmPeriod:    20,
		DecisionPeriod:   30,
		DecisionDeposit:  big.NewInt(50).Bytes(),
		MaxBalance:       big.NewInt(1000).Bytes(),
	}
	_ = m.Init(members, false, sudo, &track)

	rt.caller = []byte{2}
	err := m.SubmitProposal(CallContent{Amount: big.NewInt(10).Bytes()}, 0)
	require.NoError(t, err)

	// Check proposal created
	q := DaoQuery{DAO: m.DAO}
	proposals, err := q.Proposals()
	require.NoError(t, err)
	require.Len(t, proposals, 1)
	require.Equal(t, uint32(0), proposals[0].ID)
}

func TestDaoMutation_CancelProposal_NotOwner(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	members := []Member{{Account: []byte{2}, Balance: big.NewInt(100).Bytes()}}
	track := TrackData{
		Name:            "test",
		PreparePeriod:    10,
		MaxDeciding:      100,
		ConfirmPeriod:    20,
		DecisionPeriod:   30,
		DecisionDeposit:  big.NewInt(50).Bytes(),
		MaxBalance:       big.NewInt(1000).Bytes(),
	}
	_ = m.Init(members, false, sudo, &track)

	rt.caller = []byte{2}
	_ = m.SubmitProposal(CallContent{Amount: big.NewInt(10).Bytes()}, 0)

	rt.caller = []byte{1} // different caller
	err := m.CancelProposal(0)
	require.ErrorIs(t, err, ErrInvalidProposalCaller)
}

func TestDaoMutation_CancelProposal_Success(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	members := []Member{{Account: []byte{2}, Balance: big.NewInt(100).Bytes()}}
	track := TrackData{
		Name:            "test",
		PreparePeriod:    10,
		MaxDeciding:      100,
		ConfirmPeriod:    20,
		DecisionPeriod:   30,
		DecisionDeposit:  big.NewInt(50).Bytes(),
		MaxBalance:       big.NewInt(1000).Bytes(),
	}
	_ = m.Init(members, false, sudo, &track)

	rt.caller = []byte{2}
	_ = m.SubmitProposal(CallContent{Amount: big.NewInt(10).Bytes()}, 0)

	err := m.CancelProposal(0)
	require.NoError(t, err)

	q := DaoQuery{DAO: m.DAO}
	prop, err := q.Proposal(0)
	require.NoError(t, err)
	require.True(t, prop.IsSome())
	p, err := prop.UnWrap()
	require.NoError(t, err)
	require.Equal(t, ProposalCanceled, p.Status.State)
}

func TestDaoMutation_DepositProposal_InvalidStatus(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	members := []Member{{Account: []byte{2}, Balance: big.NewInt(100).Bytes()}}
	track := TrackData{
		Name:            "test",
		PreparePeriod:    10,
		MaxDeciding:      100,
		ConfirmPeriod:    20,
		DecisionPeriod:   30,
		DecisionDeposit:  big.NewInt(50).Bytes(),
		MaxBalance:       big.NewInt(1000).Bytes(),
	}
	_ = m.Init(members, false, sudo, &track)

	rt.caller = []byte{2}
	_ = m.SubmitProposal(CallContent{Amount: big.NewInt(10).Bytes()}, 0)
	_ = m.CancelProposal(0)

	err := m.DepositProposal(0, big.NewInt(50).Bytes())
	require.ErrorIs(t, err, ErrInvalidProposalStatus)
}

func TestDaoMutation_DepositProposal_InvalidDeposit(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	members := []Member{{Account: []byte{2}, Balance: big.NewInt(100).Bytes()}}
	track := TrackData{
		Name:            "test",
		PreparePeriod:    0, // 0 so we can deposit immediately
		MaxDeciding:      100,
		ConfirmPeriod:    20,
		DecisionPeriod:   30,
		DecisionDeposit:  big.NewInt(50).Bytes(),
		MaxBalance:       big.NewInt(1000).Bytes(),
	}
	_ = m.Init(members, false, sudo, &track)

	rt.caller = []byte{2}
	_ = m.SubmitProposal(CallContent{Amount: big.NewInt(10).Bytes()}, 0)

	// Advance height past prepare period to allow deposit
	rt.height = 1

	err := m.DepositProposal(0, big.NewInt(10).Bytes()) // less than DecisionDeposit
	require.ErrorIs(t, err, ErrInvalidDeposit)
}

func TestDaoMutation_ExecProposal_NotOngoing(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	members := []Member{{Account: []byte{2}, Balance: big.NewInt(100).Bytes()}}
	track := TrackData{
		Name:            "test",
		PreparePeriod:    10,
		MaxDeciding:      100,
		ConfirmPeriod:    20,
		DecisionPeriod:   30,
		DecisionDeposit:  big.NewInt(50).Bytes(),
		MaxBalance:       big.NewInt(1000).Bytes(),
	}
	_ = m.Init(members, false, sudo, &track)

	rt.caller = []byte{2}
	_ = m.SubmitProposal(CallContent{Amount: big.NewInt(10).Bytes()}, 0)

	err := m.ExecProposal(0) // still pending, not ongoing
	require.ErrorIs(t, err, ErrPropNotOngoing)
}


