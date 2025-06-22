package model

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/cockroachdb/pebble"
	"github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/gogoproto/proto"
	"wetee.app/dsecret/internal/model/protoio"
)

var DBINS *DB

const (
	dbPath = "./chain_data/wetee"
)

type DB struct {
	*pebble.DB
}

func (db *DB) NewTransaction() *Txn {
	return &Txn{in: db.DB.NewIndexedBatch()}
}

func NewDB() (*DB, error) {
	// Open DB
	db, err := pebble.Open(dbPath, &pebble.Options{})
	if err != nil {
		return nil, err
	}

	// Create a new DB instance and initialize with DB
	dbInstance := &DB{}
	dbInstance.DB = db

	DBINS = dbInstance

	return DBINS, nil
}

func Set(key string, value []byte) error {
	return SetKey("", key, value)
}

func Get(key string) ([]byte, error) {
	return GetKey("", key)
}

func SetKey(namespace, key string, value []byte) error {
	val, err := SealWithProductKey(value, nil)
	if err != nil {
		return err
	}

	return DBINS.Set([]byte(namespace+"_"+key), val, pebble.Sync)
}

func GetKey(namespace, key string) ([]byte, error) {
	value, _, err := DBINS.Get([]byte(namespace + "_" + key))
	if err != nil {
		return nil, err
	}

	return Unseal(value, nil)
}

func GetJson[T any](namespace, key string) (*T, error) {
	v, err := GetKey(namespace, key)
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

func SetJson[T any](namespace, key string, val *T) error {
	bt, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return SetKey(namespace, key, bt)
}

func keyUpperBound(b []byte) []byte {
	end := make([]byte, len(b))
	copy(end, b)
	for i := len(end) - 1; i >= 0; i-- {
		end[i] = end[i] + 1
		if end[i] != 0 {
			return end[:i+1]
		}
	}
	return nil
}

func GetProtoMessageList[T any](namespace, key string) (list []*T, err error) {
	rkey := []byte(namespace + "_" + key)
	iter, err := DBINS.NewIter(&pebble.IterOptions{
		LowerBound: rkey,
		UpperBound: keyUpperBound(rkey),
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

func GetProtoMessage[T any](namespace, key string) (*T, error) {
	v, err := GetKey(namespace, key)
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

func SetProtoMessage[T proto.Message](namespace, key string, value T) error {
	buf := new(bytes.Buffer)
	err := types.WriteMessage(value, buf)
	if err != nil {
		return err
	}
	return SetKey(namespace, key, buf.Bytes())
}

func DeleteKey(namespace, key string) error {
	return DBINS.Delete([]byte(namespace+"_"+key), pebble.Sync)
}

func DeletekeysByPrefix(namespace, key string) error {
	rkey := []byte(namespace + "_" + key)
	iter, err := DBINS.NewIter(&pebble.IterOptions{
		LowerBound: rkey,
		UpperBound: keyUpperBound(rkey),
	})
	if err != nil {
		return err
	}
	defer iter.Close()

	txn := DBINS.DB.NewBatch()
	for iter.First(); iter.Valid(); iter.Next() {
		txn.Delete(iter.Key(), pebble.Sync)
	}

	return txn.Commit(pebble.Sync)
}
