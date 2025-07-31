package model

import (
	"errors"
	"fmt"

	"github.com/cockroachdb/pebble"
	"github.com/wetee-dao/tee-dsecret/pkg/util"
)

func AddToList(namespace string, key string, val []byte) error {
	indexBt, err := GetKey(namespace, "index_"+key)
	if err != nil && !errors.Is(err, pebble.ErrNotFound) {
		return err
	}

	var index uint64 = 0
	if indexBt != nil {
		index = util.BytesToUint64(indexBt)
	}

	txn := DBINS.NewTransaction()
	txn.SetKey(namespace, "index_"+key, util.Uint64ToBytes(index+1))
	txn.SetKey(namespace, key+fmt.Sprint(index), val)
	return txn.Commit()
}

func GetList(namespace string, subkey string, page int, size int) ([][]byte, error) {
	key := []byte(comboKey(namespace, subkey))
	iter, err := DBINS.NewIter(&pebble.IterOptions{
		LowerBound: key,
		UpperBound: keyUpperBound(key),
	})
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var list [][]byte
	for iter.First(); iter.Valid(); iter.Next() {
		v := iter.Value()
		value, err := util.Unseal(v, nil)
		if err != nil {
			return nil, err
		}

		list = append(list, value)
	}

	return list, nil
}

func DeleteList(namespace string, key string) error {
	return DeletekeysByPrefix(namespace, key)
}
