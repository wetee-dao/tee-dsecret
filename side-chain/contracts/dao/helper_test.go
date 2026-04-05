package dao

import (
	"math/big"
	"testing"
)

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
