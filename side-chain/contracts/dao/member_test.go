package dao

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDaoQuery_Members(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	d := DaoQuery{DAO: *NewDAO(rt)}
	members, err := d.Members()
	require.NoError(t, err)
	require.Empty(t, members)
}

func TestDaoQuery_PublicJoin_Default(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	d := DaoQuery{DAO: *NewDAO(rt)}
	enabled, err := d.GetPublicJoin()
	require.NoError(t, err)
	require.False(t, enabled)
}

func TestDaoMutation_Init_Basic(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	err := m.Init()
	require.NoError(t, err)

	// Check public join enabled by default
	enabled, err := DaoQuery{DAO: m.DAO}.GetPublicJoin()
	require.NoError(t, err)
	require.True(t, enabled)
}

func TestDaoMutation_Init_WithDefaultTrack(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	err := m.Init()
	require.NoError(t, err)

	// Check default track set
	dt, err := DaoQuery{DAO: m.DAO}.DefaultTrack()
	require.NoError(t, err)
	require.NotNil(t, dt)
	require.Equal(t, uint32(0), *dt)

	// Check track exists
	tr, err := DaoQuery{DAO: m.DAO}.Track(0)
	require.NoError(t, err)
	require.True(t, tr.IsSome())
}

func TestDaoMutation_PublicJoin_NotAllowed(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	rt.caller = []byte{1}
	err := m.PublicJoin()
	require.ErrorIs(t, err, ErrPublicJoinNotAllowed)
}

func TestDaoMutation_SetPublicJoin_NotGov(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	rt.caller = []byte{1}
	err := m.SetPublicJoin(true)
	require.ErrorIs(t, err, ErrMustCallByGov)
}

func TestDaoMutation_SetPublicJoin_ByGov(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	_ = m.Init()
	rt.sudoAccount = sudo

	rt.caller = sudo
	err := m.SetPublicJoin(true)
	require.NoError(t, err)

	enabled, err := DaoQuery{DAO: m.DAO}.GetPublicJoin()
	require.NoError(t, err)
	require.True(t, enabled)
}

func TestDaoMutation_Join_NotGov(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	rt.caller = []byte{1}
	err := m.Join([]byte{2}, big.NewInt(100).Bytes())
	require.ErrorIs(t, err, ErrMustCallByGov)
}

func TestDaoMutation_Join_ByGov(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	_ = m.Init()
	rt.sudoAccount = sudo

	rt.caller = sudo
	err := m.Join([]byte{2}, big.NewInt(100).Bytes())
	require.NoError(t, err)

	member, err := m.member([]byte{2})
	require.NoError(t, err)
	require.NotNil(t, member)
	require.Equal(t, 0, cmp(member.Balance, big.NewInt(100).Bytes()))
}

func TestDaoMutation_Join_MemberExisted(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	_ = m.Init()
	rt.sudoAccount = sudo

	// First join a member
	rt.caller = sudo
	_ = m.Join([]byte{2}, big.NewInt(100).Bytes())

	// Try to join again
	err := m.Join([]byte{2}, big.NewInt(50).Bytes())
	require.ErrorIs(t, err, ErrMemberExisted)
}

func TestDaoMutation_Leave_NotMember(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	rt.caller = []byte{1}
	err := m.Leave()
	require.ErrorIs(t, err, ErrMemberNotExisted)
}

func TestDaoMutation_Leave_BalanceNotZero(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	_ = m.Init()
	rt.sudoAccount = sudo

	// Join a member with balance
	rt.caller = sudo
	_ = m.Join([]byte{2}, big.NewInt(100).Bytes())

	rt.caller = []byte{2}
	err := m.Leave()
	require.ErrorIs(t, err, ErrMemberBalanceNotZero)
}

func TestDaoMutation_Leave_Success(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	_ = m.Init()
	rt.sudoAccount = sudo

	// Join a member with zero balance
	rt.caller = sudo
	_ = m.Join([]byte{2}, nil)

	rt.caller = []byte{2}
	err := m.Leave()
	require.NoError(t, err)

	member, err := m.member([]byte{2})
	require.NoError(t, err)
	require.Nil(t, member)
}

func TestDaoMutation_PublicJoin_Success(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	_ = m.Init()
	rt.sudoAccount = sudo

	rt.caller = []byte{2}
	err := m.PublicJoin()
	require.NoError(t, err)

	member, err := m.member([]byte{2})
	require.NoError(t, err)
	require.NotNil(t, member)
}

func TestDaoMutation_Delete_NotMember(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	err := m.Delete([]byte{1})
	require.ErrorIs(t, err, ErrMemberNotExisted)
}

func TestDaoMutation_Delete_ByGov(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	_ = m.Init()
	rt.sudoAccount = sudo

	// Join a member
	rt.caller = sudo
	_ = m.Join([]byte{2}, big.NewInt(100).Bytes())

	err := m.Delete([]byte{2})
	require.NoError(t, err)

	member, err := m.member([]byte{2})
	require.NoError(t, err)
	require.Nil(t, member)

	supply, err := DaoQuery{DAO: m.DAO}.TotalSupply()
	require.NoError(t, err)
	require.True(t, isZero(supply))
}

func TestDaoMutation_LeaveWithBurn_Success(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	_ = m.Init()
	rt.sudoAccount = sudo

	// Join a member with balance
	rt.caller = sudo
	_ = m.Join([]byte{2}, big.NewInt(100).Bytes())

	rt.caller = []byte{2}
	err := m.LeaveWithBurn()
	require.NoError(t, err)

	member, err := m.member([]byte{2})
	require.NoError(t, err)
	require.Nil(t, member)

	supply, err := DaoQuery{DAO: m.DAO}.TotalSupply()
	require.NoError(t, err)
	require.True(t, isZero(supply))
}
