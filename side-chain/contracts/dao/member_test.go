package dao

import (
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

const testDBSubdir = "chain_data/wetee"

// testRuntime implements model.ContractApi for tests.
type testRuntime struct {
	height int64
	caller []byte
	txn    *model.Txn
}

func (r *testRuntime) GetHeight() int64   { return r.height }
func (r *testRuntime) GetTxn() *model.Txn { return r.txn }
func (r *testRuntime) GetCaller() []byte  { return r.caller }

// setupTestDB creates a temporary database for testing.
func setupTestDB(t *testing.T) (rt *testRuntime, cleanup func()) {
	t.Helper()
	tmp := t.TempDir()
	oldWD, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmp))
	require.NoError(t, os.MkdirAll(testDBSubdir, 0o755))

	db, err := model.NewDB()
	require.NoError(t, err)

	txn := db.NewTransaction()
	rt = &testRuntime{txn: txn}

	cleanup = func() {
		_ = txn.Rollback()
		_ = db.Close()
		model.DBINS = nil
		_ = os.Chdir(oldWD)
	}
	return rt, cleanup
}

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
	enabled, err := d.PublicJoin()
	require.NoError(t, err)
	require.False(t, enabled)
}

func TestDaoMutation_Init_Basic(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	members := []Member{
		{Account: []byte{1}, Balance: big.NewInt(100).Bytes()},
		{Account: []byte{2}, Balance: big.NewInt(200).Bytes()},
	}
	err := m.Init(members, true, []byte{1}, nil)
	require.NoError(t, err)

	// Check total supply
	supply, err := DaoQuery{DAO: m.DAO}.TotalSupply()
	require.NoError(t, err)
	require.Equal(t, 0, cmp(supply, big.NewInt(300).Bytes()))

	// Check public join enabled
	enabled, err := DaoQuery{DAO: m.DAO}.PublicJoin()
	require.NoError(t, err)
	require.True(t, enabled)
}

func TestDaoMutation_Init_WithDefaultTrack(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	track := TrackData{
		Name:            "test",
		PreparePeriod:   10,
		MaxDeciding:     100,
		ConfirmPeriod:   20,
		DecisionPeriod:  30,
		DecisionDeposit: big.NewInt(50).Bytes(),
		MaxBalance:      big.NewInt(1000).Bytes(),
	}
	err := m.Init([]Member{}, false, []byte{1}, &track)
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
	_ = m.Init([]Member{}, false, sudo, nil)

	rt.caller = sudo
	err := m.SetPublicJoin(true)
	require.NoError(t, err)

	enabled, err := DaoQuery{DAO: m.DAO}.PublicJoin()
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
	_ = m.Init([]Member{}, false, sudo, nil)

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
	members := []Member{{Account: []byte{2}, Balance: big.NewInt(100).Bytes()}}
	_ = m.Init(members, false, sudo, nil)

	rt.caller = sudo
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
	members := []Member{{Account: []byte{2}, Balance: big.NewInt(100).Bytes()}}
	_ = m.Init(members, false, sudo, nil)

	rt.caller = []byte{2}
	err := m.Leave()
	require.ErrorIs(t, err, ErrMemberBalanceNotZero)
}

func TestDaoMutation_Leave_Success(t *testing.T) {
	rt, cleanup := setupTestDB(t)
	defer cleanup()

	m := DaoMutation{DAO: *NewDAO(rt)}
	sudo := []byte{1}
	members := []Member{{Account: []byte{2}, Balance: nil}}
	_ = m.Init(members, false, sudo, nil)

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
	_ = m.Init([]Member{}, true, sudo, nil)

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
	members := []Member{{Account: []byte{2}, Balance: big.NewInt(100).Bytes()}}
	_ = m.Init(members, false, sudo, nil)

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
	members := []Member{{Account: []byte{2}, Balance: big.NewInt(100).Bytes()}}
	_ = m.Init(members, false, sudo, nil)

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
