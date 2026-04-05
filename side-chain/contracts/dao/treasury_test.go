package dao

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDaoMutation_Spend_NotMember(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
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

	rt.caller = []byte{2}
	err := m.Spend([]byte{3}, big.NewInt(100).Bytes(), 0)
	require.ErrorIs(t, err, ErrMemberNotExisted)
}

func TestDaoMutation_Spend_Success(t *testing.T) {
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
	err := m.Spend([]byte{3}, big.NewInt(100).Bytes(), 0)
	require.NoError(t, err)

	// Check proposal created
	q := DaoQuery{DAO: m.DAO}
	proposals, err := q.Proposals()
	require.NoError(t, err)
	require.Len(t, proposals, 1)
}

func TestDaoMutation_Payout_NotGov(t *testing.T) {
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
	_ = m.Spend([]byte{3}, big.NewInt(100).Bytes(), 0)

	err := m.Payout(0)
	require.ErrorIs(t, err, ErrMustCallByGov)
}

func TestDaoMutation_Payout_NotFound(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	_ = m.Init([]Member{}, false, sudo, nil)

	rt.caller = sudo
	err := m.Payout(999)
	require.ErrorIs(t, err, ErrSpendNotFound)
}

func TestDaoMutation_Payout_AlreadyExecuted(t *testing.T) {
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
	_ = m.Spend([]byte{3}, big.NewInt(100).Bytes(), 0)

	rt.caller = sudo
	_ = m.Payout(0)

	err := m.Payout(0)
	require.ErrorIs(t, err, ErrSpendAlreadyExecuted)
}

func TestDaoMutation_Payout_Success(t *testing.T) {
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
	_ = m.Spend([]byte{3}, big.NewInt(100).Bytes(), 0)

	rt.caller = sudo
	err := m.Payout(0)
	require.NoError(t, err)
}
