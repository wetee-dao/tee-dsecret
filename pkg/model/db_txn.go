package model

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/cockroachdb/pebble"
	"github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/wetee-dao/tee-dsecret/pkg/model/protoio"
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

func (txn *Txn) Delete(key []byte) error {
	return txn.in.Delete(key, pebble.Sync)
}

func (txn *Txn) DeletekeysByPrefix(prefix []byte) error {
	iter, err := DBINS.NewIter(&pebble.IterOptions{
		LowerBound: prefix,
		UpperBound: keyUpperBound(prefix),
	})
	if err != nil {
		return err
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		if err := txn.Delete(iter.Key()); err != nil {
			return err
		}
	}

	return nil
}

func TxnGetJson[T any](txn *Txn, key []byte) (*T, error) {
	v, err := txn.Get(key)
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if len(v) == 0 {
		return new(T), nil
	}

	val := new(T)
	err = json.Unmarshal(v, val)

	return val, err
}

func TxnSetJson[T any](txn *Txn, key []byte, val *T) error {
	bt, err := json.Marshal(val)
	if err != nil {
		return err
	}

	// fmt.Println("set", key, string(bt))
	return txn.Set(key, bt)
}

func TxnGetProtoMessageList[T any](txn *Txn, key []byte) (list []*T, err error) {
	iter, err := DBINS.NewIter(&pebble.IterOptions{
		LowerBound: key,
		UpperBound: keyUpperBound(key),
	})
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		v := iter.Value()
		value, err := Unseal(v, nil)
		if err != nil {
			return nil, err
		}

		val := new(T)
		err = protoio.ReadMessage(bytes.NewBuffer(value), val)
		if err == nil {
			list = append(list, val)
		}
	}

	return
}

func TxnGetProtoMessage[T any](txn *Txn, key []byte) (*T, error) {
	v, err := txn.Get(key)
	if err != nil {
		return nil, err
	}

	if len(v) == 0 {
		return nil, nil
	}

	value, err := Unseal(v, nil)
	if err != nil {
		return nil, err
	}

	val := new(T)
	err = protoio.ReadMessage(bytes.NewBuffer(value), val)
	return val, err
}

func TxnSetProtoMessage[T proto.Message](txn *Txn, key []byte, value T) error {
	buf := new(bytes.Buffer)
	err := types.WriteMessage(value, buf)
	if err != nil {
		return err
	}
	return txn.Set(key, buf.Bytes())
}

func (txn *Txn) Commit() error {
	return txn.in.Commit(pebble.Sync)
}
