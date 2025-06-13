package model

import (
	"github.com/cockroachdb/pebble"
)

type Txn struct {
	in *pebble.Batch
}

func (txn *Txn) Set(key, value []byte) error {
	val, err := SealWithProductKey(value, nil)
	if err != nil {
		return err
	}

	return txn.in.Set(key, val, pebble.Sync)
}

func (txn *Txn) Get(key []byte) ([]byte, error) {
	v, _, err := txn.in.Get(key)
	if err != nil {
		return nil, err
	}

	return Unseal(v, nil)
}

func (txn *Txn) Commit() error {
	return txn.in.Commit(pebble.Sync)
}
