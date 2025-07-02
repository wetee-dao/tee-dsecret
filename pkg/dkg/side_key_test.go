package dkg

import (
	"encoding/hex"
	"testing"
)

func TestGenerateSideKey(t *testing.T) {
	shares, pub, err := NewSr25519Split(10, 6)
	if err != nil {
		t.Fatal(err)
	}

	key1, err := CombineSr25519(shares[:5])
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(pub.ToBytes()) == hex.EncodeToString(key1.Public()) {
		t.Fatal("key1 except error, but get success")
	}

	key2, err := CombineSr25519(shares[:6])
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(pub.ToBytes()) != hex.EncodeToString(key2.Public()) {
		t.Fatal("key2 except success")
	}
}
