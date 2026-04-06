package dao

import (
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

// testDBSubdir is the subdirectory where the test database is stored.
const testDBSubdir = "chain_data/wetee"

// testRuntime implements model.ContractApi for tests.
type testRuntime struct {
	height      int64
	caller      []byte
	txn         *model.Txn
	sudoAccount []byte
}

func (r *testRuntime) GetHeight() int64       { return r.height }
func (r *testRuntime) GetTxn() *model.Txn     { return r.txn }
func (r *testRuntime) GetCaller() []byte      { return r.caller }
func (r *testRuntime) GetSudoAccount() []byte { return r.sudoAccount }

// setupTestDB creates a temporary database for testing with an initialized runtime.
// It returns a testRuntime with the database transaction set up.
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

func TestDecodeEncodeAmount(t *testing.T) {
	if decodeAmount(nil).Sign() != 0 {
		t.Fatal("nil -> 0")
	}
	if decodeAmount([]byte{}).Sign() != 0 {
		t.Fatal("empty -> 0")
	}
	v := big.NewInt(12345)
	b := encodeAmount(v)
	if decodeAmount(b).Cmp(v) != 0 {
		t.Fatalf("roundtrip: got %s", decodeAmount(b).String())
	}
}

func TestAddSubCmp(t *testing.T) {
	a := big.NewInt(100).Bytes()
	b := big.NewInt(30).Bytes()
	if cmp(add(a, b), big.NewInt(130).Bytes()) != 0 {
		t.Fatal("add")
	}
	if cmp(sub(a, b), big.NewInt(70).Bytes()) != 0 {
		t.Fatal("sub")
	}
	if sub(b, a) != nil {
		t.Fatal("sub underflow -> nil")
	}
	if !isZero(nil) || !isZero([]byte{}) {
		t.Fatal("isZero")
	}
	if isZero(big.NewInt(1).Bytes()) {
		t.Fatal("isZero non-zero")
	}
}

func TestAllowanceKeyStable(t *testing.T) {
	k := allowanceKey([]byte{0xab}, []byte{0xcd})
	if k != "ab:cd" {
		t.Fatalf("got %q", k)
	}
}

func TestCloneBytes(t *testing.T) {
	if cloneBytes(nil) != nil {
		t.Fatal("nil clone")
	}
	x := []byte{1, 2}
	y := cloneBytes(x)
	x[0] = 9
	if y[0] == 9 {
		t.Fatal("deep copy expected")
	}
}
