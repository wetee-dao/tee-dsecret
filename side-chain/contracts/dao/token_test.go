package dao

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDaoQuery_TotalSupply_Empty(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	q := DaoQuery{DAO: *NewDAO(rt)}
	supply, err := q.TotalSupply()
	require.NoError(t, err)
	require.Nil(t, supply)
}

func TestDaoQuery_TotalSupply_AfterInit(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	members := []Member{
		{Account: []byte{1}, Balance: big.NewInt(100).Bytes()},
		{Account: []byte{2}, Balance: big.NewInt(200).Bytes()},
	}
	_ = m.Init(members, false, []byte{1}, nil)

	q := DaoQuery{DAO: m.DAO}
	supply, err := q.TotalSupply()
	require.NoError(t, err)
	require.Equal(t, 0, cmp(supply, big.NewInt(300).Bytes()))
}

func TestDaoQuery_BalanceOf_NotMember(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	q := DaoQuery{DAO: *NewDAO(rt)}
	balance, err := q.BalanceOf([]byte{1})
	require.NoError(t, err)
	require.Nil(t, balance)
}

func TestDaoQuery_BalanceOf_Success(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	members := []Member{{Account: []byte{1}, Balance: big.NewInt(100).Bytes()}}
	_ = m.Init(members, false, []byte{2}, nil)

	q := DaoQuery{DAO: m.DAO}
	balance, err := q.BalanceOf([]byte{1})
	require.NoError(t, err)
	require.Equal(t, 0, cmp(balance, big.NewInt(100).Bytes()))
}

func TestDaoQuery_LockBalanceOf_NotMember(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	q := DaoQuery{DAO: *NewDAO(rt)}
	lock, err := q.LockBalanceOf([]byte{1})
	require.NoError(t, err)
	require.Nil(t, lock)
}

func TestDaoQuery_Allowance_NotSet(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	q := DaoQuery{DAO: *NewDAO(rt)}
	allowance, err := q.Allowance([]byte{1}, []byte{2})
	require.NoError(t, err)
	require.Nil(t, allowance)
}

func TestDaoMutation_Transfer_Disabled(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	members := []Member{
		{Account: []byte{1}, Balance: big.NewInt(100).Bytes()},
		{Account: []byte{2}, Balance: big.NewInt(100).Bytes()},
	}
	_ = m.Init(members, false, []byte{3}, nil)

	rt.caller = []byte{1}
	err := m.Transfer([]byte{2}, big.NewInt(10).Bytes())
	require.ErrorIs(t, err, ErrTransferDisabled)
}

func TestDaoMutation_Transfer_Success(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{3}
	members := []Member{
		{Account: []byte{1}, Balance: big.NewInt(100).Bytes()},
		{Account: []byte{2}, Balance: big.NewInt(100).Bytes()},
	}
	_ = m.Init(members, false, sudo, nil)

	// Enable transfer
	_ = m.transferEnabled.Set(rt.txn, singletonKey, true)

	rt.caller = []byte{1}
	err := m.Transfer([]byte{2}, big.NewInt(10).Bytes())
	require.NoError(t, err)

	q := DaoQuery{DAO: m.DAO}
	balance1, _ := q.BalanceOf([]byte{1})
	balance2, _ := q.BalanceOf([]byte{2})
	require.Equal(t, 0, cmp(balance1, big.NewInt(90).Bytes()))
	require.Equal(t, 0, cmp(balance2, big.NewInt(110).Bytes()))
}

func TestDaoMutation_Transfer_LowBalance(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{3}
	members := []Member{
		{Account: []byte{1}, Balance: big.NewInt(100).Bytes()},
		{Account: []byte{2}, Balance: big.NewInt(100).Bytes()},
	}
	_ = m.Init(members, false, sudo, nil)

	_ = m.transferEnabled.Set(rt.txn, singletonKey, true)

	rt.caller = []byte{1}
	err := m.Transfer([]byte{2}, big.NewInt(200).Bytes())
	require.ErrorIs(t, err, ErrLowBalance)
}

func TestDaoMutation_Transfer_MemberNotExisted(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{3}
	members := []Member{{Account: []byte{1}, Balance: big.NewInt(100).Bytes()}}
	_ = m.Init(members, false, sudo, nil)

	_ = m.transferEnabled.Set(rt.txn, singletonKey, true)

	rt.caller = []byte{1}
	err := m.Transfer([]byte{2}, big.NewInt(10).Bytes()) // member 2 not exists
	require.ErrorIs(t, err, ErrMemberNotExisted)
}

func TestDaoMutation_Approve_Success(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	members := []Member{{Account: []byte{1}, Balance: big.NewInt(100).Bytes()}}
	_ = m.Init(members, false, []byte{2}, nil)

	rt.caller = []byte{1}
	err := m.Approve([]byte{2}, big.NewInt(50).Bytes())
	require.NoError(t, err)

	q := DaoQuery{DAO: m.DAO}
	allowance, err := q.Allowance([]byte{1}, []byte{2})
	require.NoError(t, err)
	require.Equal(t, 0, cmp(allowance, big.NewInt(50).Bytes()))
}

func TestDaoMutation_TransferFrom_InsufficientAllowance(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{3}
	members := []Member{
		{Account: []byte{1}, Balance: big.NewInt(100).Bytes()},
		{Account: []byte{2}, Balance: big.NewInt(100).Bytes()},
	}
	_ = m.Init(members, false, sudo, nil)

	_ = m.transferEnabled.Set(rt.txn, singletonKey, true)

	rt.caller = []byte{1}
	_ = m.Approve([]byte{2}, big.NewInt(10).Bytes())

	rt.caller = []byte{2}
	err := m.TransferFrom([]byte{1}, []byte{2}, big.NewInt(20).Bytes())
	require.ErrorIs(t, err, ErrInsufficientAllowance)
}

func TestDaoMutation_TransferFrom_Success(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{3}
	members := []Member{
		{Account: []byte{1}, Balance: big.NewInt(100).Bytes()},
		{Account: []byte{2}, Balance: big.NewInt(100).Bytes()},
	}
	_ = m.Init(members, false, sudo, nil)

	_ = m.transferEnabled.Set(rt.txn, singletonKey, true)

	rt.caller = []byte{1}
	_ = m.Approve([]byte{2}, big.NewInt(50).Bytes())

	rt.caller = []byte{2}
	err := m.TransferFrom([]byte{1}, []byte{2}, big.NewInt(10).Bytes())
	require.NoError(t, err)

	q := DaoQuery{DAO: m.DAO}
	allowance, _ := q.Allowance([]byte{1}, []byte{2})
	require.Equal(t, 0, cmp(allowance, big.NewInt(40).Bytes()))

	balance1, _ := q.BalanceOf([]byte{1})
	balance2, _ := q.BalanceOf([]byte{2})
	require.Equal(t, 0, cmp(balance1, big.NewInt(90).Bytes()))
	require.Equal(t, 0, cmp(balance2, big.NewInt(110).Bytes()))
}

func TestDaoMutation_TransferFrom_Disabled(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	members := []Member{
		{Account: []byte{1}, Balance: big.NewInt(100).Bytes()},
		{Account: []byte{2}, Balance: big.NewInt(100).Bytes()},
	}
	_ = m.Init(members, false, []byte{3}, nil)

	rt.caller = []byte{2}
	err := m.TransferFrom([]byte{1}, []byte{2}, big.NewInt(10).Bytes())
	require.ErrorIs(t, err, ErrTransferDisabled)
}
