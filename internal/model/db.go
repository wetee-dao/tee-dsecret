package model

import (
	"bytes"
	"encoding/json"
	"fmt"

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

func (db *DB) NewTransaction(update bool) *Txn {
	return &Txn{in: db.DB.NewBatch()}
}

func NewDB() error {
	// Open DB
	db, err := pebble.Open(dbPath, &pebble.Options{})
	if err != nil {
		return err
	}

	// Create a new DB instance and initialize with DB
	dbInstance := &DB{}
	dbInstance.DB = db

	DBINS = dbInstance

	return nil
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

	return DBINS.Set([]byte(namespace+"__"+key), val, pebble.NoSync)
}

func GetKey(namespace, key string) ([]byte, error) {
	fmt.Println("xxxxxxxxxxxxxxxxxxx")
	value, _, err := DBINS.Get([]byte(namespace + "__" + key))
	fmt.Println("xxxxxxxxxxxxxxxxxxx2")
	fmt.Println("xxxxxxxxxxxxxxxxxxx3")
	if err != nil {
		return nil, err
	}

	return Unseal(value, nil)
}

func GetJson[T any](namespace, key string) (*T, error) {
	v, err := GetKey(namespace, key)
	if err != nil {
		return nil, err
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

func GetAbciMessageList[T any](namespace, key string) (list []T, err error) {
	iter, err := DBINS.NewIter(&pebble.IterOptions{
		LowerBound: []byte(namespace + "__" + key),
		// UpperBound: []byte("prefix_upper_bound"),
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
			list = append(list, *val)
		}
	}

	return
}

func GetAbciMessage[T any](namespace, key string) (*T, error) {
	v, err := GetKey(namespace, key)
	if err != nil {
		return nil, err
	}

	value, err := Unseal(v, nil)
	if err != nil {
		return nil, err
	}

	val := new(T)
	err = protoio.ReadMessage(bytes.NewBuffer(value), val)
	return val, err
}

func SetAbciMessage[T proto.Message](namespace, key string, value T) error {
	buf := new(bytes.Buffer)
	err := types.WriteMessage(value, buf)
	if err != nil {
		return err
	}
	return SetKey(namespace, key, buf.Bytes())
}
