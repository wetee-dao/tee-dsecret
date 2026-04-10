package gov

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func addr(b []byte) model.UniAddr {
	return model.UniAddr{T: 0, V: b}
}

func amount(i *big.Int) model.Amount {
	return model.Amount{Int: i}
}

func TestGovQuery_Members(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	d := GovQuery{Gov: *NewGov(rt)}
	members, err := d.Members()
	require.NoError(t, err)
	require.Empty(t, members)
}

func TestGovQuery_PublicJoin_Default(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	d := GovQuery{Gov: *NewGov(rt)}
	enabled, err := d.GetPublicJoin()
	require.NoError(t, err)
	require.False(t, enabled)
}

func TestGovMutation_Init_Basic(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := GovMutation{Gov: *NewGov(rt)}
	err := m.Init()
	require.NoError(t, err)

	enabled, err := GovQuery{Gov: m.Gov}.GetPublicJoin()
	require.NoError(t, err)
	require.True(t, enabled)
}

func TestGovMutation_Init_WithDefaultTrack(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := GovMutation{Gov: *NewGov(rt)}
	err := m.Init()
	require.NoError(t, err)

	dt, err := GovQuery{Gov: m.Gov}.DefaultTrack()
	require.NoError(t, err)
	require.Equal(t, uint32(0), dt)

	tr, err := GovQuery{Gov: m.Gov}.Track(0)
	require.NoError(t, err)
	require.True(t, tr.IsSome())
}

func TestGovMutation_PublicJoin_NotAllowed(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := GovMutation{Gov: *NewGov(rt)}
	rt.caller = addr([]byte{1})
	err := m.PublicJoin()
	require.ErrorIs(t, err, ErrPublicJoinNotAllowed)
}

func TestGovMutation_SetPublicJoin_NotGov(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := GovMutation{Gov: *NewGov(rt)}
	rt.caller = addr([]byte{1})
	err := m.SetPublicJoin(true)
	require.ErrorIs(t, err, ErrMustCallByGov)
}

func TestGovMutation_SetPublicJoin_ByGov(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := GovMutation{Gov: *NewGov(rt)}
	sudo := addr([]byte{1})
	_ = m.Init()
	rt.sudoAccount = sudo

	rt.caller = sudo
	err := m.SetPublicJoin(true)
	require.NoError(t, err)

	enabled, err := GovQuery{Gov: m.Gov}.GetPublicJoin()
	require.NoError(t, err)
	require.True(t, enabled)
}

func TestGovMutation_Join_NotGov(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := GovMutation{Gov: *NewGov(rt)}
	rt.caller = addr([]byte{1})
	err := m.Join(addr([]byte{2}), amount(big.NewInt(100)))
	require.ErrorIs(t, err, ErrMustCallByGov)
}

func TestGovMutation_Join_ByGov(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := GovMutation{Gov: *NewGov(rt)}
	sudo := addr([]byte{1})
	_ = m.Init()
	rt.sudoAccount = sudo

	rt.caller = sudo
	err := m.Join(addr([]byte{2}), amount(big.NewInt(100)))
	require.NoError(t, err)

	member, err := m.members.GetOrDefault(rt.txn, addr([]byte{2}), model.ZeroAmount)
	require.NoError(t, err)
	require.Equal(t, 0, member.Int.Cmp(big.NewInt(100)))
}

func TestGovMutation_Join_MemberExisted(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := GovMutation{Gov: *NewGov(rt)}
	sudo := addr([]byte{1})
	_ = m.Init()
	rt.sudoAccount = sudo

	rt.caller = sudo
	_ = m.Join(addr([]byte{2}), amount(big.NewInt(100)))

	err := m.Join(addr([]byte{2}), amount(big.NewInt(50)))
	require.ErrorIs(t, err, ErrMemberExisted)
}

func TestGovMutation_Leave_NotMember(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := GovMutation{Gov: *NewGov(rt)}
	rt.caller = addr([]byte{1})
	err := m.Leave()
	require.ErrorIs(t, err, ErrMemberNotExisted)
}

func TestGovMutation_Leave_BalanceNotZero(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := GovMutation{Gov: *NewGov(rt)}
	sudo := addr([]byte{1})
	_ = m.Init()
	rt.sudoAccount = sudo

	rt.caller = sudo
	_ = m.Join(addr([]byte{2}), amount(big.NewInt(100)))

	rt.caller = addr([]byte{2})
	err := m.Leave()
	require.ErrorIs(t, err, ErrMemberBalanceNotZero)
}

func TestGovMutation_Leave_Success(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := GovMutation{Gov: *NewGov(rt)}
	sudo := addr([]byte{1})
	_ = m.Init()
	rt.sudoAccount = sudo

	rt.caller = sudo
	_ = m.Join(addr([]byte{2}), model.ZeroAmount)

	rt.caller = addr([]byte{2})
	err := m.Leave()
	require.NoError(t, err)

	member, err := m.members.Get(rt.txn, addr([]byte{2}))
	require.NoError(t, err)
	require.Nil(t, member)
}

func TestGovMutation_PublicJoin_Success(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := GovMutation{Gov: *NewGov(rt)}
	sudo := addr([]byte{1})
	_ = m.Init()
	rt.sudoAccount = sudo

	rt.caller = sudo
	_ = m.SetPublicJoin(true)

	rt.caller = addr([]byte{2})
	err := m.PublicJoin()
	require.NoError(t, err)

	member, err := m.members.Get(rt.txn, addr([]byte{2}))
	require.NoError(t, err)
	require.NotNil(t, member)
}

func TestGovMutation_Delete_NotMember(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := GovMutation{Gov: *NewGov(rt)}
	sudo := addr([]byte{1})
	_ = m.Init()
	rt.sudoAccount = sudo

	rt.caller = sudo
	err := m.Delete(addr([]byte{2}))
	require.ErrorIs(t, err, ErrMemberNotExisted)
}

func TestGovMutation_Delete_ByGov(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := GovMutation{Gov: *NewGov(rt)}
	sudo := addr([]byte{1})
	_ = m.Init()
	rt.sudoAccount = sudo

	rt.caller = sudo
	_ = m.Join(addr([]byte{2}), amount(big.NewInt(100)))

	err := m.Delete(addr([]byte{2}))
	require.NoError(t, err)

	member, err := m.members.Get(rt.txn, addr([]byte{2}))
	require.NoError(t, err)
	require.Nil(t, member)

	supply, err := GovQuery{Gov: m.Gov}.TotalSupply()
	require.NoError(t, err)
	require.Equal(t, 0, supply.Int.Sign())
}

func TestGovMutation_LeaveWithBurn_Success(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := GovMutation{Gov: *NewGov(rt)}
	sudo := addr([]byte{1})
	_ = m.Init()
	rt.sudoAccount = sudo

	rt.caller = sudo
	_ = m.Join(addr([]byte{2}), amount(big.NewInt(100)))

	rt.caller = addr([]byte{2})
	err := m.LeaveWithBurn()
	require.NoError(t, err)

	member, err := m.members.Get(rt.txn, addr([]byte{2}))
	require.NoError(t, err)
	require.Nil(t, member)

	supply, err := GovQuery{Gov: m.Gov}.TotalSupply()
	require.NoError(t, err)
	require.Equal(t, 0, supply.Int.Sign())
}
