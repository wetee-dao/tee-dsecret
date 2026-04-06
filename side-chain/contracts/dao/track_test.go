package dao

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDaoQuery_Tracks_Empty(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	q := DaoQuery{DAO: *NewDAO(rt)}
	tracks, err := q.Tracks()
	require.NoError(t, err)
	require.Empty(t, tracks)
}

func TestDaoQuery_Track_NotFound(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	q := DaoQuery{DAO: *NewDAO(rt)}
	result, err := q.Track(999)
	// Track returns ErrNoTrack when not found, not None
	require.ErrorIs(t, err, ErrNoTrack)
	require.True(t, result.IsNone())
}

func TestDaoQuery_DefaultTrack_NotSet(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	q := DaoQuery{DAO: *NewDAO(rt)}
	dt, err := q.DefaultTrack()
	require.NoError(t, err)
	require.Nil(t, dt)
}

func TestDaoMutation_AddTrack_NotGov(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	rt.caller = []byte{1}

	track := TrackData{
		Name:            "test",
		PreparePeriod:   10,
		MaxDeciding:     100,
		ConfirmPeriod:   20,
		DecisionPeriod:  30,
		DecisionDeposit: big.NewInt(50).Bytes(),
		MaxBalance:      big.NewInt(1000).Bytes(),
	}
	err := m.AddTrack(track)
	require.ErrorIs(t, err, ErrMustCallByGov)
}

func TestDaoMutation_AddTrack_Success(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	_ = m.Init()
	rt.sudoAccount = sudo

	rt.caller = sudo
	track := TrackData{
		Name:            "test",
		PreparePeriod:   10,
		MaxDeciding:     100,
		ConfirmPeriod:   20,
		DecisionPeriod:  30,
		DecisionDeposit: big.NewInt(50).Bytes(),
		MaxBalance:      big.NewInt(1000).Bytes(),
	}
	err := m.AddTrack(track)
	require.NoError(t, err)

	q := DaoQuery{DAO: m.DAO}
	result, err := q.Track(0)
	require.NoError(t, err)
	require.True(t, result.IsSome())
	tr, err := result.UnWrap()
	require.NoError(t, err)
	require.Equal(t, "test", tr.Name)
}

func TestDaoMutation_SetDefaultTrack_NotGov(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	_ = m.Init()
	rt.sudoAccount = []byte{1}

	rt.caller = []byte{2}
	err := m.SetDefaultTrack(0)
	require.ErrorIs(t, err, ErrMustCallByGov)
}

func TestDaoMutation_SetDefaultTrack_NoTrack(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	_ = m.Init()
	rt.sudoAccount = sudo

	rt.caller = sudo
	err := m.SetDefaultTrack(999)
	require.ErrorIs(t, err, ErrNoTrack)
}

func TestDaoMutation_SetDefaultTrack_Success(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	_ = m.Init()
	rt.sudoAccount = sudo

	rt.caller = sudo
	err := m.SetDefaultTrack(0)
	require.NoError(t, err)

	q := DaoQuery{DAO: m.DAO}
	dt, err := q.DefaultTrack()
	require.NoError(t, err)
	require.NotNil(t, dt)
	require.Equal(t, uint32(0), *dt)
}

func TestDaoQuery_Track_Success(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	_ = m.Init()
	rt.sudoAccount = sudo

	q := DaoQuery{DAO: m.DAO}
	result, err := q.Track(0)
	require.NoError(t, err)
	require.True(t, result.IsSome())
	unwrapped, _ := result.UnWrap()
	require.Equal(t, "default", unwrapped.Name)
	require.Equal(t, uint32(10), unwrapped.PreparePeriod)
}

func TestDaoQuery_Tracks_Multiple(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	_ = m.Init()
	rt.sudoAccount = sudo

	rt.caller = sudo
	track1 := TrackData{Name: "track1", PreparePeriod: 10}
	track2 := TrackData{Name: "track2", PreparePeriod: 20}
	_ = m.AddTrack(track1)
	_ = m.AddTrack(track2)

	q := DaoQuery{DAO: m.DAO}
	tracks, err := q.Tracks()
	require.NoError(t, err)
	require.Len(t, tracks, 2)
}
